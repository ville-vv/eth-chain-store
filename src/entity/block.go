package entity

import (
	"github.com/ville-vv/eth-chain-store/src/common/go-ethereum/common"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

type TxWriteFun func(*model.TransactionData) error

func (tw TxWriteFun) TxWrite(txData *model.TransactionData) error {
	return tw(txData)
}

type BlockCaptor struct {
	ethRpcCli ethrpc.EthRpcClient
	txWrite   TxWriter
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
		vlog.ERROR("get block by number %s error %s", bkNumber, err.Error())
		return
	}
	sel.dealBlock(block)
}

// 处理块数据
func (sel *BlockCaptor) dealBlock(block *ethrpc.RpcBlock) {
	var err error
	for _, trfData := range block.Transactions {
		block.TimeStamp = common.HexToHash(block.TimeStamp).Big().String()
		if trfData.IsContract() {
			// 合约交易
			if err = sel.contractTransaction(&block.RpcBlockHeader, trfData); err != nil {
				vlog.ERROR("处理以太坊合约交易错误：hash=%s %s", trfData.Hash, err.Error())
			}
			continue
		}
		// 非合约交易
		if err = sel.TxWrite(&model.TransactionData{
			TimeStamp:   block.TimeStamp,
			BlockHash:   trfData.BlockHash,
			BlockNumber: trfData.BlockNumber,
			From:        trfData.From,
			Hash:        trfData.Hash,
			To:          trfData.To,
			Value:       common.HexToHash(trfData.Value).Big().String(),
		}); err != nil {
			vlog.ERROR("处理以太坊原生交易错误：hash=%s %s", trfData.Hash, err.Error())
		}
	}
}

// contractTransaction 合约交易
func (sel *BlockCaptor) contractTransaction(header *ethrpc.RpcBlockHeader, tfData *ethrpc.RpcTransaction) error {
	if !tfData.IsTransfer() {
		// 如果不是直接的转账交易，就获取合约的交易收据信息
		tfReceipt, err := sel.ethRpcCli.GetTransactionReceipt(tfData.Hash)
		if err != nil {
			return err
		}
		if tfReceipt == nil {
			vlog.WARN("合约交易凭证查询为空：hash=%s", tfData.Hash)
			return nil
		}
		for _, lg := range tfReceipt.Logs {
			if !lg.IsTransfer() {
				// 非转账凭证不处理
				continue
			}
			if err = sel.TxWrite(&model.TransactionData{
				ContractAddress: lg.Address,
				TimeStamp:       header.TimeStamp,
				BlockHash:       lg.BlockHash,
				BlockNumber:     lg.BlockNumber,
				From:            lg.From(),
				Hash:            lg.TransactionHash,
				To:              lg.To(),
				Value:           lg.Value(),
				IsContract:      true,
			}); err != nil {
				return err
			}
		}
	}
	to, val := ethrpc.TransferParser(tfData.Input).TransferParse()
	return sel.TxWrite(&model.TransactionData{
		ContractAddress: tfData.To,
		TimeStamp:       common.HexToHash(header.TimeStamp).Big().String(),
		BlockHash:       tfData.BlockHash,
		BlockNumber:     tfData.BlockNumber,
		From:            tfData.From,
		Hash:            tfData.Hash,
		To:              to,
		Value:           val,
		IsContract:      true,
	})
}

func (sel *BlockCaptor) TxWrite(txData *model.TransactionData) error {
	//var err error
	//blockNumber := common.HexToHash(txData.BlockNumber).Big().Int64()
	//// 获取当前交易的余额
	//if txData.IsContract {
	//	txData.Balance, err = sel.ethRpcCli.GetContractBalanceByBlockNumber(txData.ContractAddress, txData.From, blockNumber)
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	txData.Balance, err = sel.ethRpcCli.GetBalanceByBlockNumber(txData.From, blockNumber)
	//	if err != nil {
	//		return err
	//	}
	//}
	return sel.txWrite.TxWrite(txData)
}
