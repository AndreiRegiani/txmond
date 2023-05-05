package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/AndreiRegiani/txmond/cmd/storage"
)

type JsonRpcRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

type BlockNumberResponse struct {
	Result string `json:"result"`
}

type BlockResult struct {
	Result BlockResultTransactions `json:"result"`
}

type BlockResultTransactions struct {
	Transactions []storage.Transaction `json:"transactions"`
}

func BlockNumber() (uint64, error) {
	jsonReq := &JsonRpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []string{},
		ID:      1,
	}

	reqBytes, err := json.Marshal(jsonReq)
	if err != nil {
		return 0, errors.New("failure marshaling block number request")
	}

	ethereumRPCURL := os.Getenv("ETHEREUM_RPC_URL")

	resp, err := http.Post(ethereumRPCURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return 0, errors.New("failure posting to Ethereum RPC")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("failure reading response")
	}

	var jsonResp BlockNumberResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return 0, errors.New("failure unmarshaling response")
	}

	// Convert hex string to uint64: "0x4b7" -> 1207
	blockNum, err := strconv.ParseUint(jsonResp.Result[2:], 16, 64)
	if err != nil {
		return 0, errors.New("failure parsing block number")
	}

	return blockNum, nil
}

func GetTransactions(blockNumber uint64) ([]storage.Transaction, error) {
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)

	jsonReq := &JsonRpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []string{blockNumberHex, strconv.FormatBool(true)},
		ID:      1,
	}

	reqBytes, err := json.Marshal(jsonReq)
	if err != nil {
		return nil, errors.New("failure marshaling block transactions request")
	}

	ethereumRPCURL := os.Getenv("ETHEREUM_RPC_URL")

	resp, err := http.Post(ethereumRPCURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errors.New("failure posting to Ethereum RPC")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failure reading response")
	}

	var jsonResp BlockResult
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, errors.New("failure unmarshaling response")
	}

	return jsonResp.Result.Transactions, nil
}
