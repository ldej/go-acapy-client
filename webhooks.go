package acapy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type WebhookHandlers struct {
	ConnectionsEventHandler            func(event Connection)
	BasicMessagesEventHandler          func(event BasicMessagesEvent)
	ProblemReportEventHandler          func(event ProblemReportEvent)
	CredentialExchangeEventHandler     func(event CredentialExchangeRecord)
	CredentialExchangeV2EventHandler   func(event CredentialExchangeRecordV2)
	CredentialExchangeDIFEventHandler  func(event CredentialExchangeDIF)
	CredentialExchangeIndyEventHandler func(event CredentialExchangeIndy)
	RevocationRegistryEventHandler     func(event RevocationRegistry)
	PresentationExchangeEventHandler   func(event PresentationExchangeRecord)
	CredentialRevocationEventHandler   func(event CredentialRevocationRecord)
	PingEventHandler                   func(event PingEvent)
	OutOfBandEventHandler              func(event OutOfBandEvent)
}

func CreateWebhooksHandler(handlers WebhookHandlers) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		topic := path[len(path)-1]

		defer r.Body.Close()

		switch topic {
		case "connections":
			if handlers.ConnectionsEventHandler != nil {
				var connectionsEvent Connection
				json.NewDecoder(r.Body).Decode(&connectionsEvent)
				handlers.ConnectionsEventHandler(connectionsEvent)
			}
		case "basicmessages":
			if handlers.BasicMessagesEventHandler != nil {
				var basicMessagesEvent BasicMessagesEvent
				json.NewDecoder(r.Body).Decode(&basicMessagesEvent)
				handlers.BasicMessagesEventHandler(basicMessagesEvent)
			}
		case "problem_report":
			if handlers.ProblemReportEventHandler != nil {
				var problemReportEvent ProblemReportEvent
				json.NewDecoder(r.Body).Decode(&problemReportEvent)
				handlers.ProblemReportEventHandler(problemReportEvent)
			}
		case "issue_credential":
			if handlers.CredentialExchangeEventHandler != nil {
				var credentialExchangeEvent CredentialExchangeRecord
				json.NewDecoder(r.Body).Decode(&credentialExchangeEvent)
				handlers.CredentialExchangeEventHandler(credentialExchangeEvent)
			}
		case "issuer_cred_rev":
			if handlers.CredentialRevocationEventHandler != nil {
				var credentialRevocationEvent CredentialRevocationRecord
				json.NewDecoder(r.Body).Decode(&credentialRevocationEvent)
				handlers.CredentialRevocationEventHandler(credentialRevocationEvent)
			}
		case "issue_credential_v2_0":
			if handlers.CredentialExchangeV2EventHandler != nil {
				var credentialExchangeV2Event CredentialExchangeRecordV2
				json.NewDecoder(r.Body).Decode(&credentialExchangeV2Event)
				handlers.CredentialExchangeV2EventHandler(credentialExchangeV2Event)
			}
		case "issue_credential_v2_0_dif":
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Println(string(body))
			if handlers.CredentialExchangeDIFEventHandler != nil {
				var credentialExchangeDIFEvent CredentialExchangeDIF
				if err := json.Unmarshal(body, &credentialExchangeDIFEvent); err != nil {
					log.Fatal(err)
				}
				handlers.CredentialExchangeDIFEventHandler(credentialExchangeDIFEvent)
			}
		case "issue_credential_v2_0_indy":
			if handlers.CredentialExchangeIndyEventHandler != nil {
				var credentialExchangeIndyEvent CredentialExchangeIndy
				json.NewDecoder(r.Body).Decode(&credentialExchangeIndyEvent)
				handlers.CredentialExchangeIndyEventHandler(credentialExchangeIndyEvent)
			}
		case "revocation_registry":
			if handlers.RevocationRegistryEventHandler != nil {
				var revocationRegistryEvent RevocationRegistry
				json.NewDecoder(r.Body).Decode(&revocationRegistryEvent)
				handlers.RevocationRegistryEventHandler(revocationRegistryEvent)
			}
		case "oob_invitation":
			if handlers.OutOfBandEventHandler != nil {
				var outOfBandEvent OutOfBandEvent
				json.NewDecoder(r.Body).Decode(&outOfBandEvent)
				handlers.OutOfBandEventHandler(outOfBandEvent)
			}
		case "present_proof":
			if handlers.PresentationExchangeEventHandler != nil {
				var presentationExchangeEvent PresentationExchangeRecord
				json.NewDecoder(r.Body).Decode(&presentationExchangeEvent)
				handlers.PresentationExchangeEventHandler(presentationExchangeEvent)
			}
		case "ping":
			if handlers.PingEventHandler != nil {
				var pingEvent PingEvent
				json.NewDecoder(r.Body).Decode(&pingEvent)
				handlers.PingEventHandler(pingEvent)
			}
		default:
			log.Printf("Webhook topic not supported: %q\n", topic)
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(200)
	}
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
