package main

import (
	"context"
	"log"
	"time"

	"github.com/crlnravel/go-fiber-template/internal/config"
	"github.com/crlnravel/go-fiber-template/internal/models"
	"github.com/crlnravel/go-fiber-template/platform/database"
)

// @title LemonForm API
// @version 1.0
// @description API Endpoints for LemonForm — a form builder application.
// @license.name MIT
// @host localhost:8080
// @BasePath /
func main() {
	// Connect to database (Supabase PostgreSQL via GORM)
	database.ConnectPostgres()

	// Auto-migrate models
	if err := database.DB.AutoMigrate(&models.User{}, &models.Form{}, &models.Question{}, &models.Response{}, &models.Answer{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	prod := config.GetStageStatus() == config.EnvironmentProduction

	app := NewApp(&appConfig{
		prod: prod,
		db:   database.DB,
	})

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"database": func(ctx context.Context) error {
			sqlDB, err := database.DB.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
		"http-server": func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	port := ":" + config.GetEnv("SERVER_PORT", "8080")

	if err := app.Listen(port); err != nil {
		log.Panic(err)
	}

	<-wait
}
