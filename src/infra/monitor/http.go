package monitor

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/shirou/gopsutil/cpu"
	"github.com/ville-vv/vilgo/vtask"
	"net/http"
)

var (
	MqSize                  vtask.AtomicInt64
	AccountWriteProcessNum  vtask.AtomicInt64
	ContractWriteProcessNum vtask.AtomicInt64
	TxWriteProcessNum       vtask.AtomicInt64
)

func StartMonitor() {
	http.HandleFunc("/mq_size", func(writer http.ResponseWriter, request *http.Request) {
		bodyStr := fmt.Sprintf(`
MQ Pool Size: %d
Account Write Process Numbers: %d
Contract Write Process Numbers: %d
Transaction Write Process Numbers: %d
CPU Info : %s
`, MqSize.Load(), AccountWriteProcessNum.Load(), ContractWriteProcessNum.Load(), TxWriteProcessNum.Load(), Cpu())
		writer.Write([]byte(bodyStr))
	})
	http.ListenAndServe("0.0.0.0:5489", nil)
}

func Cpu() string {
	info, err := cpu.Times(false)
	if err != nil {
		return ""
	}
	str, _ := jsoniter.MarshalToString(info)
	return str
}
