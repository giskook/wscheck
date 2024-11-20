package erc20balance

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	url            = ""
	transferTopics = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

var (
	usdtContract = common.HexToAddress("0x1e4a5963abfd975d8c9021ce480b42188849d41d")
)

func CheckBalance() {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	newHeight := make(chan *big.Int, 20)
	go func() {
		var h uint64
		for {
			select {
			case <-time.After(300 * time.Millisecond):
				cur, err := client.BlockNumber(context.Background())
				if err != nil {
					log.Printf("Error getting block number: %v", err)
					continue
				}
				if cur != h {
					newHeight <- big.NewInt(int64(cur))
					h = cur
					log.Println("block number:", cur)
				}
			}
		}
	}()

	for height := range newHeight {
		filterQuery := ethereum.FilterQuery{}
		filterQuery.FromBlock = height
		filterQuery.ToBlock = height
		filterQuery.Addresses = []common.Address{usdtContract}
		filterQuery.Topics = [][]common.Hash{{common.HexToHash(transferTopics)}}
		logs, _ := client.FilterLogs(context.Background(), filterQuery)
		for _, l := range logs {
			from := common.BytesToAddress(l.Topics[1][12:32])
			callMsg := ethereum.CallMsg{
				From: from,
				To:   &usdtContract,
				Data: manualEncodeBalanceOf(from),
			}
			res, err := client.CallContract(context.Background(), callMsg, nil)
			if err != nil {
				log.Println("Error calling contract:", err)
				continue
			}
			curBalance := new(big.Int).SetBytes(res)
			resprv, err := client.CallContract(context.Background(), callMsg, big.NewInt(height.Int64()-1))
			preBalance := new(big.Int).SetBytes(resprv)
			if preBalance == curBalance {
				log.Printf("------------------block number %v, address %v, cur balance %v", l.BlockNumber, from.String(), curBalance)
			}
			log.Printf("block number %v, address %v, cur balance %v", l.BlockNumber, from.String(), curBalance)
		}
	}
}

func manualEncodeBalanceOf(address common.Address) []byte {
	// Function selector for balanceOf(address): 0x70a08231
	selector := []byte{0x70, 0xa0, 0x82, 0x31}

	// Address padded to 32 bytes
	paddedAddress := common.LeftPadBytes(address.Bytes(), 32)

	// Concatenate the selector and the padded address
	return append(selector, paddedAddress...)
}
