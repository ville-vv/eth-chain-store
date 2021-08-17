package monitor

import (
	"fmt"
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
`, MqSize.Load(), AccountWriteProcessNum.Load(), ContractWriteProcessNum.Load(), TxWriteProcessNum.Load())
		writer.Write([]byte(bodyStr))
	})
	http.ListenAndServe("0.0.0.0:5489", nil)
}
