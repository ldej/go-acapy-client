# go-acapy-client

A library for interacting with ACA-py in Go.

## Context

You can create your own Self-Sovereign Identity solution using the Hyperledger Ursa, Indy, Aries stack. Learn more about the background by watching these videos: 

- [The Story of Open SSI Standards](https://www.youtube.com/watch?v=RllH91rcFdE)
- [Decentralized Identifiers (DIDs) - The Fundamental Building Block of Self Sovereign Identity](https://www.youtube.com/watch?v=Jcfy9wd5bZI&)

To become an Aries developer, attend these courses by the Linux Foundation on [edx.org](https://edx.org):

- [Introduction to Hyperledger Sovereign Identity Blockchain Solutions: Indy, Aries & Ursa](https://courses.edx.org/courses/course-v1:LinuxFoundationX+LFS172x+3T2019/course/)
- [Becoming a Hyperledger Aries Developer](https://courses.edx.org/courses/course-v1:LinuxFoundationX+LFS173x+1T2020/course/) 

## Installation

```shell
$ go get -u github.com/ldej/go-acapy-client
```

## Compatibility

Both [ACA-py](https://github.com/hyperledger/aries-cloudagent-python) and `go-acapy-client` are under active development and might be incompatible. Currently `go-acapy-client` supports v0.6.0-pre of ACA-py.

## Development

Start a local Indy ledger network VON-network. Make a checkout of [github.com/bcgov/von-network](https://github.com/bcgov/von-network). Then run:

```shell script
./manage start --logs
```

This starts 4 Indy nodes and a von-webserver. The von-webserver has a web interface at [localhost:9000](http://localhost:9000) which allows you to browse the transactions in the blockchain.

Start a Tails server for the revocation registry tails files: Make a checkout of [github.com/bcgov/indy-tails-server](https://github.com/bcgov/indy-tails-server). Then run:

```shell script
./docker/manage start
```

Start an Aries-Cloud-Agent-Python (ACA-py) instance and configure the right command line parameters. Read about ACA-py and the command line parameters on my blog:

- [Becoming a Hyperledger Aries Developer - Part 1: Terminology](https://ldej.nl/post/becoming-aries-developer-part-1-terminology/)
- [Becoming a Hyperledger Aries Developer - Part 2: Development Environment](https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-2-development-environment/)
- [Becoming a Hyperledger Aries Developer - Part 3: Connecting using Swagger](https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-3-connecting-using-swagger/)
- [Becoming a Hyperledger Aries Developer - Part 3: Connecting using DIDComm Exchange](https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-3-connecting-using-didcomm-exchange/)
- [Becoming a Hyperledger Aries Developer - Part 4: Connecting using go-acapy-client](https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-4-connecting-using-go-acapy-client/)
- [Becoming a Hyperledger Aries Developer - Part 5: Issue Credentials](https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-5-issue-credentials/)
- [Becoming a Hyperledger Aries Developer - Part 6: Revocation](https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-6-revocation/)
- [Becoming a Hyperledger Aries Developer - Part 7: Present Proof](https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-7-present-proof/)
- [Connecting ACA-py to Development Ledgers](https://ldej.nl/post/connecting-acapy-to-development-ledgers/)
- [Aries Cloud Agent Python (ACA-py) Webhooks](https://ldej.nl/post/aries-cloudagent-python-webhooks/)

## Examples

Create a client, register a DID in the ledger and create an invitation.

```go
package main

import "github.com/ldej/go-acapy-client"

func main() {
    var ledgerURL = "http://localhost:9000"
    var acapyURL = "http://localhost:8000"
    client := acapy.NewClient(acapyURL)
    
    didResponse, err := acapy.RegisterDID(ledgerURL, "Alice", "000000000000000000000000MySeed01", acapy.Endorser)
    if err != nil {
    	// handle error
    }

    // Start aca-py with registered DID
    
    invitation, err := client.CreateInvitation("Bob", false, false, false)
    if err != nil {
        // handle error
    }
}
```

Examples can be found in the [examples](./examples) folder.

## Implemented Endpoints

### Action Menu

`{id}` = connection identifier

| Function Name | Method | Endpoint                    | Implemented   |
| ------------- | ------ | --------------------------- | ------------- |
| -             | POST   | /action-menu/{id}/close     | :exclamation: |
| -             | POST   | /action-menu/{id}/fetch     | :exclamation: |
| -             | POST   | /action-menu/{id}/perform   | :exclamation: |
| -             | POST   | /action-menu/{id}/request   | :exclamation: |
| -             | POST   | /action-menu/{id}/send-menu | :exclamation: |

### Basic Message

`{id}` = connection identifier

| Function Name    | Method | Endpoint                       | Implemented        |
| ---------------- | ------ | ------------------------------ | ------------------ |
| SendBasicMessage | POST   | /connections/{id}/send-message | :heavy_check_mark: |

### Connection

`{id}` = connection identifier

`{ref_id}` = inbound connection identifier

| Function Name     | Method | Endpoint                                     | Implemented        |
| ----------------- | ------ | -------------------------------------------- | ------------------ |
| QueryConnections  | GET    | /connections                                 | :heavy_check_mark: |
| CreateInvitation  | POST   | /connections/create-invitation               | :heavy_check_mark: |
| -                 | POST   | /connections/create-static                   | :exclamation:      |
| ReceiveInvitation | POST   | /connections/receive-invitation              | :heavy_check_mark: |
| GetConnection     | GET    | /connections/{id}                            | :heavy_check_mark: |
| RemoveConnection  | DELETE | /connections/{id}                            | :heavy_check_mark: |
| AcceptInvitation  | POST   | /connections/{id}/accept-invitation          | :heavy_check_mark: |
| AcceptRequest     | POST   | /connections/{id}/accept-request             | :heavy_check_mark: |
| -                 | POST   | /connections/{id}/establish-inbound/{ref_id} | :exclamation:      |
| -                 | GET    | /connections/{id}/metadata                   | :exclamation:      |
| -                 | POST   | /connections/{id}/metadata                   | :exclamation:      |

### Credential Definitions

`{id}` = credential definition identifier

| Function Name               | Method | Endpoint                        | Implemented        |
| --------------------------- | ------ | ------------------------------- | ------------------ |
| CreateCredentialDefinitions | POST   | /credential-definitions         | :heavy_check_mark: |
| QueryCredentialDefinitions  | GET    | /credential-definitions/created | :heavy_check_mark: |
| GetCredentialDefinition     | GET    | /credential-definitions/{id}    | :heavy_check_mark: |

### Credentials

`{id}` = credential identifier, also known as referent

| Function Name       | Method | Endpoint                    | Implemented        |
| ------------------- | ------ | --------------------------- | ------------------ |
| CredentialMimeTypes | GET    | /credential/mime-types/{id} | :heavy_check_mark: |
| IsCredentialRevoked | GET    | /credential/revoked/{id}    | :heavy_check_mark: |
| GetCredential       | GET    | /credential/{id}            | :heavy_check_mark: |
| RemoveCredential    | DELETE | /credential/{id}            | :heavy_check_mark: |
| GetCredentials      | GET    | /credentials                | :heavy_check_mark: |

### DID Exchange

`{id}` = connection identifier

| Function Name               | Method | Endpoint                            | Implemented        |
| --------------------------- | ------ | ----------------------------------- | ------------------ |
| DIDExchangeAcceptInvitation | POST   | /didexchange/{id}/accept-invitation | :heavy_check_mark: |
| DIDExchangeAcceptRequest    | POST   | /didexchange/{id}/accept-request    | :heavy_check_mark: |

### Introduction

`{id}` = connection identifier

| Function Name | Method | Endpoint                             | Implemented   |
| ------------- | ------ | ------------------------------------ | ------------- |
| -             | POST   | /connections/{id}/start-introduction | :exclamation: |

### Issue Credentials (Credential Exchange v1.0)

`{id}` = credential exchange identifier

| Function Name                   | Method | Endpoint                                      | Implemented        |
| ------------------------------- | ------ | --------------------------------------------- | ------------------ |
| CreateCredentialExchange        | POST   | /issue-credential/create                      | :heavy_check_mark: |
| QueryCredentialExchange         | GET    | /issue-credential/records                     | :heavy_check_mark: |
| GetCredentialExchange           | GET    | /issue-credential/records/{id}                | :heavy_check_mark: |
| RemoveCredentialExchange        | DELETE | /issue-credential/records/{id}                | :heavy_check_mark: |
| IssueCredentialByID             | POST   | /issue-credential/records/{id}/issue          | :heavy_check_mark: |
| ReportCredentialExchangeProblem | POST   | /issue-credential/records/{id}/problem-report | :heavy_check_mark: |
| SendCredentialOfferByID         | POST   | /issue-credential/records/{id}/send-offer     | :heavy_check_mark: |
| SendCredentialRequestByID       | POST   | /issue-credential/records/{id}/send-request   | :heavy_check_mark: |
| StoreReceivedCredential         | POST   | /issue-credential/records/{id}/store          | :heavy_check_mark: |
| SendCredential                  | POST   | /issue-credential/send                        | :heavy_check_mark: |
| SendCredentialOffer             | POST   | /issue-credential/send-offer                  | :heavy_check_mark: |
| SendCredentialProposal          | POST   | /issue-credential/send-proposal               | :heavy_check_mark: |

### Issue Credentials (Credential Exchange v2.0)

`{id}` = credential exchange identifier

| Function Name                     | Method | Endpoint                                          | Implemented        |
| --------------------------------- | ------ | ------------------------------------------------- | ------------------ |
| CreateCredentialExchangeV2        | POST   | /issue-credential-2.0/create                      | :heavy_check_mark: |
| QueryCredentialExchangeV2         | GET    | /issue-credential-2.0/records                     | :heavy_check_mark: |
| GetCredentialExchangeV2           | GET    | /issue-credential-2.0/records/{id}                | :heavy_check_mark: |
| RemoveCredentialExchangeV2        | DELETE | /issue-credential-2.0/records/{id}                | :heavy_check_mark: |
| IssueCredentialByIDV2             | POST   | /issue-credential-2.0/records/{id}/issue          | :heavy_check_mark: |
| ReportCredentialExchangeProblemV2 | POST   | /issue-credential-2.0/records/{id}/problem-report | :heavy_check_mark: |
| SendCredentialOfferByIDV2         | POST   | /issue-credential-2.0/records/{id}/send-offer     | :heavy_check_mark: |
| SendCredentialRequestByIDV2       | POST   | /issue-credential-2.0/records/{id}/send-request   | :heavy_check_mark: |
| StoreReceivedCredentialV2         | POST   | /issue-credential-2.0/records/{id}/store          | :heavy_check_mark: |
| SendCredentialV2                  | POST   | /issue-credential-2.0/send                        | :heavy_check_mark: |
| SendCredentialOfferV2             | POST   | /issue-credential-2.0/send-offer                  | :heavy_check_mark: |
| SendCredentialProposalV2          | POST   | /issue-credential-2.0/send-proposal               | :heavy_check_mark: |


### Ledger

| Function Name            | Method | Endpoint                          | Implemented        |
| ------------------------ | ------ | --------------------------------- | ------------------ |
| GetDIDEndpointFromLedger | GET    | /ledger/did-endpoint              | :heavy_check_mark: |
| GetDIDVerkeyFromLedger   | GET    | /ledger/did-verkey                | :heavy_check_mark: |
| GetDIDRoleFromLedger     | GET    | /ledger/get-nym-role              | :heavy_check_mark: |
| -                        | POST   | /ledger/register-nym              | :exclamation:      |
| -                        | PATCH  | /ledger/rotate-public-did-keypair | :exclamation:      |
| -                        | GET    | /ledger/taa                       | :exclamation:      |
| -                        | POST   | /ledger/taa/accept                | :exclamation:      |

### Mediation

`{id}` = connection identifier

`{mid}` = mediation identifier

| Function Name | Method | Endpoint                                     | Implemented   |
| ------------- | ------ | -------------------------------------------- | ------------- |
| -             | GET    | /mediation/default-mediator                  | :exclamation: |
| -             | DELETE | /mediation/default-mediator                  | :exclamation: |
| -             | GET    | /mediation/keylists                          | :exclamation: |
| -             | POST   | /mediation/keylists/{mid}/send-keylist-query | :exclamation: |
| -             | POST   | /mediation/keylists{mid}/send-keylist-update | :exclamation: |
| -             | POST   | /mediation/request/{id}                      | :exclamation: |
| -             | GET    | /mediation/requests                          | :exclamation: |
| -             | GET    | /mediation/requests/{mid}                    | :exclamation: |
|               | DELETE | /mediation/requests/{mid}                    | :exclamation: |
| -             | POST   | /mediation/requests/{mid}/deny               | :exclamation: |
| -             | POST   | /mediation/requests/{mid}/grant              | :exclamation: |
| -             | PUT    | /mediation/{mid}/default-mediator            | :exclamation: |

### Out-of-Band

| Function Name              | Method | Endpoint                        | Implemented        |
| -------------------------- | ------ | ------------------------------- | ------------------ |
| CreateOutOfBandInvitation  | POST   | /out-of-band/create-invitation  | :heavy_check_mark: |
| ReceiveOutOfBandInvitation | POST   | /out-of-band/receive-invitation | :heavy_check_mark: |

### Present Proof

`{id}` = presentation exchange identifier

| Function Name                  | Method | Endpoint                                        | Implemented        |
| ------------------------------ | ------ | ----------------------------------------------- | ------------------ |
| CreatePresentationRequest      | POST   | /present-proof/create-request                   | :heavy_check_mark: |
| QueryPresentationExchange      | GET    | /present-proof/records                          | :heavy_check_mark: |
| GetPresentationExchangeByID    | GET    | /present-proof/records/{id}                     | :heavy_check_mark: |
| RemovePresentationExchangeByID | DELETE | /present-proof/records/{id}                     | :heavy_check_mark: |
| GetPresentationCredentialsByID | GET    | /present-proof/records/{id}/credentials         | :heavy_check_mark: |
| SendPresentationRequestByID    | POST   | /present-proof/records/{id}/send-request        | :heavy_check_mark: |
| SendPresentationByID           | POST   | /present-proof/records/{id}/send-presentation   | :heavy_check_mark: |
| VerifyPresentationByID         | POST   | /present-proof/records/{id}/verify-presentation | :heavy_check_mark: |
| SendPresentationProposal       | POST   | /present-proof/send-proposal                    | :heavy_check_mark: |
| SendPresentationRequest        | POST   | /present-proof/send-request                     | :heavy_check_mark: |

### Revocation

`{id}` = revocation registry identifier, `{cred_def_id}` = credential definition identifier

| Function Name                       | Method | Endpoint                                  | Implemented        |
| ----------------------------------- | ------ | ----------------------------------------- | ------------------ |
| GetActiveRevocationRegistry         | GET    | /revocation/active-registry/{cred_def_id} | :heavy_check_mark: |
| ClearPendingRevocations             | POST   | /revocation/clear-pending-revocations     | :heavy_check_mark: |
| CreateRevocationRegistry            | POST   | /revocation/create-registry               | :heavy_check_mark: |
| GetCredentialRevocationStatus       | GET    | /revocation/credential-record             | :heavy_check_mark: |
| PublishRevocations                  | POST   | /revocation/publish-revocations           | :heavy_check_mark: |
| QueryRevocationRegistries           | GET    | /revocation/registries/created            | :heavy_check_mark: |
| GetRevocationRegistry               | GET    | /revocation/registry/{id}                 | :heavy_check_mark: |
| UpdateRevocationRegistryTailsURI    | PATCH  | /revocation/registry/{id}                 | :heavy_check_mark: |
| PublishRevocationRegistryDefinition | POST   | /revocation/registry/{id}/definition      | :heavy_check_mark: |
| PublishRevocationRegistryEntry      | POST   | /revocation/registry/{id}/entry           | :heavy_check_mark: |
| GetNumberOfIssuedCredentials        | GET    | /revocation/registry/{id}/issued          | :heavy_check_mark: |
| SetRevocationRegistryState          | PATCH  | /revocation/registry/{id}/set-state       | :heavy_check_mark: |
| UploadRegistryTailsFile             | PUT    | /revocation/registry/{id}/tails-file      | :heavy_check_mark: |
| DownloadRegistryTailsFile           | GET    | /revocation/registry/{id}/tails-file      | :heavy_check_mark: |
| RevokeIssuedCredential              | POST   | /revocation/revoke                        | :heavy_check_mark: |

### Schema

`{id}` = schema identifier

| Function Name  | Method | Endpoint         | Implemented        |
| -------------- | ------ | ---------------- | ------------------ |
| RegisterSchema | POST   | /schemas         | :heavy_check_mark: |
| QuerySchemas   | GET    | /schemas/created | :heavy_check_mark: |
| GetSchema      | GET    | /schemas/{id}    | :heavy_check_mark: |

### Server

| Function Name   | Method | Endpoint      | Implemented        |
| --------------- | ------ | ------------- | ------------------ |
| Features        | GET    | /features     | :heavy_check_mark: |
| Plugins         | GET    | /plugins      | :heavy_check_mark: |
| Shutdown        | GET    | /shutdown     | :heavy_check_mark: |
| Status          | GET    | /status       | :heavy_check_mark: |
| IsAlive         | GET    | /status/live  | :heavy_check_mark: |
| IsReady         | GET    | /status/ready | :heavy_check_mark: |
| ResetStatistics | POST   | /status/reset | :heavy_check_mark: |

### Trust ping

| Function Name | Method | Endpoint                    | Implemented        |
| ------------- | ------ | --------------------------- | ------------------ |
| SendPing      | POST   | /connections/{id}/send-ping | :heavy_check_mark: |


### Wallet

| Function Name            | Method | Endpoint                         | Implemented        |
| ------------------------ | ------ | -------------------------------- | ------------------ |
| QueryDIDs                | GET    | /wallet/did                      | :heavy_check_mark: |
| CreateLocalDID           | POST   | /wallet/did/create               | :heavy_check_mark: |
| RotateKeypair            | PATCH  | /wallet/did/local/rotate-keypair | :heavy_check_mark: |
| GetPublicDID             | GET    | /wallet/did/public               | :heavy_check_mark: |
| SetPublicDID             | POST   | /wallet/did/public               | :heavy_check_mark: |
| GetDIDEndpointFromWallet | GET    | /wallet/get-public-did           | :heavy_check_mark: |
| SetDIDEndpointInWallet   | POST   | /wallet/set-public-did           | :heavy_check_mark: |

### JSON-LD (unlisted in Swagger)

| Function Name | Method | Endpoint       | Implemented   |
| ------------- | ------ | -------------- | ------------- |
| SignJSONLD    | GET    | /jsonld/sign   | :exclamation: |
| VerifyJSONLD  | GET    | /jsonld/verify | :exclamation: |

## Webhooks

When an event occurs in ACA-py, for example a connection request has been received, a webhook is called on your controller on a certain topic. `go-acapy-client` provides a webhook handler where you can register your own functions to handle these events. Based on an event happening you can update your UI or inform the user about the event.

Read more about [ACA-py webhooks](https://ldej.nl/post/aries-cloudagent-python-webhooks/) on my blog.

```go
func main() {
    r := mux.NewRouter()
	webhookHandler := acapy.CreateWebhooksHandler(acapy.WebhookHandlers{
		ConnectionsEventHandler:            ConnectionsEventHandler,
		BasicMessagesEventHandler:          BasicMessagesEventHandler,
		ProblemReportEventHandler:          ProblemReportEventHandler,
		CredentialExchangeEventHandler:     CredentialExchangeEventHandler,
		CredentialExchangeV2EventHandler:   CredentialExchangeV2EventHandler,
		CredentialExchangeDIFEventHandler:  CredentialExchangeDIFEventHandler,
		CredentialExchangeIndyEventHandler: CredentialExchangeIndyEventHandler,
		RevocationRegistryEventHandler:     RevocationRegistryEventHandler,
		PresentationExchangeEventHandler:   PresentationExchangeEventHandler,
		CredentialRevocationEventHandler:   CredentialRevocationEventHandler,
		PingEventHandler:                   PingEventHandler,
		OutOfBandEventHandler:              OutOfBandEventHandler,
	})
    
    r.HandleFunc("/webhooks/topic/{topic}/", webhookHandler).Methods(http.MethodPost)
    
    // and so on
}

func ConnectionsEventHandler(event acapy.ConnectionsEvent) {
    fmt.Printf("\n -> Connection %q (%s), update to state %q rfc23 state %q\n", event.Alias, event.ConnectionID, event.State, event.RFC23State)
}

func BasicMessagesEventHandler(event acapy.BasicMessagesEvent) {
    fmt.Printf("\n -> Received message on connection %s: %s\n", event.ConnectionID, event.Content)
}

func ProblemReportEventHandler(event acapy.ProblemReportEvent) {
    fmt.Printf("\n -> Received problem report: %+v\n", event)
}

func CredentialExchangeEventHandler(event acapy.CredentialExchange) {
    fmt.Printf("\n -> Credential Exchange update: %s - %s\n", event.CredentialExchangeID, event.State)
}

func CredentialExchangeV2EventHandler(event acapy.CredentialExchangeRecordV2) {
	fmt.Printf("\n -> Credential Exchange V2 update: %s - %s\n", event.CredentialExchangeID, event.State)
}

func CredentialExchangeDIFEventHandler(event acapy.CredentialExchangeDIF) {
	fmt.Printf("\n -> Credential Exchange DIF Event: %s - %s", event.CredentialExchangeID, event.State)
}

func CredentialExchangeIndyEventHandler(event acapy.CredentialExchangeIndy) {
	fmt.Printf("\n -> Credential Exchange Indy Event: %s - %s", event.CredentialExchangeID, event.CredentialExchangeIndyID)
}

func RevocationRegistryEventHandler(event acapy.RevocationRegistry) {
    fmt.Printf("\n -> Revocation Registry update: %s - %s\n", event.RevocationRegistryID, event.State)
}

func PresentationExchangeEventHandler(event acapy.PresentationExchange) {
    fmt.Printf("\n -> Presentation Exchange update: %s - %s\n", event.PresentationExchangeID, event.State)
}

func CredentialRevocationEventHandler(event acapy.IssuerCredentialRevocationEvent) {
    fmt.Printf("\n -> Issuer Credential Revocation: %s - %s - %s\n", event.CredentialExchangeID, event.RecordID, event.State)
}

func PingEventHandler(event acapy.PingEvent) {
    fmt.Printf("\n -> Ping Event: %q state: %q responded: %t\n", event.ConnectionID, event.State, event.Responded)
}

func OutOfBandEventHandler(event acapy.OutOfBandEvent) {
    fmt.Printf("\n -> Out of Band Event: %q state %q\n", event.InvitationID, event.State)
}
```

You are free to choose the URL for your webhooks. Don't forget to set the command-line parameter for ACA-py: `--webhook-url http://localhost:{port}/webhooks`. The URL you provide to ACA-py is the base URL which will be extended with `/topic/{topic}` by default. So whatever URL you choose, make sure that:

- if the `--webhook-url` is `http://myhost:{port}/webhooks` 
- then the webhooks handler should listen on `http://myhost:{port}/webhooks/topic/{topic}`

The `acapy.WebhookHandler` is web framework agnostic and reads the topic from the URL by itself. The handler returned by `acapy.WebhookHandler` has the standard handler signature `func (w http.ResponseWriter, r *http.Request) {}`.

## TODO

- [ ] godoc
- [ ] Proper error handling
- [x] Admin API Key
- [x] Tracing via global config
- [ ] Automation of steps via global config
- [ ] Payment decorators https://github.com/hyperledger/aries-rfcs/tree/master/features/0075-payment-decorators
- [ ] Constructors for JSON-LD types
- [ ] Add types for roles, states, predicates
- [ ] Allow for a connection-less credential exchange
- [ ] Allow for a connection-less proof by making a QR code of a payload below. The base64 payload is the result of your call to  `/present-proof/create-request`.
```json
{
    "@id": "3b67c4bf-3953-4ace-94ef-28e0969288c5",
    "@type": "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/present-proof/1.0/request-presentation",
    "request_presentations~attach": [
        {
            "@id": "libindy-request-presentation-0",
            "mime-type": "application/json",
            "data": {
                "base64": "eyJuYW1lIjoiQ29udGFjdCBBZGRyZXNzIiwidmVyc2lvbiI6IjEuMC4wIiwibm9uY2UiOiIzOTIyNTEwMjk1NjY5MzcxNDIxNTIzMDgiLCJyZXF1ZXN0ZWRfYXR0cmlidXRlcyI6eyJHaXZlbiBOYW1lcyI6eyJuYW1lIjoiZ2l2ZW5fbmFtZXMiLCJyZXN0cmljdGlvbnMiOlt7InNjaGVtYV9pc3N1ZXJfZGlkIjoiODU0NTlHeGpOeVNKOEh3VFRRNHZxNyIsInNjaGVtYV9uYW1lIjoidmVyaWZpZWRfcGVyc29uIiwic2NoZW1hX3ZlcnNpb24iOiIxLjQuMCJ9LHsic2NoZW1hX25hbWUiOiJ1bnZlcmlmaWVkX3BlcnNvbiIsInNjaGVtYV92ZXJzaW9uIjoiMC4xLjAiLCJpc3N1ZXJfZGlkIjoiOFlxN0VoS0JNdWpoMjVOa0xHR2IydCJ9XX0sIkZhbWlseSBOYW1lIjp7Im5hbWUiOiJmYW1pbHlfbmFtZSIsInJlc3RyaWN0aW9ucyI6W3sic2NoZW1hX2lzc3Vlcl9kaWQiOiI4NTQ1OUd4ak55U0o4SHdUVFE0dnE3Iiwic2NoZW1hX25hbWUiOiJ2ZXJpZmllZF9wZXJzb24iLCJzY2hlbWFfdmVyc2lvbiI6IjEuNC4wIn0seyJzY2hlbWFfbmFtZSI6InVudmVyaWZpZWRfcGVyc29uIiwic2NoZW1hX3ZlcnNpb24iOiIwLjEuMCIsImlzc3Vlcl9kaWQiOiI4WXE3RWhLQk11amgyNU5rTEdHYjJ0In1dfSwiRGF0ZSBvZiBCaXJ0aCI6eyJuYW1lIjoiYmlydGhkYXRlIiwicmVzdHJpY3Rpb25zIjpbeyJzY2hlbWFfaXNzdWVyX2RpZCI6Ijg1NDU5R3hqTnlTSjhId1RUUTR2cTciLCJzY2hlbWFfbmFtZSI6InZlcmlmaWVkX3BlcnNvbiIsInNjaGVtYV92ZXJzaW9uIjoiMS40LjAifSx7InNjaGVtYV9uYW1lIjoidW52ZXJpZmllZF9wZXJzb24iLCJzY2hlbWFfdmVyc2lvbiI6IjAuMS4wIiwiaXNzdWVyX2RpZCI6IjhZcTdFaEtCTXVqaDI1TmtMR0diMnQifV19LCJTdHJlZXQgQWRkcmVzcyI6eyJuYW1lIjoic3RyZWV0X2FkZHJlc3MiLCJyZXN0cmljdGlvbnMiOlt7InNjaGVtYV9pc3N1ZXJfZGlkIjoiODU0NTlHeGpOeVNKOEh3VFRRNHZxNyIsInNjaGVtYV9uYW1lIjoidmVyaWZpZWRfcGVyc29uIiwic2NoZW1hX3ZlcnNpb24iOiIxLjQuMCJ9LHsic2NoZW1hX25hbWUiOiJ1bnZlcmlmaWVkX3BlcnNvbiIsInNjaGVtYV92ZXJzaW9uIjoiMC4xLjAiLCJpc3N1ZXJfZGlkIjoiOFlxN0VoS0JNdWpoMjVOa0xHR2IydCJ9XX0sIlBvc3RhbCBDb2RlIjp7Im5hbWUiOiJwb3N0YWxfY29kZSIsInJlc3RyaWN0aW9ucyI6W3sic2NoZW1hX2lzc3Vlcl9kaWQiOiI4NTQ1OUd4ak55U0o4SHdUVFE0dnE3Iiwic2NoZW1hX25hbWUiOiJ2ZXJpZmllZF9wZXJzb24iLCJzY2hlbWFfdmVyc2lvbiI6IjEuNC4wIn0seyJzY2hlbWFfbmFtZSI6InVudmVyaWZpZWRfcGVyc29uIiwic2NoZW1hX3ZlcnNpb24iOiIwLjEuMCIsImlzc3Vlcl9kaWQiOiI4WXE3RWhLQk11amgyNU5rTEdHYjJ0In1dfSwiQ2l0eSI6eyJuYW1lIjoibG9jYWxpdHkiLCJyZXN0cmljdGlvbnMiOlt7InNjaGVtYV9pc3N1ZXJfZGlkIjoiODU0NTlHeGpOeVNKOEh3VFRRNHZxNyIsInNjaGVtYV9uYW1lIjoidmVyaWZpZWRfcGVyc29uIiwic2NoZW1hX3ZlcnNpb24iOiIxLjQuMCJ9LHsic2NoZW1hX25hbWUiOiJ1bnZlcmlmaWVkX3BlcnNvbiIsInNjaGVtYV92ZXJzaW9uIjoiMC4xLjAiLCJpc3N1ZXJfZGlkIjoiOFlxN0VoS0JNdWpoMjVOa0xHR2IydCJ9XX0sIlByb3ZpbmNlIjp7Im5hbWUiOiJyZWdpb24iLCJyZXN0cmljdGlvbnMiOlt7InNjaGVtYV9pc3N1ZXJfZGlkIjoiODU0NTlHeGpOeVNKOEh3VFRRNHZxNyIsInNjaGVtYV9uYW1lIjoidmVyaWZpZWRfcGVyc29uIiwic2NoZW1hX3ZlcnNpb24iOiIxLjQuMCJ9LHsic2NoZW1hX25hbWUiOiJ1bnZlcmlmaWVkX3BlcnNvbiIsInNjaGVtYV92ZXJzaW9uIjoiMC4xLjAiLCJpc3N1ZXJfZGlkIjoiOFlxN0VoS0JNdWpoMjVOa0xHR2IydCJ9XX0sIkNvdW50cnkiOnsibmFtZSI6ImNvdW50cnkiLCJyZXN0cmljdGlvbnMiOlt7InNjaGVtYV9pc3N1ZXJfZGlkIjoiODU0NTlHeGpOeVNKOEh3VFRRNHZxNyIsInNjaGVtYV9uYW1lIjoidmVyaWZpZWRfcGVyc29uIiwic2NoZW1hX3ZlcnNpb24iOiIxLjQuMCJ9LHsic2NoZW1hX25hbWUiOiJ1bnZlcmlmaWVkX3BlcnNvbiIsInNjaGVtYV92ZXJzaW9uIjoiMC4xLjAiLCJpc3N1ZXJfZGlkIjoiOFlxN0VoS0JNdWpoMjVOa0xHR2IydCJ9XX19LCJyZXF1ZXN0ZWRfcHJlZGljYXRlcyI6e319"
            }
        }
    ],
    "comment": null,
    "~service": {
        "recipientKeys": [
            "F2fFPEXABoPKt8mYjNAavBwbsmQKYqNTcv3HKqBgqpLw"
        ],
        "routingKeys": null,
        "serviceEndpoint": "https://my-url.test.org"
    }
}
```