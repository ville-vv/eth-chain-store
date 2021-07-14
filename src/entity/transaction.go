package entity

import (
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/eth-chain-store/src/repo"
)

type TxFlowProcessor interface {
	TxWriter
	UpdateAccount(*model.AccountContent) error
}

type TransactionWriter struct {
	contractFlow TxFlowProcessor
	nativeFlow   TxFlowProcessor
}

func (sel *TransactionWriter) getTxFlowProcessor(txData *model.TransactionData) TxFlowProcessor {
	if txData.IsContract {
		return sel.contractFlow
	}
	return sel.nativeFlow
}

func (sel *TransactionWriter) TxWrite(txData *model.TransactionData) error {
	prc := sel.getTxFlowProcessor(txData)
	if err := prc.UpdateAccount(&model.AccountContent{
		Contract: txData.ContractAddress,
		Address:  txData.From,
		TxTime:   txData.TimeStamp,
		TxHash:   txData.Hash,
	}); err != nil {
		return err
	}

	if err := prc.UpdateAccount(&model.AccountContent{
		Contract: txData.ContractAddress,
		Address:  txData.To,
		TxTime:   txData.TimeStamp,
		TxHash:   txData.Hash,
	}); err != nil {
		return err
	}
	return prc.TxWrite(txData)
}

type NativeFlowProcessor struct {
	ethCli  ethrpc.EthRpcClient
	actRepo repo.AccountRepo
}

func (sel *NativeFlowProcessor) TxWrite(txData *model.TransactionData) error {
	// 查询该地址是否存在，写入地址
	return nil
}

func (sel *NativeFlowProcessor) UpdateAccount(act *model.AccountContent) error {
	// 获取余额
	balance, err := sel.ethCli.GetBalance(act.Address)
	if err != nil {
		return err
	}
	act.Balance = balance
	//
	return nil
}

type ContractFlowProcessor struct {
	ethCli  ethrpc.EthRpcClient
	actRepo repo.AccountRepo
}

func (sel *ContractFlowProcessor) TxWrite(txData *model.TransactionData) error {

	return nil
}

func (sel *ContractFlowProcessor) UpdateAccount(act *model.AccountContent) error {
	// 获取转账者的余额
	balance, err := sel.ethCli.GetContractBalance(act.Contract, act.Address)
	if err != nil {
		return err
	}
	act.Balance = balance
	// 在数据库中是否存在该合约账户
	// 如果存在就直接更新余额
	// 如果不存在就需要获取该合约相关的东西
	// 写入合约信息
	return nil
}
