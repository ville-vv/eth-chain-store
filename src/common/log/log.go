package log

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/utils"
	"github.com/ville-vv/vilgo/vlog"
	"os"
)

func init()  {
	logLevel := vlog.LogLevelInfo
	logFile := fmt.Sprintf("eth_chain_store_%s.log", utils.RandStringBytesMask(8))
	if os.Getenv("ETH_CHAIN_STORE_ENV") == "local"{
		logLevel = vlog.LogLevelDebug
		logFile = "stdout"
		return
	}
	cnf := &vlog.LogCnf{
		OutPutErrFile: []string{logFile},
		ProgramName:   "eth-chain-store",
		Level:         logLevel,
	}
	vlog.SetLogger(vlog.NewGoLogger(cnf))
	return
}

func DEBUG(format string, args ...interface{}) {
	vlog.LogD(format, args...)
}

func ERROR(format string, args ...interface{}) {
	vlog.LogE(format, args...)
}

func INFO(format string, args ...interface{}) {
	vlog.LogI(format, args...)
}
func WARN(format string, args ...interface{}) {
	vlog.LogW(format, args...)
}