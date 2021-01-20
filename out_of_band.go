package acapy

import "strconv"

// {
//  "attachments": [
//    {
//      "id": "string",
//      "type": "string"
//    }
//  ],
//  "include_handshake": true,
//  "metadata": {},
//  "use_public_did": true
//}

type Attachment struct {
	ID   string `json:"id"`   // either CredentialExchangeID or PresentationExchangeID
	Type string `json:"type"` // either credential-offer or present-proof
}

type CreateOutOfBandInvitationRequest struct {
	Attachments      []Attachment `json:"attachments,omitempty"` // When I put something in here it crashes, the CredentialExchange or PresentationExchange should probably be in the right state
	IncludeHandshake bool         `json:"include_handshake"`     // Invitation must include handshake protocols, request attachments, or both
	Metadata         struct{}     `json:"metadata,omitempty"`    // TODO
	UsePublicDID     bool         `json:"use_public_did,omitempty"`
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
	Type               string   `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/out-of-band/1.0/invitation
	ID                 string   `json:"@id"`
	Label              string   `json:"label"`
	HandshakeProtocols []string `json:"handshake_protocols"`
	Service            []struct {
		ID              string   `json:"id"`
		Type            string   `json:"type"`
		RecipientKeys   []string `json:"recipientKeys"`
		ServiceEndpoint string   `json:"serviceEndpoint"`
	} `json:"service"`
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

func (c *Client) ReceiveOutOfBandInvitation() {}
