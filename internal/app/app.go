package app

import (
	"fmt"
	"github.com/Wucop228/marketplace/internal/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

func Run(cfg *config.Config) {
	e := echo.New()

	if err := e.Start(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
