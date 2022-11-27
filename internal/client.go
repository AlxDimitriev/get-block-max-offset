package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"
)

const greater = 1

type Client struct {
	httpClient        *http.Client
	apiKey           string
	AddressesBalance map[string]big.Int
	mu               sync.Mutex
}

func NewClient(apiKey string) *Client {
	return &Client{
		&http.Client{Timeout: 30 * time.Second},
		apiKey,
		make(map[string]big.Int),
		sync.Mutex{},
	}
}

func (c *Client) GetMaxBalanceOffset() (string, string) {
	var (
		maxOffsetAddr string
		maxBalanceOffset big.Int
		absMaxBalanceOffset big.Int
	)
	for addr, balance := range c.AddressesBalance {
		absBalance := new(big.Int)
		absBalance.Abs(&balance)
		if absBalance.Cmp(&absMaxBalanceOffset) == greater {
			absMaxBalanceOffset = *absBalance
			maxBalanceOffset = balance
			maxOffsetAddr = addr
		}
	}
	return maxOffsetAddr, maxBalanceOffset.String()
}

func (c *Client) GetLastBlock() int {
	reqBody := []byte(`{"jsonrpc": "2.0",
                        "method": "eth_blockNumber",
                        "params": [],
                        "id": "getblock.io"}`)
	reqBodyReader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, "https://eth.getblock.io/mainnet/", reqBodyReader)
	if err != nil {
		log.Fatalf("failed to create request")
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		r := LastBlockResp{}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			log.Fatalf("failed to decode last block number response: %v", err)
		}
		return hex2Int(r.Result)
	}
	return 0
}

func (c *Client) GetBlockTransactions(lastBlockNum, offset int, wg *sync.WaitGroup)  {
	defer wg.Done()
	blockNumHex := fmt.Sprintf("0x%x", lastBlockNum-offset)
	reqBody := []byte(fmt.Sprintf(`{"jsonrpc": "2.0",
                                           "method": "eth_getBlockByNumber",
                                           "params": ["%s", true],
                                           "id": "getblock.io"}`, blockNumHex))
	reqBodyReader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, "https://eth.getblock.io/mainnet/", reqBodyReader)
	if err != nil {
		log.Fatalf("failed to create request")
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	r := TransactionsBlockResp{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("failed to decode %s block info response: %v", blockNumHex, err)
	}

	for _, tr := range r.Result.Transactions {
		val := hex2BigInt(tr.Value)
		c.mu.Lock()
		if prevVal, ok :=  c.AddressesBalance[tr.From]; ok {
			newVal := new(big.Int)
			newVal.Sub(&prevVal, &val)
			c.AddressesBalance[tr.From] = *newVal
		} else {
			newVal := new(big.Int)
			newVal.Sub(big.NewInt(0), &val)
			c.AddressesBalance[tr.From] = *newVal
		}

		if prevVal, ok :=  c.AddressesBalance[tr.To]; ok {
			newVal := new(big.Int)
			newVal.Add(&prevVal, &val)
			c.AddressesBalance[tr.To] = *newVal
		} else {
			newVal := new(big.Int)
			newVal.Add(big.NewInt(0), &val)
			c.AddressesBalance[tr.To] = *newVal
		}
		c.mu.Unlock()
	}
}