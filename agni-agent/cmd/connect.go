package cmd

import (
	"fmt"
	"log"

	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/rpc"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect and authenticate to the agent tunnel",
	Long:  `This command establishes and authenticates a connection to the agent tunnel, allowing you to interact with the IndraNet network.`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Println("📝 Registering agent with the registry...")
		agent, gatewayIdentity, err := bridge.AgentRegistry()
		if err != nil {
			log.Fatalf("❌ Failed to register agent: %v", err)
		}
		log.Printf("✅ Agent registered successfully: ID=%s, Domain=%s fingurePrint=%s", agent.ID, agent.Domain, agent.Identity)
		log.Println("🔌 Connecting to the agent tunnel...")
		gatewayConntion := fmt.Sprintf("%s:%d", agent.GatewayIP, agent.WssPort)

		log.Println("Connecting to gatewayURL", gatewayConntion)
		_ = rpc.InitateConnection(gatewayConntion, gatewayIdentity)

		rpc.SendConnection(agent)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
