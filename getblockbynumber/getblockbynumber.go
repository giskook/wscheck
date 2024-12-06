package getblockbynumber

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	url    = "https://rpc.xlayer.tech"
	ticker = 10 * time.Millisecond
)

func GetBlockByNumber() {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	var h, cur uint64
	for {
		select {
		case <-time.After(ticker):
			cur, err = client.BlockNumber(context.Background())
			if err != nil {
				log.Printf("Error getting block number: %v", err)
				continue
			}
			if cur > h {
				h = cur
				log.Println("block number:", cur)
			}
			if cur < h {
				log.Printf("block number decreased: %v -> %v", h, cur)
			}
			_, err = client.BlockByNumber(context.Background(), big.NewInt(int64(cur)))
			if err != nil {
				log.Printf("Error getting block by number: %v", err)
				continue
			}
		}
	}
}
