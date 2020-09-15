package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	acapy "github.com/ldej/go-acapy-client"
)

type App struct {
	client *acapy.Client
	server *http.Server
	port   string
}

func (app *App) ReadCommands() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println(`Choose:
	(1) Register DID
	(2) Register Schema
	(3) Create invitation
	(4) Receive invitation
	(5) Accept invitation
	(6) Accept request
	(7) Send basic message
	(8) Query connections
	(exit) Exit`)
		fmt.Print("Enter Command: ")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "exit":
			app.Exit()
			return
		case "1":
			fmt.Print("Seed: ")
			scanner.Scan()
			seed := scanner.Text()
			app.RegisterDID(seed)
		case "2":
			app.RegisterSchema()
		case "3":
			fmt.Print("Alias: ")
			scanner.Scan()
			alias := scanner.Text()
			invitationResponse, _ := app.CreateInvitation(alias, false, false, true)
			invitation, _ := json.Marshal(invitationResponse.Invitation)
			fmt.Printf("Invitation json: %s\n", string(invitation))
		case "4":
			fmt.Print("Invitation json: ")
			scanner.Scan()
			invitation := scanner.Bytes()
			receiveInvitation, _ := app.ReceiveInvitation(invitation)
			fmt.Printf("Connection id: %s\n", receiveInvitation.ConnectionID)
		case "5":
			fmt.Print("Connection id: ")
			scanner.Scan()
			connectionID := scanner.Text()
			app.AcceptInvitation(connectionID)
		case "6":
			fmt.Print("Connection id: ")
			scanner.Scan()
			connectionID := scanner.Text()
			app.AcceptRequest(connectionID)
		case "7":
			fmt.Print("Connection id: ")
			scanner.Scan()
			connectionID := scanner.Text()
			fmt.Print("Message: ")
			scanner.Scan()
			message := scanner.Text()
			app.SendBasicMessage(connectionID, message)
		case "8":
			app.QueryConnections(acapy.QueryConnectionsParams{})
		}
	}
}

func (app *App) StartWebserver() {
	r := mux.NewRouter()
	webhooksHandler := acapy.Webhooks(
		app.ConnectionsEventHandler,
		app.BasicMessagesEventHandler,
		app.ProblemReportEventHandler,
	)

	r.HandleFunc("/webhooks/topic/{topic}/", webhooksHandler).Methods(http.MethodPost)
	fmt.Printf("Listening on %v\n", app.port)

	app.server = &http.Server{
		Addr:    ":" + app.port,
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

func (app *App) ConnectionsEventHandler(event acapy.ConnectionsEvent) {
	fmt.Printf("\n -> Connection %q (%s), update to state %q\n", event.Alias, event.ConnectionID, event.State)
}

func (app *App) BasicMessagesEventHandler(event acapy.BasicMessagesEvent) {
	connection, _ := app.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Received message from %q (%s): %s\n", connection.Alias, event.ConnectionID, event.Content)
}

func (app *App) ProblemReportEventHandler(event acapy.ProblemReportEvent) {
	fmt.Printf("\n -> Received problem report: %+v\n", event)
}

func main() {
	var port = "4455"
	var ledgerURL = "http://localhost:9000"
	var acapyURL = "http://localhost:11000"

	flag.StringVar(&port, "port", "4455", "port")
	flag.StringVar(&ledgerURL, "ledger", "http://localhost:9000", "Ledger URL")
	flag.StringVar(&acapyURL, "acapy", "http://localhost:11000", "ACA-py URL")
	flag.Parse()

	app := App{
		client: acapy.NewClient(ledgerURL, acapyURL),
		port:   port,
	}
	app.StartWebserver()
	app.ReadCommands()
}

func (app *App) RegisterDID(seed string) (acapy.RegisterDIDResponse, error) {
	didResponse, err := app.client.RegisterDID(
		"Laurence de Jong",
		seed,
		"TRUST_ANCHOR",
	)
	if err != nil {
		return acapy.RegisterDIDResponse{}, err
	}
	return didResponse, nil
}

func (app *App) RegisterSchema() acapy.SchemaResponse {
	schemaResponse, err := app.client.RegisterSchema(
		"Laurence",
		"1.0",
		[]string{"name"},
	)
	if err != nil {
		log.Fatal("Failed to register schema: ", err)
	}
	fmt.Printf("Registered schema: %+v\n", schemaResponse)
	return schemaResponse
}

func (app *App) CreateInvitation(alias string, autoAccept bool, multiUse bool, public bool) (acapy.CreateInvitationResponse, error) {
	invitationResponse, err := app.client.CreateInvitation(alias, autoAccept, multiUse, public)
	if err != nil {
		log.Fatal("Failed to create invitation: ", err)
		return acapy.CreateInvitationResponse{}, err
	}
	return invitationResponse, nil
}

func (app *App) ReceiveInvitation(inv []byte) (acapy.ReceiveInvitationResponse, error) {
	var invitation acapy.Invitation
	err := json.Unmarshal(inv, &invitation)
	if err != nil {
		return acapy.ReceiveInvitationResponse{}, err
	}
	return app.client.ReceiveInvitation(invitation)
}

func (app *App) AcceptInvitation(connectionID string) (acapy.Connection, error) {
	return app.client.AcceptInvitation(connectionID)
}

func (app *App) AcceptRequest(connectionID string) (acapy.Connection, error) {
	return app.client.AcceptRequest(connectionID)
}

func (app *App) SendPing(connectionID string) (acapy.Thread, error) {
	return app.client.SendPing(connectionID)
}

func (app *App) SendBasicMessage(connectionID string, message string) error {
	return app.client.SendBasicMessage(connectionID, message)
}

func (app *App) GetConnection(connectionID string) (acapy.Connection, error) {
	return app.client.GetConnection(connectionID)
}

func (app *App) QueryConnections(params acapy.QueryConnectionsParams) ([]acapy.Connection, error) {
	connections, err := app.client.QueryConnections(params)
	if err != nil {
		log.Fatal("Failed to list connections: ", err)
		return nil, err
	}
	indent, err := json.MarshalIndent(connections, "", "  ")
	if err != nil {
		return nil, err
	}
	fmt.Printf(string(indent) + "\n")
	return connections, nil
}
