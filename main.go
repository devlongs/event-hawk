package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/devlongs/event-hawk/cmd"
	"github.com/devlongs/event-hawk/notifications"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func main() {
	cmd.Execute()

	client, err := ethclient.Dial(viper.GetString("ethereum.wss_url"))
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	// Prepare filter query
	contractAddr := common.HexToAddress(viper.GetString("ethereum.contract_address"))
	eventSignatures := viper.GetStringSlice("ethereum.event_signatures")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics:    [][]common.Hash{getTopics(eventSignatures)},
	}

	// Subscribe to event logs
	logsChan := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsChan)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("Monitoring events in real-time...")

	// Main loop: watch for new logs or errors
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case vLog := <-logsChan:
			handleLog(client, vLog)
		}
	}
}

func getTopics(signatures []string) []common.Hash {
	var topicList []common.Hash
	for _, sig := range signatures {
		// The Keccak-256 hash of the ABI signature string
		hash := crypto.Keccak256Hash([]byte(sig))
		topicList = append(topicList, hash)
	}
	return topicList
}

func handleLog(client *ethclient.Client, vLog types.Log) {
	// Construct event model
	event := notifications.Event{
		BlockNumber: vLog.BlockNumber,
		TxHash:      vLog.TxHash.Hex(),
		Data:        make(map[string]interface{}),
	}

	transferAbi, err := abi.JSON(strings.NewReader(`[{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "from", "type": "address"},
			{"indexed": true, "name": "to", "type": "address"},
			{"indexed": false, "name": "value", "type": "uint256"}
		],
		"name": "Transfer",
		"type": "event"
	}]`))
	if err != nil {
		log.Printf("Error parsing transfer ABI: %v", err)
		return
	}

	// Unpack the log’s data into our map
	err = transferAbi.UnpackIntoMap(event.Data, "Transfer", vLog.Data)
	if err != nil {
		log.Printf("Error unpacking log data: %v", err)
		return
	}

	// Output the event according to CLI flag
	switch cmd.OutputFormat {
	case "json":
		if out, err := json.Marshal(event); err == nil {
			fmt.Println(string(out))
		} else {
			log.Printf("Failed to marshal JSON: %v", err)
		}
	default:
		fmt.Printf("Block: %d | Tx: %s\n", event.BlockNumber, event.TxHash)
		for k, v := range event.Data {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	// Send a Slack notification if configured
	if webhook := viper.GetString("notifications.slack_webhook"); webhook != "" {
		notifications.SendSlack(webhook, event)
	}
}
