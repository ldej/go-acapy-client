package acapy

import (
	"fmt"
)

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
	MimeType string `json:"mime-type"` // optional
	Value    string `json:"value"`
}

type createCredentialExchangeRecordRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"`
	IssuerDID              string            `json:"issuer_did"`
	Comment                string            `json:"comment"`
	SchemaID               string            `json:"schema_id"`
	Trace                  bool              `json:"trace"`
	AutoRemove             bool              `json:"auto_remove"`

	// not supported
	SchemaName      string `json:"schema_name"`
	SchemaVersion   string `json:"schema_version"`
	SchemaIssuerDID string `json:"schema_issuer_did"`
}

func (c *Client) CreateCredentialExchangeRecord(
	credentialPreview CredentialPreview, // required
	credentialDefinitionID string, // optional
	issuerDID string, // optional
	schemaID string, // optional
	comment string, // optional
) (CredentialExchangeRecord, error) {

	var request = createCredentialExchangeRecordRequest{
		CredentialDefinitionID: credentialDefinitionID,
		IssuerDID:              issuerDID,
		SchemaID:               schemaID,
		CredentialPreview:      credentialPreview,
		Comment:                comment,
		Trace:                  c.tracing,
		AutoRemove:             !c.preserveExchangeRecords,
	}
	var credentialExchange CredentialExchangeRecord
	err := c.post("/issue-credential/create", nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchange, nil
}

type credentialProposalRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id"`
	ConnectionID           string            `json:"connection_id"`
	IssuerDID              string            `json:"issuer_did"`
	Comment                string            `json:"comment"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"`
	SchemaID               string            `json:"schema_id"`
	Trace                  bool              `json:"trace"`
	AutoRemove             bool              `json:"auto_remove"`

	// not supported
	SchemaName      string `json:"schema_name"`
	SchemaVersion   string `json:"schema_version"`
	SchemaIssuerDID string `json:"schema_issuer_did"`
}

// ProposeCredential tells the issuer what the holder hopes to receive
func (c *Client) ProposeCredential(
	connectionID string, // required
	credentialPreview CredentialPreview, // required
	comment string, // optional
	credentialDefinitionID string, // optional
	issuerDID string, // optional
	schemaID string, // optional
) (CredentialExchangeRecord, error) {

	var response CredentialExchangeRecord

	var request = credentialProposalRequest{
		CredentialDefinitionID: credentialDefinitionID,
		ConnectionID:           connectionID,
		IssuerDID:              issuerDID,
		Comment:                comment,
		CredentialPreview:      credentialPreview,
		SchemaID:               schemaID,
		Trace:                  c.tracing,
		AutoRemove:             !c.preserveExchangeRecords,
	}

	err := c.post("/issue-credential/send-proposal", nil, request, &response)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return response, nil
}

type credentialOfferRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id"`
	ConnectionID           string            `json:"connection_id"`
	CredentialPreview      CredentialPreview `json:"credential_preview"`
	Comment                string            `json:"comment"`

	// filled automatically
	Trace      bool `json:"trace"`
	AutoRemove bool `json:"auto_remove"`
	AutoIssue  bool `json:"auto_issue"`
}

// OfferCredential sends to the holder what the issuer can offer
// TODO support payment decorator
func (c *Client) OfferCredential(
	connectionID string, // required
	credentialPreview CredentialPreview, // required
	comment string, // optional
	credentialDefinitionID string, // optional
) (CredentialExchangeRecord, error) {
	var offer = credentialOfferRequest{
		CredentialDefinitionID: credentialDefinitionID,
		ConnectionID:           connectionID,
		CredentialPreview:      credentialPreview,
		Comment:                comment,
		Trace:                  c.tracing,
		AutoRemove:             !c.preserveExchangeRecords,
		AutoIssue:              c.autoRespondCredentialOffer,
	}
	var credentialExchangeRecord CredentialExchangeRecord
	err := c.post("/issue-credential/send-offer", nil, offer, &credentialExchangeRecord)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchangeRecord, nil
}

// OfferCredentialByID sends an offer to the holder based on a previously received proposal
func (c *Client) OfferCredentialByID(credentialExchangeID string) (CredentialExchangeRecord, error) {
	var credentialExchange CredentialExchangeRecord
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/send-offer", credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchange, nil
}

// RequestCredentialByID sends a credential request to the issuer based on a previously received offer
func (c *Client) RequestCredentialByID(credentialExchangeID string) (CredentialExchangeRecord, error) {
	var credentialExchange CredentialExchangeRecord
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/send-request", credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchange, nil
}

type issueCredentialRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id"`
	ConnectionID           string            `json:"connection_id"`
	IssuerDID              string            `json:"issuer_did"`
	Comment                string            `json:"comment"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"`
	SchemaID               string            `json:"schema_id"`
	Trace                  bool              `json:"trace"`
	AutoRemove             bool              `json:"auto_remove"`

	// not supported
	SchemaName      string `json:"schema_name"`
	SchemaVersion   string `json:"schema_version"`
	SchemaIssuerDID string `json:"schema_issuer_did"`
}

// IssueCredential issues a credential to the holder
func (c *Client) IssueCredential(
	connectionID string, // required
	credentialPreview CredentialPreview, // required
	comment string, // optional
	credentialDefinitionID string, // optional
	issuerDID string, // optional
	schemaID string, // optional
) (CredentialExchangeRecord, error) {

	var credentialExchange CredentialExchangeRecord

	var request = issueCredentialRequest{
		CredentialDefinitionID: credentialDefinitionID,
		ConnectionID:           connectionID,
		IssuerDID:              issuerDID,
		Comment:                comment,
		CredentialPreview:      credentialPreview,
		SchemaID:               schemaID,
		Trace:                  c.tracing,
		AutoRemove:             !c.preserveExchangeRecords,
	}

	err := c.post("/issue-credential/send", nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchange, nil
}

// IssueCredentialByID issues a credential to the holder based on a previously received request
func (c *Client) IssueCredentialByID(credentialExchangeID string, comment string) (CredentialExchangeRecord, error) {
	var credentialExchange CredentialExchangeRecord
	var body = struct {
		Comment string `json:"comment"`
	}{
		Comment: comment,
	}
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/issue", credentialExchangeID), nil, body, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchange, nil
}

// StoreCredentialByID the holder stores a credential based on a previously received issued credential
// credentialID is optional: https://github.com/hyperledger/aries-cloudagent-python/issues/594#issuecomment-656113125
func (c *Client) StoreCredentialByID(credentialExchangeID string, credentialID string) (CredentialExchangeRecord, error) {
	var credentialExchange CredentialExchangeRecord
	var body = struct {
		CredentialID string `json:"credential_id"`
	}{
		CredentialID: credentialID,
	}
	err := c.post(fmt.Sprintf("/issue-credential/records/%s/store", credentialExchangeID), nil, body, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchange, nil
}

type QueryCredentialExchangeParams struct {
	ConnectionID string `json:"connection_id"`
	Role         string `json:"role"`
	State        string `json:"state"`
	ThreadID     string `json:"thread_id"`
}

func (c *Client) QueryCredentialExchange(params QueryCredentialExchangeParams) ([]CredentialExchangeRecord, error) {
	var result = struct {
		Results []CredentialExchangeRecord `json:"result"`
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

func (c *Client) GetCredentialExchange(credentialExchangeID string) (CredentialExchangeRecord, error) {
	var credentialExchange CredentialExchangeRecord
	err := c.get(fmt.Sprintf("/issue-credential/records/%s", credentialExchangeID), nil, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
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
