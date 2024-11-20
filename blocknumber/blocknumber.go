package blocknumber

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	url    = ""
	ticker = 10 * time.Millisecond
)

func BlockNumber() {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	var h uint64
	for {
		select {
		case <-time.After(ticker):
			cur, err := client.BlockNumber(context.Background())
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
		}
	}
}
