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
	client       *acapy.Client
	server       *http.Server
	ledgerURL    string
	port         int
	label        string
	seed         string
	rand         string
	myDID        string
	connectionID string
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
	(3) Accept invitation
	(4) Accept request
	(5) Send ping
	(6) Send basic message
	(7) Query connections
	(exit) Exit
`)
		fmt.Print("Choose: ")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "exit":
			app.Exit(nil)
			return
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
				false,
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
			app.connectionID = connection.ConnectionID
			fmt.Printf("Connection id: %s\n", connection.ConnectionID)
		case "3":
			_, err := app.client.DIDExchangeAcceptInvitation(app.connectionID, "", "")
			if err != nil {
				app.Exit(err)
			}
		case "4":
			_, err := app.client.DIDExchangeAcceptRequest(app.connectionID, "")
			if err != nil {
				app.Exit(err)
			}
		case "5":
			_, err := app.client.SendPing(app.connectionID)
			if err != nil {
				app.Exit(err)
			}
		case "6":
			fmt.Print("Message: ")
			scanner.Scan()
			message := scanner.Text()

			err := app.client.SendBasicMessage(app.connectionID, message)
			if err != nil {
				app.Exit(err)
			}
		case "7":
			connections, err := app.client.QueryConnections(nil)
			if err != nil {
				app.Exit(err)
			}
			for _, connection := range connections {
				fmt.Printf("%s - %s - %s - %s\n", connection.TheirLabel, connection.ConnectionID, connection.State, connection.TheirDID)
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
		"--monitor-ping",
		"--auto-respond-messages",
		"--wallet-type", "indy",
		"--wallet-name", id,
		"--wallet-key", id,
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
		ConnectionsEventHandler:   app.ConnectionsEventHandler,
		BasicMessagesEventHandler: app.BasicMessagesEventHandler,
		ProblemReportEventHandler: app.ProblemReportEventHandler,
		PingEventHandler:          app.PingEventHandler,
		OutOfBandEventHandler:     app.OutOfBandEventHandler,
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
	app.connectionID = event.ConnectionID
	fmt.Printf("\n -> Connection %q (%s), update to state %q rfc23 state %q\n", event.Alias, event.ConnectionID, event.State, event.RFC23State)
}

func (app *App) BasicMessagesEventHandler(event acapy.BasicMessagesEvent) {
	connection, _ := app.client.GetConnection(event.ConnectionID)
	fmt.Printf("\n -> Received message from %q (%s): %s\n", connection.TheirLabel, event.ConnectionID, event.Content)
}

func (app *App) OutOfBandEventHandler(event acapy.OutOfBandEvent) {
	fmt.Printf("\n -> Out of Band Event: %q state %q\n", event.InvitationID, event.State)
}

func (app *App) PingEventHandler(event acapy.PingEvent) {
	fmt.Printf("\n -> Ping Event: %q state: %q responded: %t\n", event.ConnectionID, event.State, event.Responded)
}

func (app *App) ProblemReportEventHandler(event acapy.ProblemReportEvent) {
	fmt.Printf("\n -> Received problem report: %+v\n", event)
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
	return app.client.ReceiveOutOfBandInvitation(invitation, false)
}
