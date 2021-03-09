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

type CreateCredentialExchangeRecordRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id,omitempty"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"` // required
	IssuerDID              string            `json:"issuer_did,omitempty"`
	Comment                string            `json:"comment,omitempty"`
	SchemaID               string            `json:"schema_id,omitempty"`
	SchemaVersion          string            `json:"schema_version,omitempty"`
	SchemaIssuerDID        string            `json:"schema_issuer_did,omitempty"`
	SchemaName             string            `json:"schema_name,omitempty"`

	Trace      bool `json:"trace,omitempty"`
	AutoRemove bool `json:"auto_remove,omitempty"`
}

func (c *Client) CreateCredentialExchangeRecord(request CreateCredentialExchangeRecordRequest) (CredentialExchangeRecord, error) {
	var credentialExchange CredentialExchangeRecord
	err := c.post("/issue-credential/create", nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecord{}, err
	}
	return credentialExchange, nil
}

type credentialProposalRequest struct {
	CredentialDefinitionID string            `json:"cred_def_id,omitempty"`
	ConnectionID           string            `json:"connection_id"` // required
	IssuerDID              string            `json:"issuer_did,omitempty"`
	Comment                string            `json:"comment,omitempty"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"` // required
	SchemaID               string            `json:"schema_id,omitempty"`
	Trace                  bool              `json:"trace,omitempty"`
	AutoRemove             bool              `json:"auto_remove,omitempty"`

	// not supported
	SchemaName      string `json:"schema_name,omitempty"`
	SchemaVersion   string `json:"schema_version,omitempty"`
	SchemaIssuerDID string `json:"schema_issuer_did,omitempty"`
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
	CredentialDefinitionID string            `json:"cred_def_id,omitempty"`
	ConnectionID           string            `json:"connection_id"`      // required
	CredentialPreview      CredentialPreview `json:"credential_preview"` // required
	Comment                string            `json:"comment,omitempty"`

	// filled automatically
	Trace      bool `json:"trace,omitempty"`
	AutoRemove bool `json:"auto_remove,omitempty"`
	AutoIssue  bool `json:"auto_issue,omitempty"`
}

// OfferCredential sends to the holder what the issuer can offer
// TODO support payment decorator
func (c *Client) OfferCredential(
	connectionID string, // required
	credentialPreview CredentialPreview, // required
	credentialDefinitionID string, // optional
	comment string, // optional
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
	CredentialDefinitionID string            `json:"cred_def_id,omitempty"`
	ConnectionID           string            `json:"connection_id"` // required
	IssuerDID              string            `json:"issuer_did,omitempty"`
	Comment                string            `json:"comment,omitempty"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"` // required
	SchemaID               string            `json:"schema_id,omitempty"`
	Trace                  bool              `json:"trace,omitempty"`
	AutoRemove             bool              `json:"auto_remove,omitempty"`

	// not supported
	SchemaName      string `json:"schema_name,omitempty"`
	SchemaVersion   string `json:"schema_version,omitempty"`
	SchemaIssuerDID string `json:"schema_issuer_did,omitempty"`
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

type OutOfBandCredential struct {
	ID                string            `json:"@id"`
	Type              string            `json:"@type"`
	Comment           string            `json:"comment"`
	Service           Service           `json:"~service"`
	CredentialPreview CredentialPreview `json:"credential_preview"`
	OffersAttach      []OfferAttach     `json:"offers~attach"`
}

func (c *Client) CreateOutOfBandCredential(request CreateCredentialExchangeRecordRequest) (OutOfBandCredential, error) {

	record, err := c.CreateCredentialExchangeRecord(request)
	if err != nil {
		return OutOfBandCredential{}, err
	}

	invitation, err := c.CreateInvitation("", false, false, false)
	if err != nil {
		return OutOfBandCredential{}, err
	}

	return OutOfBandCredential{
		Comment:           request.Comment,
		Type:              record.CredentialOfferMap.Type,
		CredentialPreview: record.CredentialOfferMap.CredentialPreview,
		OffersAttach:      record.CredentialOfferMap.OffersAttach,
		Service: Service{
			RecipientKeys:   invitation.Invitation.RecipientKeys,
			RoutingKeys:     invitation.Invitation.RoutingKeys,
			ServiceEndpoint: invitation.Invitation.ServiceEndpoint,
		},
	}, nil
}
