package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var socketFile string

var pbcopyRegex = regexp.MustCompile("/?pbcopy$")
var pbpasteRegex = regexp.MustCompile("/?pbpaste$")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lemonade",
	Short: "Copy and paste between client and server machines.",
	Long: `Use Unix sockets to share a clipboard with a remote
machine, such as over SSH.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmdName := os.Args[0]
		slog.Error(cmdName)

		if pbcopyRegex.MatchString(cmdName) {
			// macOS copy command
			argText := ""
			if len(args) > 0 {
				argText = args[0]
			}
			runCopy(argText)
			return
		}

		if pbpasteRegex.MatchString(cmdName) {
			// macOS paste command
			runPaste()
			return
		}

		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s built at %s from commit %s", version, date, commit)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&socketFile, "socket", "/tmp/lemonade.sock", "socket file (default is /var/run/lemonade.sock)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
