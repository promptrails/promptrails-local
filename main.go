package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/promptrails/promptrails-local/internal/seed"
	"github.com/promptrails/promptrails-local/internal/server"
	"github.com/promptrails/promptrails-local/internal/store"
	"go.uber.org/zap"
)

var Version = "dev"

func main() {
	port := flag.Int("port", envInt("PORT", 8080), "Server port")
	seedData := flag.Bool("seed", envBool("SEED", true), "Load seed data on startup")
	fixturesDir := flag.String("fixtures", envStr("FIXTURES", ""), "Load additional fixtures from directory (JSON files)")
	logLevel := flag.String("log-level", envStr("LOG_LEVEL", "info"), "Log level (debug, info, warn, error)")
	corsOrigins := flag.String("cors-origins", envStr("CORS_ORIGINS", "*"), "CORS allowed origins")
	flag.Parse()

	logger := initLogger(*logLevel)
	defer logger.Sync()

	s := store.New()

	if *seedData {
		if err := seed.Load(s, logger); err != nil {
			logger.Fatal("failed to load seed data", zap.Error(err))
		}
		logger.Info("seed data loaded successfully")
	}

	if *fixturesDir != "" {
		if err := seed.LoadFromDir(s, *fixturesDir, logger); err != nil {
			logger.Fatal("failed to load fixtures from directory", zap.String("dir", *fixturesDir), zap.Error(err))
		}
		logger.Info("fixtures loaded from directory", zap.String("dir", *fixturesDir))
	}

	srv := server.New(s, logger, *corsOrigins, Version)

	go func() {
		printBanner(*port)
		addr := fmt.Sprintf(":%d", *port)
		if err := srv.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown error", zap.Error(err))
	}
}

func printBanner(port int) {
	fmt.Println()
	fmt.Println("  ╔═══════════════════════════════════════════╗")
	fmt.Println("  ║         PromptRails Local Emulator        ║")
	fmt.Printf("  ║         Version: %-25s║\n", Version)
	fmt.Println("  ╠═══════════════════════════════════════════╣")
	fmt.Printf("  ║  API:  http://localhost:%-19s║\n", fmt.Sprintf("%d/api/v1", port))
	fmt.Printf("  ║  Docs: http://localhost:%-19s║\n", fmt.Sprintf("%d/docs", port))
	fmt.Println("  ╚═══════════════════════════════════════════╝")
	fmt.Println()
}

func initLogger(level string) *zap.Logger {
	var cfg zap.Config
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, _ := cfg.Build()
	return logger
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		fmt.Sscanf(v, "%d", &i)
		return i
	}
	return def
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "true" || v == "1" {
		return true
	}
	if v == "false" || v == "0" {
		return false
	}
	return def
}
