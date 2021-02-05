package acapy

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func WebhookHandler(
	connectionsEventHandler func(event Connection),
	basicMessagesEventHandler func(event BasicMessagesEvent),
	problemReportEventHandler func(event ProblemReportEvent),
	credentialExchangeEventHandler func(event CredentialExchangeRecord),
	revocationRegistryEventHandler func(event RevocationRegistry),
	presentationExchangeEventHandler func(event PresentationExchangeRecord),
	credentialRevocationEventHandler func(event CredentialRevocationRecord),
	pingEventHandler func(event PingEvent),
	outOfBandEventHandler func(event OutOfBandEvent),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		topic := path[len(path)-1]

		switch topic {
		case "connections":
			if connectionsEventHandler != nil {
				var connectionsEvent Connection
				json.NewDecoder(r.Body).Decode(&connectionsEvent)
				connectionsEventHandler(connectionsEvent)
			}
		case "basicmessages":
			if basicMessagesEventHandler != nil {
				var basicMessagesEvent BasicMessagesEvent
				json.NewDecoder(r.Body).Decode(&basicMessagesEvent)
				basicMessagesEventHandler(basicMessagesEvent)
			}
		case "problem_report":
			if problemReportEventHandler != nil {
				var problemReportEvent ProblemReportEvent
				json.NewDecoder(r.Body).Decode(&problemReportEvent)
				problemReportEventHandler(problemReportEvent)
			}
		case "issue_credential":
			if credentialExchangeEventHandler != nil {
				var credentialExchangeEvent CredentialExchangeRecord
				json.NewDecoder(r.Body).Decode(&credentialExchangeEvent)
				credentialExchangeEventHandler(credentialExchangeEvent)
			}
		case "issuer_cred_rev":
			if credentialRevocationEventHandler != nil {
				var credentialRevocationEvent CredentialRevocationRecord
				json.NewDecoder(r.Body).Decode(&credentialRevocationEvent)
				credentialRevocationEventHandler(credentialRevocationEvent)
			}
		case "revocation_registry":
			if revocationRegistryEventHandler != nil {
				var revocationRegistryEvent RevocationRegistry
				json.NewDecoder(r.Body).Decode(&revocationRegistryEvent)
				revocationRegistryEventHandler(revocationRegistryEvent)
			}
		case "oob_invitation":
			if outOfBandEventHandler != nil {
				var outOfBandEvent OutOfBandEvent
				json.NewDecoder(r.Body).Decode(&outOfBandEvent)
				outOfBandEventHandler(outOfBandEvent)
			}
		case "present_proof":
			if presentationExchangeEventHandler != nil {
				var presentationExchangeEvent PresentationExchangeRecord
				json.NewDecoder(r.Body).Decode(&presentationExchangeEvent)
				presentationExchangeEventHandler(presentationExchangeEvent)
			}
		case "ping":
			if pingEventHandler != nil {
				var pingEvent PingEvent
				json.NewDecoder(r.Body).Decode(&pingEvent)
				pingEventHandler(pingEvent)
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
