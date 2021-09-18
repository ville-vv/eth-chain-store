package entity

import (
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"strconv"
)

type Contract struct {
	rpcCli ethrpc.EthRPC
	model.ContractContent
	contractRepo repo.ContractRepository
}

func NewContract(ethCli ethrpc.EthRPC, contractRepo repo.ContractRepository) *Contract {
	return &Contract{
		rpcCli:          ethCli,
		contractRepo:    contractRepo,
		ContractContent: model.ContractContent{},
	}
}

func (sel *Contract) SetPublishTime(tm string) {
	sel.PublishTime = tm
}

func (sel *Contract) SetAddress(addr string) {
	sel.Address = addr
}

// GetContentInGrpc 从 rpc 接口中获取合约信息
func (sel *Contract) SetErc20ContentFromRpc() error {
	addr := sel.Address
	supply, err := sel.rpcCli.GetContractTotalSupply(addr)
	if err != nil {
		// 没有总发行量就直接退出
		return err
	}
	// 存在总发行量就一定是（Peter 说的）erc20

	supply = common.HexToHash(supply).Big().String()
	// 获取合约名字,暂时不要
	//name, _ := sel.rpcCli.GetContractName(addr)

	symbol, _ := sel.rpcCli.GetContractSymbol(addr)

	decimal, err := sel.rpcCli.GetContractDecimals(addr)
	if err != nil {
		decimal = "0x0"
	}
	decimal = common.HexToHash(decimal).Big().String()
	var decimalInt int64
	decimalInt, _ = strconv.ParseInt(decimal, 10, 64)

	sel.Symbol = symbol
	sel.Address = addr
	sel.IsErc20 = true
	sel.DecimalBit = int(decimalInt)

	return nil
}

func (sel *Contract) CreateRecord() error {
	return nil
}
