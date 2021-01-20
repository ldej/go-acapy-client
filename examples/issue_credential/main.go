package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	acapy "github.com/ldej/go-acapy-client"
)

type App struct {
	client *acapy.Client
	server *http.Server
	port   int
	label  string
	seed   string
	rand   string
}

func (app *App) ReadCommands() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Your name: ")
	scanner.Scan()
	app.label = scanner.Text()

	fmt.Printf("Seed: ")
	scanner.Scan()
	app.seed = scanner.Text()

	didResponse, _ := app.RegisterDID(app.label, app.seed)
	fmt.Printf("Registered DID: %s\n", didResponse.DID)

	var schema acapy.Schema

	for {
		fmt.Println(`Choose:
	(1) Create invitation
	(2) Receive invitation
	(3) Register schema
	(4) Create credential definition
	(5) Send credential offer
	(6) Send credential request
	(7) Issue credential
	(8) Store credential
	(9) List credentials
	(exit) Exit
`)
		fmt.Print("Enter Command: ")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "exit":
			app.Exit()
			return
		case "1":
			fmt.Print("Alias: ")
			scanner.Scan()
			alias := scanner.Text()
			invitationResponse, _ := app.CreateInvitation(alias, true, false, true)
			invitation, _ := json.Marshal(invitationResponse.Invitation)
			fmt.Printf("Invitation json: %s\n", string(invitation))
		case "2":
			fmt.Print("Invitation json: ")
			scanner.Scan()
			invitation := scanner.Bytes()
			connection, _ := app.ReceiveInvitation(invitation)
			fmt.Printf("Connection ID: %s\n", connection.ConnectionID)
		case "3":
			fmt.Print("Schema name: ")
			scanner.Scan()
			schemaName := scanner.Text()
			fmt.Printf("Version: ")
			scanner.Scan()
			version := scanner.Text()
			fmt.Printf("Attributes: ")
			scanner.Scan()
			attributesString := scanner.Text()
			attributes := strings.Split(attributesString, ",")
			schema, _ = app.RegisterSchema(schemaName, version, attributes)
			fmt.Printf("Schema: %+v\n", schema)
		case "4":
			fmt.Printf("Tag: ")
			scanner.Scan()
			tag := scanner.Text()
			fmt.Printf("Schema ID: ")
			scanner.Scan()
			schemaID := scanner.Text()
			definition, _ := app.client.CreateCredentialDefinition(tag, false, 0, schemaID)
			fmt.Printf("Credential Definition ID: %s\n", definition)
		case "5":
			fmt.Printf("Credential Definition ID: ")
			scanner.Scan()
			credentialDefinitionID := scanner.Text()

			fmt.Printf("Connection ID: ")
			scanner.Scan()
			connectionID := scanner.Text()

			var attributes []acapy.CredentialPreviewAttribute

			for _, attr := range schema.AttributeNames {
				fmt.Printf("Attribute %q value: ", attr)
				scanner.Scan()
				val := scanner.Text()

				attributes = append(attributes, acapy.CredentialPreviewAttribute{
					Name:     attr,
					MimeType: "text/plain",
					Value:    val,
				})
			}

			fmt.Printf("Comment: ")
			scanner.Scan()
			comment := scanner.Text()

			var offer = acapy.CredentialOfferRequest{
				CredentialDefinitionID: credentialDefinitionID,
				ConnectionID:           connectionID,
				CredentialPreview: acapy.CredentialPreview{
					Type:       "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/1.0/credential-preview",
					Attributes: attributes,
				},
				Comment:    comment,
				Trace:      false,
				AutoRemove: false,
				AutoIssue:  false,
			}
			exchangeOffer, _ := app.client.SendCredentialOffer(offer)
			fmt.Printf("Credential Exchange ID: %s\n", exchangeOffer.CredentialExchangeID)
		case "6":
			fmt.Printf("Credential Exchange ID: ")
			scanner.Scan()
			credentialExchangeID := scanner.Text()
			_, _ = app.client.SendCredentialRequestByID(credentialExchangeID)
		case "7":
			fmt.Printf("Credential Exchange ID: ")
			scanner.Scan()
			credentialExchangeID := scanner.Text()

			fmt.Printf("Comment: ")
			scanner.Scan()
			comment := scanner.Text()

			_, _ = app.client.IssueCredentialByID(credentialExchangeID, comment)
		case "8":
			fmt.Printf("Credential Exchange ID: ")
			scanner.Scan()
			credentialExchangeID := scanner.Text()

			fmt.Printf("Credential ID: ")
			scanner.Scan()
			credentialID := scanner.Text()

			_, _ = app.client.StoreCredentialByID(credentialExchangeID, credentialID)
		case "9":
			credentials, _ := app.client.GetCredentials(10, 0, "")
			for _, cred := range credentials {
				fmt.Printf("%s - %s", cred.Referent, cred.Attributes)
			}
		}
	}
}

