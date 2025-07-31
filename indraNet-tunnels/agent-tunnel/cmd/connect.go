package cmd

import (
	"agent-tunnel/internal/connections/connects"
	"context"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var (
	gatewayURL  string
	token       string
	secret      string
	agentID     string
	portforward string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect and authenticate to the agent tunnel",
	Long:  `This command establishes and authenticates a connection to the agent tunnel, allowing you to interact with the IndraNet network.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("🔌 Connecting to the agent tunnel...")

		client := &connects.TunnelClient{
			Cfg: connects.ClientConfig{
				GatewayURL:  gatewayURL,
				Token:       token,
				Secret:      secret,
				AgentID:     agentID,
				Portforward: portforward,
			},
			Close:   make(chan struct{}),
			Streams: make(map[string]net.Conn),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := client.Start(ctx); err != nil {
			log.Fatalf("❌ Tunnel client exited: %v", err)
		}

	},
}

func init() {

	connectCmd.Flags().StringVar(&gatewayURL, "gateway", "", "Gateway WebSocket URL (wss://)")
	connectCmd.Flags().StringVar(&token, "token", "", "Auth token")
	connectCmd.Flags().StringVar(&secret, "secret", "", "HMAC secret key")
	connectCmd.Flags().StringVar(&agentID, "id", "", "Agent ID")
	connectCmd.Flags().StringVar(&portforward, "port", "", "local port to forward")

	connectCmd.MarkFlagRequired("gateway")
	connectCmd.MarkFlagRequired("token")
	connectCmd.MarkFlagRequired("secret")
	connectCmd.MarkFlagRequired("id")
	connectCmd.MarkFlagRequired("port")

	rootCmd.AddCommand(connectCmd)

}
