package conf

import (
	"flag"
	"os"
	"strings"
)

var (
	//syncInterval      string
	//fastSyncInterval  string
	RpcUrl              string
	DbUser              string
	DbPassword          string
	DbHost              string
	DbPort              string
	LogFile             string
	SaveToSqlFileDbName string
	Debug               bool
	MaxPullNum          int
	MaxWriteNum         int
	MaxBatchInsertNum   int
	MaxSqlFileSize      int
	IsMaxProcs          bool
	IsHelp              bool
	WriteToDbInterval   int
	TxDataInHive        bool
	WithTxBalance       bool
	StartBlockNumber    int64 // 开始区块
	EndBlockNumber      int64 // 结束区块
	SaveAccount         bool
	SaveContract        bool
	SaveTransaction     bool
	SaveType            string
)

func ReadEnv(val *string, key string) {
	if envVal := os.Getenv(key); strings.TrimSpace(envVal) != "" {
		*val = envVal
	}
}

func ReadFlag(val *string, key string) {
	f := flag.Lookup(key)
	if f == nil {
		return
	}
	if v := f.Value.String(); v != "" {
		*val = v
	}
}

var (
	GlobalExitSignal          = make(chan int)
	GlobalProgramFinishSigmal = make(chan int)
)

var globalConfig *GlobalConfig

type GlobalConfig struct {
	MysqlCfg   *MysqlConf      `json:"mysql_cfg" name:"" toml:"Mysql"`
	BlkSyncCfg BlockSyncConfig `json:"blk_sync_cfg" name:""`
}

type BlockSyncConfig struct {
}
