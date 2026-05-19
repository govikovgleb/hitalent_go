package main

import (
	"context"
	"hitalent_go/internal/config"
	"hitalent_go/internal/handlers"
	"hitalent_go/internal/repository"
	"hitalent_go/internal/router"
	"hitalent_go/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	config := config.Load()

	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	deptRepo, err := repository.NewDepartmentRepository(db)
	if err != nil {
		slog.Error("failed to initialize department repository", "error", err)
		os.Exit(1)
	}

	empRepo, err := repository.NewEmployeeRepository(db)
	if err != nil {
		slog.Error("failed to initialize employee repository", "error", err)
		os.Exit(1)
	}

	srv := service.New(deptRepo, empRepo)

	handlers := handlers.New(srv)
	router := router.New(handlers)
	port := config.Port

	slog.Info("starting server", "port", port)

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
}
