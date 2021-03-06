package ethrpc

import (
	"github.com/ville-vv/eth-chain-store/src/common/go-eth"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common/hexutil"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/rpc"
	"math/big"
	"strings"
)

const (
	ERC20MethodIDForBalanceOf    = "0x70a08231"                                                         // balanceOf(address)
	ERC20MethodIDForTransfer     = "0xa9059cbb"                                                         // transfer(address,uint256)
	ERC20MethodIDForDecimals     = "0x313ce567"                                                         // decimals()
	ERC20MethodIDForAllowance    = "0xdd62ed3e"                                                         // allowance(address,address)
	ERC20MethodIDForSymbol       = "0x95d89b41"                                                         // symbol()
	ERC20MethodIDForTotalSupply  = "0x18160ddd"                                                         // totalSupply()
	ERC20MethodIDForName         = "0x06fdde03"                                                         // name()
	ERC20MethodIDForApprove      = "0x095ea7b3"                                                         // approve(address,uint256)
	ERC20MethodIDForTransferFrom = "0x23b872dd"                                                         // transferFrom(address,address,uint256)
	ERC20EventIDForTransfer      = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" // Transfer(address,address,uint256)
)

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

var _ EthRPC = &Client{}

// Client
type Client struct {
	ethCli *rpc.Client
}

func NewClient(endpoint string) *Client {
	ethCli, err := rpc.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	return &Client{ethCli: ethCli}
}
func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}
func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

// GetBalance ??????ETH????????????
func (sel *Client) GetBalance(addr string) (string, error) {
	var result hexutil.Big
	err := sel.ethCli.Call(&result, "eth_getBalance", common.HexToAddress(addr), "latest")
	if err != nil {
		return "0", err
	}
	return result.ToInt().String(), nil
}

// GetBalanceByBlockNumber ??????ETH??????????????????
func (sel *Client) GetBalanceByBlockNumber(addr string, blockNumber int64) (string, error) {
	var result hexutil.Big
	bigInt := hexutil.Big(*big.NewInt(blockNumber))
	err := sel.ethCli.Call(&result, "eth_getBalance", common.HexToAddress(addr), bigInt.String())
	if err != nil {
		return "0", err
	}
	return result.ToInt().String(), nil
}

// GetContractBalance ??????ERC20????????????????????????
func (sel *Client) GetContractBalance(contract, addr string) (string, error) {
	toAddr := common.HexToAddress(contract)
	//msg := ethereum.CallMsg{
	//	From: common.HexToAddress(contract),
	//	To:   &toAddr,
	//	Data: (ContractCallParam{
	//		MethodID: ERC20MethodIDForBalanceOf,
	//		Params:   []string{addr},
	//	}).HexByte(),
	//}
	arg := map[string]interface{}{
		//"from": common.HexToAddress(contract),
		"to": &toAddr,
		"data": (ContractCallParam{
			MethodID: ERC20MethodIDForBalanceOf,
			Params:   []string{addr},
		}).String(),
	}
	var hex hexutil.Bytes
	err := sel.ethCli.Call(&hex, "eth_call", arg, "latest")
	if err != nil {
		return "0", err
	}
	return common.HexToHash(hex.String()).Big().String(), nil
}

// GetContractBalanceByBlockNumber ??????ERC20?????????????????????????????????
func (sel *Client) GetContractBalanceByBlockNumber(contract, addr string, blockNumber int64) (string, error) {
	toAddr := common.HexToAddress(contract)
	//msg := ethereum.CallMsg{
	//	From: common.HexToAddress(addr),
	//	To:   &toAddr,
	//	Data: (ContractCallParam{
	//		MethodID: ERC20MethodIDForBalanceOf,
	//		Params:   []string{addr},
	//	}).HexByte(),
	//}
	arg := map[string]interface{}{
		//"from": common.HexToAddress(contract),
		"to": &toAddr,
		"data": (ContractCallParam{
			MethodID: ERC20MethodIDForBalanceOf,
			Params:   []string{addr},
		}).String(),
	}

	var hex hexutil.Bytes

	bigInt := hexutil.Big(*big.NewInt(blockNumber))

	err := sel.ethCli.Call(&hex, "eth_call", arg, bigInt.String())
	if err != nil {
		return "0", err
	}
	return common.HexToHash(hex.String()).Big().String(), nil
}

