package acapy

import "strconv"

type Attachment struct {
	ID   string `json:"id"`   // either CredentialExchangeID or PresentationExchangeID
	Type string `json:"type"` // either credential-offer or present-proof
}

// CreateOutOfBandInvitationRequest must have IncludeHandshake true or, Attachments should be filled, or both
type CreateOutOfBandInvitationRequest struct {
	// When I put something in Attachments it crashes,
	// the CredentialExchangeRecord or PresentationExchangeRecord should probably be in the right state
	Attachments      []Attachment `json:"attachments,omitempty"`
	IncludeHandshake bool         `json:"include_handshake"`
	Metadata         struct{}     `json:"metadata,omitempty"` // TODO
	UsePublicDID     bool         `json:"use_public_did"`
}

type OutOfBandInvitationResponse struct {
	AutoAccept          bool                `json:"auto_accept"`
	InvitationMessageID string              `json:"invi_msg_id"`
	UpdatedAt           string              `json:"updated_at"`
	State               string              `json:"state"`
	InvitationID        string              `json:"invitation_id"`
	InvitationURL       string              `json:"invitation_url"`
	Trace               bool                `json:"trace"`
	MultiUse            bool                `json:"multi_use"`
	CreatedAt           string              `json:"created_at"`
	Invitation          OutOfBandInvitation `json:"invitation"`
}

type OutOfBandInvitation struct {
	Type               string    `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/out-of-band/1.0/invitation
	ID                 string    `json:"@id"`
	Label              string    `json:"label"`
	HandshakeProtocols []string  `json:"handshake_protocols"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/didexchange/v1.0
	Service            []Service `json:"service"`
	ServiceBlocks      []Service `json:"service_blocks,omitempty"`
	ServiceDIDs        []string  `json:"service_dids,omitempty"`
}

type Service struct {
	DID             string   `json:"did,omitempty"`
	ID              string   `json:"id"`
	Type            string   `json:"type"` // did-communication
	RecipientKeys   []string `json:"recipientKeys,omitempty"`
	RoutingKeys     []string `json:"routingKeys,omitempty"`
	ServiceEndpoint string   `json:"serviceEndpoint,omitempty"`
}

func (c *Client) CreateOutOfBandInvitation(request CreateOutOfBandInvitationRequest, autoAccept bool, multiUse bool) (OutOfBandInvitationResponse, error) {
	var result OutOfBandInvitationResponse
	var queryParams = map[string]string{
		"auto_accept": strconv.FormatBool(autoAccept),
		"multi_use":   strconv.FormatBool(multiUse),
	}
	err := c.post("/out-of-band/create-invitation", queryParams, request, &result)
	if err != nil {
		return OutOfBandInvitationResponse{}, err
	}
	return result, nil
}

func (c *Client) ReceiveOutOfBandInvitation(invitation OutOfBandInvitation, autoAccept bool) (Connection, error) {
	var result Connection
	var queryParams = map[string]string{
		"auto_accept": strconv.FormatBool(autoAccept),
		"alias":       invitation.Label,
	}
	err := c.post("/out-of-band/receive-invitation", queryParams, invitation, &result)
	if err != nil {
		return Connection{}, err
	}
	return result, nil
}
