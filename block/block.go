package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/blockchain/block"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func ChainClient(endpoint string) iotexapi.APIServiceClient {
	opt := grpc.WithInsecure()
	conn, err := grpc.Dial(endpoint, opt)
	if err != nil {
		panic(err)
	}
	return iotexapi.NewAPIServiceClient(conn)
}
func main() {
	config.SetEVMNetworkID(4690)
	endpoint := "api.testnet.iotex.one:80"
	//endpoint := "api.mainnet.iotex.one:80"
	chainClient := ChainClient(endpoint)
	count := uint64(1)
	startHeight := uint64(8715290)
	//startHeight := uint64(10618157)
	rawRequest := &iotexapi.GetRawBlocksRequest{
		StartHeight:  startHeight,
		Count:        count,
		WithReceipts: true,
	}
	getRawBlocksRes, err := chainClient.GetRawBlocks(context.Background(), rawRequest)
	if err != nil {
		panic(err)
	}

	for _, blkInfo := range getRawBlocksRes.GetBlocks() {
		blk := &block.Block{}
		if err := blk.ConvertFromBlockPb(blkInfo.GetBlock()); err != nil {
			panic(err)
		}
		receipts := map[hash.Hash256]*action.Receipt{}
		for _, receiptPb := range blkInfo.GetReceipts() {
			log.Printf("receiptPb: %s \n actHash: %x\n", receiptPb, string(receiptPb.GetActHash()))
			receipt := &action.Receipt{}
			receipt.ConvertFromReceiptPb(receiptPb)
			receipts[receipt.ActionHash] = receipt
			blk.Receipts = append(blk.Receipts, receipt)
		}
		transactionLogs, err := chainClient.GetTransactionLogByBlockHeight(
			context.Background(),
			&iotexapi.GetTransactionLogByBlockHeightRequest{
				BlockHeight: blk.Header.Height(),
			},
		)
		if err != nil {
			panic(err)
		}
		for _, tlogs := range transactionLogs.TransactionLogs.Logs {
			log.Printf("tlogs %s\n", tlogs)
			logs := make([]*action.TransactionLog, len(tlogs.Transactions))
			for i, txn := range tlogs.Transactions {
				amount, ok := new(big.Int).SetString(txn.Amount, 10)
				if !ok {
					panic(errors.Errorf("failed to parse %s", txn.Amount))
				}
				logs[i] = &action.TransactionLog{
					Type:      txn.Type,
					Amount:    amount,
					Sender:    txn.Sender,
					Recipient: txn.Recipient,
				}
			}
			actHash := hash.BytesToHash256(tlogs.ActionHash)
			receipts[actHash].AddTransactionLogs(logs...)
		}

		//debugBlock(blk)
		debugTransaction(blk)
		//msgPack(blk)
	}

}
func msgPack(blk *block.Block) {
	var err error
	data, err := blk.Serialize()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%x", data)
	unblk := &block.Block{}
	err = unblk.Deserialize(data)
	if err != nil {
		log.Panicln(err)
	}
	debugBlock(unblk)
}

func calculateTxRoot(acts []action.SealedEnvelope) hash.Hash256 {
	h := make([]hash.Hash256, 0, len(acts))
	for _, act := range acts {
		h = append(h, act.Hash())
	}
	if len(h) == 0 {
		return hash.ZeroHash256
	}
	return crypto.NewMerkleTree(h).HashTree()
}

type income struct {
	inFlow        *big.Int
	outFlow       *big.Int
	inNumActions  int
	outNumActions int
}

