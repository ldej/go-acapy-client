# Running the example

## Start a VON-network ledger

Make a checkout of [github.com/bcgov/von-network](https://github.com/bcgov/von-network). Then run:

```shell script
./manage start --logs
```

This starts 4 Indy nodes and a von-webserver. The von-webserver has a web interface at [localhost:9000](http://localhost:9000) which allows you to browse the transactions in the blockchain.

Start a Tails server for the revocation registry tails files: Make a checkout of [github.com/bcgov/indy-tails-server](https://github.com/bcgov/indy-tails-server). Then run:

```shell script
./docker/manage start
```

## Starting instances

Make sure the command `aca-py` is available in the terminal where you execute the example.

Starting Alice

```shell script
$ go run examples/present_proof/main.go -port 4000 -name Alice
```

It will start a webserver for receiving webhook calls on port 4000. ACA-py will be started and will be available on port 4001 (that's the port ACA-py instances communicate over), and the ACA-py admin interface will be available on port 4002 (that's the port where Swagger is available and where your controller connects to).

Starting Bob

```shell script
$ go run examples/present_proof/main.go -port 4003 -name Bob
```

The same goes for Bob, but then with port 4003, 4004 and 4005.

## Present Proof

You can choose a couple of options:

```text
Choose:
	(1) Create invitation
	(2) Receive invitation
	(3) Register schema
	(4) Create credential definition
	(5) Issue credential
	(6) Send presentation proposal
	(7) Send presentation request
	(8) Send presentation
	(9) Verify presentation
	(10) List presentation proofs
	(exit) Exit

Enter Command: 

```

First, Alice creates an invitation. Then Bob should receive the invitation, copy Alice' invitation json to Bob. Then Alice registers a new schema, creates a credential definition, issues a credential.

Bob can send a presentation proposal, then Alice can send a presentation request, Bob can send the presentation, and Alice can verify it.