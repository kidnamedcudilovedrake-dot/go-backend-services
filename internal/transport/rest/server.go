package rest

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	simulatorcontroller "antiscam-simulator/internal/simulator/controller/http"
	usercontroller "antiscam-simulator/internal/user/controller/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, userController *usercontroller.UserController, simController *simulatorcontroller.SimulatorController) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/register", userController.Register)
	mux.HandleFunc("GET /api/v1/users/{user_id}/history", userController.GetHistory)

	mux.HandleFunc("POST /api/v1/game/start", simController.StartGame())
	mux.HandleFunc("POST /api/v1/game/step", simController.ProcessStep())
	mux.HandleFunc("GET /api/v1/scenarios", simController.GetScenarios())
	mux.HandleFunc("POST /api/v1/game/generate_ai", simController.GenerateScenario())

	return &Server{
		httpServer: &http.Server{
			Addr:    port,
			Handler: corsMiddleware(mux),
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			os.Exit(1)
		}
	}()

	select {
	case <-ctx.Done():
	case <-shutdownChan:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(shutdownCtx)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
