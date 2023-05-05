package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AndreiRegiani/txmond/cmd/storage"
)

func blockCurrentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("REST: GET /v1/ethereum/block/current/")

	lastProcessedBlock, _ := storage.Db.GetLastProcessedBlock()

	responseData := map[string]interface{}{
		"lastProcessedBlock": lastProcessedBlock,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func walletHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("REST: POST /v1/ethereum/wallet/")

	var requestBody struct {
		Address string `json:"address"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wallet := storage.Wallet{
		Address: requestBody.Address,
	}

	// Insert the wallet into the database
	err = storage.Db.InsertWallet(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success responseData
	responseData := map[string]interface{}{
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func walletTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("REST: GET /v1/ethereum/wallet/transactions/")
}

func Listen() {
	http.HandleFunc("/v1/ethereum/block/current/", blockCurrentHandler)
	http.HandleFunc("/v1/ethereum/wallet/", walletHandler)
	http.HandleFunc("/v1/ethereum/wallet/transactions/", walletTransactionsHandler)

	port, err := strconv.Atoi(os.Getenv("TXMOND_REST_API_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("REST: serving API on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
