package contransfer

import (
	"log"
	"github.com/ethereum/go-ethereum/ethclient"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"fmt"
	"math/big"
	"crypto/ecdsa"
	"time"
)

const (
	rpcURL   = "wss://testws.xlayer.tech"
	gasPrice = 30000000000
	gasLimit = 21000
	chainID  = 1
)

var privateKeys = []string{"0x1234"}

func Boom() {
	// Connect to Ethereum client (local node or Infura)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	nonces := make([]uint64, 0, len(privateKeys))
	for _, pk := range privateKeys {
		nonces = append(nonces, getNonce(client, pk))
	}
	for i, pk := range privateKeys {
		privateKey := pk
		index := i
		go func() {
			err = client.SendTransaction(context.Background(), buildTx(privateKey, nonces[index]))
			if err != nil {
				log.Printf("Failed to send transaction:%v %v", privateKey, err)
			}
		}()
	}
}

func buildTx(pk string, nonce uint64) *types.Transaction {
	// Sender's private key
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Error casting public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("From address: ", fromAddress.Hex())

	// Specify the gas price
	gp := big.NewInt(gasPrice) // in wei (30 gwei)

	// Build the transaction
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &fromAddress,
		Gas:      gasLimit,
		GasPrice: gp,
		Value:    big.NewInt(0),
	})

	signedTx, err := types.SignTx(rawTx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	return signedTx
}

func getNonce(client *ethclient.Client, pk string) uint64 {
	// Sender's private key
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Error casting public key")
	}

	nonce := uint64(0)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	for {
		nonce, err = client.NonceAt(context.Background(), address, nil)
		if err != nil {
			log.Printf("Failed to get nonce: %v, %v", address, err)
			time.Sleep(time.Second)
		}
		break
	}

	return nonce
}
