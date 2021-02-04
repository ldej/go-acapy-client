package acapy

import (
	"fmt"
)

type CredentialExchange struct {
	CredentialExchangeID      string                    `json:"credential_exchange_id"`
	CredentialDefinitionID    string                    `json:"credential_definition_id"`
	ConnectionID              string                    `json:"connection_id"`
	ThreadID                  string                    `json:"thread_id"`
	ParentThreadID            string                    `json:"parent_thread_id"`
	SchemaID                  string                    `json:"schema_id"`
	RevocationID              string                    `json:"revocation_id"`
	RevocationRegistryID      string                    `json:"revoc_reg_id"`
	State                     string                    `json:"state"`
	CredentialOffer           CredentialOffer           `json:"credential_offer"`
	CredentialOfferMap        CredentialOfferMap        `json:"credential_offer_dict"`
	CredentialProposalMap     CredentialProposal        `json:"credential_proposal_dict"`
	CredentialRequest         CredentialRequest         `json:"credential_request"`
	CredentialRequestMetadata CredentialRequestMetadata `json:"credential_request_metadata"`
	Credential                Credential                `json:"credential"`
	RawCredential             RawCredential             `json:"raw_credential"`
	Role                      string                    `json:"role"`
	Initiator                 string                    `json:"initiator"`
	CreatedAt                 string                    `json:"created_at"`
	UpdatedAt                 string                    `json:"updated_at"`
	ErrorMessage              string                    `json:"error_msg"`
	Trace                     bool                      `json:"trace"`
	AutoOffer                 bool                      `json:"auto_offer"`
	AutoIssue                 bool                      `json:"auto_issue"`
	AutoRemove                bool                      `json:"auto_remove"`
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
	Type              string            `json:"@type"`
	ID                string            `json:"@id"`
	Thread            Thread            `json:"~thread"`
	CredentialPreview CredentialPreview `json:"credential_preview"`
	Comment           string            `json:"comment"`
	OffersAttach      []struct {
		ID       string `json:"@id"`
		MimeType string `json:"mime-type"`
		Data     struct {
			Base64 string `json:"base64"`
		} `json:"data"`
	} `json:"offers~attach"`
}

type CredentialProposal struct {
	ID                     string            `json:"@id"`
	Type                   string            `json:"@type"`
	CredentialDefinitionID string            `json:"cred_def_id"`
	SchemaID               string            `json:"schema_id"`
	IssuerDID              string            `json:"issuer_did"`
	SchemaName             string            `json:"schema_name"`
	SchemaIssuerDID        string            `json:"schema_issuer_did"`
	SchemaVersion          string            `json:"schema_version"`
	Comment                string            `json:"comment"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"`
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

type CredentialRequestMetadata struct {
	Description string `json:"description"`
}

type RawCredential struct {
	Description string `json:"raw_credential"`
}

func NewCredentialPreview(attributes []CredentialPreviewAttribute) CredentialPreview {
	return CredentialPreview{
		Type:       "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/1.0/credential-preview",
		Attributes: attributes,
	}
}

type CredentialPreview struct {
	Type       string                       `json:"@type"`
	Attributes []CredentialPreviewAttribute `json:"attributes"`
}

type CredentialPreviewAttribute struct {
	Name     string `json:"name"`
	MimeType string `json:"mime-type"`
	Value    string `json:"value"`
}

type CredentialOfferRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id"`
	ConnectionID           string            `json:"connection_id"`
	CredentialPreview      CredentialPreview `json:"credential_preview"`
	Comment                string            `json:"comment"`
	Trace                  bool              `json:"trace"`
	AutoRemove             bool              `json:"auto_remove"`
	AutoIssue              bool              `json:"auto_issue"`
}

type CredentialProposalRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id"`
	ConnectionID           string            `json:"connection_id"`
	IssuerDID              string            `json:"issuer_did"`
	Comment                string            `json:"comment"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"`
	SchemaName             string            `json:"schema_name"`
	SchemaVersion          string            `json:"schema_version"`
	SchemaID               string            `json:"schema_id"`
	SchemaIssuerDID        string            `json:"schema_issuer_did"`
	Trace                  bool              `json:"trace"`
	AutoRemove             bool              `json:"auto_remove"`
}

func (c *Client) SendCredentialProposal(proposal CredentialProposalRequest) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post("/issue-credential/send-proposal", nil, proposal, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) SendCredentialOffer(offer CredentialOfferRequest) (CredentialExchange, error) {
	var credentialExchangeRecord CredentialExchange
	err := c.post("/issue-credential/send-offer", nil, offer, &credentialExchangeRecord)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchangeRecord, nil
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
	err := c.get("/issue-credential/records", queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (c *Client) GetCredentialExchange(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.get(fmt.Sprintf("/issue-credential/records/%s", credentialExchangeID), nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

type CredentialCreateRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id"`
	IssuerDID              string            `json:"issuer_did"`
	Comment                string            `json:"comment"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"`
	SchemaName             string            `json:"schema_name"`
	SchemaVersion          string            `json:"schema_version"`
	SchemaID               string            `json:"schema_id"`
	SchemaIssuerDID        string            `json:"schema_issuer_did"`
	Trace                  bool              `json:"trace"`
	AutoRemove             bool              `json:"auto_remove"`
}

func (c *Client) CreateCredentialExchange(request CredentialCreateRequest) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post("/issue-credential/create", nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

type CredentialSendRequest CredentialProposalRequest

func (c *Client) SendCredential(request CredentialSendRequest) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post("/issue-credential/send", nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) SendCredentialOfferByID(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/send-offer", credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) SendCredentialRequestByID(credentialExchangeID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/send-request", credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) IssueCredentialByID(credentialExchangeID string, comment string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	var body = struct {
		Comment string `json:"comment"`
	}{
		Comment: comment,
	}
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/issue", credentialExchangeID), nil, body, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

// credentialID is optional: https://github.com/hyperledger/aries-cloudagent-python/issues/594#issuecomment-656113125
func (c *Client) StoreCredentialByID(credentialExchangeID string, credentialID string) (CredentialExchange, error) {
	var credentialExchange CredentialExchange
	var body = struct {
		CredentialID string `json:"credential_id"`
	}{
		CredentialID: credentialID,
	}
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/store", credentialExchangeID), nil, body, &credentialExchange)
	if err != nil {
		return CredentialExchange{}, err
	}
	return credentialExchange, nil
}

func (c *Client) RemoveCredentialExchange(credentialExchangeID string) error {
	return c.delete(fmt.Sprintf("/issue-credential/records/%s", credentialExchangeID))
}

func (c *Client) ReportCredentialExchangeProblem(credentialExchangeID string, message string) error {
	var body = struct {
		Message string `json:"explain_ltxt"`
	}{
		Message: message,
	}
	return c.post(fmt.Sprintf("/issue-credential/records/%s/problem-report", credentialExchangeID), nil, body, nil)
}
