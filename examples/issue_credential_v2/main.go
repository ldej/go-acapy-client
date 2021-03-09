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
	"github.com/ldej/go-acapy-client"
)

type App struct {
	client    *acapy.Client
	server    *http.Server
	ledgerURL string
	port      int
	label     string
	seed      string
	rand      string
	myDID     string

	connection             acapy.Connection
	schema                 acapy.Schema
	credentialDefinitionID string
	credentialExchange     acapy.CredentialExchangeRecordResult
}

func (app *App) ReadCommands() {
	scanner := bufio.NewScanner(os.Stdin)

	didResponse, err := app.RegisterDID(app.label, app.label+app.rand)
	if err != nil {
		app.Exit(err)
	}
	app.myDID = didResponse.DID
	fmt.Printf("Hi %s, your registered DID is %s\n", app.label, didResponse.DID)

	for {
		fmt.Println(`Options:
	(1) Create invitation
	(2) Receive invitation
	(3) Register schema
	(4) Create credential definition
	(5) Send credential offer
	(6) Send credential request
	(7) Issue credential
	(8) Store credential
	(9) List credentials
	(10) Problem report
	(exit) Exit
`)
		fmt.Print("Choose: ")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "exit":
			app.Exit(nil)
		case "1":
			fmt.Println("Who/What is the invitation for?")
			scanner.Scan()
			theirLabel := scanner.Text()

			invitationResponse, err := app.client.CreateOutOfBandInvitation(
				acapy.CreateOutOfBandInvitationRequest{
					Alias:              theirLabel,
					HandshakeProtocols: acapy.DefaultHandshakeProtocols,
					MyLabel:            app.label,
				},
				true,
				false,
			)
			if err != nil {
				app.Exit(err)
			}
			invitation, _ := json.Marshal(invitationResponse.Invitation)
			fmt.Printf("Invitation json: %s\n", string(invitation))
		case "2":
			fmt.Println("Invitation json: ")
			scanner.Scan()
			invitation := scanner.Bytes()
			connection, err := app.ReceiveInvitation(invitation)
			if err != nil {
				app.Exit(err)
			}
			app.connection = connection
			fmt.Printf("Connection id: %s\n", connection.ConnectionID)
		case "3":
			fmt.Print("Schema name: ")
			scanner.Scan()
			schemaName := scanner.Text()

			fmt.Printf("Version: ")
			scanner.Scan()
			version := scanner.Text()

			fmt.Printf("Attributes (comma separated, i.e.: name,age): ")
			scanner.Scan()
			attributesString := scanner.Text()
			attributes := strings.Split(attributesString, ",")

			app.schema, err = app.RegisterSchema(schemaName, version, attributes)
			if err != nil {
				app.Exit(err)
			}
			fmt.Printf("Schema: %+v\n", app.schema)
		case "4":
			fmt.Println("This is slow, it takes a couple of seconds.")
			app.credentialDefinitionID, err = app.client.CreateCredentialDefinition("tag", false, 0, app.schema.ID)
			if err != nil {
				app.Exit(err)
			}
			fmt.Printf("Credential Definition ID: %s\n", app.credentialDefinitionID)
		case "5":
			fmt.Printf("Comment: ")
			scanner.Scan()
			comment := scanner.Text()

			var attributes []acapy.CredentialPreviewAttributeV2

			for _, attr := range app.schema.AttributeNames {
				fmt.Printf("Attribute %q value: ", attr)
				scanner.Scan()
				val := scanner.Text()

				attributes = append(attributes, acapy.CredentialPreviewAttributeV2{
					Name:     attr,
					MimeType: "text/plain",
					Value:    val,
				})
			}

			if credentialExchange, err := app.client.OfferCredentialV2(
				app.connection.ConnectionID,
				acapy.NewCredentialPreviewV2(attributes),
				app.credentialDefinitionID,
				comment,
			); err != nil {
				app.Exit(err)
			} else {
				app.credentialExchange = credentialExchange
			}
			fmt.Printf("Credential Exchange ID: %s\n", app.credentialExchange.CredentialExchangeRecord.CredentialExchangeID)
		case "6":
			_, err := app.client.RequestCredentialByIDV2(app.credentialExchange.CredentialExchangeRecord.CredentialExchangeID)
			if err != nil {
				app.Exit(err)
			}
		case "7":
			fmt.Printf("Comment: ")
			scanner.Scan()
			comment := scanner.Text()

			_, err := app.client.IssueCredentialByIDV2(app.credentialExchange.CredentialExchangeRecord.CredentialExchangeID, comment)
			if err != nil {
				app.Exit(err)
			}
		case "8":
			_, err = app.client.StoreCredentialByIDV2(app.credentialExchange.CredentialExchangeRecord.CredentialExchangeID, "")
			if err != nil {
				app.Exit(err)
			}
		case "9":
			credentials, err := app.client.GetCredentials(10, 0, "")
			if err != nil {
				app.Exit(err)
			}
			for _, cred := range credentials {
				fmt.Printf("%s - %s", cred.Referent, cred.Attributes)
			}
		case "10":
			fmt.Printf("Message: ")
			scanner.Scan()
			message := scanner.Text()

			err := app.client.ReportCredentialExchangeProblemV2(app.credentialExchange.CredentialExchangeRecord.CredentialExchangeID, message)
			if err != nil {
				app.Exit(err)
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
		"--genesis-url", fmt.Sprintf("%s/genesis", app.ledgerURL),
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
	webhookHandler := acapy.CreateWebhooksHandler(acapy.WebhookHandlers{
		ConnectionsEventHandler:            app.ConnectionsEventHandler,
		ProblemReportEventHandler:          app.ProblemReportEventHandler,
		CredentialExchangeV2EventHandler:   app.CredentialExchangeV2EventHandler,
		CredentialExchangeDIFEventHandler:  app.CredentialExchangeDIFEventHandler,
		CredentialExchangeIndyEventHandler: app.CredentialExchangeIndyEventHandler,
		OutOfBandEventHandler:              app.OutOfBandEventHandler,
	})

	r.HandleFunc("/webhooks/topic/{topic}/", webhookHandler).Methods(http.MethodPost)
	fmt.Printf("Listening on http://localhost:%d\n", app.port)
	fmt.Printf("ACA-py Admin API on http://localhost:%d\n", app.port+2)

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

func (app *App) Exit(err error) {
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func (app *App) ConnectionsEventHandler(event acapy.Connection) {
	if event.Alias == "" {
		connection, _ := app.client.GetConnection(event.ConnectionID)
		event.Alias = connection.TheirLabel
	}
	app.connection = event
	fmt.Printf("\n -> Connection %q (%s), update to state %q rfc23 state %q\n", event.Alias, event.ConnectionID, event.State, event.RFC23State)
}

func (app *App) CredentialExchangeV2EventHandler(event acapy.CredentialExchangeRecordV2) {
	connection, _ := app.client.GetConnection(event.ConnectionID)
	app.credentialExchange.CredentialExchangeRecord = event
	fmt.Printf("\n -> Credential Exchange V2 update: %s - %s - %s\n", event.CredentialExchangeID, connection.TheirLabel, event.State)
}

func (app *App) CredentialExchangeDIFEventHandler(event acapy.CredentialExchangeDIF) {
	record, _ := app.client.GetCredentialExchangeV2(event.CredentialExchangeID)
	connection, _ := app.client.GetConnection(record.CredentialExchangeRecord.ConnectionID)
	app.credentialExchange.DIF = event
	fmt.Printf("\n -> Credential Exchange DIF Event: %s - %s - %s", connection.TheirLabel, event.CredentialExchangeID, event.State)
}

func (app *App) CredentialExchangeIndyEventHandler(event acapy.CredentialExchangeIndy) {
	record, _ := app.client.GetCredentialExchangeV2(event.CredentialExchangeID)
	connection, _ := app.client.GetConnection(record.CredentialExchangeRecord.ConnectionID)
	app.credentialExchange.Indy = event
	fmt.Printf("\n -> Credential Exchange Indy Event: %s - %s - %s", connection.TheirLabel, event.CredentialExchangeID, event.CredentialExchangeIndyID)
}

func (app *App) ProblemReportEventHandler(event acapy.ProblemReportEvent) {
	fmt.Printf("\n -> Received problem report: %+v\n", event)
}

func (app *App) OutOfBandEventHandler(event acapy.OutOfBandEvent) {
	fmt.Printf("\n -> Out of Band Event: %q state %q\n", event.InvitationID, event.State)
}

func main() {
	var port = 4455
	var ledgerURL = "http://localhost:9000"
	var name = ""

	flag.IntVar(&port, "port", 4455, "port")
	flag.StringVar(&name, "name", "Alice", "alice")
	flag.Parse()

	acapyURL := fmt.Sprintf("http://localhost:%d", port+2)

	app := App{
		client:    acapy.NewClient(acapyURL),
		ledgerURL: ledgerURL,
		port:      port,
		label:     name,
		rand:      strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000)),
	}
	app.StartWebserver()
	app.ReadCommands()
}

func (app *App) RegisterDID(alias string, seed string) (acapy.RegisterDIDResponse, error) {
	didResponse, err := acapy.RegisterDID(
		app.ledgerURL+"/register",
		alias,
		seed,
		acapy.Endorser,
	)
	if err != nil {
		return acapy.RegisterDIDResponse{}, err
	}
	app.label = alias
	app.seed = didResponse.Seed
	app.StartACApy()
	return didResponse, nil
}

func (app *App) ReceiveInvitation(inv []byte) (acapy.Connection, error) {
	var invitation acapy.OutOfBandInvitation
	err := json.Unmarshal(inv, &invitation)
	if err != nil {
		return acapy.Connection{}, err
	}
	return app.client.ReceiveOutOfBandInvitation(invitation, true)
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
