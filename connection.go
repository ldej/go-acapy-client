package acapy

import (
	"fmt"
	"strconv"
)

type Connection struct {
	Accept              string `json:"accept"` // auto / manual
	Alias               string `json:"alias"`
	ConnectionID        string `json:"connection_id"`
	CreatedAt           string `json:"created_at"`
	ErrorMsg            string `json:"error_msg"`
	InboundConnectionID string `json:"inbound_connection_id"`
	InvitationKey       string `json:"invitation_key"`
	InvitationMode      string `json:"invitation_mode"` // once / multi
	InvitationMessageID string `json:"invitation_msg_id"`
	MyDID               string `json:"my_did"`
	RequestID           string `json:"request_id"`
	RFC23State          string `json:"rfc23_state"`
	RoutingState        string `json:"routing_state"`
	State               string `json:"state"`
	TheirDID            string `json:"their_did"`
	TheirLabel          string `json:"their_label"`
	TheirRole           string `json:"their_role"`
	UpdatedAt           string `json:"updated_at"`
}

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
	//var createInvitationRequest CreateInvitationRequest
	err := c.post("/connections/create-invitation", queryParams, nil, &createInvitationResponse)
	if err != nil {
		return CreateInvitationResponse{}, err
	}
	return createInvitationResponse, nil
}

func (c *Client) ReceiveInvitation(invitation Invitation, autoAccept bool) (Connection, error) {
	var connection Connection
	err := c.post("/connections/receive-invitation", map[string]string{
		"alias":       invitation.Label,
		"auto_accept": strconv.FormatBool(autoAccept),
	}, invitation, &connection)
	if err != nil {
		return Connection{}, err
	}
	return connection, nil
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
	err := c.post(fmt.Sprintf("/connections/%s/accept-invitation", connectionID), nil, nil, &connection)
	if err != nil {
		return Connection{}, err
	}
	return connection, nil
}

func (c *Client) AcceptRequest(connectionID string) (Connection, error) {
	var connection Connection
	err := c.post(fmt.Sprintf("/connections/%s/accept-request", connectionID), nil, nil, &connection)
	if err != nil {
		return Connection{}, err
	}
	return connection, nil
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

func (c *Client) QueryConnections(params *QueryConnectionsParams) ([]Connection, error) {
	var results = struct {
		Connections []Connection `json:"results"`
	}{}

	var queryParams = map[string]string{}
	if params != nil {
		queryParams = map[string]string{
			"alias":            params.Alias,
			"initiator":        params.Initiator,
			"invitation_key":   params.InvitationKey,
			"my_did":           params.MyDID,
			"connection_state": params.State,
			"their_did":        params.TheirDID,
			"their_role":       params.TheirRole,
		}
	}

	err := c.get("/connections", queryParams, &results)
	if err != nil {
		return nil, err
	}
	return results.Connections, nil
}

func (c *Client) GetConnection(connectionID string) (Connection, error) {
	var connection Connection
	err := c.get(fmt.Sprintf("/connections/%s", connectionID), nil, &connection)
	if err != nil {
		return Connection{}, err
	}
	return connection, nil
}

func (c *Client) RemoveConnection(connectionID string) error {
	return c.delete(fmt.Sprintf("/connections/%s", connectionID))
}

type Thread struct {
	ThreadID string `json:"thread_id"`
}

// Trust Ping
func (c *Client) SendPing(connectionID string) (Thread, error) {
	ping := struct {
		Comment string `json:"comment"`
	}{
		Comment: "ping",
	}
	var thread Thread
	err := c.post(fmt.Sprintf("/connections/%s/send-ping", connectionID), nil, ping, &thread)
	if err != nil {
		return Thread{}, err
	}
	return thread, nil
}

// TODO CreateStaticConnection EstablishInboundConnection
