package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initializationCmd = &cobra.Command{
	Use:   "initialize",
	Short: "Initialize the Agni CLI",
	Long:  `This command initializes the Agni CLI, setting up necessary configurations and environment variables`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Agent Tunnel Details:")
		// Here you would typically fetch and display the tunnel details
		fmt.Println("Configuration: [details here]")
		fmt.Println("Status: [status here]")
	},
}

func init() {
	rootCmd.AddCommand(initializationCmd)
}
