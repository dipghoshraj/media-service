package cmd

import (
	"fmt"

	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/rpc"
	"github.com/spf13/cobra"
)

var connectConfigFile string

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect and authenticate to the agent tunnel",
	Long:  `This command establishes and authenticates a connection to the agent tunnel, allowing you to interact with the IndraNet network.`,
	Run: func(cmd *cobra.Command, args []string) {
		bridge.LoadConfig(connectConfigFile)

		bridge.Logger.Info("registering agent with the registry")
		agent, gatewayIdentity, err := bridge.AgentRegistry()
		if err != nil {
			bridge.Logger.Error("failed to register agent", "error", err)
			return
		}
		bridge.Logger.Info("agent registered",
			"agent_id", agent.ID,
			"domain", agent.Domain,
			"fingerprint", agent.Identity,
		)
		gatewayConntion := fmt.Sprintf("%s:%d", agent.GatewayIP, agent.WssPort)
		bridge.Logger.Info("connecting to gateway", "gateway", gatewayConntion)
		_ = rpc.InitateConnection(gatewayConntion, gatewayIdentity)

		session, err := rpc.NewTunnelSession(agent)
		if err != nil {
			bridge.Logger.Error("failed to build tunnel session", "error", err)
			return
		}
		session.SendConnection()
	},
}

func init() {
	connectCmd.Flags().StringVarP(&connectConfigFile, "file", "f", "", "Path to config file (default: agni-config.yaml)")
	rootCmd.AddCommand(connectCmd)
}
