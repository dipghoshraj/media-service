package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of agent-tunnel",
	Long:  `All software has versions. This is agent-tunnel's version.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print the version of the application
		// This should be set during the build process
		version := "0.1.2" // Replace with actual version variable if set during build
		log.Println("agent-tunnel version:", version)
	},
}

func init() {
	// Add the version command to the root command
	rootCmd.AddCommand(versionCmd)
}
