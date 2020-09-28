package acapy

import (
	"fmt"
)

type PresentationExchange struct {
	PresentationExchangeID   string                  `json:"presentation_exchange_id"`
	ConnectionID             string                  `json:"connection_id"`
	ThreadID                 string                  `json:"thread_id"`
	State                    string                  `json:"state"`
	Initiator                string                  `json:"initiator"`
	Role                     string                  `json:"role"`
	PresentationProposalDict PresentationProposalMap `json:"presentation_proposal_dict"`
	PresentationRequest      PresentationRequest     `json:"presentation_request"`
	PresentationRequestDict  PresentationRequestMap  `json:"presentation_request_dict"`
	Presentation             Presentation            `json:"presentation"`
	Verified                 string                  `json:"verified"`
	CreatedAt                string                  `json:"created_at"`
	UpdatedAt                string                  `json:"updated_at"`
	ErrorMsg                 string                  `json:"error_msg"`
	AutoPresent              bool                    `json:"auto_present"`
	Trace                    bool                    `json:"trace"`
}

type PresentationProposalMap struct {
}

type Presentation struct {
}

type PresentationRequestMap struct {
}

type PresentationProposalRequest struct {
	Comment             string              `json:"comment"`
	AutoPresent         bool                `json:"auto_present"`
	PresentationPreview PresentationPreview `json:"presentation_proposal"`
	ConnectionID        string              `json:"connection_id"`
	Trace               bool                `json:"trace"`
}

type PresentationAttribute struct {
	Name                   string `json:"name"`
	CredentialDefinitionID string `json:"cred_def_id"`
	MimeType               string `json:"mime-type"`
	Value                  string `json:"value"`
	Referent               string `json:"referent"`
}

type Predicate struct {
	Name                   string `json:"name"`
	CredentialDefinitionID string `json:"cred_def_id"`
	Predicate              string `json:"predicate"`
	Threshold              int    `json:"threshold"`
}

type PresentationPreview struct {
	Type       string                  `json:"@type"`
	Attributes []PresentationAttribute `json:"attributes"`
	Predicates []Predicate             `json:"predicates"`
}

type PresentationRequest struct {
	Trace        bool         `json:"trace"`
	Comment      string       `json:"comment"`
	ConnectionID string       `json:"connection_id"`
	ProofRequest ProofRequest `json:"proof_request"`
}

type NonRevoked struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type AdditionalProp struct {
	Restrictions []map[string]string `json:"restrictions"`
	Name         string              `json:"name"`
	Names        []string            `json:"names"`
	PType        string              `json:"p_type"`
	PValue       int                 `json:"p_value"`
	NonRevoked   NonRevoked          `json:"non_revoked"`
}

type ProofRequest struct {
	Name                string                    `json:"name"`
	Nonce               string                    `json:"nonce"`
	RequestedPredicates map[string]AdditionalProp `json:"requested_predicates"`
	RequestedAttributes map[string]AdditionalProp `json:"requested_attributes"`
	Version             string                    `json:"version"`
	NonRevoked          NonRevoked                `json:"non_revoked"`
}

func (c *Client) SendPresentationProposal(request PresentationProposalRequest) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.post(fmt.Sprintf("%s/present-proof/send-proposal", c.ACApyURL), nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, err
	}
	return presentationExchange, nil
}

func (c *Client) CreatePresentationRequest() (PresentationExchange, error) {
	return PresentationExchange{}, nil
}

func (c *Client) SendPresentationRequest(request PresentationRequest) (PresentationExchange, error) {
	return PresentationExchange{}, nil
}

func (c *Client) QueryPresentationExchange() ([]PresentationExchange, error) {
	return nil, nil
}

func (c *Client) GetPresentationExchangeByID(presentationExchangeID string) (PresentationExchange, error) {
	return PresentationExchange{}, nil
}

func (c *Client) SendPresentationRequestByID(presentationExchangeID string) (PresentationExchange, error) {
	return PresentationExchange{}, nil
}

func (c *Client) SendPresentationByID(presentationExchangeID string) (PresentationExchange, error) {
	return PresentationExchange{}, nil
}

func (c *Client) VerifyPresentationByID(presentationExchangeID string) (PresentationExchange, error) {
	return PresentationExchange{}, nil
}

func (c *Client) GetPresentationCredentialsByID() ([]Credential, error) {
	return nil, nil
}

func (c *Client) RemovePresentationExchangeByID(presentationExchangeID string) (PresentationExchange, error) {
	return PresentationExchange{}, nil
}
