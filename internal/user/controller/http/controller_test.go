package usercontroller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"antiscam-simulator/internal/user/controller/http"
	"antiscam-simulator/internal/user/domain"
	"github.com/stretchr/testify/assert"
)

type mockUserUsecase struct {
	registerFunc   func(ctx context.Context, username string) (string, error)
	saveGameFunc   func(ctx context.Context, result *userdomain.TrainingResult) error
	getHistoryFunc func(ctx context.Context, userID string) ([]userdomain.TrainingHistoryItem, error)
}

func (m *mockUserUsecase) Register(ctx context.Context, username string) (string, error) {
	return m.registerFunc(ctx, username)
}

func (m *mockUserUsecase) SaveTrainingResult(ctx context.Context, result *userdomain.TrainingResult) error {
	return m.saveGameFunc(ctx, result)
}

func (m *mockUserUsecase) GetHistory(ctx context.Context, userID string) ([]userdomain.TrainingHistoryItem, error) {
	return m.getHistoryFunc(ctx, userID)
}

func TestUserController_Register(t *testing.T) {
	tests := []struct {
		name         string
		payload      map[string]interface{}
		mockRegister func(ctx context.Context, username string) (string, error)
		expectedCode int
	}{
		{
			name:    "success",
			payload: map[string]interface{}{"username": "testuser"},
			mockRegister: func(_ context.Context, _ string) (string, error) {
				return "uuid-123", nil
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "empty username",
			payload:      map[string]interface{}{"username": ""},
			mockRegister: nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid json",
			payload:      nil,
			mockRegister: nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockUserUsecase{registerFunc: tt.mockRegister}
			handler := usercontroller.NewUserController(svc)

			var body []byte
			if tt.payload != nil {
				body, _ = json.Marshal(tt.payload)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestUserController_GetHistory(t *testing.T) {
	tests := []struct {
		name        string
		pathUserID  string
		mockHistory func(ctx context.Context, userID string) ([]userdomain.TrainingHistoryItem, error)
		expectCode  int
	}{
		{
			name:       "success",
			pathUserID: "123",
			mockHistory: func(_ context.Context, _ string) ([]userdomain.TrainingHistoryItem, error) {
				return []userdomain.TrainingHistoryItem{{ScenarioID: "scen1"}}, nil
			},
			expectCode: http.StatusOK,
		},
		{
			name:       "user not found",
			pathUserID: "404",
			mockHistory: func(_ context.Context, _ string) ([]userdomain.TrainingHistoryItem, error) {
				return nil, userdomain.ErrUserNotFound
			},
			expectCode: http.StatusNotFound,
		},
		{
			name:        "empty path value",
			pathUserID:  "",
			mockHistory: nil,
			expectCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockUserUsecase{getHistoryFunc: tt.mockHistory}
			handler := usercontroller.NewUserController(svc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.pathUserID+"/history", http.NoBody)
			if tt.pathUserID != "" {
				req.SetPathValue("user_id", tt.pathUserID)
			}

			w := httptest.NewRecorder()

			handler.GetHistory(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}
