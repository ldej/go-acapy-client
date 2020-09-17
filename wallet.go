package acapy

import (
	"fmt"
	"strconv"
)

type QueryDIDsParams DID

type DID struct {
	DID    string `json:"did"`
	Public bool   `json:"public"`
	VerKey string `json:"verkey"`
}

func (c *Client) QueryDIDs(params QueryDIDsParams) ([]DID, error) {
	type results struct {
		DIDs []DID `json:"results"`
	}
	var r results
	queryParams := map[string]string{
		"did":    params.DID,
		"public": strconv.FormatBool(params.Public),
		"verkey": params.VerKey,
	}
	err := c.get(fmt.Sprintf("%s/wallet/did", c.ACApyURL), queryParams, &r)
	if err != nil {
		return nil, err
	}
	return r.DIDs, nil
}

type didResult struct {
	DID `json:"didResult"`
}

func (c *Client) CreateLocalDID() (DID, error) {
	var r didResult
	err := c.post(fmt.Sprintf("%s/wallet/didResult/create", c.ACApyURL), nil, nil, &r)
	if err != nil {
		return DID{}, err
	}
	return r.DID, nil
}

func (c *Client) GetPublicDID() (DID, error) {
	var r didResult
	err := c.get(fmt.Sprintf("%s/wallet/did/public", c.ACApyURL), nil, &r)
	if err != nil {
		return DID{}, err
	}
	return r.DID, nil
}

func (c *Client) SetPublicDID(did string) (DID, error) {
	var r didResult
	queryParams := map[string]string{
		"did": did,
	}
	err := c.post(fmt.Sprintf("%s/wallet/did/public", c.ACApyURL), queryParams, nil, &r)
	if err != nil {
		return DID{}, err
	}
	return r.DID, nil
}

func (c *Client) SetDIDEndpoint(did string, endpoint string, endpointType string) error {
	var setDIDEndpointRequest = struct {
		DID          string `json:"did"`
		Endpoint     string `json:"endpoint"`
		EndpointType string `json:"endpoint_type"`
	}{
		DID:          did,
		Endpoint:     endpoint,
		EndpointType: endpointType,
	}
	return c.post(fmt.Sprintf("%s/wallet/set-did-endpoint", c.ACApyURL), nil, setDIDEndpointRequest, nil)
}

func (c *Client) GetDIDEndpoint(did string) (string, error) {
	var r = struct {
		DID      string `json:"did"`
		Endpoint string `json:"endpoint"`
	}{}
	queryParams := map[string]string{
		"did": did,
	}
	err := c.get(fmt.Sprintf("%s/wallet/get-did-endpoint", c.ACApyURL), queryParams, &r)
	if err != nil {
		return "", err
	}
	return r.Endpoint, nil
}

func (c *Client) RotateKeypair(did string) error {
	queryParams := map[string]string{
		"did": did,
	}
	return c.patch(fmt.Sprintf("%s/wallet/did/local/rotate-keypair", c.ACApyURL), queryParams, nil, nil)
}
