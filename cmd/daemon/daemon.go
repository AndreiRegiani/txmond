package daemon

import (
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
	log.Println("Daemon: polling")

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
	log.Printf("Daemon: processing block: %d", blockNumber)

	wallets, err := storage.Db.GetWallets()
	if err != nil {
		log.Printf("Daemon: db error: GetWallets: %v", err)
		return
	}

	for _, wallet := range wallets {
		processWallet(blockNumber, wallet)
	}
}

func processWallet(blockNumber uint64, wallet storage.Wallet) {
	log.Printf("Daemon: processing block: %d: wallet: %s", blockNumber, wallet.Address)
}
