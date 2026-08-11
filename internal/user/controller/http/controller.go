package usercontroller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"antiscam-simulator/internal/user/domain"
	"antiscam-simulator/internal/user/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, username string) (string, error)
	GetHistory(ctx context.Context, userID string) ([]userdomain.TrainingHistoryItem, error)
}

type UserController struct {
	usecase UserUsecase
}

func NewUserController(usecase UserUsecase) *UserController {
	return &UserController{usecase: usecase}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req userdto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" {
		c.writeError(w, http.StatusBadRequest, "username is required")
		return
	}

	userID, err := c.usecase.Register(r.Context(), req.Username)
	if err != nil {
		c.writeError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(userdto.RegisterResponse{UserID: userID})
}



func (c *UserController) GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		c.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := r.PathValue("user_id")
	if userID == "" {
		c.writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	history, err := c.usecase.GetHistory(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userdomain.ErrUserNotFound) {
			c.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		c.writeError(w, http.StatusInternalServerError, "failed to get user statistics")
		return
	}

	dtoHistory := userdto.MapHistoryFromDomain(history)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(userdto.UserHistoryResponse{
		UserID:  userID,
		History: dtoHistory,
	})
}

func (c *UserController) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(userdto.ErrorResponse{Error: msg})
}
