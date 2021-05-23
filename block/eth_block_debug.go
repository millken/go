package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	log.SetFlags(log.Llongfile)
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/ca34940b25684b49b1511bf73340e87c")
	if err != nil {
		log.Fatal(err)
	}
	blockNumber := big.NewInt(12478368)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	a, err := json.Marshal(block)
	fmt.Println(string(a))
	return
	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time())                // 1527211625
	fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println(len(block.Transactions()))   // 144

	// for _, tx := range block.Transactions() {
	// 	fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	// 	fmt.Println(tx.Value().String())    // 10000000000000000
	// 	fmt.Println(tx.Gas())               // 105000
	// 	fmt.Println(tx.GasPrice().Uint64()) // 102000000000
	// 	fmt.Println(tx.Nonce())             // 110644
	// 	fmt.Println(tx.Data())              // []
	// 	fmt.Println(tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
	// }
}
