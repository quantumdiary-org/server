package main

import (
	"flag"
	"fmt"
	"os"

	"netschool-proxy/api/api/internal/app"
	"netschool-proxy/api/api/internal/config"
	"netschool-proxy/api/api/internal/pkg/logger"
)

var (
	configFile = flag.String("config", "config/dev.yaml", "Path to config file")
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Config validation error: %v\n", err)
		os.Exit(1)
	}

	// Create and start application
	appInstance, err := app.New(cfg)
	if err != nil {
		logger.Fatal("Failed to create application", "error", err)
		os.Exit(1)
	}

	if err := appInstance.Start(); err != nil {
		logger.Fatal("Application failed to start", "error", err)
		os.Exit(1)
	}
}