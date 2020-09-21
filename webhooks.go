package acapy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func WebhookHandler(
	connectionsEventHandler func(event ConnectionsEvent),
	basicMessagesEventHandler func(event BasicMessagesEvent),
	problemReportEventHandler func(event ProblemReportEvent),
	credentialExchangeEventHandler func(event CredentialExchange),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		topic := path[len(path)-1]

		switch topic {
		case "connections":
			var connectionsEvent ConnectionsEvent
			json.NewDecoder(r.Body).Decode(&connectionsEvent)
			connectionsEventHandler(connectionsEvent)
		case "basicmessages":
			var basicMessagesEvent BasicMessagesEvent
			json.NewDecoder(r.Body).Decode(&basicMessagesEvent)
			basicMessagesEventHandler(basicMessagesEvent)
		case "problem_report":
			var problemReportEvent ProblemReportEvent
			json.NewDecoder(r.Body).Decode(&problemReportEvent)
			problemReportEventHandler(problemReportEvent)
		case "issue_credential":
			var credentialExchangeEvent CredentialExchange
			json.NewDecoder(r.Body).Decode(&credentialExchangeEvent)
			credentialExchangeEventHandler(credentialExchangeEvent)
		case "oob-invitation":
			// TODO
		case "present_proof":
			// TODO
		default:
			fmt.Printf("Topic not supported: %q\n", topic)
			w.WriteHeader(404)
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Printf(string(body))
			return
		}
		w.WriteHeader(200)
	}
}

type ConnectionsEvent struct {
	Initiator      string `json:"initiator"`
	CreatedAt      string `json:"created_at"`
	State          string `json:"state"`
	ConnectionID   string `json:"connection_id"`
	Accept         string `json:"accept"`
	Alias          string `json:"alias"`
	InvitationMode string `json:"invitation_mode"`
	UpdatedAt      string `json:"updated_at"`
	RoutingState   string `json:"routing_state"`
	InvitationKey  string `json:"invitation_key"`
}

type BasicMessagesEvent struct {
	ConnectionID string `json:"connection_id"`
	MessageID    string `json:"message_id"`
	State        string `json:"state"`
	Content      string `json:"content"`
}

type ProblemReportEvent struct {
	Type   string `json:"@type"`
	ID     string `json:"@id"`
	Thread struct {
		Thid string `json:"thid"`
	} `json:"~thread"`
	ExplainLtxt string `json:"explain-ltxt"`
}
