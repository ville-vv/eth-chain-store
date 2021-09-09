package dao

import (
	"bufio"
	"context"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vfile"
	"github.com/ville-vv/vilgo/vlog"
	"os"
	"path"
	"sync"
	"time"
)

type TransactionHiveDataFile struct {
	normalTxFs   *txHiveDataFile
	contractTxFs *txHiveDataFile
}

func (sel *TransactionHiveDataFile) Scheme() string {
	return "TransactionHiveDataFile"
}

func (sel *TransactionHiveDataFile) Init() error {
	return nil
}

func (sel *TransactionHiveDataFile) Start() error {
	go sel.normalTxFs.Start()
	go sel.contractTxFs.Start()
	return nil
}

func (sel *TransactionHiveDataFile) Exit(ctx context.Context) error {
	sel.normalTxFs.Exit()
	sel.contractTxFs.Exit()
	return nil
}

func NewHiveDataFile(normalTxFile string, contractTxFile string, interval int64) *TransactionHiveDataFile {
	return &TransactionHiveDataFile{
		normalTxFs:   newTxHiveDataFile(normalTxFile, interval),
		contractTxFs: newTxHiveDataFile(contractTxFile, interval),
	}
}

func (sel *TransactionHiveDataFile) CreateTransactionRecord(txData *model.TransactionData) error {
	// vlog.INFO("创建交易数据")
	//return nil
	txRecord := &model.TransactionRecord{
		CreatedAt:       time.Now(),
		BlockNumber:     txData.BlockNumber,
		BlockHash:       txData.BlockHash,
		TxHash:          txData.Hash,
		Timestamp:       txData.TimeStamp,
		ContractAddress: txData.ContractAddress,
		FromAddr:        txData.From,
		ToAddr:          txData.To,
		GasPrice:        txData.GasPrice,
		Value:           txData.Value,
		FromAddrBalance: txData.FromBalance,
		ToAddrBalance:   txData.ToBalance,
	}

	if txData.IsContractToken {
		return sel.contractTxFs.InsertData(txRecord)
	}
	return sel.normalTxFs.InsertData(txRecord)
}

type txHiveDataFile struct {
	interval  int64
	fileName  string
	dataFile  *os.File
	waitExit  chan int
	stopCh    chan int
	dataCache []string
	maxLen    int
	seq       int64
	sync.Mutex
}

func newTxHiveDataFile(fileName string, interval int64) *txHiveDataFile {
	dirPath := path.Dir(fileName)
	if !vfile.PathExists(dirPath) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	seq, err := ReadLineNum(fileName)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error() + fileName)
	}
	return &txHiveDataFile{
		interval:  interval,
		dataFile:  f,
		waitExit:  make(chan int),
		stopCh:    make(chan int),
		maxLen:    100000,
		seq:       seq,
		dataCache: make([]string, 0, 100000),
	}
}

func ReadLineNum(fileName string) (int64, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error() + fileName)
	}
	defer f.Close()
	fd := bufio.NewReader(f)
	count := int64(0)
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	return count, nil
}

func (sel *txHiveDataFile) Start() {
	tmr := time.NewTicker(time.Second * time.Duration(sel.interval))
	for {
		select {
		case <-sel.stopCh:
			sel.save()
			sel.waitExit <- 1
			return
		case <-tmr.C:
			sel.save()
		}
	}
}

func (sel *txHiveDataFile) save() {
	sel.Lock()
	dataCache := sel.dataCache
	l := len(dataCache)
	if sel.maxLen < l {
		sel.maxLen = l
	}
	sel.dataCache = make([]string, 0, sel.maxLen)
	sel.Unlock()

	if l > 0 {
		//vlog.INFO("当前条数：%d", l)
		data := ""
		for i := 0; i < l; i++ {
			data = data + dataCache[i]
		}
		_, err := sel.dataFile.WriteString(data)
		if err != nil {
			vlog.ERROR("write file data failed fileName%s \n%s", sel.fileName, data)
		}
	}
}

func (sel *txHiveDataFile) Exit() {
	close(sel.stopCh)
	vlog.INFO("Hive数据文件写入线程退出等待中")
	<-sel.waitExit
	close(sel.waitExit)
	_ = sel.dataFile.Close()
	vlog.INFO("Hive数据文件写入线程退出正常")
}

func (sel *txHiveDataFile) InsertData(txData *model.TransactionRecord) error {
	sel.Lock()
	sel.seq++
	txData.ID = sel.seq
	sel.dataCache = append(sel.dataCache, txData.String())
	sel.Unlock()
	return nil
}
