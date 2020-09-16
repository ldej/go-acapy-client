# Running the example

## Start a VON-network ledger

Make a checkout of [github.com/bcgov/von-network](https://github.com/bcgov/von-network). Then run:

```shell script
./manage start --logs
```

This starts 4 Indy nodes and a von-webserver. The von-webserver has a web interface at [localhost:9000](http://localhost:9000) which allows you to browse the transactions in the blockchain.

## Starting instances

Make sure the command `aca-py` is available in the terminal where you execute the example.

Starting Alice

```shell script
$ go run examples/connecting/main.go --port 4000
```

It will start a webserver for receiving webhook calls on port 4000. ACA-py will be started and will be available on port 4001 (that's the port ACA-py instances communicate over), and the ACA-py admin interface will be available on port 4002 (that's the port where Swagger is available and where your controller connects to).

Starting Bob

```shell script
$ go run examples/connectin/main.go --port 4003
```

The same goes for Bob, but then with port 4003, 4004 and 4005.

## Connecting instances

After entering the name for your instance, you can choose a couple of options:

```text
Choose:
	(1) Register DID and start ACA-py
	(2) Register Schema
	(3) Create invitation
	(4) Receive invitation
	(5) Accept invitation
	(6) Accept request
	(7) Send basic message
	(8) Query connections
	(exit) Exit
```

First, register a DID and start ACA-py. Then, create an invitation with Alice. After that receive the invitation with Bob, then accept the invitation with Bob. Finally, accept the request with Alice. When the connection is established, a basic message can be sent from one agent to the other.