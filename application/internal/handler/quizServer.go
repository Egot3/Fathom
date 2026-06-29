package handler

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	acceptutils "github.com/egot3/fathom/internal/acceptUtils"
	archiveutlis "github.com/egot3/fathom/internal/archiveUtlis"
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
		switch {
		case errors.Is(err, carefulness.ErrMalformedRequest):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.ErrMalformedRequest.JSONError()) // no err check as рукописи не горят

		case errors.Is(err, carefulness.ErrUnprocessableRequest):
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(carefulness.ErrUnprocessableRequest.JSONError())

		case errors.Is(err, io.EOF):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Empty body"})

		case errors.Is(err, io.ErrUnexpectedEOF):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Data loss"})

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	accept, err := acceptutils.BestAccept(r.Header.Get("Accept"),
		"application/zip", "application/tar", "application/gzip",
	)
	if (err != nil) || (accept == "") {
		logger.Info("Got an unaccaptable accept header",
			slog.String("accept", accept),
			slog.String("Error", err.Error()),
		)

		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	type quizFile struct {
		uuid string
		path string
		fi   os.FileInfo // used by tar
	}
	var files []quizFile
	for _, uuid := range req.UUIDs {
		path, err := c.quizRepo.QuizPath(ctx, uuid)
		if err != nil {
			logger.Error("couldn't get path", "uuid", uuid, "error", err)
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(carefulness.JSONError{Error: fmt.Sprintf("%v not found", uuid)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(carefulness.JSONError{Error: fmt.Sprintf("unable to process %v", uuid)})
			}
			return
		}
		fi, err := os.Stat(path)
		if err != nil {
			logger.Error("quiz file not accessible", "path", path, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		files = append(files, quizFile{uuid: uuid.String(), path: path, fi: fi})
	}

	switch accept {
	case "application/zip":
		w.Header().Set("Content-Type", "application/zip")
		zipWriter := zip.NewWriter(w)
		logger = logger.With(slog.String("strategy", "zip"))

		for _, qf := range files {
			if err := archiveutlis.AddFileToZip(zipWriter, qf.path); err != nil {
				logger.Error("error writing to zip",
					slog.String("path", qf.path),
					slog.String("Error", err.Error()),
				)
				zipWriter.Close()
				return
			}
		}
		if err := zipWriter.Close(); err != nil {
			logger.Error("error finalising zip", "error", err)
		}
		return

	case "application/tar":
		logger = logger.With(slog.String("strategy", "tar"))
		w.Header().Set("Content-Type", "application/tar")
		tarWriter := tar.NewWriter(w)

		for _, qf := range files {
			if err := archiveutlis.AddFileToTar(tarWriter, qf.path, qf.fi); err != nil {
				logger.Error("error writing to tar",
					slog.String("path", qf.path),
					slog.String("Error", err.Error()),
				)

				tarWriter.Close()
				return
			}
		}
		if err := tarWriter.Close(); err != nil {
			logger.Error("error finalising tar", "error", err)
		}
		return

	case "application/gzip":
		logger = logger.With(slog.String("strategy", "tar.gz"))
		w.Header().Set("Content-Type", "application/tar")

		gzipWriter := gzip.NewWriter(w)
		tarWriter := tar.NewWriter(gzipWriter)

		for _, qf := range files {
			if err := archiveutlis.AddFileToTar(tarWriter, qf.path, qf.fi); err != nil {
				logger.Error("error writing to tar",
					slog.String("path", qf.path),
					slog.String("Error", err.Error()),
				)

				tarWriter.Close()
				return
			}
		}

		if err := gzipWriter.Close(); err != nil {
			logger.Error("error finalising gz", "error", err)
		}
		if err := tarWriter.Close(); err != nil {
			logger.Error("error finalising tar", "error", err)
		}
		return

	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return

	}
}

func (c *chiService) Import(w http.ResponseWriter, r *http.Request) {
	panic("unimp")
}
