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
		case "issuer_cred_rev":
			// TODO
			// {
			//		"created_at": "2021-01-20 10:20:31.317083Z",
			//		"cred_def_id": "Vbs7XwP3vD1xc2JSCs4Hdx:3:CL:507:ForBob",
			//		"cred_rev_id": "1",
			//		"record_id": "a09c38fb-1daf-4923-a0c3-841b1696f6d1",
			//		"rev_reg_id": "Vbs7XwP3vD1xc2JSCs4Hdx:4:Vbs7XwP3vD1xc2JSCs4Hdx:3:CL:507:ForBob:CL_ACCUM:cfadbbf9-30c0-4a81-b46c-e723c61af4a9",
			//		"cred_ex_id": "b50b032f-9391-4e25-99b2-98cf45bb601d",
			//		"state": "issued",
			//		"updated_at": "2021-01-20 10:20:31.317083Z"
			// }
		case "revocation_registry":
			var revocationRegistryEvent RevocationRegistry
			json.NewDecoder(r.Body).Decode(&revocationRegistryEvent)
			revocationRegistryEventHandler(revocationRegistryEvent)
		case "oob_invitation":
			// TODO
		case "present_proof":
			var presentationExchangeEvent PresentationExchange
			json.NewDecoder(r.Body).Decode(&presentationExchangeEvent)
			presentationExchangeEventHandler(presentationExchangeEvent)
		case "ping":
			// TODO
			// {
			//   "comment": context.message.comment,
			//   "connection_id": context.message_receipt.connection_id,
			//   "responded": context.message.response_requested,
			//   "state": "received",
			//   "thread_id": context.message._thread_id,
			// }
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
