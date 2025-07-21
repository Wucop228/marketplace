package main

import (
	"github.com/Wucop228/marketplace/internal/app"
	"github.com/Wucop228/marketplace/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	app.Run(cfg)
}
