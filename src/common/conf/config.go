package conf

import (
	"flag"
	"os"
	"strings"
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
	GlobalExitSignal = make(chan int)
)

var globalConfig *GlobalConfig

type GlobalConfig struct {
	MysqlCfg   *MysqlConf      `json:"mysql_cfg" name:"" toml:"Mysql"`
	BlkSyncCfg BlockSyncConfig `json:"blk_sync_cfg" name:""`
}

type BlockSyncConfig struct {
}
