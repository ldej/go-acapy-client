package acapy

import (
	"fmt"
)

func (c *Client) Plugins() ([]string, error) {
	var result = struct {
		Result []string `json:"result"`
	}{}
	err := c.get(fmt.Sprintf("%s/plugins", c.ACApyURL), nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

type Status struct {
	Version   string `json:"version"`
	Label     string `json:"label"`
	Conductor struct {
		InSessions  int `json:"in_sessions"`
		OutEncode   int `json:"out_encode"`
		OutDeliver  int `json:"out_deliver"`
		TaskActive  int `json:"task_active"`
		TaskDone    int `json:"task_done"`
		TaskFailed  int `json:"task_failed"`
		TaskPending int `json:"task_pending"`
	} `json:"conductor"`
}

func (c *Client) Status() (Status, error) {
	var status Status
	err := c.get(fmt.Sprintf("%s/status", c.ACApyURL), nil, &status)
	if err != nil {
		return Status{}, err
	}
	return status, nil
}

func (c *Client) ResetStatistics() error {
	return c.post(fmt.Sprintf("%s/status/reset", c.ACApyURL), nil, nil, nil)
}

func (c *Client) IsAlive() (bool, error) {
	var result = struct {
		Alive bool `json:"alive"`
	}{}
	err := c.get(fmt.Sprintf("%s/status/alive", c.ACApyURL), nil, &result)
	if err != nil {
		return false, err
	}
	return result.Alive, nil
}

func (c *Client) IsReady() (bool, error) {
	var result = struct {
		Ready bool `json:"ready"`
	}{}
	err := c.get(fmt.Sprintf("%s/status/ready", c.ACApyURL), nil, &result)
	if err != nil {
		return false, err
	}
	return result.Ready, nil
}

func (c *Client) Shutdown() error {
	return c.get(fmt.Sprintf("%s/shutdown", c.ACApyURL), nil, nil)
}

func (c *Client) Features() ([]string, error) {
	var result = struct {
		Results map[string]interface{} `json:"results"`
	}{}
	err := c.get(fmt.Sprintf("%s/features", c.ACApyURL), nil, &result)
	if err != nil {
		return nil, err
	}
	var features []string
	for feature, _ := range result.Results {
		features = append(features, feature)
	}
	return features, nil
}
