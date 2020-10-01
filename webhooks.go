package acapy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func WebhookHandler(
	connectionsEventHandler func(event Connection),
	basicMessagesEventHandler func(event BasicMessagesEvent),
	problemReportEventHandler func(event ProblemReportEvent),
	credentialExchangeEventHandler func(event CredentialExchange),
	revocationRegistryEventHandler func(event RevocationRegistry),
	presentationExchangeEventHandler func(event PresentationExchange),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		topic := path[len(path)-1]

		switch topic {
		case "connections":
			var connectionsEvent Connection
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
		case "revocation_registry":
			var revocationRegistryEvent RevocationRegistry
			json.NewDecoder(r.Body).Decode(&revocationRegistryEvent)
			revocationRegistryEventHandler(revocationRegistryEvent)
		case "oob-invitation":
			// TODO
		case "present_proof":
			var presentationExchangeEvent PresentationExchange
			json.NewDecoder(r.Body).Decode(&presentationExchangeEvent)
			presentationExchangeEventHandler(presentationExchangeEvent)
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
