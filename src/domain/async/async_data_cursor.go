package async

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/domain/valobj"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type DataCursorAggregate struct {
	cursorType string
	valobj.DataCursor
	cursorRepo repo.DataCursorRepository
}

func NewDataCursorAggregate(cursorType string, cursorRepo repo.DataCursorRepository) *DataCursorAggregate {
	return &DataCursorAggregate{cursorType: cursorType, cursorRepo: cursorRepo}
}

func (sel *DataCursorAggregate) Init() error {
	var dtc = &valobj.DataCursor{}
	cursorInfo, err := sel.cursorRepo.QueryCursorInfo(sel.cursorType)
	if err != nil {
		return err
	}
	if cursorInfo == "" {
		dtc.BlockSize = 50
		cursorInfo, err = jsoniter.MarshalToString(dtc)
		if err != nil {
			return err
		}
		return sel.cursorRepo.CreateCursorInfo(sel.cursorType, cursorInfo)
	}
	return jsoniter.UnmarshalFromString(cursorInfo, dtc)
}

func (sel *DataCursorAggregate) GetTxData() ([]*model.TransactionRecord, error) {
	return sel.cursorRepo.GetTxRecordAroundBlockNo(sel.BlockNo, sel.BlockSize)
}

func (sel *DataCursorAggregate) Finish() error {
	return sel.cursorRepo.UpdateFinishInfo(sel.BlockNo)
}
