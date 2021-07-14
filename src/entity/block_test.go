package entity

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"testing"
)

func TestBlockCaptor_PullBlockByNumber(t *testing.T) {
	log.Init()
	ethClt := ethrpc.NewClient("https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119")
	count := 0
	txWrite := func(txData *model.TransactionData) error {
		count++
		fmt.Println(count)
		return nil
	}

	bc := &BlockCaptor{ethRpcCli: ethClt, txWrite: TxWriteFun(txWrite)}
	bc.PullBlockByNumber(12696216)
}
