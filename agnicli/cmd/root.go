package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "agnicli",
	Short: "Agnicli is a command line interface for Agni",
	Long:  `Agnicli is a command line interface for Agni, providing various commands to interact with the Agni platform.`,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
