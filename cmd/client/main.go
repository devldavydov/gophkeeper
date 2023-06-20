// Package represents main application for client.
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

	"github.com/devldavydov/gophkeeper/internal/client"
	"github.com/devldavydov/gophkeeper/internal/common/info"
	gkLog "github.com/devldavydov/gophkeeper/internal/common/log"
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

	cltSettings, err := ClientSettingsAdapt(config)
	if err != nil {
		return fmt.Errorf("failed to create application settings: %w", err)
	}

	logger.Info(appVer)
	client := client.NewClient(cltSettings, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	return client.Start(ctx)
}
