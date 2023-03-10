package chain

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Test() {
	fmt.Println("gogo ")
	client, err := ethclient.Dial("wss://polygon-mainnet.g.alchemy.com/v2/4pp0CAODUCRqF3NueIbPh3MNj_R0wYJv")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xe290d1be511b42c6fc92b0714b79e00a0b6f3af9")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
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
			fmt.Println(vLog) // pointer to event log
		}
	}
}
