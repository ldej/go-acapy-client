package acapy

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type registerDIDRequest struct {
	Alias string `json:"alias"`
	Seed  string `json:"seed"`
	Role  string `json:"role"`
	DID   string `json:"did"`
}

type RegisterDIDResponse struct {
	DID    string `json:"did"`
	Seed   string `json:"seed"`
	Verkey string `json:"verkey"`
}

type DIDRole string

const (
	Endorser       DIDRole = "ENDORSER"
	Steward        DIDRole = "STEWARD"
	Trustee        DIDRole = "TRUSTEE"
	NetworkMonitor DIDRole = "NETWORK_MONITOR"
)

func RegisterDID(ledgerURL string, alias string, seed string, role DIDRole) (RegisterDIDResponse, error) {
	var request registerDIDRequest
	var response RegisterDIDResponse

	request = registerDIDRequest{
		Alias: alias,
		Seed:  seed, // Should be random in develop mode
		Role:  string(role),
	}
	body, err := json.Marshal(request)
	if err != nil {
		return RegisterDIDResponse{}, err
	}
	resp, err := http.Post(ledgerURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return RegisterDIDResponse{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return RegisterDIDResponse{}, err
	}

	return response, nil
}
