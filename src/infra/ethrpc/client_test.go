package ethrpc

import (
	"bytes"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common/hexutil"
	"golang.org/x/crypto/sha3"
	"strconv"
	"testing"
)

func TestClient_GetBalanceByBlockNumber(t *testing.T) {
	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	cli = NewClient("http://172.16.16.119:8545")
	balance, err := cli.GetContractBalance("0xdAC17F958D2ee523a2206206994597C13D831ec7",
		"0xdAC17F958D2ee523a2206206994597C13D831ec7")
	fmt.Println(balance, err)

	kek := sha3.NewLegacyKeccak256()
	kek.Reset()
	kek.Write([]byte("balanceOf(address)"))

	fmt.Println(hexutil.Bytes(kek.Sum(nil)).String()[:10])
}

func TestClient_GetContractTotalSupply(t *testing.T) {
	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	totalSup, err := cli.GetContractTotalSupply("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	fmt.Println(totalSup, err)
	fmt.Println(strconv.ParseInt(totalSup, 0, 64))
}

func TestClient_GetCode(t *testing.T) {
	//	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	cli := NewClient("http://172.16.16.119:8545")
	totalSup, err := cli.GetCode("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	fmt.Println(totalSup, err)

	fmt.Println(bytes.Index(hexutil.Bytes(totalSup), common.FromHex(ERC20MethodIDForBalanceOf)))

}

func TestClient_GetBlockByNumber(t *testing.T) {
	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	block, err := cli.GetBlock()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("head", block.EthBlockHeader)
	for _, val := range block.Transactions {
		fmt.Println(val, err)
	}
}

func TestClient_GetBlockNumber(t *testing.T) {
	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	block, err := cli.GetBlockNumber()
	fmt.Println(block, err)
	//fmt.Println(Keccak256Hash("Transfer(address,address,uint256)"))
}

func TestClient_GetTransactionReceipt(t *testing.T) {
	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	receipt, err := cli.GetTransactionReceipt("0xe6600a046d6ba96d475aa7bf9ee98b3218a713aaf89e1d968651dfe1599280f7")
	if err != nil {
		t.Error(err)
		return
	}
	for _, lg := range receipt.Logs {
		if lg.IsTransfer() {
			fmt.Println(lg.Value(), lg.From(), lg.To())
		}
	}
}

//
func TestClient_GetSymbol(t *testing.T) {
	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	symbol, err := cli.GetContractSymbol("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	fmt.Println(symbol, err)
}

func TestClient_GetContractDecimals(t *testing.T) {
	cli := NewClient("http://172.16.16.111:8545")
	symbol, err := cli.GetContractDecimals("0x9582C4ADACB3BCE56Fea3e590F05c3ca2fb9C477")
	fmt.Println(symbol, err)
}
