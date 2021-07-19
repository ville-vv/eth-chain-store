package ethrpc

import "github.com/ville-vv/eth-chain-store/src/common/go-ethereum/common"

type Input string

// IsContract 判断是否合约交易
// return 如果是合约交易就返回 true
func (t Input) IsContract() bool {
	dataLen := len(t)
	if dataLen < 10 && t == "0x" {
		return false
	}
	return true
}

func (t Input) IsTransfer() bool {
	dataLen := len(t)
	if dataLen < 10 {
		return false
	}
	if t[:10] == ERC20MethodIDForTransfer {
		return true
	}
	return false
}

type TransferParser string

func (sel TransferParser) TransferParse() (to string, val string) {
	if len(sel) < 138 {
		return "", "0"
	}
	to = common.HexToAddress("0x" + string(sel[32:71])).String()
	val = common.HexToHash("0x" + string(sel[72:])).Big().String()
	return
}
