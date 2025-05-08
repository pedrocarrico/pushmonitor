package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/pedrocarrico/pushmonitor/internal/config"
	"github.com/pedrocarrico/pushmonitor/internal/logger"
)

var (
	cfg config.Config
)

func writePIDFile(pidFile string) error {
	pid := os.Getpid()
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

func sendReloadSignal(pidFile string) error {
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("failed to read PID file: %v", err)
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		return fmt.Errorf("invalid PID in file: %v", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	if err := process.Signal(syscall.SIGHUP); err != nil {
		return fmt.Errorf("failed to send SIGHUP signal: %v", err)
	}

	return nil
}

func main() {
	// Initialize logger with default settings (stdout and info level)
	logger.Init("info", os.Stdout)

	reloadFlag := flag.Bool("reload", false, "Reload the service configuration")
	flag.Parse()

	logger.Info("Loading configuration...")
	if err := cfg.Load(); err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded successfully.")

	if *reloadFlag {
		logger.Info("Reloading configuration...")
		if err := sendReloadSignal(cfg.PIDFile); err != nil {
			logger.Error("Failed to reload configuration: %v", err)
			os.Exit(1)
		}
		logger.Info("Reloaded configuration successfully")
		return
	}

	logger.Info("Starting Push Monitor...")
	logger.Debug("Found %d push tests", len(cfg.PushTests))

	logger.Debug("Writing PID file to: %s", cfg.PIDFile)
	if err := writePIDFile(cfg.PIDFile); err != nil {
		logger.Error("Failed to write PID file: %v", err)
		os.Exit(1)
	}

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
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
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

	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP:
			logger.Info("Received reload signal, reloading configuration...")
			if err := cfg.Reload(); err != nil {
				logger.Error("Failed to reload configuration: %v", err)
				os.Exit(1)
			}
			logger.Debug("Setting up logging on file: %s", cfg.Logging.File)
			logFile, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logger.Error("error opening log file: %v", err)
				os.Exit(1)
			}

			logger.Init(cfg.Logging.Level, os.Stdout, logFile)
			logger.Debug("Logging setup completed")

			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			for _, test := range cfg.PushTests {
				wg.Add(1)
				go test.Run(httpClient, &wg, ctx)
			}
			logger.Info("Configuration reloaded and tests restarted")
		case syscall.SIGINT, syscall.SIGTERM:
			logger.Info("Received shutdown signal, initiating graceful shutdown...")
			cancel()
			wg.Wait()
			logger.Info("All tests stopped, shutting down...")
			return
		}
	}
}
