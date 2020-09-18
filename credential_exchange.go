package acapy

import (
	"fmt"
	"strconv"
)

type CredentialExchange struct {
	CredentialExchangeID   string             `json:"credential_exchange_id"`
	CredentialDefinitionID string             `json:"credential_definition_id"`
	ConnectionID           string             `json:"connection_id"`
	ThreadID               string             `json:"thread_id"`
	SchemaID               string             `json:"schema_id"`
	State                  string             `json:"state"`
	CredentialOffer        CredentialOffer    `json:"credential_offer"`
	CredentialOfferMap     CredentialOfferMap `json:"credential_offer_dict"`
	CredentialProposalMap  CredentialProposal `json:"credential_proposal_dict"`
	CredentialRequest      CredentialRequest  `json:"credential_request"`
	Initiator              string             `json:"initiator"`
	CreatedAt              string             `json:"created_at"`
	UpdatedAt              string             `json:"updated_at"`
	Trace                  bool               `json:"trace"`
	AutoOffer              bool               `json:"auto_offer"`
	AutoIssue              bool               `json:"auto_issue"`
	AutoRemove             bool               `json:"auto_remove"`
}

type CredentialOffer struct {
	SchemaID               string `json:"schema_id"`
	CredentialDefinitionID string `json:"cred_def_id"`
	KeyCorrectnessProof    struct {
		C     string     `json:"c"`
		XzCap string     `json:"xz_cap"`
		XrCap [][]string `json:"xr_cap"`
	} `json:"key_correctness_proof"`
	Nonce string `json:"nonce"`
}

type CredentialOfferMap struct {
	Type              string     `json:"@type"`
	ID                string     `json:"@id"`
	Thread            Thread     `json:"~thread"`
	CredentialPreview Credential `json:"credential_preview"`
	Comment           string     `json:"comment"`
	OffersAttach      []struct {
		ID       string `json:"@id"`
		MimeType string `json:"mime-type"`
		Data     struct {
			Base64 string `json:"base64"`
		} `json:"data"`
	} `json:"offers~attach"`
}

type CredentialProposal struct {
	ID                     string     `json:"@id"`
	Type                   string     `json:"@type"`
	CredentialDefinitionID string     `json:"cred_def_id"`
	SchemaID               string     `json:"schema_id"`
	IssuerDID              string     `json:"issuer_did"`
	SchemaName             string     `json:"schema_name"`
	SchemaIssuerDID        string     `json:"schema_issuer_did"`
	SchemaVersion          string     `json:"schema_version"`
	Comment                string     `json:"comment"`
	Credential             Credential `json:"credential_proposal"`
}

type CredentialRequest struct {
	CredentialDefinitionID string `json:"cred_def_id"`
	ProverDID              string `json:"prover_did"`
	BlindedMs              struct {
		U                   string   `json:"u"`
		Ur                  string   `json:"ur"`
		HiddenAttributes    []string `json:"hidden_attributes"`
		CommittedAttributes struct {
		} `json:"committed_attributes"` // TODO
	} `json:"blinded_ms"`
	BlindedMsCorrectnessProof struct {
		C        string `json:"c"`
		VDashCap string `json:"v_dash_cap"`
		MCaps    struct {
			MasterSecret string `json:"master_secret"`
		} `json:"m_caps"`
		RCaps struct {
		} `json:"r_caps"`
	} `json:"blinded_ms_correctness_proof"`
	Nonce string `json:"nonce"`
}

type Credential struct {
	Type       string      `json:"@type"`
	Attributes []Attribute `json:"attributes"`
}

type CredentialOfferRequest struct {
	CredentialDefinitionID string     `json:"cred_def_id"`
	ConnectionID           string     `json:"connection_id"`
	CredentialPreview      Credential `json:"credential_preview"`
	Comment                string     `json:"comment"`
	Trace                  bool       `json:"trace"`
	AutoRemove             bool       `json:"auto_remove"`
	AutoIssue              bool       `json:"auto_issue"`
}

type Attribute struct {
	Name     string `json:"name"`
	MimeType string `json:"mime-type"`
	Value    string `json:"value"`
}

func (c *Client) SendCredentialExchangeOffer(offer CredentialOfferRequest) (CredentialExchange, error) {
	var credentialExchangeRecord CredentialExchange
	err := c.post(fmt.Sprintf("%s/issue-credential/send-offer", c.ACApyURL), nil, offer, &credentialExchangeRecord)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchangeRecord, nil
}

type CredentialProposalRequest struct {
	CredentialDefinitionID string     `json:"cred_def_id"`
	ConnectionID           string     `json:"connection_id"`
	IssuerDID              string     `json:"issuer_did"`
	Comment                string     `json:"comment"`
	CredentialProposal     Credential `json:"credential_proposal"`
	SchemaName             string     `json:"schema_name"`
	SchemaVersion          string     `json:"schema_version"`
	SchemaID               string     `json:"schema_id"`
	SchemaIssuerDID        string     `json:"schema_issuer_did"`
	Trace                  bool       `json:"trace"`
	AutoRemove             bool       `json:"auto_remove"`
}

