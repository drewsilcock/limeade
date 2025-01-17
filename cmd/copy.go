package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/drewsilcock/limeade/client"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy text from stdin or argument to server clipboard.",
	Long: `Copy text from stdin or argument to server clipboard.

If text is provided via stdin, it will be used. If nothing is provided in
stdin, the first argument to this command will be used.`,
	Run: func(cmd *cobra.Command, args []string) {
		argText := ""
		if len(args) > 0 {
			argText = args[0]
		}
		runCopy(argText)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCopy(argText string) {
	text := argText

	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		slog.Error(fmt.Sprintf("unable to stat stdin: %s", err.Error()))
		os.Exit(1)
	}

	if stdinStat.Size() > 0 {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error(fmt.Sprintf("unable to read from stdin: %s", err.Error()))
			os.Exit(1)
		}

		text = string(stdin)
	}

	slog.Debug(fmt.Sprintf("copying text to '%s': %s", socketFile, text))

	c := client.New(socketFile)
	if err := c.Copy(text); err != nil {
		slog.Error(fmt.Sprintf("unable to copy text to server clipboard: %s", err.Error()))
		os.Exit(1)
	}
}
