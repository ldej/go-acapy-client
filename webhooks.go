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
	issuerCredentialReceivedEventHandler func(event IssuerCredentialReceivedEvent),
	pingEventHandler func(event PingEvent),
	outOfBandEventHandler func(event OutOfBandEvent),
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
			body, _ := ioutil.ReadAll(r.Body) // TODO
			fmt.Println(string(body))
			var problemReportEvent ProblemReportEvent
			json.NewDecoder(r.Body).Decode(&problemReportEvent)
			problemReportEventHandler(problemReportEvent)
		case "issue_credential":
			body, _ := ioutil.ReadAll(r.Body) // TODO
			fmt.Println(string(body))
			var credentialExchangeEvent CredentialExchange
			json.NewDecoder(r.Body).Decode(&credentialExchangeEvent)
			credentialExchangeEventHandler(credentialExchangeEvent)
		case "issuer_cred_rev":
			body, _ := ioutil.ReadAll(r.Body) // TODO
			fmt.Println(string(body))
			var issuerCredentialReceivedEvent IssuerCredentialReceivedEvent
			json.NewDecoder(r.Body).Decode(&issuerCredentialReceivedEvent)
			issuerCredentialReceivedEventHandler(issuerCredentialReceivedEvent)
		case "revocation_registry":
			body, _ := ioutil.ReadAll(r.Body) // TODO
			fmt.Println(string(body))
			var revocationRegistryEvent RevocationRegistry
			json.NewDecoder(r.Body).Decode(&revocationRegistryEvent)
			revocationRegistryEventHandler(revocationRegistryEvent)
		case "oob_invitation":
			var outOfBandEvent OutOfBandEvent
			json.NewDecoder(r.Body).Decode(&outOfBandEvent)
			outOfBandEventHandler(outOfBandEvent)
		case "present_proof":
			body, _ := ioutil.ReadAll(r.Body) // TODO
			fmt.Println(string(body))
			var presentationExchangeEvent PresentationExchange
			json.NewDecoder(r.Body).Decode(&presentationExchangeEvent)
			presentationExchangeEventHandler(presentationExchangeEvent)
		case "ping":
			var pingEvent PingEvent
			json.NewDecoder(r.Body).Decode(&pingEvent)
			pingEventHandler(pingEvent)
		default:
			fmt.Printf("Topic not supported: %q\n", topic)
			w.WriteHeader(404)
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Println(string(body))
			return
		}
		w.WriteHeader(200)
	}
}

type IssuerCredentialReceivedEvent struct {
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
	CredentialDefinitionID string `json:"cred_def_id"`
	CredentialRevisionID   string `json:"1"`
	RecordID               string `json:"record_id"`
	RevocationRegistryID   string `json:"rev_reg_id"`
	CredentialExchangeID   string `json:"cred_ex_id"`
	State                  string `json:"state"`
}

type PingEvent struct {
	Comment      string `json:"comment"`
	ConnectionID string `json:"connection_id"`
	Responded    bool   `json:"responded"`
	State        string `json:"state"`
	ThreadID     string `json:"thread_id"`
}

type OutOfBandEvent struct {
	InvitationID        string              `json:"invitation_id"`
	InvitationMessageID string              `json:"invi_msg_id"`
	Invitation          OutOfBandInvitation `json:"invitation"`
	State               string              `json:"state"`
	InvitationURL       string              `json:"invitation_url"`
	UpdatedAt           string              `json:"updated_at"`
	CreatedAt           string              `json:"created_at"`
	AutoAccept          bool                `json:"auto_accept"`
	MultiUse            bool                `json:"multi_use"`
	Trace               bool                `json:"trace"`
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
