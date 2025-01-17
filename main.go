package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"

	"github.com/drewsilcock/lemonade/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
		}),
	))

	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
