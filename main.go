package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"

	"github.com/drewsilcock/lemonade/cmd"
)

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
		}),
	))

	cmd.Execute()
}
