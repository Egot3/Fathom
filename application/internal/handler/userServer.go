package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/passwordutils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// DeleteUser implements [TestService].
func (c *chiTestService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	userUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}

	logger = logger.With(
		slog.String("userUUID", userUUID.String()),
	)
	ctx = logging.WithLogger(ctx, logger)

	err := c.userRepo.DeleteUser(ctx, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "User not found"})

			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUser implements [TestService].
func (c *chiTestService) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	userUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}

	logger.With(slog.String("userUUID", userUUID.String()))

	user, err := c.userRepo.User(ctx, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "User not found"})
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
	json.NewEncoder(w).Encode(contracts.GetUserResponse{
		User: *user,
	})
}

// Login implements [TestService].
func (c *chiTestService) Login(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
	logger.With(slog.String("nickname", req.Nickname))

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Hashing error"})
		return
	}

	user, err := c.userRepo.Login(ctx, req.Nickname, passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "No active user with this username"})
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newToken, err := jwtutils.GenerateToken(user.UUID, user.IsTeacher)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		cookie := &http.Cookie{
			Name:     "jwt_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Bad token"})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    newToken,
		Path:     "/",
		Expires:  time.Now().Add(jwtutils.JWTTTL),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contracts.RegisterResponse{
		User: *user,
	})
}

// PatchUser implements [TestService].
func (c *chiTestService) PatchUser(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	userUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to retrieve uuid"})
		return
	}
	var req contracts.PatchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
	logger.With(slog.String("uuid", userUUID.String()))

	var passwordHash []byte = nil
	if req.Password != nil {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Hashing error"})
			return
		}
	}

	var isTeacher = req.IsTeacher
	if isTeacher != nil {
		is, err := c.userRepo.IsTeacher(ctx, userUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Unable to authorize given user"})
			return
		}
		if !is {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	err = c.userRepo.UpdateUser(ctx, models.PatchUser{
		UUID:         userUUID,
		Nickname:     req.Nickname,
		PasswordHash: passwordHash,
		IsTeacher:    isTeacher,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Register implements [TestService].
func (c *chiTestService) Register(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
	logger.With(slog.String("nickname", req.Nickname))

	err = passwordutils.CheckPasswordSafety(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Hashing error"})
		return
	}

	user, err := c.userRepo.Register(ctx, req.Nickname, passwordHash)
	if err != nil {
		if conflict, ok := errors.AsType[*carefulness.Conflict](err); ok {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(conflict.JSONError())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newToken, err := jwtutils.GenerateToken(user.UUID, user.IsTeacher)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		cookie := &http.Cookie{
			Name:     "jwt_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Bad token"})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    newToken,
		Path:     "/",
		Expires:  time.Now().Add(jwtutils.JWTTTL),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contracts.RegisterResponse{
		User: *user,
	})
}

// ListUsers implements [TestService].
func (c *chiTestService) ListUsers(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.ListUsersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
	logger.With(slog.Int("page", req.Page), slog.Int("size", req.Size))

	if req.Size <= 0 {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(carefulness.JSONError{Error: "Size must be > 0"})
		return
	}

	users, total, err := c.userRepo.List(ctx, req.Page, req.Size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Error: "User not found"})
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
	json.NewEncoder(w).Encode(contracts.ListUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	})
}
