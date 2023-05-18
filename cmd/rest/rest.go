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

	// Processed by the daemon
	lastProcessedBlock, _ := storage.Db.GetLastProcessedBlock()

	responseData := map[string]any{
		"lastProcessedBlock": lastProcessedBlock,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

	// Daemon will be processing this wallet
	err = storage.Db.InsertWallet(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := map[string]any{
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func walletTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("REST: GET /v1/ethereum/wallet/transactions/")

	// GET param: ?address=
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address parameter is required", http.StatusBadRequest)
		return
	}

	transactions, err := storage.Db.GetTransactionsByAddress(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := map[string]any{
		"transactions": transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Listen() {
	http.HandleFunc("/v1/ethereum/block/current/", blockCurrentHandler)
	http.HandleFunc("/v1/ethereum/wallet/", walletHandler)
	http.HandleFunc("/v1/ethereum/wallet/transactions/", walletTransactionsHandler)

	port, err := strconv.Atoi(os.Getenv("TXMOND_REST_API_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("REST: serving API on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