func (app *App) StartACApy() {
	id := strings.Replace(app.label+app.rand, " ", "_", -1)
	cmd := exec.Command("aca-py",
		"start",
		"--auto-provision",
		"-it", "http", "0.0.0.0", strconv.Itoa(app.port+1),
		"-ot", "http",
		"--admin", "0.0.0.0", strconv.Itoa(app.port+2),
		"--admin-insecure-mode",
		"--genesis-url", fmt.Sprintf("%s/genesis", app.client.LedgerURL),
		"--seed", app.seed,
		"--endpoint", fmt.Sprintf("http://localhost:%d/", app.port+1),
		"--webhook-url", fmt.Sprintf("http://localhost:%d/webhooks", app.port),
		"--label", app.label,
		"--public-invites",
		"--wallet-type", "indy",
		"--wallet-name", id,
		"--wallet-key", id,
		"--auto-accept-invites",
		"--auto-accept-requests",
		"--auto-ping-connection",
	)
	cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	go func() {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (app *App) StartWebserver() {
	r := mux.NewRouter()
	webhookHandler := acapy.WebhookHandler(
		app.ConnectionsEventHandler,
		app.BasicMessagesEventHandler,
		app.ProblemReportEventHandler,
		app.CredentialExchangeEventHandler,
		nil,
		nil,
	)

	r.HandleFunc("/webhooks/topic/{topic}/", webhookHandler).Methods(http.MethodPost)
	fmt.Printf("Listening on %v\n", app.port)

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: r,
	}

	go func() {
		if err := app.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (app *App) Exit() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func (app *App) ConnectionsEventHandler(event acapy.Connection) {
	if event.Alias == "" {
		connection, _ := app.client.GetConnection(event.ConnectionID)
		event.Alias = connection.TheirLabel
	}
	fmt.Printf("\n -> Connection %q (%s), update to state %q\n", event.Alias, event.ConnectionID, event.State)
}

func (app *App) BasicMessagesEventHandler(event acapy.BasicMessagesEvent) {
	connection, _ := app.client.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Received message from %q (%s): %s\n", connection.TheirLabel, event.ConnectionID, event.Content)
}

func (app *App) CredentialExchangeEventHandler(event acapy.CredentialExchange) {
	connection, _ := app.client.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Credential Exchange update: %s - %s - %s\n", event.CredentialExchangeID, connection.TheirLabel, event.State)
}

func (app *App) ProblemReportEventHandler(event acapy.ProblemReportEvent) {
	fmt.Printf("\n -> Received problem report: %+v\n", event)
}

func main() {
	var port = 4455
	var ledgerURL = "http://localhost:9000"

	flag.IntVar(&port, "port", 4455, "port")
	flag.Parse()

	acapyURL := fmt.Sprintf("http://localhost:%d", port+2)

	app := App{
		client: acapy.NewClient(ledgerURL, "", acapyURL),
		port:   port,
		rand:   strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000)),
	}
	app.StartWebserver()
	app.ReadCommands()
}

func (app *App) RegisterDID(alias string, seed string) (acapy.RegisterDIDResponse, error) {
	didResponse, err := app.client.RegisterDID(
		alias,
		seed,
		"ENDORSER", // TODO
	)
	if err != nil {
		return acapy.RegisterDIDResponse{}, err
	}
	app.label = alias
	app.seed = didResponse.Seed
	app.StartACApy()
	return didResponse, nil
}

func (app *App) CreateInvitation(alias string, autoAccept bool, multiUse bool, public bool) (acapy.CreateInvitationResponse, error) {
	invitationResponse, err := app.client.CreateInvitation(alias, autoAccept, multiUse, public)
	if err != nil {
		log.Printf("Failed to create invitation: %+v", err)
		return acapy.CreateInvitationResponse{}, err
	}
	return invitationResponse, nil
}

func (app *App) ReceiveInvitation(inv []byte) (acapy.Connection, error) {
	var invitation acapy.Invitation
	err := json.Unmarshal(inv, &invitation)
	if err != nil {
		return acapy.Connection{}, err
	}
	return app.client.ReceiveInvitation(invitation, true)
}

func (app *App) RegisterSchema(name string, version string, attributes []string) (acapy.Schema, error) {
	schema, err := app.client.RegisterSchema(
		name,
		version,
		attributes,
	)
	if err != nil {
		log.Printf("Failed to register schema: %+v", err)
		return acapy.Schema{}, err
	}
	fmt.Printf("Registered schema: %+v\n", schema)
	return schema, nil
}