func debugTransaction(blk *block.Block) {
	incomes := make(map[string]*income)
	for _, receipt := range blk.Receipts {

		//transaction
		for _, transation := range receipt.TransactionLogs() {
			if transation.Sender != "" {
				if _, ok := incomes[transation.Sender]; !ok {
					incomes[transation.Sender] = &income{
						outFlow:       big.NewInt(0).Set(transation.Amount),
						outNumActions: 1,
						inFlow:        big.NewInt(0),
						inNumActions:  0,
					}
				} else {
					incomes[transation.Sender].outFlow = incomes[transation.Sender].outFlow.Add(incomes[transation.Sender].outFlow, transation.Amount)
					incomes[transation.Sender].outNumActions += 1
				}
				fmt.Println(transation.Sender, incomes[transation.Sender].outFlow.String())
			}
			if transation.Recipient != "" {
				if _, ok := incomes[transation.Recipient]; !ok {
					incomes[transation.Recipient] = &income{
						inFlow:        big.NewInt(0).Set(transation.Amount),
						inNumActions:  1,
						outFlow:       big.NewInt(0),
						outNumActions: 0,
					}
				} else {
					incomes[transation.Recipient].inFlow = incomes[transation.Recipient].inFlow.Add(incomes[transation.Recipient].inFlow, transation.Amount)
					incomes[transation.Recipient].inNumActions += 1
				}
				fmt.Println(transation.Recipient, incomes[transation.Recipient].inFlow.String())
			}

		}
	}
	fmt.Println("==")
	for accountAddress, income := range incomes {

		fmt.Println(accountAddress, income.inFlow.String(), income.inNumActions, income.outFlow.String(), income.outNumActions)
		// log.L().Info("insert",
		// 	zap.Duration("timeSpent", time.Since(timeStart)),
		// )

	}
}
func debugBlock(blk *block.Block) {
	fmt.Printf("===== header =====\n")
	txtRoot := blk.Header.TxRoot()
	fmt.Printf("txRoot : %s\n", hex.EncodeToString(txtRoot[:]))
	calTxtRoot := blk.CalculateTxRoot()
	fmt.Printf("CalculateTxRoot : %s\n", hex.EncodeToString(calTxtRoot[:]))

	testCalTxtRoot := calculateTxRoot(blk.Actions)
	fmt.Printf("TestCalculateTxRoot : %s\n", hex.EncodeToString(testCalTxtRoot[:]))
	for i, selp := range blk.Actions {
		fmt.Printf("===== action: #%d =====\n", i)
		actionHash := selp.Hash()
		fmt.Printf("actionHash : %s\n", hex.EncodeToString(actionHash[:]))
		sender, _ := address.FromBytes(selp.SrcPubkey().Hash())
		fmt.Printf("from : %s\n", sender)

		dst, _ := selp.Destination()
		fmt.Printf("to : %s\n", dst)
		gasPrice := selp.GasPrice().String()
		fmt.Printf("gasPrice : %s\n", gasPrice)
		gasLimit := selp.GasLimit()
		fmt.Printf("gasLimit : %d\n", gasLimit)
		fmt.Printf("gasPrice : %s\n", gasPrice)
		nonce := selp.Nonce()
		fmt.Printf("nonce : %d\n", nonce)

		act := selp.Action()
		actionType := fmt.Sprintf("%T", act)
		fmt.Printf("actionType : %s\n", actionType)
		amount := "0"

		switch a := act.(type) {
		case *action.Transfer:
			amount = a.Amount().String()
		case *action.Execution:
			amount = a.Amount().String()
		case *action.DepositToRewardingFund:
			amount = a.Amount().String()
		case *action.ClaimFromRewardingFund:
			amount = a.Amount().String()
		case *action.CreateStake:
			amount = a.Amount().String()
		case *action.DepositToStake:
			amount = a.Amount().String()
		case *action.CandidateRegister:
			amount = a.Amount().String()
		}
		fmt.Printf("amount : %s\n", amount)
	}
	for j, receipt := range blk.Receipts {
		fmt.Printf("===== receipt: #%d =====\n", j)
		fmt.Printf("receipt.ActionHash : %s\n", hex.EncodeToString(receipt.ActionHash[:]))
		receiptHash := receipt.Hash()
		fmt.Printf("receipt.Hash : %s\n", hex.EncodeToString(receiptHash[:]))
		fmt.Printf("receipt.Status : %d\n", receipt.Status)
		fmt.Printf("receipt.GasConsumed : %d\n", receipt.GasConsumed)
		fmt.Printf("receipt.BlockHeight : %d\n", receipt.BlockHeight)
		fmt.Printf("receipt.ContractAddress : %s\n", receipt.ContractAddress)
		fmt.Printf("receipt.executionRevertMsg : %s\n", receipt.ExecutionRevertMsg())
		for j1, transLog := range receipt.TransactionLogs() {
			fmt.Printf("===== receipt: #%d transaction: #%d =====\n", j, j1)
			fmt.Printf("transaction.Type : %s\n", transLog.Type)
			fmt.Printf("transaction.Amount : %s\n", transLog.Amount.String())
			fmt.Printf("transaction.Sender : %s\n", transLog.Sender)
			fmt.Printf("transaction.Recipient : %s\n", transLog.Recipient)
		}
		for j2, log := range receipt.Logs() {
			fmt.Printf("===== receipt: #%d log: #%d =====\n", j, j2)
			if len(log.Topics) > 0 {
				bucketIndex := new(big.Int).SetBytes(log.Topics[1][:])

				fmt.Printf("bucketID : %v \n", bucketIndex.String())
			}
			fmt.Printf("log.Address : %s\n", log.Address)
			fmt.Printf("log.Topics : %x\n", log.Topics)
			fmt.Printf("log.Data : %x\n", log.Data)
			fmt.Printf("log.BlockHeight : %d\n", log.BlockHeight)
			fmt.Printf("log.ActionHash : %s\n", hex.EncodeToString(log.ActionHash[:]))
			fmt.Printf("log.Index : %d\n", log.Index)
			fmt.Printf("log.NotFixTopicCopyBug : %v\n", log.NotFixTopicCopyBug)

		}

	}
}
