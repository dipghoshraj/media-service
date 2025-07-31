package persist

import (
	"fmt"
	"log"
)

func SetAgent(agentID string, socketHost string) error {
	// Set the agent data in Redis
	err := rdb.HSet(ctx, agentID, "router_host", socketHost).Err()
	if err != nil {
		fmt.Printf("Failed to set agent %s: %v", agentID, err)
	}
	log.Printf("Agent %s set with router host %s", agentID, socketHost)
	return nil
}
