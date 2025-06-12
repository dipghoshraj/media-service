package cmd

import "github.com/spf13/cobra"

var (
	gatewayURL string
	token      string
	secret     string
	agentID    string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect and authenticate to the agent tunnel",
	Long:  `This command establishes and authenticates a connection to the agent tunnel, allowing you to interact with the IndraNet network.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Here you would typically implement the logic to connect to the agent tunnel
		// For now, we will just print a message
		println("Connecting to the agent tunnel...")
		// Add connection logic here
	},
}

func init() {

	connectCmd.Flags().StringVar(&gatewayURL, "gateway", "", "Gateway WebSocket URL (wss://)")
	connectCmd.Flags().StringVar(&token, "token", "", "Auth token")
	connectCmd.Flags().StringVar(&secret, "secret", "", "HMAC secret key")
	connectCmd.Flags().StringVar(&agentID, "id", "", "Agent ID")

	connectCmd.MarkFlagRequired("gateway")
	connectCmd.MarkFlagRequired("token")
	connectCmd.MarkFlagRequired("secret")
	connectCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(connectCmd)

}
