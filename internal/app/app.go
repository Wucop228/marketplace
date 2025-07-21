package app

import (
	"database/sql"
	"fmt"
	"github.com/Wucop228/marketplace/internal/config"
	"github.com/Wucop228/marketplace/internal/delivery/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func Run(cfg *config.Config) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s port=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.DBPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error", err)
	}
	defer db.Close()

	authCfg := config.AuthConfig{AccessTokenTTL: cfg.AccessTokenTTL, JWTSecret: cfg.JWTSecret}

	e := echo.New()
	h := http.NewAuthHandler(db, &authCfg)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")
	{
		api.POST("/login", h.Login)
		api.POST("/register", h.Register)
	}

	if err := e.Start(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
