package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type DataCursorRepository interface {
	GetTxRecordAroundBlockNo(blockNo int64, blockSize int64) ([]*model.TransactionRecord, error)

	UpdateFinishInfo(blockNo int64) error

	QueryCursorInfo(typ string) (string, error)

	CreateCursorInfo(typ string, content string) error
}
