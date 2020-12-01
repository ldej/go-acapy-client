package acapy

import (
	"fmt"
)

func (c *Client) SendBasicMessage(connectionID string, message string) error {
	type BasicMessage struct {
		Content string `json:"content"`
	}
	var basicMessage = BasicMessage{
		Content: message,
	}

	return c.post(fmt.Sprintf("/connections/%s/send-message", connectionID), nil, basicMessage, nil)
}
