package main

import (
	"log/slog"
	"os"
	"tuntun/internal/app"
)

const (
	ShowPackageByteLimit = 10
)

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slog.LevelDebug,
        // AddSource: true,
    }))
	slog.SetDefault(logger)

	slog.Info("Starting TUN interface...")

    tun := app.NewTUN(logger)

    tun.MustRun(ShowPackageByteLimit)
}