package chain

import (
	"accounts/api/pkg/config"
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func GetClient(net string) (*ethclient.Client, error) {
	var client *ethclient.Client
	if net == "kla" || net == "klay" || net == "klaytn" {
		headers := http.Header{}
		headers.Add("Authorization", "Basic S0FTS1VHRVkwT0dLVUlNM0lUTURCVzk1OmFHRUhCQ2pjWU5HMjlvRC1MSDhBV25mRUpzRDgtLTVYTTF2NmxiZE8=")
		headers.Add("x-chain-id", "1001")
		rpcClient, err := rpc.DialOptions(context.Background(), "https://node-api.klaytnapi.com/v1/klaytn", rpc.WithHeaders(headers))
		//rpc.DialHTTPWithClient("https://node-api.klaytnapi.com/v1/klaytn", httpClient)

		if err != nil {
			return nil, err
		}
		client = ethclient.NewClient(rpcClient)
	} else {
		url, err := config.GetRPC(net)
		if err != nil {
			return nil, errors.New("invalid net")
		}
		client, err = ethclient.DialContext(context.Background(), url)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func Transfer(pk string, to string, amount int64, service_id int32, net string, contractAddr string, isNative bool) (string, error) {
	if isNative {
		txHash, err := transferNativeBalance(pk, to, amount, net)
		if err != nil {
			return "", err
		}
		return txHash, nil
	} else {
		transferTokenBalance(pk, to, amount, net, contractAddr)
		return "", errors.New("ERC20 is not supported yet")
	}
}

func transferTokenBalance(pk string, to string, amount int64, net string, contract string) error {
	_, err := GetClient(net)

	if err != nil {
		return err
	}
	return nil
}

func transferNativeBalance(pk string, to string, amount int64, net string) (string, error) {
	client, err := GetClient(net)
	if err != nil {
		return "", err
	}

	fromAddress, privateKey, err := GetPrivateKey(pk)
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(to)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	gasLimit := uint64(21000)
	value := big.NewInt(amount)
	nonce, err := client.PendingNonceAt(context.Background(), *fromAddress)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	fmt.Printf("Transaction sent: %s\n %s", signedTx.Hash().Hex(), net)
	return signedTx.Hash().Hex(), nil
}
