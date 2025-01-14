package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/watzon/0x45/internal/config"
	"github.com/watzon/0x45/internal/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title 0x45 API
// @version 1.0
// @description API for 0x45
// @license.name MIT
// @license.url https://github.com/watzon/0x45/blob/main/LICENSE
// @host localhost:3000
// @BasePath /
func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize logger with environment-based level
	logLevel := getLogLevel()
	logConfig := zap.NewDevelopmentConfig()
	logConfig.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info("logger initialized", zap.String("level", logLevel.String()))

	// Initialize server with storage manager
	srv := server.New(cfg, logger)

	// Create and setup server
	srv.SetupRoutes()

	go func() {
		// Start server
		log.Printf("Starting server on %s", cfg.Server.Address)
		log.Fatal(srv.Start(cfg.Server.Address))
	}()

	// Wait for shutdown signal and initiate graceful shutdown once received.
	<-ctx.Done()
	defer func() {
		if err := srv.Cleanup(); err != nil {
			log.Printf("failed cleaning up server: %v", err)
		}
	}()
	log.Print("shutdown signal received, initiate shutdown")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed shutting down gracefully: %v", err)
	}
}

func getLogLevel() zapcore.Level {
	env := os.Getenv("0X_LOG_LEVEL")
	if env == "" {
		return zapcore.InfoLevel
	}

	switch strings.ToLower(env) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		log.Fatalf("unknown log level: %s", env)
		return zapcore.InfoLevel
	}
}
