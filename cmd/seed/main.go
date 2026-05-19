package main

import (
	"hitalent_go/internal/config"
	"hitalent_go/internal/repository"
	"hitalent_go/internal/seed"
	"hitalent_go/internal/service"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting database seeder...")

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
	seeder := seed.New(srv)

	if err := seeder.Run(); err != nil {
		slog.Error("failed to seed database", "error", err)
		os.Exit(1)
	}

	slog.Info("Seeding completed successfully!")
}
