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

func (c *Client) RegisterDID(alias string, seed string, role string) (RegisterDIDResponse, error) {
	var registerDID registerDIDRequest
	var registerDIDResponse RegisterDIDResponse

	registerDID = registerDIDRequest{
		Alias: alias,
		Seed:  seed, // Should be random in Develop mode
		Role:  role,
	}
	err := c.post(c.LedgerURL+"/register", nil, registerDID, &registerDIDResponse)
	if err != nil {
		return RegisterDIDResponse{}, err
	}
	return registerDIDResponse, nil
}
