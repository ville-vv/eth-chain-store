package async

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/domain/valobj"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

const (
	CursorTypeContractTx = "ContractTxScanCursor"
	CursorTypeEthereumTx = "EthereumTxScanCursor"
)

type DataCursorAggregate struct {
	cursorType string
	dC         valobj.DataCursor
	cursorRepo repo.DataCursorRepository
}

func NewDataCursorAggregate(cursorType string, cursorRepo repo.DataCursorRepository) *DataCursorAggregate {
	return &DataCursorAggregate{cursorType: cursorType, cursorRepo: cursorRepo}
}

func (sel *DataCursorAggregate) Init() error {
	cursorInfo, err := sel.cursorRepo.QueryCursorInfo(sel.cursorType)
	if err != nil {
		return err
	}
	if cursorInfo == "" {
		sel.dC.BlockSize = 50
		cursorInfo, err = jsoniter.MarshalToString(sel.dC)
		if err != nil {
			return err
		}
		return sel.cursorRepo.CreateCursorInfo(sel.cursorType, cursorInfo)
	}
	return jsoniter.UnmarshalFromString(cursorInfo, &sel.dC)
}

func (sel *DataCursorAggregate) tBName() string {
	switch sel.cursorType {
	case CursorTypeContractTx:
		return "contract_transaction_records_orc"
	case CursorTypeEthereumTx:
		return "transaction_records_orc"
	default:
		panic("sel.cursorType not support")
	}
	return ""
}

func (sel *DataCursorAggregate) GetTxData() ([]*model.TransactionRecord, error) {
	return sel.cursorRepo.GetTxRecordAroundBlockNo(sel.tBName(), sel.dC.BlockNo, sel.dC.BlockSize)
}

func (sel *DataCursorAggregate) Finish() error {
	sel.dC.BlockNo = sel.dC.BlockNo + sel.dC.BlockSize - 1
	vlog.INFO("完成数据块 [%s] [%d]", sel.cursorType, sel.dC.BlockNo)
	return sel.cursorRepo.UpdateFinishInfo(sel.cursorType, sel.dC.ToString())
}
