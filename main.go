package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nzai/go-trader/config"
	"github.com/nzai/go-trader/exchange"
	"github.com/nzai/go-trader/strategy"
)

const (
	appName    = "go-trader"
	appVersion = "0.1.0"
)

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	version := flag.Bool("version", false, "Print version and exit")
	// Default dry-run to true so I don't accidentally place real orders while experimenting
	dryRun := flag.Bool("dry-run", true, "Run in dry-run mode (no real orders)")
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s\n", appName, appVersion)
		os.Exit(0)
	}

	log.Printf("Starting %s v%s", appName, appVersion)

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", *configFile, err)
	}

	if *dryRun {
		log.Println("Running in dry-run mode — no real orders will be placed")
		cfg.DryRun = true
	}

	// Initialize exchange client
	client, err := exchange.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize exchange client: %v", err)
	}
	defer client.Close()

	// Initialize and start trading strategy
	strat, err := strategy.New(cfg, client)
	if err != nil {
		log.Fatalf("Failed to initialize strategy: %v", err)
	}

	if err := strat.Start(); err != nil {
		log.Fatalf("Failed to start strategy: %v", err)
	}
	defer strat.Stop()

	log.Printf("Trader is running. Press Ctrl+C to stop.")

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh

	log.Printf("Received signal %s, shutting down...", sig)
}
