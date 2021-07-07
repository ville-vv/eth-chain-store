package ethrpc

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/go-ethereum"
	"github.com/ville-vv/eth-chain-store/src/common/go-ethereum/common"
	"github.com/ville-vv/eth-chain-store/src/common/go-ethereum/common/hexutil"
	"github.com/ville-vv/eth-chain-store/src/common/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
)

const (
	ERC20MethodIDForBalanceOf    = "0x70a08231" // balanceOf(address)
	ERC20MethodIDForTransfer     = "0xa9059cbb" // transfer(address,uint256)
	ERC20MethodIDForDecimals     = "0x313ce567" // decimals()
	ERC20MethodIDForAllowance    = "0xdd62ed3e" // allowance(address,address)
	ERC20MethodIDForSymbol       = "0x95d89b41" // symbol()
	ERC20MethodIDForTotalSupply  = "0x18160ddd" // totalSupply()
	ERC20MethodIDForName         = "0x06fdde03" // name()
	ERC20MethodIDForApprove      = "0x095ea7b3" // approve(address,uint256)
	ERC20MethodIDForTransferFrom = "0x23b872dd" // transferFrom(address,address,uint256)
)

// GenMethodId 可以使这个方法生成 method id， 比如：GenMethodId("balanceOf(address)") 得到 0x70a08231
func GenMethodId(method string) string {
	kek := sha3.NewLegacyKeccak256()
	kek.Reset()
	kek.Write([]byte("balanceOf(address)"))
	return hexutil.Encode(kek.Sum(nil))[:10]
}

type ContractCallParam struct {
	MethodID string
	Params   []string
}

func (sel ContractCallParam) HexByte() []byte {
	return common.FromHex(sel.String())
}

func (sel ContractCallParam) String() string {
	param := sel.MethodID + "000000000000000000000000"
	for _, val := range sel.Params {
		tmp := strings.Replace(val, "0x", "", 1)
		param = param + tmp
	}
	return param
}

// Client
type Client struct {
	ethCli *ethclient.Client
}

func NewClient(endpoint string) *Client {
	ethCli, err := ethclient.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	return &Client{ethCli: ethCli}
}

// GetBalance 获取ETH最新余额
func (sel *Client) GetBalance(addr string) (string, error) {
	balance, err := sel.ethCli.BalanceAt(context.Background(), common.HexToAddress(addr), nil)
	if err != nil {
		return "0", err
	}
	return balance.String(), nil
}

// GetBalanceByBlockNumber 获取ETH指定区块余额
func (sel *Client) GetBalanceByBlockNumber(addr string, blockNumber int64) (string, error) {
	balance, err := sel.ethCli.BalanceAt(context.Background(), common.HexToAddress(addr), big.NewInt(blockNumber))
	if err != nil {
		return "0", err
	}
	return balance.String(), nil
}

// GetContractBalance 获取ERC20合约代币最新余额
func (sel *Client) GetContractBalance(contract, addr string) (string, error) {
	toAddr := common.HexToAddress(contract)
	msg := ethereum.CallMsg{
		From: common.HexToAddress(contract),
		To:   &toAddr,
		Data: (ContractCallParam{
			MethodID: ERC20MethodIDForBalanceOf,
			Params:   []string{addr},
		}).HexByte(),
	}
	respData, err := sel.ethCli.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "0x0", err
	}
	return hexutil.Encode(respData), nil
}

// GetContractBalanceByBlockNumber 获取ERC20合约代币指定区块时余额
func (sel *Client) GetContractBalanceByBlockNumber(contract, addr string, blockNumber int64) (string, error) {
	toAddr := common.HexToAddress(contract)
	msg := ethereum.CallMsg{
		From: common.HexToAddress(addr),
		To:   &toAddr,
		Data: (ContractCallParam{
			MethodID: ERC20MethodIDForBalanceOf,
			Params:   []string{addr},
		}).HexByte(),
	}
	respData, err := sel.ethCli.CallContract(context.Background(), msg, big.NewInt(blockNumber))
	if err != nil {
		return "0x0", err
	}
	return hexutil.Encode(respData), nil
}

// GetContractTotalSupply 获取ERC 20 代币 发行总量
func (sel *Client) GetContractTotalSupply(contract string) (string, error) {
	toAddr := common.HexToAddress(contract)
	msg := ethereum.CallMsg{
		From: common.HexToAddress(contract),
		To:   &toAddr,
		Data: (ContractCallParam{
			MethodID: ERC20MethodIDForTotalSupply,
			Params:   []string{contract},
		}).HexByte(),
	}
	respData, err := sel.ethCli.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "0x0", err
	}
	return hexutil.Encode(respData), nil
}

// 获取合约地址的编译后的代码，如果是非合约地址，返回 0x
func (sel *Client) GetCode(addr string) (string, error) {
	respData, err := sel.ethCli.CodeAt(context.Background(), common.HexToAddress(addr), nil)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(respData), nil
}

// GetTransactionByBlockNumber
func (sel *Client) GetBlockByNumber(blockNumber int64) error {
	blockInfo, err := sel.ethCli.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return err
	}
	blockInfo.Body()
	return nil
}
