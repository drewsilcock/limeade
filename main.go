package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	"github.com/lmittmann/tint"

	"github.com/drewsilcock/limeade/cmd"
)

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
		}),
	))

	if err := setVersionInfo(); err != nil {
		slog.Error(fmt.Sprintf("failed to set version info: %s", err.Error()))
	}

	cmd.Execute()
}

func setVersionInfo() error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("failed to read build info")
	}

	version := buildInfo.Main.Version
	commit := "unknown"
	commitDate := "unknown"
	goVersion := buildInfo.GoVersion

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			commit = setting.Value
		case "vcs.time":
			commitDate = setting.Value
		}
	}

	cmd.SetVersionInfo(version, commit, commitDate, goVersion)
	return nil
}
