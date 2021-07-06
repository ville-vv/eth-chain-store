package ethrpc

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
	"strconv"
	"testing"
)

func TestClient_GetBalanceByBlockNumber(t *testing.T) {
	cli := NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
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
