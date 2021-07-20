package repo

type BlockNumberRepo interface {
	UpdateBlockNumber(bkNum int64) error
}
