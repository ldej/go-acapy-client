package acapy

import (
	"fmt"
	"strconv"
)

type CreateInvitationResponse struct {
	ConnectionID  string `json:"connection_id,omitempty"`
	InvitationURL string `json:"invitation_url,omitempty"`
	Invitation    struct {
		Type            string   `json:"@type,omitempty"`
		ID              string   `json:"@id,omitempty"`
		DID             string   `json:"did,omitempty"`
		ImageURL        string   `json:"imageUrl,omitempty"`
		Label           string   `json:"label,omitempty"`
		RecipientKeys   []string `json:"recipientKeys,omitempty"`
		RoutingKeys     []string `json:"routingKeys,omitempty"`
		ServiceEndpoint string   `json:"serviceEndpoint,omitempty"`
	} `json:"invitation,omitempty"`
}

func (c *Client) CreateInvitation(alias string, autoAccept bool, multiUse bool, public bool) (CreateInvitationResponse, error) {
	var createInvitationResponse CreateInvitationResponse
	queryParams := map[string]string{
		"alias":       alias,
		"auto_accept": strconv.FormatBool(autoAccept),
		"multi_use":   strconv.FormatBool(multiUse),
		"public":      strconv.FormatBool(public),
	}
	err := c.post(c.ACApyURL+"/connections/create-invitation", queryParams, nil, &createInvitationResponse)
	if err != nil {
		return CreateInvitationResponse{}, err
	}
	return createInvitationResponse, nil
}

type ReceiveInvitationResponse struct {
	Accept              string `json:"accept,omitempty"`
	Alias               string `json:"alias,omitempty"`
	ConnectionID        string `json:"connection_id,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
	ErrorMsg            string `json:"error_msg,omitempty"`
	InboundConnectionID string `json:"inbound_connection_id,omitempty"`
	Initiator           string `json:"initiator,omitempty"`
	InvitationKey       string `json:"invitation_key,omitempty"`
	InvitationMode      string `json:"invitation_mode,omitempty"`
	MyDID               string `json:"my_did,omitempty"`
	RequestID           string `json:"request_id,omitempty"`
	RoutingState        string `json:"routing_state,omitempty"`
	State               string `json:"state,omitempty"`
	TheirDID            string `json:"their_did,omitempty"`
	TheirLabel          string `json:"their_label,omitempty"`
	TheirRole           string `json:"their_role,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`
}

func (c *Client) ReceiveInvitation(invitation Invitation) (ReceiveInvitationResponse, error) {
	var receiveInvitationResponse ReceiveInvitationResponse
	err := c.post(c.ACApyURL+"/connections/receive-invitation", map[string]string{
		"alias":       invitation.Label,
		"auto_accept": strconv.FormatBool(false),
	}, invitation, &receiveInvitationResponse)
	if err != nil {
		return ReceiveInvitationResponse{}, err
	}
	return receiveInvitationResponse, nil
}

type Invitation struct {
	ID              string   `json:"@id,omitempty"`
	DID             string   `json:"did,omitempty"`
	ImageURL        string   `json:"imageUrl,omitempty"`
	Label           string   `json:"label,omitempty"`
	RecipientKeys   []string `json:"recipientKeys,omitempty"`
	RoutingKeys     []string `json:"routingKeys,omitempty"`
	ServiceEndpoint string   `json:"serviceEndpoint,omitempty"`
}

func (c *Client) AcceptInvitation(connectionID string) (Connection, error) {
	var connection Connection
	err := c.post(fmt.Sprintf("%s/connections/%s/accept-invitation", c.ACApyURL, connectionID), nil, nil, &connection)
	if err != nil {
		return Connection{}, err
	}
	return connection, nil
}

func (c *Client) AcceptRequest(connectionID string) (Connection, error) {
	var connection Connection
	err := c.post(fmt.Sprintf("%s/connections/%s/accept-request", c.ACApyURL, connectionID), nil, nil, &connection)
	if err != nil {
		return Connection{}, err
	}
	return connection, nil
}

type Connection struct {
	Accept              string `json:"accept"`
	Alias               string `json:"alias"`
	ConnectionID        string `json:"connection_id"`
	CreatedAt           string `json:"created_at"`
	ErrorMsg            string `json:"error_msg"`
	InboundConnectionID string `json:"inbound_connection_id"`
	Initiator           string `json:"initiator"`
	InvitationKey       string `json:"invitation_key"`
	InvitationMode      string `json:"invitation_mode"`
	MyDID               string `json:"my_did"`
	RequestID           string `json:"request_id"`
	RoutingState        string `json:"routing_state"`
	State               string `json:"state"`
	TheirDID            string `json:"their_did"`
	TheirLabel          string `json:"their_label"`
	TheirRole           string `json:"their_role"`
	UpdatedAt           string `json:"updated_at"`
}

// QueryConnectionsParams model
//
// Parameters for querying connections
//
type QueryConnectionsParams struct {

	// Alias of connection invitation
	Alias string `json:"alias,omitempty"`

	// Initiator is Connection invitation initiator
	Initiator string `json:"initiator,omitempty"`

	// Invitation key
	InvitationKey string `json:"invitation_key,omitempty"`

	// MyDID is DID of the agent
	MyDID string `json:"my_did,omitempty"`

	// State of the connection invitation
	State string `json:"state"`

	// TheirDID is other party's DID
	TheirDID string `json:"their_did,omitempty"`

	// TheirRole is other party's role
	TheirRole string `json:"their_role,omitempty"`
}

func (c *Client) QueryConnections(params QueryConnectionsParams) ([]Connection, error) {
	var results = struct {
		Connections []Connection `json:"results"`
	}{}

	var queryParams = map[string]string{
		"alias":            params.Alias,
		"initiator":        params.Initiator,
		"invitation_key":   params.InvitationKey,
		"my_did":           params.MyDID,
		"connection_state": params.State,
		"their_did":        params.TheirDID,
		"their_role":       params.TheirRole,
	}
	err := c.get(c.ACApyURL+"/results", queryParams, &results)
	if err != nil {
		return nil, err
	}
	return results.Connections, nil
}

func (c *Client) GetConnection(connectionID string) (Connection, error) {
	var connection Connection
	err := c.get(fmt.Sprintf("%s/connections/%s", c.ACApyURL, connectionID), nil, &connection)
	if err != nil {
		return Connection{}, err
	}
	return connection, nil
}

func (c *Client) RemoveConnection(connectionID string) error {
	return c.post(fmt.Sprintf("%s/connections/%s", c.ACApyURL, connectionID), nil, nil, nil)
}

type Thread struct {
	ThreadID string `json:"thread_id"`
}

func (c *Client) SendPing(connectionID string) (Thread, error) {
	ping := struct {
		Comment string `json:"comment"`
	}{
		Comment: "ping",
	}
	var thread Thread
	err := c.post(fmt.Sprintf("%s/connections/%s/send-ping", c.ACApyURL, connectionID), nil, ping, &thread)
	if err != nil {
		return Thread{}, err
	}
	return thread, nil
}
