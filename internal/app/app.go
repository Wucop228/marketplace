package app

import (
	"database/sql"
	"fmt"
	"github.com/Wucop228/marketplace/internal/config"
	"github.com/Wucop228/marketplace/internal/delivery/http"
	appMiddleware "github.com/Wucop228/marketplace/internal/middleware"
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
	authHandler := http.NewAuthHandler(db, &authCfg)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authMiddleware := appMiddleware.AuthMiddleware(appMiddleware.DefaultConfig(authCfg.JWTSecret))

	auth := e.Group("")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	productHandler := http.NewProductHandler(db)

	api := e.Group("/api")
	api.Use(authMiddleware)
	api.POST("/create-product", productHandler.CreateProduct)

	if err := e.Start(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
