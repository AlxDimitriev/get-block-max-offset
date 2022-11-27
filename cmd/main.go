package main

import (
	"get_block_test/internal"
	"github.com/jessevdk/go-flags"
	"log"
	"sync"
)

func main()  {
	var opts struct {
		ApiKey string `long:"api_key" env:"API_KEY"`
		NumBlocksToProcess int `long:"num_blocks_to_process" env:"NUM_BLOCKS_TO_PROCESS" default:"100"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("Failed to parse args: %s", err)
	}
	client := internal.NewClient(opts.ApiKey)
	lastBlock := client.GetLastBlock()
	wg := sync.WaitGroup{}
	wg.Add(opts.NumBlocksToProcess)

	for i := 0; i < opts.NumBlocksToProcess; i++ {
		go client.GetBlockTransactions(lastBlock, i, &wg)
	}
	wg.Wait()
	addr, balanceOffset := client.GetMaxBalanceOffset()
	log.Printf("Address: %s\n", addr)
	log.Printf("Offset: %s", balanceOffset)
}