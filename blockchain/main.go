package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BlockNumber(address string) string {

	client, err := ethclient.Dial(address)
	if err != nil {
		log.Fatalf("Oops! There was a problem", err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return ""
	}
	return string(header.Number.String())
}

func Transaction(block_number, address string) (*types.Transaction, int) {
	client, err := ethclient.Dial(address)
	if err != nil {
		log.Fatalf("Oops! There was a problem", err)
	}
	block_int, _ := strconv.ParseInt(block_number, 10, 64)
	blockNumber := big.NewInt(block_int)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	var txr *types.Transaction
	var status uint64
	for _, tx := range block.Transactions() {

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(receipt.Status) // 1
		txr, status = tx, receipt.Status
		break
	}
	return txr, int(status)

}

type Post struct {
	Id      string `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

func ChainID(address string) string {
	posturl := address

	// JSON body
	body := []byte(`{
		"jsonrpc": "2.0",
		"method": "net_version",
		"params": [],
		"id": "getblock.io"
	}`)

	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	post := &Post{}
	err = json.NewDecoder(res.Body).Decode(post)
	if err != nil {
		panic(err)
	}
	return post.Result
}
