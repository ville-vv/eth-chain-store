package conf

import (
	"fmt"
	"testing"
)

func TestGetEthBusinessDbConfig(t *testing.T) {
	mysqlCfg := GetEthBusinessDbConfig()
	fmt.Println(mysqlCfg)
}
