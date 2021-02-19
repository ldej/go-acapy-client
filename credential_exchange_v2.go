package acapy

import (
	"fmt"
)

type createCredentialExchangeRecordRequestV2 struct {
	CredentialPreview CredentialPreviewV2 `json:"credential_preview"` // required
	Filter            Filter              `json:"filter"`
	Comment           string              `json:"comment,omitempty"`

	Trace      bool `json:"trace,omitempty"`
	AutoRemove bool `json:"auto_remove,omitempty"`
}

func (c *Client) CreateCredentialExchangeRecordV2(
	credentialPreview CredentialPreviewV2, // required
	credentialDefinitionID string, // optional
	issuerDID string, // optional
	schemaID string, // optional
	comment string, // optional
) (CredentialExchangeRecordResult, error) {

	var request = createCredentialExchangeRecordRequestV2{
		CredentialPreview: credentialPreview,
		Filter: Filter{
			DIF: map[string]string{},
			Indy: IndyFilter{
				CredentialDefinitionID: credentialDefinitionID,
				IssuerDID:              issuerDID,
				SchemaID:               schemaID,
			},
		},
		Comment:    comment,
		Trace:      c.tracing,
		AutoRemove: !c.preserveExchangeRecords,
	}
	var credentialExchange CredentialExchangeRecordResult
	err := c.post("/issue-credential-2.0/create", nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchange, nil
}

type credentialProposalRequestV2 struct {
	ConnectionID      string              `json:"connection_id"` // required
	Comment           string              `json:"comment,omitempty"`
	CredentialPreview CredentialPreviewV2 `json:"credential_preview"` // required
	Filter            Filter              `json:"filter"`

	Trace      bool `json:"trace,omitempty"`
	AutoRemove bool `json:"auto_remove,omitempty"`
}

// ProposeCredential tells the issuer what the holder hopes to receive
func (c *Client) ProposeCredentialV2(
	connectionID string, // required
	credentialPreview CredentialPreviewV2, // required
	comment string, // optional
	credentialDefinitionID string, // optional
	issuerDID string, // optional
	schemaID string, // optional
) (CredentialExchangeRecordResult, error) {

	var response CredentialExchangeRecordResult

	var request = credentialProposalRequestV2{
		ConnectionID:      connectionID,
		Comment:           comment,
		CredentialPreview: credentialPreview,
		Filter: Filter{
			DIF: map[string]string{},
			Indy: IndyFilter{
				CredentialDefinitionID: credentialDefinitionID,
				IssuerDID:              issuerDID,
				SchemaID:               schemaID,
			},
		},
		Trace:      c.tracing,
		AutoRemove: !c.preserveExchangeRecords,
	}

	err := c.post("/issue-credential-2.0/send-proposal", nil, request, &response)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return response, nil
}

type credentialOfferRequestV2 struct {
	CredentialDefinitionID string              `json:"cred_def_id,omitempty"`
	ConnectionID           string              `json:"connection_id"`      // required
	CredentialPreview      CredentialPreviewV2 `json:"credential_preview"` // required
	Comment                string              `json:"comment,omitempty"`
	Filter                 Filter              `json:"filter"`

	// filled automatically
	Trace      bool `json:"trace,omitempty"`
	AutoRemove bool `json:"auto_remove,omitempty"`
	AutoIssue  bool `json:"auto_issue,omitempty"`
}

type Filter struct {
	DIF  map[string]string `json:"dif"`
	Indy IndyFilter        `json:"indy"`
}

type IndyFilter struct {
	CredentialDefinitionID string `json:"cred_def_id,omitempty"`
	IssuerDID              string `json:"issuer_did,omitempty"`
	SchemaID               string `json:"schema_id,omitempty"`
	SchemaName             string `json:"schema_name,omitempty"`
	SchemaVersion          string `json:"schema_version,omitempty"`
	SchemaIssuerDID        string `json:"schema_issuer_did,omitempty"`
}

// OfferCredential sends to the holder what the issuer can offer
// TODO support payment decorator
func (c *Client) OfferCredentialV2(
	connectionID string, // required
	credentialPreview CredentialPreviewV2, // required
	credentialDefinitionID string, // optional
	comment string, // optional
) (CredentialExchangeRecordResult, error) {
	var offer = credentialOfferRequestV2{
		CredentialDefinitionID: credentialDefinitionID,
		ConnectionID:           connectionID,
		CredentialPreview:      credentialPreview,
		Comment:                comment,
		Filter: Filter{
			DIF: map[string]string{},
			Indy: IndyFilter{
				CredentialDefinitionID: credentialDefinitionID,
			},
		},
		Trace:      c.tracing,
		AutoRemove: !c.preserveExchangeRecords,
		AutoIssue:  c.autoRespondCredentialOffer,
	}
	var credentialExchangeRecord CredentialExchangeRecordResult
	err := c.post("/issue-credential-2.0/send-offer", nil, offer, &credentialExchangeRecord)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchangeRecord, nil
}

// OfferCredentialByID sends an offer to the holder based on a previously received proposal
func (c *Client) OfferCredentialByIDV2(credentialExchangeID string) (CredentialExchangeRecordResult, error) {
	var credentialExchange CredentialExchangeRecordResult
	err := c.post(fmt.Sprintf("/issue-credential-2.0/records/%s/send-offer", credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchange, nil
}

// RequestCredentialByID sends a credential request to the issuer based on a previously received offer
func (c *Client) RequestCredentialByIDV2(credentialExchangeID string) (CredentialExchangeRecordResult, error) {
	var credentialExchange CredentialExchangeRecordResult
	err := c.post(fmt.Sprintf("/issue-credential-2.0/records/%s/send-request", credentialExchangeID), nil, nil, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchange, nil
}

type issueCredentialRequestV2 struct {
	ConnectionID      string              `json:"connection_id"` // required
	Comment           string              `json:"comment,omitempty"`
	CredentialPreview CredentialPreviewV2 `json:"credential_preview"` // required
	Filter            Filter              `json:"filter"`

	Trace      bool `json:"trace,omitempty"`
	AutoRemove bool `json:"auto_remove,omitempty"`
}

// IssueCredential issues a credential to the holder
func (c *Client) IssueCredentialV2(
	connectionID string, // required
	credentialPreview CredentialPreviewV2, // required
	comment string, // optional
	credentialDefinitionID string, // optional
	issuerDID string, // optional
	schemaID string, // optional
) (CredentialExchangeRecordResult, error) {

	var credentialExchange CredentialExchangeRecordResult

	var request = issueCredentialRequestV2{
		ConnectionID:      connectionID,
		Comment:           comment,
		CredentialPreview: credentialPreview,
		Filter: Filter{
			DIF: map[string]string{},
			Indy: IndyFilter{
				CredentialDefinitionID: credentialDefinitionID,
				IssuerDID:              issuerDID,
				SchemaID:               schemaID,
			},
		},
		Trace:      c.tracing,
		AutoRemove: !c.preserveExchangeRecords,
	}

	err := c.post("/issue-credential-2.0/send", nil, request, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchange, nil
}

// IssueCredentialByID issues a credential to the holder based on a previously received request
func (c *Client) IssueCredentialByIDV2(credentialExchangeID string, comment string) (CredentialExchangeRecordResult, error) {
	var credentialExchange CredentialExchangeRecordResult
	var body = struct {
		Comment string `json:"comment"`
	}{
		Comment: comment,
	}
	err := c.post(fmt.Sprintf("/issue-credential-2.0/records/%s/issue", credentialExchangeID), nil, body, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchange, nil
}

// StoreCredentialByID the holder stores a credential based on a previously received issued credential
// credentialID is optional: https://github.com/hyperledger/aries-cloudagent-python/issues/594#issuecomment-656113125
func (c *Client) StoreCredentialByIDV2(credentialExchangeID string, credentialID string) (CredentialExchangeRecordResult, error) {
	var credentialExchange CredentialExchangeRecordResult
	var body = struct {
		CredentialID string `json:"credential_id,omitempty"`
	}{
		CredentialID: credentialID,
	}
	err := c.post(fmt.Sprintf("/issue-credential-2.0/records/%s/store", credentialExchangeID), nil, body, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchange, nil
}

type QueryCredentialExchangeParamsV2 struct {
	ConnectionID string `json:"connection_id"`
	Role         string `json:"role"`
	State        string `json:"state"`
	ThreadID     string `json:"thread_id"`
}

func (c *Client) QueryCredentialExchangeV2(params QueryCredentialExchangeParamsV2) ([]CredentialExchangeRecordResult, error) {
	var result = struct {
		Results []CredentialExchangeRecordResult `json:"results"`
	}{}
	var queryParams = map[string]string{
		"connection_id": params.ConnectionID,
		"role":          params.Role,
		"state":         params.State,
		"thread_id":     params.ThreadID,
	}
	err := c.get("/issue-credential-2.0/records", queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.Results, nil
}

func (c *Client) GetCredentialExchangeV2(credentialExchangeID string) (CredentialExchangeRecordResult, error) {
	var credentialExchange CredentialExchangeRecordResult
	err := c.get(fmt.Sprintf("/issue-credential-2.0/records/%s", credentialExchangeID), nil, &credentialExchange)
	if err != nil {
		return CredentialExchangeRecordResult{}, err
	}
	return credentialExchange, nil
}

func (c *Client) RemoveCredentialExchangeV2(credentialExchangeID string) error {
	return c.delete(fmt.Sprintf("/issue-credential-2.0/records/%s", credentialExchangeID))
}

func (c *Client) ReportCredentialExchangeProblemV2(credentialExchangeID string, message string) error {
	var body = struct {
		Message string `json:"explain_ltxt"`
	}{
		Message: message,
	}
	return c.post(fmt.Sprintf("/issue-credential-2.0/records/%s/problem-report", credentialExchangeID), nil, body, nil)
}
