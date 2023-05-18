package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AndreiRegiani/txmond/cmd/blockchain"
	"github.com/AndreiRegiani/txmond/cmd/storage"
)

func Start() {
	log.Println("Daemon: starting")

	pollingSeconds, err := strconv.Atoi(os.Getenv("ETHEREUM_RPC_POLLING_SECONDS"))
	if err != nil {
		log.Fatal("Daemon: error parsing ETHEREUM_RPC_POLLING_SECONDS")
	}

	for {
		polling()
		time.Sleep(time.Duration(pollingSeconds) * time.Second)
	}
}

func polling() {
	// Skip polling the blockchain if there are no wallets beiung monitored yet
	wallets, err := storage.Db.GetWallets()
	if err != nil {
		log.Printf("Daemon: db error: GetWallets: %v\n", err)
		return
	}

	log.Printf("Daemon: monitoring %d wallets\n", len(wallets))

	if len(wallets) == 0 {
		return
	}

	lastBlock, err := storage.Db.GetLastProcessedBlock()
	if err != nil {
		log.Printf("Daemon: db error: GetLastProcessedBlock: %v\n", err)
		return
	}

	currentBlock, err := blockchain.BlockNumber()
	if err != nil {
		log.Printf("Daemon: blockchain error: BlockNumber: %v\n", err)
		return
	}

	log.Printf("Daemon: lastBlock: %d | currentBlock: %d\n", lastBlock, currentBlock)

	// First run starts from the current block height
	if lastBlock == 0 {
		lastBlock = currentBlock - 1
	}

	// Process all new blocks since the last polling
	if lastBlock < currentBlock {
		for block := lastBlock + 1; block <= currentBlock; block++ {
			processBlock(block)
		}
	}

	err = storage.Db.SetLastProcessedBlock(currentBlock)
	if err != nil {
		log.Printf("Daemon: db error: SetLastProcessedBlock: %v\n", err)
		return
	}
}

func processBlock(blockNumber uint64) {
	log.Printf("Daemon: process block: %d\n", blockNumber)

	wallets, err := storage.Db.GetWallets()
	if err != nil {
		log.Printf("Daemon: db error: GetWallets: %v\n", err)
		return
	}

	transactions, err := blockchain.GetTransactions(blockNumber)
	if err != nil {
		log.Printf("Daemon: blockchain error: GetTransactions: %v\n", err)
		return
	}

	for _, transaction := range transactions {
		for _, wallet := range wallets {
			if transaction.From == wallet.Address || transaction.To == wallet.Address {
				log.Printf("Daemon: found transaction for wallet: %s\n", wallet.Address)

				err = storage.Db.InsertTransaction(wallet.Address, transaction)
				if err != nil {
					log.Printf("Daemon: db error: InsertTransaction: %v\n", err)
					return
				}

				// Send webhook if configured
				if webhookURL := os.Getenv("TXMOND_WEBHOOK_URL"); webhookURL != "" {
					go func(url string, tx storage.Transaction) {
						err := sendWebhook(url, tx)
						if err != nil {
							log.Printf("Daemon: webhook error: %v\n", err)
						}
					}(webhookURL, transaction)
				}
			}
		}
	}
}

func sendWebhook(url string, transaction storage.Transaction) error {
	payload, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send webhook, status code: %d", resp.StatusCode)
	}

	return nil
}
