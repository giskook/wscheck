package wscheck

import (
	"fmt"
	"log"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"sync"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	testnetWsUrl    = "wss://testws.xlayer.tech"
	testnetContract = "0xBec7859BC3d0603BeC454F7194173E36BF2Aa5C8"
	wsUrl           = testnetWsUrl
	contractAddress = testnetContract
	erc20ABI        = `[{"constant":false,"inputs":[{"indexed":false,"internalType":"string","name":"name","type":"string"},{"indexed":false,"internalType":"string","name":"symbol","type":"string"},{"indexed":false,"internalType":"uint8","name":"__decimals","type":"uint8"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","payable":false,"type":"event"},{"constant":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","payable":false,"type":"event"},{"constant":false,"inputs":[],"name":"DOMAIN_SEPARATOR","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"DOMAIN_TYPEHASH","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"PERMIT_TYPEHASH","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"VERSION","outputs":[{"internalType":"string","name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"bridgeAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"burn","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"deploymentChainId","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"mint","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"","type":"address"}],"name":"nonces","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"deadline","type":"uint256"},{"indexed":false,"internalType":"uint8","name":"v","type":"uint8"},{"indexed":false,"internalType":"bytes32","name":"r","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"s","type":"bytes32"}],"name":"permit","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
)

func WsConcurrent() {
	client, err := ethclient.Dial(wsUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	headers := make(chan *types.Header)
	_, err = client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new headers: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Fatalf("Failed to parse the ABI: %v", err)
	}
	wg := sync.WaitGroup{}
	for {
		header := <-headers
		fmt.Printf("New block: %v\n", header.Number.String())
		getMetaInfos(client, parsedABI, contractAddress, &wg)
	}
}

func getMetaInfos(client *ethclient.Client, parsedABI abi.ABI, contractAddress string, wg *sync.WaitGroup) {
	wg.Add(40)
	var decimals, symbol, name, totalSupply string
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			decimals = getDecimals(client, parsedABI, common.HexToAddress(contractAddress))
		}()
		go func() {
			defer wg.Done()
			symbol = getSymbol(client, parsedABI, common.HexToAddress(contractAddress))
		}()
		go func() {
			defer wg.Done()
			name = getName(client, parsedABI, common.HexToAddress(contractAddress))
		}()
		go func() {
			defer wg.Done()
			totalSupply = getTotalSupply(client, parsedABI, common.HexToAddress(contractAddress))
		}()
	}
	wg.Wait()
	log.Printf("Decimals: %v, Symbol: %v, Name: %v, Total supply: %v\n", decimals, symbol, name, totalSupply)
}

func getDecimals(client *ethclient.Client, parsedABI abi.ABI, tokenAddress common.Address) string {
	data, err := parsedABI.Pack("decimals")
	if err != nil {
		log.Fatalf("Failed to pack the function call: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	var decimals uint8
	err = parsedABI.UnpackIntoInterface(&decimals, "decimals", result)
	if err != nil {
		log.Fatalf("Failed to unpack the result: %v", err)
	}
	return fmt.Sprintf("%d", decimals)
}

func getSymbol(client *ethclient.Client, parsedABI abi.ABI, tokenAddress common.Address) string {
	data, err := parsedABI.Pack("symbol")
	if err != nil {
		log.Fatalf("Failed to pack the function call: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	var symbol string
	err = parsedABI.UnpackIntoInterface(&symbol, "symbol", result)
	if err != nil {
		log.Fatalf("Failed to unpack the result: %v", err)
	}

	return symbol
}

func getName(client *ethclient.Client, parsedABI abi.ABI, tokenAddress common.Address) string {
	data, err := parsedABI.Pack("name")
	if err != nil {
		log.Fatalf("Failed to pack the function call: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	var name string
	err = parsedABI.UnpackIntoInterface(&name, "name", result)
	if err != nil {
		log.Fatalf("Failed to unpack the result: %v", err)
	}

	return name
}

func getTotalSupply(client *ethclient.Client, parsedABI abi.ABI, tokenAddress common.Address) string {
	data, err := parsedABI.Pack("totalSupply")
	if err != nil {
		log.Fatalf("Failed to pack the function call: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	var totalSupply *big.Int
	err = parsedABI.UnpackIntoInterface(&totalSupply, "totalSupply", result)
	if err != nil {
		log.Fatalf("Failed to unpack the result: %v", err)
	}
	return totalSupply.String()
}
