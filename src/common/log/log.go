package log

import (
	"flag"
	"github.com/ville-vv/vilgo/vlog"
)

func Init() {
	logLevel := vlog.LogLevelInfo
	var logFile = "stdout"
	//logFile := fmt.Sprintf("logs/eth_chain_store_%s.log", utils.RandStringBytesMask(8))
	//if os.Getenv("ETH_CHAIN_STORE_ENV") == "local" {
	//	logLevel = vlog.LogLevelDebug
	//	logFile = "stdout"
	//	return
	//}

	if logVal := flag.Lookup("logFile"); logVal != nil {
		valStr := logVal.Value.String()
		if valStr != "" {
			logFile = logVal.Value.String()
			//logFile = fmt.Sprintf("%s%s.log", logFile, utils.RandStringBytesMask(8))
		}
	}

	if logVal := flag.Lookup("debug"); logVal != nil {
		if logVal.Value.String() == "true" {
			logLevel = vlog.LogLevelDebug
		}
	}
	cnf := &vlog.LogCnf{
		OutPutFile:  []string{logFile},
		ProgramName: "eth-chain-store-mysql",
		Level:       logLevel,
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
