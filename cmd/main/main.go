package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supchaser/LO_test_task/internal/app/delivery"
	"github.com/supchaser/LO_test_task/internal/app/repository"
	"github.com/supchaser/LO_test_task/internal/app/usecase"
	"github.com/supchaser/LO_test_task/internal/config"
	"github.com/supchaser/LO_test_task/internal/middleware/logging"
	recovery "github.com/supchaser/LO_test_task/internal/middleware/panic"
	"github.com/supchaser/LO_test_task/internal/utils/logger"
)

func main() {
	logger.Init()
	defer logger.Close()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatal("failed to load config", err, nil)
	}

	logger.Info("configuration loaded successfully", nil)

	repo := repository.CreateTaskRepository()
	uc := usecase.CreateTaskUsecase(repo)
	delivery := delivery.CreateTaskDelivery(uc)

	handlerChain := func(h http.Handler) http.Handler {
		return recovery.RecoveryMiddleware(logging.LoggingMiddleware(h))
	}

	mux := http.NewServeMux()

	mux.Handle("POST /tasks", handlerChain(http.HandlerFunc(delivery.CreateTask)))
	mux.Handle("GET /tasks/{id}", handlerChain(http.HandlerFunc(delivery.GetTask)))
	mux.Handle("GET /tasks", handlerChain(http.HandlerFunc(delivery.ListTasks)))
	mux.Handle("PUT /tasks/{id}", handlerChain(http.HandlerFunc(delivery.UpdateTask)))
	mux.Handle("DELETE /tasks/{id}", handlerChain(http.HandlerFunc(delivery.DeleteTask)))
	mux.Handle("GET /health", handlerChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})))

	port := ":" + cfg.ServerPort
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	serverErr := make(chan error, 1)

	go func() {
		logger.Info("starting HTTP server", map[string]any{
			"address": "http://localhost" + port,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", err, nil)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		logger.Error("failed to start server", err, nil)
		os.Exit(1)
	case sig := <-quit:
		logger.Info("server is shutting down", map[string]any{
			"signal": sig.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", err, nil)
			os.Exit(1)
		}

		logger.Info("server stopped", nil)
	}
}
