package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/test/identityset"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
)

const (
	//18160ddd -> totalSupply()
	totalSupplyString = "18160ddd"
	//70a08231 -> balanceOf(address)
	balanceOfString = "70a08231000000000000000000000000fea7d8ac16886585f1c232f13fefc3cfa26eb4cc"
	//dd62ed3e -> allowance(address,address)
	allowanceString = "dd62ed3e000000000000000000000000fea7d8ac16886585f1c232f13fefc3cfa26eb4cc000000000000000000000000fea7d8ac16886585f1c232f13fefc3cfa26eb4cc"
	//095ea7b3 -> approve(address,uint256)
	approveString = "095ea7b3000000000000000000000000fea7d8ac16886585f1c232f13fefc3cfa26eb4cc0000000000000000000000000000000000000000000000000000000000000001"

	// transferSha3 is sha3 of xrc20's transfer event,keccak('Transfer(address,address,uint256)')
	transferSha3 = "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	topicsPlusDataLen = 256
	sha3Len           = 64
	contractParamsLen = 64
	addressLen        = 40
	successStatus     = uint64(1)
	revertStatus      = uint64(106)
)

var (
	totalSupply, _   = hex.DecodeString(totalSupplyString)
	balanceOf, _     = hex.DecodeString(balanceOfString)
	allowance, _     = hex.DecodeString(allowanceString)
	approve, _       = hex.DecodeString(approveString)
	erc20Contract    = make(map[string]bool)
	nonErc20Contract = make(map[string]bool)
	nonce            = uint64(1)
	transferAmount   = big.NewInt(0)
	gasLimit         = uint64(100000)
	gasPrice         = big.NewInt(10000000)
	callerAddress    = identityset.Address(30).String()
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
	//endpoint := "api.testnet.iotex.one:80"
	endpoint := "api.mainnet.iotex.one:80"
	cli := ChainClient(endpoint)

	addr := "io18ndhcj88pwz5a5h68yhzz6r4q8vykwhugq45ns"
	callData := totalSupply
	execution, err := action.NewExecution(addr, nonce, transferAmount, gasLimit, gasPrice, callData)
	if err != nil {
		log.Fatal(err)
	}
	request := &iotexapi.ReadContractRequest{
		Execution:     execution.Proto(),
		CallerAddress: callerAddress,
	}

	res, err := cli.ReadContract(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	decoded, err := hex.DecodeString(res.GetData())
	if err != nil {
		log.Fatal(err)
	}
	balance := new(big.Int).SetBytes(decoded)
	//https://iotexscan.io/token/io18ndhcj88pwz5a5h68yhzz6r4q8vykwhugq45ns
	log.Printf("%v", balance)
}
