// Package represents main application for server.
//
//nolint:gochecknoglobals // build inline variable
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devldavydov/gophkeeper/internal/common/info"
	gkLog "github.com/devldavydov/gophkeeper/internal/common/log"
	"github.com/devldavydov/gophkeeper/internal/server"
	_ "github.com/lib/pq"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	appVer := info.FormatVersion(buildVersion, buildDate, buildCommit)
	fmt.Println(appVer)

	config, err := LoadConfig(*flag.CommandLine, os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to load flag and ENV settings: %w", err)
	}

	logger, err := gkLog.NewLogger(config.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	serverSettings, err := ServiceSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create server settings: %w", err)
	}

	logger.Info(appVer)
	serverService := server.NewService(serverSettings, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	return serverService.Start(ctx)
}
