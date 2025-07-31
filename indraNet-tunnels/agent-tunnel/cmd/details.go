package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Print details about the agent tunnel",
	Long:  `This command provides detailed information about the agent tunnel configuration and status.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Agent Tunnel Details:")
		// Here you would typically fetch and display the tunnel details
		fmt.Println("Configuration: [details here]")
		fmt.Println("Status: [status here]")
	},
}

func init() {
	// Add the details command to the root command
	rootCmd.AddCommand(detailsCmd)

	// 	// Add global flags
	// 	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agent-tunnel.yaml)")

	// 	// Add local flags
	// 	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
