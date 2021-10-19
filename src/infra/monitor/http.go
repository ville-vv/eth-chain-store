package monitor

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/shirou/gopsutil/cpu"
	"github.com/ville-vv/vilgo/vtask"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

var (
	MqSize                  vtask.AtomicInt64
	AccountWriteProcessNum  vtask.AtomicInt64
	ContractWriteProcessNum vtask.AtomicInt64
	TxWriteProcessNum       vtask.AtomicInt64
)

func StartMonitor() {
	//	http.HandleFunc("/mq_size", func(writer http.ResponseWriter, request *http.Request) {
	//		bodyStr := fmt.Sprintf(`
	//MQ Pool Size: %d
	//Account Write Process Numbers: %d
	//Contract Write Process Numbers: %d
	//Transaction Write Process Numbers: %d
	//CPU Info : %s
	//`, MqSize.Load(), AccountWriteProcessNum.Load(), ContractWriteProcessNum.Load(), TxWriteProcessNum.Load(), Cpu())
	//		writer.Write([]byte(bodyStr))
	//	})
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	http.ListenAndServe("0.0.0.0:6060", nil)
}

func Cpu() string {
	info, err := cpu.Times(false)
	if err != nil {
		return ""
	}
	str, _ := jsoniter.MarshalToString(info)
	return str
}
