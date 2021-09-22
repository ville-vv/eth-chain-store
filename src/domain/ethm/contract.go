package ethm

import (
	"github.com/pkg/errors"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"strconv"
	"sync"
	"time"
)

type Erc20Manager interface {
	IsErc20(addr string) bool
}

type RingStrList struct {
	sync.RWMutex
	list   []string
	index  int
	length int
}

func NewRingStrList() *RingStrList {
	lng := 500000
	return &RingStrList{
		list:   make([]string, lng),
		index:  0,
		length: lng,
	}
}

func (sel *RingStrList) Exist(str string) bool {
	sel.RLock()
	defer sel.RUnlock()
	for i := 0; i < sel.length; i++ {
		if str == sel.list[i] {
			return true
		}
	}
	return false
}

func (sel *RingStrList) Set(str string) {
	sel.Lock()
	sel.list[sel.index] = str
	sel.index++
	if sel.index >= sel.length {
		sel.index = 0
	}
	sel.Unlock()
}

func (sel *RingStrList) Del(str string) {
	sel.Lock()
	for i := 0; i < sel.length; i++ {
		if str == sel.list[i] {
			sel.list[i] = ""
		}
	}
	sel.Unlock()
}

type RingStrListV2 struct {
	sync.RWMutex
	list   map[string]int
	length int
}

func NewRingStrListV2() *RingStrListV2 {
	r := &RingStrListV2{
		list:   make(map[string]int),
		length: 0,
	}

	go func() {
		tmr := time.NewTicker(time.Minute * 5)
		for {
			select {
			case <-tmr.C:
				if r.length > 10000000 {
					cp := 0
					list := make(map[string]int)
					for k, val := range r.list {
						if val > 5000 {
							list[k] = 10
						}
						cp++
						if cp > 500000 {
							break
						}
					}
					r.list = list
				}
			case <-conf.GlobalExitSignal:
				return
			}
		}
	}()

	return r
}

func (sel *RingStrListV2) Exist(str string) bool {
	sel.RLock()
	defer sel.RUnlock()
	n, ok := sel.list[str]
	sel.list[str] = n + 1
	return ok
}

func (sel *RingStrListV2) Set(str string) {
	sel.Lock()
	sel.list[str] = 1
	sel.length++
	sel.Unlock()
}

func (sel *RingStrListV2) Del(str string) {
	sel.Lock()
	delete(sel.list, str)
	sel.length--
	sel.Unlock()
}

type Erc20Contract struct {
	Address     string // 合约地址
	Name        string // 名字
	Symbol      string // 标识
	TotalSupply string // 发行量
	Decimal     int
	Balance     string
}

func (sel *Erc20Contract) IsErc20() bool {
	//return sel.Name != "" && sel.Symbol != "" && sel.TotalSupply != "" && sel.DecimalBit != 0
	// 自要存在发行总量就认为是 token
	return sel.TotalSupply != ""
}

type ContractManager struct {
	rpcCli        ethrpc.EthRPC
	contractRepo  repo.ContractRepository
	haveWriteList *RingStrListV2
	sync.Mutex
}

func NewContractManager(rpcCli ethrpc.EthRPC, contractRepo repo.ContractRepository) *ContractManager {
	return &ContractManager{rpcCli: rpcCli, contractRepo: contractRepo, haveWriteList: NewRingStrListV2()}
}

// GetErc20ContractInfo ERC20 协议的合约有固定的合约接口来获取合约的基本信息
// return *Erc20Contract 合约的基本信息
func (sel *ContractManager) GetErc20ContractInfo(contractAddr string) (*Erc20Contract, error) {
	supply, err := sel.rpcCli.GetContractTotalSupply(contractAddr)
	if err != nil {
		return nil, err
	}
	supply = common.HexToHash(supply).Big().String()
	//name, _ := sel.rpcCli.GetContractName(contractAddr)
	symbol, _ := sel.rpcCli.GetContractSymbol(contractAddr)
	decimal, err := sel.rpcCli.GetContractDecimals(contractAddr)
	if err != nil {
		// 获取小数位错误
		decimal = "0x0"
	}
	decimal = common.HexToHash(decimal).Big().String()
	var decimalInt int64
	decimalInt, _ = strconv.ParseInt(decimal, 10, 64)

	return &Erc20Contract{
		Address: contractAddr,
		//Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Decimal:     int(decimalInt),
	}, nil
}

// TxWrite 合约信息写入，一笔交易存在两个地址，一个是from 地址，一个是 to 地址，两个地址都有可能是合约地址，
// 如果是  token transfer 交易，那么 to 地址一定是合约地址
func (sel *ContractManager) TxWrite(txData *model.TransactionData) (err error) {
	// vlog.DEBUG("tx writer to contract information %s", txData.Hash)

	if txData.IsContractToken {
		return sel.writeTokenContractInfo(txData.ContractAddress, txData.TimeStamp)
	}

	// 检查 to 地址
	if err = sel.writeTokenContractInfo(txData.To, txData.TimeStamp); err != nil {
		return nil
	}

	return nil
}

// writeTokenContractInfo 代币合约信息
func (sel *ContractManager) writeTokenContractInfo(addr string, timeStamp string) (err error) {
	sel.Lock()
	defer sel.Unlock()
	if sel.haveWriteList.Exist(addr) {
		//vlog.WARN("合约已经存在")
		return nil
	}
	sel.haveWriteList.Set(addr)
	// 如果存在合约地址，也要到主链中判断该地址是不是合约地址
	codeData, err := sel.rpcCli.GetCode(addr)
	if err != nil {
		sel.haveWriteList.Del(addr)
		return err
	}
	if codeData == "0x" || codeData == "" {
		// 不是合约地址直接返回
		return nil
	}

	// 查询是否已经存在记录
	if sel.contractRepo.IsContractExist(addr) {
		return nil
	}

	// erc20 合约
	erc20ContractInfo, _ := sel.GetErc20ContractInfo(addr)

	contractInfo := &model.ContractContent{
		Address:     addr,
		PublishTime: timeStamp,
	}
	if erc20ContractInfo != nil {
		contractInfo.TotalSupply = erc20ContractInfo.TotalSupply
		contractInfo.Symbol = erc20ContractInfo.Symbol
		contractInfo.IsErc20 = erc20ContractInfo.IsErc20()
		contractInfo.DecimalBit = erc20ContractInfo.Decimal // 这里加上小数位
	}
	// 如果不存在就创建
	if err = sel.contractRepo.CreateContract(contractInfo); err != nil {
		sel.haveWriteList.Del(addr)
		vlog.ERROR("ContractManager.writeTokenContractInfo create contract info failed address:%s error:%s", addr, err.Error())
		return errors.Wrap(err, "create contract info")
	}
	return
}
