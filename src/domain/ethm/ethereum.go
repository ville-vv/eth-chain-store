package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

type TxWriteFun func(*model.TransactionData) error

func (tw TxWriteFun) TxWrite(txData *model.TransactionData) error {
	return tw(txData)
}

type EthereumManager struct {
	ethRpcCli ethrpc.EthRPC
	txWrite   TxWriter
}

func NewEthereumManager(ethRpcCli ethrpc.EthRPC, txWrite TxWriter) *EthereumManager {
	return &EthereumManager{ethRpcCli: ethRpcCli, txWrite: txWrite}
}

func (sel *EthereumManager) PullBlock() error {
	// 获取块信息
	block, err := sel.ethRpcCli.GetBlock()
	if err != nil {
		vlog.ERROR("get latest block information")
		return err
	}
	return sel.dealBlock(block)
}

func (sel *EthereumManager) PullBlockByNumber(bkNumber int64) error {
	// 获取块信息
	block, err := sel.ethRpcCli.GetBlockByNumber(bkNumber)
	if err != nil {
		vlog.ERROR("get block by number %d error %s", bkNumber, err.Error())
		return err
	}
	//vlog.DEBUG("have get %d", bkNumber)
	return sel.dealBlock(block)
}

// 处理块数据
func (sel *EthereumManager) dealBlock(block *ethrpc.EthBlock) error {
	var err error
	for _, trfData := range block.Transactions {
		block.TimeStamp = common.HexToHash(block.TimeStamp).Big().String()
		if trfData.IsContractToken() {
			// 合约代币交易,需要获取合约里面的交易内容
			if err = sel.contractTransaction(&block.EthBlockHeader, trfData); err != nil {
				vlog.ERROR("处理以太坊合约交易错误：hash=%s %s", trfData.Hash, err.Error())
				return err
			}
		}
		// 无论合约交易还是非合约交易都需要记录以太坊交易的信息
		if err = sel.txWrites(&model.TransactionData{
			TimeStamp:   block.TimeStamp,
			BlockHash:   trfData.BlockHash,
			BlockNumber: trfData.BlockNumber,
			From:        trfData.From,
			GasPrice:    trfData.GasPrice,
			Hash:        trfData.Hash,
			To:          trfData.To,
			Value:       common.HexToHash(trfData.Value).Big().String(),
		}); err != nil {
			vlog.ERROR("处理以太坊原生交易错误：hash=%s %s", trfData.Hash, err.Error())
			return err
		}
	}
	return err
}

func (sel *EthereumManager) dealTransaction(header *ethrpc.EthBlockHeader, trfData *ethrpc.EthTransaction) (err error) {
	header.TimeStamp = common.HexToHash(header.TimeStamp).Big().String()
	if trfData.IsContractToken() {
		// 合约交易
		return sel.contractTransaction(header, trfData)
	}
	// 非合约交易
	if err = sel.txWrites(&model.TransactionData{
		TimeStamp:   header.TimeStamp,
		BlockHash:   trfData.BlockHash,
		BlockNumber: trfData.BlockNumber,
		From:        trfData.From,
		Hash:        trfData.Hash,
		To:          trfData.To,
		Value:       common.HexToHash(trfData.Value).Big().String(),
		TxType:      model.TxTypeNormal,
		EventType:   model.TxEventTypeTransfer,
	}); err != nil {
		vlog.ERROR("处理以太坊原生交易错误：hash=%s %s", trfData.Hash, err.Error())
	}
	return nil
}

// contractTransaction 合约代币交易
func (sel *EthereumManager) contractTransaction(header *ethrpc.EthBlockHeader, tfData *ethrpc.EthTransaction) error {
	if tfData.IsTransfer() {
		// 这个是 ERC20单笔的 token 转账
		to, val := ethrpc.TransferParser(tfData.Input).TransferParse()
		return sel.txWrites(&model.TransactionData{
			ContractAddress: tfData.To, // 单笔合约交易一般都是 To 为合约地址
			TimeStamp:       header.TimeStamp,
			BlockHash:       tfData.BlockHash,
			BlockNumber:     tfData.BlockNumber,
			From:            tfData.From,
			Hash:            tfData.Hash,
			To:              to,
			Value:           val,
			IsContractToken: true,
			TxType:          model.TxTokenTransfer,
			EventType:       model.TxEventTypeTransfer,
		})
	}
	return sel.txReceipt(header.TimeStamp, tfData.Hash)
}

func (sel *EthereumManager) txReceipt(timeStamp string, hash string) error {
	// 如果不是直接的转账交易，就获取合约的交易收据信息
	tfReceipt, err := sel.ethRpcCli.GetTransactionReceipt(hash)
	if err != nil {
		return err
	}
	if tfReceipt == nil {
		vlog.WARN("合约交易凭证查询为空：hash=%s", hash)
		return nil
	}
	if len(tfReceipt.Logs) > 0 {
		for _, lg := range tfReceipt.Logs {
			if lg.IsTransfer() {
				// 转账凭证处理
				if err = sel.txWrites(&model.TransactionData{
					ContractAddress: lg.Address, // 代币交易是存在合约地址的
					TimeStamp:       timeStamp,
					BlockHash:       lg.BlockHash,
					BlockNumber:     lg.BlockNumber,
					From:            lg.From(),
					Hash:            lg.TransactionHash,
					To:              lg.To(),
					Value:           lg.Value(),
					IsContractToken: true,
					TxType:          model.TxTokenTransfer,
					EventType:       model.TxEventTypeTransfer,
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (sel *EthereumManager) txWrites(txData *model.TransactionData) error {
	//codeData, err := sel.ethRpcCli.GetCode(txData.From)
	//if err != nil {
	//	return err
	//}
	//txData.FromIsContract = codeData != "0x" && codeData != ""
	//
	//codeData, err = sel.ethRpcCli.GetCode(txData.To)
	//if err != nil {
	//	return err
	//}
	//txData.ToIsContract = codeData != "0x" && codeData != ""

	//var err error
	//blockNumber := common.HexToHash(txData.BlockNumber).Big().Int64()
	//// 获取当前交易的余额
	//if txData.IsContractToken {
	//	txData.Balance, err = sel.ethRpcCli.GetContractBalanceByBlockNumber(txData.ContractAddressRecord, txData.From, blockNumber)
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

type EthereumWriter struct {
	filter            TxWriterFilter
	accountWriter     TxWriter
	contractWriter    TxWriter
	transactionWriter TxWriter
}

func NewEthereumWriter(filter TxWriterFilter, accountWriter TxWriter, contractWriter TxWriter, transactionWriter TxWriter) *EthereumWriter {
	return &EthereumWriter{filter: filter, accountWriter: accountWriter, contractWriter: contractWriter, transactionWriter: transactionWriter}
}

func (sel *EthereumWriter) TxWrite(txData *model.TransactionData) (err error) {
	if err = sel.filter.Filter(txData.From); err != nil {
		// 过滤发生错误不返回
		return nil
	}
	// 写入合约信息
	if err = sel.contractWriter.TxWrite(txData); err != nil {
		return
	}

	// 写入账户信息
	if err = sel.accountWriter.TxWrite(txData); err != nil {
		return
	}

	// 写入交易流水
	return sel.transactionWriter.TxWrite(txData)
}
