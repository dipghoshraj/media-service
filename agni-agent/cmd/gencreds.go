package cmd

import (
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"github.com/spf13/cobra"
)

var genCredsConfig string

var buildCredsCmd = &cobra.Command{
	Use:   "gen-creds",
	Short: "Generate TLS credentials for the agent tunnel",

	Long: `This command generates TLS credentials (certificates and keys) for secure communication in the IndraNet agent tunnel.`,
	Run: func(cmd *cobra.Command, args []string) {
		bridge.LoadConfig(genCredsConfig)
		dns := bridge.YamlConfig.Agent.Domain
		name := bridge.YamlConfig.Agent.Name
		bridge.Logger.Info("generating TLS credentials", "dns", dns, "name", name)
		err := bridge.BuildCreds(dns, name)
		if err != nil {
			bridge.Logger.Error("failed to generate TLS credentials", "error", err)
		}
	},
}

func init() {
	buildCredsCmd.Flags().StringVarP(&genCredsConfig, "file", "f", "", "Path to config file (default: agni-config.yaml)")
	rootCmd.AddCommand(buildCredsCmd)
}
