package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type slackMessage struct {
	Text string `json:"text"`
}

// SendSlack posts a simple text-based notification to Slack.
func SendSlack(webhookURL string, event Event) {
	msg := fmt.Sprintf(
		"New event detected!\nBlock: %d\nTx: %s",
		event.BlockNumber,
		event.TxHash,
	)

	body, _ := json.Marshal(slackMessage{Text: msg})
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to send Slack notification: %v", err)
	}
}
