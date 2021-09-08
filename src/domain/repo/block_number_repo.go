package repo

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"io"
	"io/ioutil"
)

//type BlockNumberRepo interface {
//	GetCntSyncBlockNumber() (int64, error)
//	UpdateSyncBlockNUmber(n int64) error
//}

type BlockNumberRepo struct {
	ebnDao *dao.EthereumBlockNumberDao
}

func NewBlockNumberRepo(ebnDao *dao.EthereumBlockNumberDao) *BlockNumberRepo {
	return &BlockNumberRepo{ebnDao: ebnDao}
}

// 初始最新区块
func (sel *BlockNumberRepo) InitLatestBlockNumber(bkNum int64) error {
	_, ok, err := sel.ebnDao.GetSyncBlockConfig("LatestBlockNumber")
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return sel.ebnDao.SetSyncBlockConfig("LatestBlockNumber", fmt.Sprintf("%d", bkNum))
}

func (sel *BlockNumberRepo) UpdateLatestBlockNumber(bkNum int64) error {
	return sel.ebnDao.UpdateLatestBlockNumber(bkNum)
}

func (sel *BlockNumberRepo) GetCntSyncBlockNumber() (int64, error) {
	return sel.ebnDao.GetSyncBlockNumber()
}

func (sel *BlockNumberRepo) UpdateSyncBlockNUmber(n int64) error {
	return sel.ebnDao.UpdateSyncBlockNumber(n)
}

type ConfigData struct {
	LatestBlockNumber  int64 `json:"latest_block_number" name:""`
	CntSyncBlockNumber int64 `json:"cnt_sync_block_number" name:""`
}
type BlockNumberRepoV2 struct {
	fs      io.ReadWriteSeeker
	cfgData *ConfigData
}

func NewBlockNumberRepoV2(rws io.ReadWriteSeeker) *BlockNumberRepoV2 {
	data, err := ioutil.ReadAll(rws)
	if err != nil {
		panic(err)
	}
	var cfgData = new(ConfigData)
	if len(data) != 0 {
		if err = jsoniter.Unmarshal(data, cfgData); err != nil {
			panic(err)
		}
	}
	return &BlockNumberRepoV2{
		fs:      rws,
		cfgData: cfgData,
	}
}

func (sel *BlockNumberRepoV2) InitLatestBlockNumber(bkNum int64) error {
	sel.cfgData.LatestBlockNumber = bkNum
	return nil
}

func (sel *BlockNumberRepoV2) UpdateLatestBlockNumber(bkNum int64) error {
	sel.cfgData.LatestBlockNumber = bkNum
	return sel.writeToFile()
}

func (sel *BlockNumberRepoV2) GetCntSyncBlockNumber() (int64, error) {
	return sel.cfgData.CntSyncBlockNumber, nil
}

func (sel *BlockNumberRepoV2) UpdateSyncBlockNUmber(n int64) error {
	sel.cfgData.CntSyncBlockNumber = n
	return sel.writeToFile()
}

func (sel *BlockNumberRepoV2) writeToFile() error {
	data, err := jsoniter.Marshal(sel.cfgData)
	if err != nil {
		return err
	}
	_, err = sel.fs.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = sel.fs.Write(data)
	return err
}
