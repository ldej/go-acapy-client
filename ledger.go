package acapy

import (
	"fmt"
)

// available endpoint types: Endpoint, Profile, LinkedDomains
func (c *Client) GetDIDEndpointFromLedger(did string, endpointType string) (string, error) {
	var result = struct {
		Endpoint string `json:"endpoint"`
	}{}
	var queryParams = map[string]string{
		"did":           did,
		"endpoint_type": endpointType,
	}
	err := c.get(fmt.Sprintf("%s/ledger/did-endpoint", c.ACApyURL), queryParams, &result)
	if err != nil {
		return "", err
	}
	return result.Endpoint, nil
}

func (c *Client) GetDIDVerkeyFromLedger(did string) (string, error) {
	var result = struct {
		Verkey string `json:"verkey"`
	}{}
	var queryParams = map[string]string{
		"did":           did,
	}
	err := c.get(fmt.Sprintf("%s/ledger/did-verkey", c.ACApyURL), queryParams, &result)
	if err != nil {
		return "", err
	}
	return result.Verkey, nil
}

// Use DID instead of NYM, as NYM is outdated
func (c *Client) GetDIDRoleFromLedger(did string) (string, error) {
	var result = struct {
		Role string `json:"role"`
	}{}
	var queryParams = map[string]string{
		"did":           did,
	}
	err := c.get(fmt.Sprintf("%s/ledger/get-nym-role", c.ACApyURL), queryParams, &result)
	if err != nil {
		return "", err
	}
	return result.Role, nil
}

