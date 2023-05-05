package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
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

func BlockNumber() (uint64, error) {
	jsonReq := &JsonRpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []string{},
		ID:      1,
	}

	reqBytes, err := json.Marshal(jsonReq)
	if err != nil {
		return 0, errors.New("failure getting block number")
	}

	ethereumRPCURL := os.Getenv("ETHEREUM_RPC_URL")

	resp, err := http.Post(ethereumRPCURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return 0, errors.New("failure getting block number")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("failure getting block number")
	}

	var jsonResp BlockNumberResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return 0, errors.New("failure getting block number")
	}

	// Convert hex string to uint64: "0x4b7" -> 1207
	blockNum, err := strconv.ParseUint(jsonResp.Result[2:], 16, 64)
	if err != nil {
		return 0, errors.New("failure getting block number")
	}

	return blockNum, nil
}
