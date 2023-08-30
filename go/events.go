package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/binocarlos/web3-events-example/go/contract"
)

type Config struct {
	PrivateKey string
	RpcURL     string
	ChainID    *big.Int
}

func getConfig() Config {
	privateKey := os.Getenv("PRIVATE_KEY")
	rpcURL := os.Getenv("RPC_URL")
	chainIDString := os.Getenv("CHAIN_ID")

	if privateKey == "" {
		fmt.Println("PRIVATE_KEY not set")
		os.Exit(1)
	}

	if rpcURL == "" {
		fmt.Println("RPC_URL not set")
		os.Exit(1)
	}

	if chainIDString == "" {
		fmt.Println("CHAIN_ID not set")
		os.Exit(1)
	}

	// convert chainID to an int and error if not
	chainID, success := big.NewInt(0).SetString(chainIDString, 10)
	if !success {
		fmt.Println("CHAIN_ID is not a valid integer")
		os.Exit(1)
	}

	return Config{
		PrivateKey: strings.Replace(privateKey, "0x", "", 1),
		RpcURL:     rpcURL,
		ChainID:    chainID,
	}
}

func main() {
	config := getConfig()

	client, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, config.ChainID)
	if err != nil {
		log.Fatal(err)
	}

	address, _, instance, err := contract.DeployContract(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Contract deloyed: 0x%x\n", address.Hex())

	go func() {

		num := 0

		for {
			num++
			tx, err := instance.SetNumber(auth, big.NewInt(int64(num)))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("tx: 0x%x\n", tx.Hash())
			time.Sleep(time.Second * 1)
		}

	}()

	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("saw event %+v\n", vLog.Data)
		}
	}

}
