package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"

	jwtutils "github.com/egot3/fathom/internal/JWTutils"
	"github.com/egot3/fathom/internal/carefulness"
	"github.com/egot3/fathom/internal/contracts"
	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/logging"
	"github.com/egot3/fathom/internal/models"
	"github.com/egot3/fathom/internal/passwordutils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// DeleteUser implements [UserService].
func (c *chiService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	userUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		logger.Error("Bad uuid")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve uuid"})
		return
	}

	logger = logger.With(
		slog.String("userUUID", userUUID.String()),
	)
	ctx = logging.WithLogger(ctx, logger)

	err := c.userRepo.DeleteUser(ctx, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("Delete user found no rows")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "User not found"})

			return
		}
		logger.Error("Delete user got an unexpected error",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUser implements [UserService].
func (c *chiService) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	userUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve uuid"})
		return
	}

	logger = logger.With(slog.String("userUUID", userUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	user, err := c.userRepo.User(ctx, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "User not found"})
			return
		}
		if gone, ok := errors.AsType[carefulness.Gone](err); ok {
			w.WriteHeader(http.StatusGone)
			json.NewEncoder(w).Encode(gone.JSONError())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("user encoding", slog.Any("user", user))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.GetUserResponse{
		User: user,
	})
}

// Login implements [UserService].
func (c *chiService) Login(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.LoginRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}

	if req.Nickname == "" || len(req.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Nickname or password are empty"})
		return
	}

	logger = logger.With(slog.String("nickname", req.Nickname))
	ctx = logging.WithLogger(r.Context(), logger)

	user, err := c.userRepo.Login(ctx, req.Nickname, []byte(req.Password))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("login user not found",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "No active user with this username"})
			return
		}
		logger.Error("unexpected login db error",
			slog.String("Error", err.Error()),
		)
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
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Bad token"})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    newToken,
		Path:     "/",
		Expires:  time.Now().Add(jwtutils.JWTTTL),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	mAge := int(math.Trunc(jwtutils.JWTTTL.Seconds()))
	w.Header().Set("Session-Control", fmt.Sprintf("max-age=%d", mAge))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.RegisterResponse{
		User: user,
	})
}

// PatchUser implements [UserService].
func (c *chiService) PatchUser(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	userUUID, ok := (r.Context().Value("uuid")).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to retrieve uuid"})
		return
	}

	var req contracts.PatchRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}
	logger = logger.With(slog.String("uuid", userUUID.String()))
	ctx = logging.WithLogger(ctx, logger)

	if req.Nickname == nil && req.IsTeacher == nil && req.Password == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var passwordHash []byte = nil
	if req.Password != nil {
		err := passwordutils.CheckPasswordSafety(*req.Password)
		if err != nil {
			logger.Error("Unsafe password",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()})
			return
		}
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("Error while parsing hash",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Hashing error"})
			return
		}
	}

	var isTeacher = req.IsTeacher
	if isTeacher != nil {
		is, err := c.userRepo.IsTeacher(ctx, userUUID)
		if err != nil {
			logger.Error("Error while checking if user is teacher",
				slog.String("Error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Unable to authorize given user"})
			return
		}
		if !is {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	err := c.userRepo.UpdateUser(ctx, models.PatchUser{
		UUID:         userUUID,
		Nickname:     req.Nickname,
		PasswordHash: passwordHash,
		IsTeacher:    isTeacher,
	})
	if err != nil {
		logger.Error("Error while updating user",
			slog.String("Error", err.Error()),
		)
		if conflict, ok := errors.AsType[carefulness.Conflict](err); ok {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(conflict.JSONError())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Register implements [UserService].
func (c *chiService) Register(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")
	var req contracts.RegisterRequest
	jerr := httputils.ParseJSON(r.Body, &req)
	if jerr != nil {
		logger.Error("Failed to parse body",
			slog.String("Error", jerr.Error()),
		)
		jerr.Encode(w)
		return
	}
	if req.Nickname == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("request contains no password/nickname",
			slog.String("nickname", req.Nickname),
			slog.Bool("pswd>0", len(req.Password) > 0),
		)

		return
	}

	logger = logger.With(slog.String("nickname", req.Nickname))

	err := passwordutils.CheckPasswordSafety(req.Password)
	if err != nil {
		logger.Error("bad password",
			slog.String("password", req.Password),
			slog.String("because", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("couldn't hash password",
			slog.String("error", err.Error()),
		)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Hashing error"})
		return
	}

	user, err := c.userRepo.Register(ctx, req.Nickname, passwordHash)
	if err != nil {
		logger.Error("couldn't register user",
			slog.String("error", err.Error()),
		)
		if conflict, ok := errors.AsType[carefulness.Conflict](err); ok {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(conflict.JSONError())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newToken, err := jwtutils.GenerateToken(user.UUID, user.IsTeacher)
	if err != nil {
		logger.Error("couldn't generate token",
			slog.String("error", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Couldn't generate token"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    newToken,
		Path:     "/",
		Expires:  time.Now().Add(jwtutils.JWTTTL),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	mAge := int(math.Trunc(jwtutils.JWTTTL.Seconds()))
	w.Header().Set("Session-Control", fmt.Sprintf("max-age=%d", mAge))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contracts.RegisterResponse{
		User: user,
	})
}

// ListUsers implements [UserService].
func (c *chiService) ListUsers(w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context()).With(
		slog.String("layer", "handler"),
	)
	ctx := logging.WithLogger(r.Context(), logger)

	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(carefulness.JSONError{Err: "Failed to parse form data"})
		return
	}

	page, size, err := httputils.Page(r)
	if err != nil {
		logger.Error("couldn't retrieve page/size from request",
			slog.String("Error", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.With(slog.Int("page", page), slog.Int("size", size))

	users, total, err := c.userRepo.List(ctx, page, size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(carefulness.JSONError{Err: "User not found"})
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
		Page:  page,
		Size:  size,
	})
}
