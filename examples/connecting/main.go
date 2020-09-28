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

	for {
		fmt.Println(`Choose:
	(1) Register DID and start ACA-py
	(2) Register Schema
	(3) Create invitation
	(4) Receive invitation
	(5) Accept invitation
	(6) Accept request
	(7) Send basic message
	(8) Query connections
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
			fmt.Print("Seed: ")
			scanner.Scan()
			seed := scanner.Text() + app.rand
			didResponse, _ := app.RegisterDID(app.label, seed)
			fmt.Printf("Registered DID: %s\n", didResponse.DID)
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
			_, _ = app.AcceptInvitation(connectionID)
		case "6":
			fmt.Print("Connection id: ")
			scanner.Scan()
			connectionID := scanner.Text()
			_, _ = app.AcceptRequest(connectionID)
		case "7":
			fmt.Print("Connection id: ")
			scanner.Scan()
			connectionID := scanner.Text()
			fmt.Print("Message: ")
			scanner.Scan()
			message := scanner.Text()
			_ = app.SendBasicMessage(connectionID, message)
		case "8":
			connections, _ := app.QueryConnections(acapy.QueryConnectionsParams{})
			for _, connection := range connections {
				fmt.Printf("%s - %s - %s - %s\n", connection.Alias, connection.ConnectionID, connection.State, connection.TheirDID)
			}
		}
	}
}

func (app *App) StartACApy() {
	cmd := exec.Command("aca-py",
		"start",
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
		"--wallet-name", strings.Replace(app.label+app.rand, " ", "_", -1),
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

func (app *App) ConnectionsEventHandler(event acapy.ConnectionsEvent) {
	alias := event.Alias
	if alias == "" {
		connection, _ := app.client.GetConnection(event.ConnectionID)
		alias = connection.Alias
	}
	fmt.Printf("\n -> Connection %q (%s), update to state %q\n", alias, event.ConnectionID, event.State)
}

func (app *App) BasicMessagesEventHandler(event acapy.BasicMessagesEvent) {
	connection, _ := app.client.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Received message from %q (%s): %s\n", connection.Alias, event.ConnectionID, event.Content)
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

func (app *App) RegisterSchema() (acapy.Schema, error) {
	schemaResponse, err := app.client.RegisterSchema(
		"Laurence",
		"1.0",
		[]string{"name"},
	)
	if err != nil {
		log.Printf("Failed to register schema: %+v", err)
		return acapy.Schema{}, err
	}
	fmt.Printf("Registered schema: %+v\n", schemaResponse)
	return schemaResponse, nil
}

func (app *App) CreateInvitation(alias string, autoAccept bool, multiUse bool, public bool) (acapy.CreateInvitationResponse, error) {
	invitationResponse, err := app.client.CreateInvitation(alias, autoAccept, multiUse, public)
	if err != nil {
		log.Printf("Failed to create invitation: %+v", err)
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
	return app.client.ReceiveInvitation(invitation, false)
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

func (app *App) QueryConnections(params acapy.QueryConnectionsParams) ([]acapy.Connection, error) {
	connections, err := app.client.QueryConnections(params)
	if err != nil {
		log.Printf("Failed to list connections: %+v", err)
		return nil, err
	}
	return connections, nil
}