// GetContractTotalSupply ??????ERC 20 ?????? ????????????
func (sel *Client) GetContractTotalSupply(contract string) (string, error) {
	toAddr := common.HexToAddress(contract)
	//msg := ethereum.CallMsg{
	//	From: common.HexToAddress(contract),
	//	To:   &toAddr,
	//	Data: (ContractCallParam{
	//		MethodID: ERC20MethodIDForTotalSupply,
	//		Params:   []string{contract},
	//	}).HexByte(),
	//}

	arg := map[string]interface{}{
		//"from": common.HexToAddress(contract),
		"to": &toAddr,
		"data": (ContractCallParam{
			MethodID: ERC20MethodIDForTotalSupply,
			Params:   []string{contract},
		}).String(),
	}
	var hex string
	err := sel.ethCli.Call(&hex, "eth_call", arg, "latest")
	if err != nil {
		return "0x0", err
	}
	return hex, nil
}

// ??????????????????????????????????????????????????????????????????????????? 0x
func (sel *Client) GetCode(addr string) (string, error) {
	var result hexutil.Bytes
	err := sel.ethCli.Call(&result, "eth_getCode", common.HexToAddress(addr), "latest")
	if err != nil {
		return "", err
	}
	return result.String(), err
}

// GetBlockByNumber ???????????????????????????
func (sel *Client) GetBlockByNumber(blockNumber int64) (*EthBlock, error) {
	var block *EthBlock
	err := sel.ethCli.Call(&block, "eth_getBlockByNumber", toBlockNumArg(big.NewInt(blockNumber)), true)
	return block, err
}

// GetLatestBlock ???????????????????????????
func (sel *Client) GetBlock() (*EthBlock, error) {
	var block *EthBlock
	err := sel.ethCli.Call(&block, "eth_getBlockByNumber", "latest", true)
	return block, err
}

func (sel *Client) GetBlockNumber() (uint64, error) {
	var result hexutil.Uint64
	err := sel.ethCli.Call(&result, "eth_blockNumber")
	return uint64(result), err
}

// GetTransactionReceipt ????????????????????? ?????????????????????????????????????????????????????????
func (sel *Client) GetTransactionReceipt(hash string) (*EthTransactionReceipt, error) {
	var result *EthTransactionReceipt
	err := sel.ethCli.Call(&result, "eth_getTransactionReceipt", common.HexToHash(hash))
	return result, err
}

func (sel *Client) GetContractSymbol(contract string) (string, error) {
	toAddr := common.HexToAddress(contract)
	arg := map[string]interface{}{
		"to": &toAddr,
		"data": (ContractCallParam{
			MethodID: ERC20MethodIDForSymbol,
			Params:   []string{contract},
		}).String(),
	}
	var hex string
	err := sel.ethCli.Call(&hex, "eth_call", arg, "latest")
	if err != nil {
		if strings.Contains(err.Error(), "invalid opcode") {
			return "no symbol", nil
		}
		return "0x0", err
	}
	return parseErc20StringProperty(hex), nil
}

func (sel *Client) GetContractName(contract string) (string, error) {
	toAddr := common.HexToAddress(contract)
	arg := map[string]interface{}{
		"to": &toAddr,
		"data": (ContractCallParam{
			MethodID: ERC20MethodIDForName,
			Params:   []string{contract},
		}).String(),
	}
	var hex string
	err := sel.ethCli.Call(&hex, "eth_call", arg, "latest")
	if err != nil {
		return "0x0", err
	}
	return parseErc20StringProperty(hex), nil
}

func (sel *Client) GetContractDecimals(contract string) (string, error) {
	toAddr := common.HexToAddress(contract)
	arg := map[string]interface{}{
		"to": &toAddr,
		"data": (ContractCallParam{
			MethodID: ERC20MethodIDForDecimals,
			Params:   []string{contract},
		}).String(),
	}
	var hex string
	err := sel.ethCli.Call(&hex, "eth_call", arg, "latest")
	if err != nil {
		return "0x0", err
	}
	return hex, nil
}
