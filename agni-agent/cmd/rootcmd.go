package cmd

import (
	"os"

	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent-tunnel",
	Short: "Agent tunnel for the IndraNet network",
	Long: `Agent tunnel is a CLI application that manages
	tunneling for the IndraNet network. It provides commands
	for creating, managing, and monitoring tunnels.`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		bridge.Logger.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
