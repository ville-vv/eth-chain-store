package repo

type IBlockNumberRepo interface {
	UpdateBlockNumber(bkNum uint64) error
}
