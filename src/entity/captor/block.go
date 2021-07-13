package captor

import (
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vlog"
)

type BlockCaptor struct {
	ethRpcCli ethrpc.EthRpcClient
}

func (sel *BlockCaptor) PullBlock() {
	// 获取块信息
	block, err := sel.ethRpcCli.GetBlock()
	if err != nil {
		vlog.ERROR("")
		return
	}
	sel.dealBlock(block)
}

func (sel *BlockCaptor) PullBlockByNumber(bkNumber int64) {
	// 获取块信息
	block, err := sel.ethRpcCli.GetBlockByNumber(bkNumber)
	if err != nil {
		vlog.ERROR("")
		return
	}
	sel.dealBlock(block)
}

// 处理块数据
func (sel *BlockCaptor) dealBlock(block *ethrpc.RpcBlock) {
	for _, trfData := range block.Transactions {
		if trfData.IsContract() {
			// 合约交易
			return
		}
		// 非合约交易
	}
}

// contractTransaction 合约交易
func (sel *BlockCaptor) contractTransaction(header *ethrpc.RpcBlockHeader, tfData *ethrpc.RpcTransaction) error {
	if !tfData.IsTransfer() {
		// 如果不是直接的转账交易，就获取合约的交易收据信息
		tfReceipt, err := sel.ethRpcCli.GetTransactionReceipt(string(tfData.Input))
		if err != nil {
			return err
		}
		for _, lg := range tfReceipt.Logs {
			if !lg.IsTransfer() {
				// 非转账凭证不处理
				continue
			}
		}
	}
	//to, val := ethrpc.TransferParser(tfData.Input).TransferParse()
	return nil
}
