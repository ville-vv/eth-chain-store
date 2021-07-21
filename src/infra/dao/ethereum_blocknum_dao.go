package dao

type EthereumBlockNumberDao struct {
	db DB
}

func NewEthereumBlockNumberDao(db DB) *EthereumBlockNumberDao {
	return &EthereumBlockNumberDao{db: db}
}
