package handler

import (
	"archive/tar"
	"archive/zip"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/config"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/logging"
	quizparser "github.com/egot3/fathom/internal/quizParser"
	"github.com/google/uuid"
	"github.com/zeebo/xxh3"
	"go.yaml.in/yaml/v4"
)

// DeleteQuiz implements [Service].
func (c *chiService) DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	quizUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		logger.Error("Bad uuid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}

	ctx = logging.WithLogger(ctx, logger.With(
		slog.String("quizUUID", quizUUID.String()),
	))

	err := c.quizRepo.DeallocateQuiz(ctx, quizUUID)
	if err != nil {
		logger.Error("unable to deallocate quiz",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Quiz not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if s := r.Header.Get("force-delete"); s == "true" {
		quizPath, err := c.quizRepo.QuizPath(ctx, quizUUID)
		if err != nil {
			logger.Error("unable to retrieve quizPath",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusMultiStatus)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to get test's path internally"})
			return
		}

		err = os.Remove(quizPath)
		if err != nil {
			logger.Error("unable to delete quiz by quizPath",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusMultiStatus)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to delete quiz"})
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetQuiz implements [Service].
// teacher-only path
func (c *chiService) GetQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	quizUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		logger.Error("Bad uuid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}

	ctx = logging.WithLogger(ctx, logger.With(
		slog.String("quizUUID", quizUUID.String()),
	))

	quizPath, err := c.quizRepo.QuizPath(ctx, quizUUID)
	if err != nil {
		logger.Error("unable to retrieve quizPath",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Quiz not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf, err := os.ReadFile(quizPath)
	if err != nil {
		logger.Error("Failed to read file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	quiz, err := quizparser.ParseQuizByBytes(buf)
	if err != nil {
		logger.Error("Unable to retrieve quiz by path",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.GetQuizResponse{Quiz: *quiz})
}

// ListQuiz implements [Service].
func (c *chiService) ListQuizzes(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Failed to parse form data"})
		return
	}

	pageInt, err := strconv.Atoi(r.Form.Get("page"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "given form page is not a number"})
		return
	}
	sizeInt, err := strconv.Atoi(r.Form.Get("size"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "given form size is not a number"})
		return
	}
	if sizeInt <= 0 {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "size can't be <= 0"})
		return
	}

	logger.With(slog.Int("page", pageInt), slog.Int("size", sizeInt))

	quizzes, total, err := c.quizRepo.ListQuizzes(ctx, pageInt, sizeInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Quizzes not found"})
			return
		}
		if gone, ok := errors.AsType[carefulness.Gone](err); ok {
			w.WriteHeader(http.StatusGone)
			json.NewEncoder(w).Encode(gone.JSONError())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.ListQuizResponse{
		Quizzes: quizzes,
		Total:   total,
		Page:    pageInt,
		Size:    sizeInt,
	})
}

// PostQuiz implements [Service].
func (c *chiService) PostQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.PostQuizRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("error in register during reading",
			slog.String("error", err.Error()),
		)
		if errors.Is(err, carefulness.ErrMalformedRequest) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest.JSONError()) // no err check as рукописи не горят

			return
		}
		if errors.Is(err, carefulness.ErrUnprocessableRequest) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.ErrUnprocessableRequest.JSONError())

			return
		}
		if errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Empty body"})

			return
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Data loss"})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Name == "" || req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("request contains no quiz's body or its name",
			slog.String("name", req.Name),
			slog.Int("body len", len(req.Body)),
		)

		return
	}

	logger = logger.With(slog.String("name", req.Name))
	ctx = logging.WithLogger(ctx, logger)

	frontmatter, err := yaml.Marshal(req.Meta)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("unable to process frontmatter",
			slog.String("Error", err.Error()),
		)

		return
	}
	buf := fmt.Appendf(nil, `
	---
	%v
	---
	%v
	`, string(frontmatter), req.Body)

	_, err = quizparser.ParseQuizByBytes(buf)
	if err != nil {
		logger.Error("couldn't parse quiz", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()})
		return
	}

	abs := config.TurnToAbs(req.Name)
	checksumUint := xxh3.HashString(req.Body)
	checksum := binary.BigEndian.AppendUint64(nil, checksumUint)
	err = c.quizRepo.RegisterQuiz(ctx, abs, checksum, req.Meta.Score)
	if err != nil {
		logger.Error("couldn't register quiz", slog.String("Error", err.Error()))
		if conflict, ok := errors.AsType[carefulness.Conflict](err); ok {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(conflict.JSONError())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't register quiz"})
		return
	}

	err = os.WriteFile(abs, buf, os.ModeAppend)
	if err != nil {
		w.WriteHeader(http.StatusMultiStatus)

		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to write file"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PutQuiz implements [Service].
func (c *chiService) PutQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	quizUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		logger.Error("Bad uuid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}
	logger = logger.With(slog.String("quizUUID", quizUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.PutQuizRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("error in register during reading",
			slog.String("error", err.Error()),
		)
		if errors.Is(err, carefulness.ErrMalformedRequest) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest.JSONError()) // no err check as рукописи не горят

			return
		}
		if errors.Is(err, carefulness.ErrUnprocessableRequest) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.ErrUnprocessableRequest.JSONError())

			return
		}
		if errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Empty body"})

			return
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Data loss"})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	abs, err := c.quizRepo.QuizPath(ctx, quizUUID)
	if err != nil {
		logger.Error("couldn't select quiz path",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	frontmatter, err := yaml.Marshal(req.Meta)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("unable to process frontmatter",
			slog.String("Error", err.Error()),
		)

		return
	}
	buf := fmt.Appendf(nil, `
	---
	%v
	---
	%v
	`, string(frontmatter), req.Body)

	_, err = quizparser.ParseQuizByBytes(buf)
	if err != nil {
		logger.Error("couldn't parse quiz", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()})
		return
	}

	checksumUint := xxh3.HashString(req.Body)
	checksum := binary.BigEndian.AppendUint64(nil, checksumUint)
	err = c.quizRepo.UpdateChecksum(ctx, quizUUID, checksum)
	if err != nil {
		logger.Error("couldn't update quiz", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't register quiz"})
		return
	}

	err = os.WriteFile(abs, buf, os.ModeAppend)
	if err != nil {
		w.WriteHeader(http.StatusMultiStatus)

		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to write file"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *chiService) PatchQuiz(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	quizUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		logger.Error("Bad uuid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}
	logger = logger.With(slog.String("quizUUID", quizUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.PatchQuizRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("error in register during reading",
			slog.String("error", err.Error()),
		)
		if errors.Is(err, carefulness.ErrMalformedRequest) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest.JSONError()) // no err check as рукописи не горят

			return
		}
		if errors.Is(err, carefulness.ErrUnprocessableRequest) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.ErrUnprocessableRequest.JSONError())

			return
		}
		if errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Empty body"})

			return
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Data loss"})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	abs, err := c.quizRepo.QuizPath(ctx, quizUUID)
	if err != nil {
		logger.Error("couldn't select quiz path",
			slog.String("Error", err.Error()),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var newAbs *string
	if req.Name != nil {
		v := config.TurnToAbs(*req.Name)
		newAbs = &v
	}

	err = c.quizRepo.PatchQuiz(ctx, quizUUID, newAbs, req.Score)
	if err != nil {
		logger.Error("couldn't update quiz", slog.String("Error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "couldn't register quiz"})
		return
	}

	if newAbs != nil {
		err := os.Rename(abs, *newAbs)
		if err != nil {
			w.WriteHeader(http.StatusMultiStatus)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()})
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *chiService) Export(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	var req contracts.ExportQuizRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("error in register during reading",
			slog.String("error", err.Error()),
		)
		if errors.Is(err, carefulness.ErrMalformedRequest) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest.JSONError()) // no err check as рукописи не горят

			return
		}
		if errors.Is(err, carefulness.ErrUnprocessableRequest) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.ErrUnprocessableRequest.JSONError())

			return
		}
		if errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Empty body"})

			return
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Data loss"})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acceptH := strings.Split(strings.ReplaceAll(r.Header.Get("Accept"), " ", ""), ",")
	acceptStarts := make([]string, len(acceptH))
	acceptQuality := make([]float32, len(acceptH))
	for i, accept := range acceptH {
		a := strings.Split(strings.TrimSpace(accept), ";q=")
		acceptStarts[i] = a[0]
		if len(a) > 1 {
			f64, err := strconv.ParseFloat(a[1], 32)
			if err != nil {
				logger.Error("Error during accept parsing",
					slog.String("Error", err.Error()),
				)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(carefulness.JSONError{Error: "bad accept header"})
				return
			}
			acceptQuality[i] = float32(f64)
		} else {
			acceptQuality[i] = 1
		}
	}

	for len(acceptStarts) > 0 {
		idx := slices.Index(acceptQuality, slices.Max(acceptQuality))
		switch acceptStarts[idx] {
		case "application/zip":
			w.Header().Set("Content-Type", "application/zip")
			zipWriter := zip.NewWriter(w)
			logger = logger.With(slog.String("strategy", "zip"))

			for _, quizUUID := range req.UUIDs {
				path, err := c.quizRepo.QuizPath(ctx, quizUUID)
				if err != nil {
					logger.Error("couldn't get path", slog.String("error", err.Error()))
					/* if errors.Is(err, sql.ErrNoRows) {
						w.WriteHeader(http.StatusNotFound)
						json.NewEncoder(w).Encode(carefulness.JSONError{Error: fmt.Sprintf("%v not found", quizUUID.String())})
						return
					}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: fmt.Sprintf("unable to process %v", quizUUID.String())})
					return */
					break
				}

				f, err := os.Open(path)
				if err != nil {
					logger.Error("couldn't open file", slog.String("error", err.Error()))
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to read quiz file"})
					return */
					break
				}
				defer f.Close()

				fileWriter, err := zipWriter.Create(filepath.Base(path))
				if err != nil {
					logger.Error("couldn't create entry in zip", slog.String("error", err.Error()))
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to create archive entry"})
					return */
					break
				}
				if _, err := io.Copy(fileWriter, f); err != nil {
					logger.Error("couldn't write file to zip", slog.String("error", err.Error()))
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to write quiz to archive"})
					return */
					break
				}
				f.Close()
			}

			if err := zipWriter.Close(); err != nil {
				logger.Error("couldn't close zip", slog.String("error", err.Error()))
				/* w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to finalize archive"})
				return */
				break
			}

			return
		case "application/tar":
			tarWriter := tar.NewWriter(w)

			logger = logger.With(slog.String("strategy", "tar"))
			w.Header().Set("Content-Type", "application/tar")

			var err error
			for _, quizUUID := range req.UUIDs {
				var path string
				path, err = c.quizRepo.QuizPath(ctx, quizUUID)
				if err != nil {
					logger.Error("couldn't get path for requested quiz",
						slog.String("Error", err.Error()),
					)
					if errors.Is(err, sql.ErrNoRows) {
						/* w.WriteHeader(http.StatusNotFound)
						json.NewEncoder(w).Encode(carefulness.JSONError{Error: fmt.Sprintf("%v not found", quizUUID.String())})
						return */
						break
					}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: fmt.Sprintf("unable to process %v", quizUUID.String())})
					return
				}

				f, err := os.Open(path)
				if err != nil {
					logger.Error("couldn't open file for requested quiz",
						slog.String("Error", err.Error()),
					)
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to read quiz file"})
					return */
					break
				}
				defer f.Close()

				info, err := f.Stat()
				if err != nil {
					logger.Error("couldn't get stats for file",
						slog.String("Error", err.Error()),
					)
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to get quiz's file stat"})
					return */
					break
				}

				header, err := tar.FileInfoHeader(info, "")
				if err != nil {
					logger.Error("couldn't get header for file",
						slog.String("Error", err.Error()),
					)
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to get quiz's file header"})
					return */
					break
				}

				if err := tarWriter.WriteHeader(header); err != nil {
					logger.Error("couldn't write quiz header to archive",
						slog.String("Error", err.Error()),
					)
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to create pass header to archive file"})
					return */
					break
				}

				if _, err = io.Copy(tarWriter, f); err != nil {
					logger.Error("couldn't write file for requested quiz",
						slog.String("Error", err.Error()),
					)
					/* w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(carefulness.JSONError{Error: "unable to write quiz file to archive"})
					return */
					break
				}
				f.Close()

			}

			err = tarWriter.Close()
			if err != nil {
				logger.Error("couldn't close file for requested quiz",
					slog.String("Error", err.Error()),
				)
				break
			}

			return
		default:
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		acceptStarts = slices.Delete(acceptStarts, idx, idx+1)
		acceptQuality = slices.Delete(acceptQuality, idx, idx+1)
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(carefulness.JSONError{Error: "None of the stratagy worked"})
}

func (c *chiService) Import(w http.ResponseWriter, r *http.Request) {
	panic("unimp")
}
