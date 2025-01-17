package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install limeade client on remote machine",
	Long: `Runs helper script to install limeade on a remote machine via SSH.

Given a hostname, this command will SSH into the remote machine and run the
installation script to make limeade available on the remote machine,
enabling copy and pasting between the host and the remote SSH machine via
limeade.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: implement install command")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
