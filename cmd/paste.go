package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/drewsilcock/lemonade/client"
)

// pasteCmd represents the paste command
var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Paste content to server clipboard.",
	Long:  `Paste content to server clipboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		runPaste()
	},
}

func init() {
	rootCmd.AddCommand(pasteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pasteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pasteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPaste() {
	c := client.New(socketFile)
	text, err := c.Paste()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	_, _ = fmt.Print(text)
}
