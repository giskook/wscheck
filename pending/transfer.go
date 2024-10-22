package pending

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	url          = "https://testrpc.xlayer.tech/pay"
	chainID      = 54574
	privateKey   = ""
	gasPrice     = 2000000
	pendingCount = 100000
)

func Transfer() {
	cId := big.NewInt(chainID)
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	gp := big.NewInt(gasPrice)
	privateKeyHex := privateKey

	pk, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	toAddr := fromAddress
	nonce += 1
	for i := 0; i < pendingCount; i++ {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gp,
			Gas:      30000,
			To:       &toAddr,
			Value:    big.NewInt(1),
		})
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(cId), pk)
		ts := types.Transactions{signedTx}
		b := new(bytes.Buffer)
		ts.EncodeIndex(0, b)
		rawTxBytes := b.Bytes()
		txToSend := new(types.Transaction)
		rlp.DecodeBytes(rawTxBytes, &txToSend)

		err = client.SendTransaction(context.Background(), txToSend)
		if err != nil {
			log.Println("send transaction error", err)
			continue
		}
		fmt.Printf("nonce %v, tx sent: %s\n", nonce, signedTx.Hash().Hex())
		nonce++
	}
}
