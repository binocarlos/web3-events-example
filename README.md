# web3-events-example

A simple example of a smart contract and emits events and 2 clients (on in node.js and the other in go) that will deploy and loop over the events emitted by transactions that we submit.

This is intended as a check that an EVM compatible chain is working with event filters as expected.

## install

You will need:

 * docker
 * node.js (>=16)
 * go (>=1.16)

## geth

If you want to test against a local geth node, you can run:

```bash
./stack geth
./stack geth-fund
```

This will start a local geth node and print out the environment variables that you need to export for the various clients.

## other chains

If you want to connect to other chains you can export the following variables:

 * `PRIVATE_KEY` - the private key of the wallet with funds to submit transactions
 * `RPC_URL` - the url of the RPC endpoint
 * `CHAIN_ID` - the chain id of the chain you are connecting to

## node.js

To run the node.js client:

```bash
export PRIVATE_KEY=xxx
export RPC_URL=xxx
export CHAIN_ID=xxx
cd js
npm install
node events.js
```

## go

To run the go client:

```bash
export PRIVATE_KEY=xxx
export RPC_URL=xxx
export CHAIN_ID=xxx
cd go
go run events.go
```

## re-compile

The artifacts folder contains the ABI and the bytecode of the contract.

The `go/contract/contract.go` file contains the golang bindings for said contract.

If you want to re-generate these because the smart contract has changed:

```bash
./stack compile
```