type CredentialProposalResponse struct {
	CredentialExchangeID  string             `json:"credential_exchange_id"`
	ConnectionID          string             `json:"connection_id"`
	ThreadID              string             `json:"thread_id"`
	State                 string             `json:"state"`
	Initiator             string             `json:"initiator"`
	Role                  string             `json:"role"`
	CredentialProposalMap CredentialProposal `json:"credential_proposal_dict"`
	CreatedAt             string             `json:"created_at"`
	UpdatedAt             string             `json:"updated_at"`
	Trace                 bool               `json:"trace"`
	AutoIssue             bool               `json:"auto_issue"`
	AutoRemove            bool               `json:"auto_remove"`
}

func (c *Client) SendCredentialExchangeProposal(proposal CredentialProposalRequest) (CredentialProposalResponse, error) {
	var credentialProposalResponse CredentialProposalResponse
	err := c.post(fmt.Sprintf("%s/issue-credential/send-proposal", c.ACApyURL), nil, proposal, &credentialProposalResponse)
	if err != nil {
		return CredentialProposalResponse{}, err
	}
	return credentialProposalResponse, nil
}

type QueryCredentialExchangeParams struct {
	ConnectionID string `json:"connection_id"`
	Role         string `json:"role"`
	State        string `json:"state"`
	ThreadID     string `json:"thread_id"`
}

func (c *Client) QueryCredentialExchange(params QueryCredentialExchangeParams) ([]CredentialExchange, error) {
	var result = struct {
		Results []CredentialExchange `json:"result"`
	}{}
	var queryParams = map[string]string{
		"connection_id": params.ConnectionID,
		"role":          params.Role,
		"state":         params.State,
		"thread_id":     params.ThreadID,
	}
	err := c.get(fmt.Sprintf("%s/issue-credential/records", c.ACApyURL), queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (c *Client) GetCredentialExchange(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.get(fmt.Sprintf("%s/issue-credential/records/%s", c.ACApyURL, credentialExchangeID), nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

type CredentialCreateRequest struct {
	CredentialDefinitionID string `json:"cred_def_id"`
	// ConnectionID           string     `json:"connection_id"`
	IssuerDID          string     `json:"issuer_did"`
	Comment            string     `json:"comment"`
	CredentialProposal Credential `json:"credential_proposal"`
	SchemaName         string     `json:"schema_name"`
	SchemaVersion      string     `json:"schema_version"`
	SchemaID           string     `json:"schema_id"`
	SchemaIssuerDID    string     `json:"schema_issuer_did"`
	Trace              bool       `json:"trace"`
	AutoRemove         bool       `json:"auto_remove"`
}

func (c *Client) CreateCredentialExchange(request CredentialCreateRequest) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("%s/issue-credential/create", c.ACApyURL), nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

type CredentialSendRequest CredentialProposalRequest

func (c *Client) SendCredentialExchange(request CredentialSendRequest) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("%s/issue-credential/send", c.ACApyURL), nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) SendCredentialExchangeOfferByID(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("%s/issue-credential/%s/send-offer", c.ACApyURL, credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) SendCredentialExchangeRequestByID(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("%s/issue-credential/%s/send-request", c.ACApyURL, credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) SendCredentialToHolder(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("%s/issue-credential/%s/issue", c.ACApyURL, credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) StoreReceivedCredential(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("%s/issue-credential/%s/store", c.ACApyURL, credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) RemoveCredentialExchange(credentialExchangeID string) error {
	return c.post(fmt.Sprintf("%s/issue-credential/%s/remove", c.ACApyURL, credentialExchangeID), nil, nil, nil)
}

func (c *Client) ReportCredentialExchangeProblem(credentialExchangeID string, message string) error {
	var body = struct {
		Message string `json:"explain_ltxt"`
	}{
		Message: message,
	}
	return c.post(fmt.Sprintf("%s/issue-credential/%s/problem-report", c.ACApyURL, credentialExchangeID), nil, body, nil)
}

func (c *Client) RevokeIssuedCredential(credentialRevocationID string, revocationRegistryID string, publish bool) error {
	queryParams := map[string]string{
		"cred_rev_id": credentialRevocationID,
		"rev_reg_id":  revocationRegistryID,
		"publish":     strconv.FormatBool(publish),
	}
	return c.post(fmt.Sprintf("%s/issue-credential/revoke", c.ACApyURL), queryParams, nil, nil)
}

// func (c *Client) PublishRevocations()
// func (c *Client) ClearPendingRevocations()
