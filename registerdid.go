package acapy

type registerDIDRequest struct {
	Alias string `json:"alias"`
	Seed  string `json:"seed"`
	Role  string `json:"role"`
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

func (c *Client) RegisterDID(alias string, seed string, role DIDRole) (RegisterDIDResponse, error) {
	var registerDID registerDIDRequest
	var registerDIDResponse RegisterDIDResponse

	registerDID = registerDIDRequest{
		Alias: alias,
		Seed:  seed, // Should be random in develop mode
		Role:  string(role),
	}
	err := c.post_ledger("/register", nil, registerDID, &registerDIDResponse)
	if err != nil {
		return RegisterDIDResponse{}, err
	}
	return registerDIDResponse, nil
}
