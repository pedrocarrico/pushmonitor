package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pedrocarrico/pushmonitor/internal/config"
	"github.com/pedrocarrico/pushmonitor/internal/logger"
	"github.com/pedrocarrico/pushmonitor/internal/version"
)

var (
	cfg         config.Config
	showVersion bool
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
}

func main() {
	flag.Parse()
	v := version.Get()

	if showVersion {
		fmt.Printf("Push Monitor version %s (build: %s, commit: %s)\n", v.Version, v.BuildTime, v.GitCommit)
		os.Exit(0)
	}

	logger.Init("info", os.Stdout)
	logger.Info("Starting Push Monitor version %s (build: %s, commit: %s)", v.Version, v.BuildTime, v.GitCommit)
	logger.Info("Loading configuration...")
	if err := cfg.Load(); err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded successfully.")

	logger.Info("Starting Push Monitor...")
	logger.Debug("Found %d push tests", len(cfg.PushTests))

	logger.Info("Setting up logging on file: %s", cfg.Logging.File)
	logFile, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("error opening log file: %v", err)
		os.Exit(1)
	}

	logger.Init(cfg.Logging.Level, os.Stdout, logFile)
	logger.Debug("Logging setup completed")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Debug("Setting up signal handlers...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	logger.Debug("Signal handlers configured")

	var httpClient = &http.Client{
		Timeout: time.Second * time.Duration(cfg.Timeout),
	}

	logger.Info("Starting push tests...")
	var wg sync.WaitGroup
	for _, test := range cfg.PushTests {
		wg.Add(1)
		go test.Run(httpClient, &wg, ctx)
	}
	logger.Info("All push tests started")

	<-sigChan
	logger.Info("Received shutdown signal, initiating graceful shutdown...")
	cancel()
	wg.Wait()
	logger.Info("All tests stopped, shutting down...")
}
