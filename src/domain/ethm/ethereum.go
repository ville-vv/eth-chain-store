package ethm

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

type EthereumDataPuller struct {
	ethRpcCli         ethrpc.EthRPC
	txWrite           TxWriter
	latestBlockNumber string
}

func NewEthereumManager(ethRpcCli ethrpc.EthRPC, txWrite TxWriter) *EthereumDataPuller {
	return &EthereumDataPuller{ethRpcCli: ethRpcCli, txWrite: txWrite}
}
func (sel *EthereumDataPuller) Pull(bkNumber int64, latestBkNum int64) error {
	// 获取块信息
	block, err := sel.ethRpcCli.GetBlockByNumber(bkNumber)
	if err != nil {
		vlog.ERROR("get block by number %d error %s", bkNumber, err.Error())
		return err
	}
	block.LatestBlockNumber = fmt.Sprintf("%d", latestBkNum)
	return sel.dealBlock(block)
}

// 处理块数据
func (sel *EthereumDataPuller) dealBlock(block *ethrpc.EthBlock) error {
	if block == nil {
		return nil
	}
	var err error

	//block.Number = common.HexToHash(block.Number).Big().String()
	// 把时间戳该为容易读懂的时间
	block.TimeStamp = block.TimeStampFormatTmString()
	// 一个块存n条交易数据
	for _, trfData := range block.Transactions {
		// 写入以太坊原生交易信息
		if err = sel.txWrites(&model.TransactionData{
			LatestNumber: block.LatestBlockNumber,
			TimeStamp:    block.TimeStamp,
			BlockHash:    trfData.BlockHash,
			BlockNumber:  trfData.BlockNumberToBig(),
			From:         trfData.From,
			GasPrice:     trfData.GasPriceToBig(),
			Hash:         trfData.Hash,
			To:           trfData.To,
			Value:        common.HexToHash(trfData.Value).Big().String(),
		}); err != nil {
			vlog.ERROR("处理以太坊原生交易错误：hash=%s %s", trfData.Hash, err.Error())
			return err
		}

		// 该条以太坊交易信息是否存在合约交易信息
		if trfData.IsContractToken() {
			// 这笔交易存在合约代币交易，需要获取合约里面的交易内容
			if err = sel.contractTransaction(&block.EthBlockHeader, trfData); err != nil {
				vlog.ERROR("处理以太坊合约交易错误：hash=%s %s", trfData.Hash, err.Error())
				return err
			}
		}

	}
	return err
}

// contractTransaction 合约代币交易
func (sel *EthereumDataPuller) contractTransaction(header *ethrpc.EthBlockHeader, tfData *ethrpc.EthTransaction) error {
	if tfData.IsTransfer() {
		// 这个是 ERC20单笔的 token 转账
		to, val := ethrpc.TransferParser(tfData.Input).TransferParse()
		return sel.txWrites(&model.TransactionData{
			LatestNumber:    header.LatestBlockNumber,
			ContractAddress: tfData.To, // 单笔合约交易一般都是 To 为合约地址
			TimeStamp:       header.TimeStamp,
			BlockHash:       tfData.BlockHash,
			BlockNumber:     tfData.BlockNumberToBig(),
			From:            tfData.From,
			Hash:            tfData.Hash,
			To:              to,
			Value:           val,
			IsContractToken: true,
			TxType:          model.TxTokenTransfer,
			EventType:       model.TxEventTypeTransfer,
		})
	}
	// 非单笔转账交易，就要获取凭证数据
	return sel.txReceipt(header.LatestBlockNumber, header.TimeStamp, tfData.Hash)
}

// txReceipt 多笔凭证交易
func (sel *EthereumDataPuller) txReceipt(latestNum, timeStamp string, hash string) error {
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
					LatestNumber:    latestNum,
					ContractAddress: lg.Address, // 代币交易是存在合约地址的
					TimeStamp:       timeStamp,
					BlockHash:       lg.BlockHash,
					BlockNumber:     lg.BlockNumberToBig(),
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

func (sel *EthereumDataPuller) txWrites(txData *model.TransactionData) error {
	return sel.txWrite.TxWrite(txData)
}
