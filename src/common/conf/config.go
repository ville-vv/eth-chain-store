package conf

import (
	"flag"
	"os"
	"strings"
)

func readEnv(val *string, key string) {
	if envVal := os.Getenv(key); strings.TrimSpace(envVal) != "" {
		*val = envVal
	}
}

func readFlag(val *string, key string) {
	f := flag.Lookup(key)
	v := f.Value.String()
	if v != "" {

	}
}

var globalConfig *GlobalConfig

type GlobalConfig struct {
	MysqlCfg *MysqlConf `json:"mysql_cfg" name:"" toml:"Mysql"`
}
