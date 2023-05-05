package daemon

import (
	"fmt"
	"log"
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
		log.Printf("Daemon: db error: GetWallets: %v", err)
		return
	}

	log.Printf("Daemon: monitoring %d wallets", len(wallets))

	if len(wallets) == 0 {
		return
	}

	lastBlock, err := storage.Db.GetLastProcessedBlock()
	if err != nil {
		log.Printf("Daemon: db error: GetLastProcessedBlock: %v", err)
		return
	}

	currentBlock, err := blockchain.BlockNumber()
	if err != nil {
		log.Printf("Daemon: blockchain error: BlockNumber: %v", err)
		return
	}

	log.Printf("Daemon: lastBlock: %d | currentBlock: %d", lastBlock, currentBlock)

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

	storage.Db.SetLastProcessedBlock(currentBlock)
}

func processBlock(blockNumber uint64) {
	log.Printf("Daemon: process block: %d", blockNumber)

	wallets, err := storage.Db.GetWallets()
	if err != nil {
		log.Printf("Daemon: db error: GetWallets: %v", err)
		return
	}

	transactions, err := blockchain.GetTransactions(blockNumber)
	if err != nil {
		log.Printf("Daemon: blockchain error: GetTransactions: %v", err)
		return
	}

	for _, transaction := range transactions {
		for _, wallet := range wallets {
			if transaction.From == wallet.Address || transaction.To == wallet.Address {
				fmt.Printf("Daemon: found transaction for wallet: %s", wallet.Address)

				err = storage.Db.InsertTransaction(wallet.Address, transaction)
				if err != nil {
					log.Printf("Daemon: db error: InsertTransaction: %v", err)
					return
				}
			}
		}
	}
}
