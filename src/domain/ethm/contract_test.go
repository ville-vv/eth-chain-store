package ethm

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/utils"
	"sync"
	"testing"
	"time"
)

// c=测试map + Mutex 效率问题
func TestNewRingStrListV2(t *testing.T) {
	rList := NewRingStrListV2()
	strList := make([]string, 0, 100000)
	var str = ""
	stTime := time.Now().UnixNano()

	var wg sync.WaitGroup
	wg.Add(100)
	for n := 0; n < 100; n++ {
		go func() {
			stTime := time.Now().UnixNano()
			for i := 0; i < 50000; i++ {
				str = utils.GenRandIntStr(64)
				rList.Set(str)
				if i < 1000 {
					strList = append(strList, str)
				}
			}
			emdTime := time.Now().UnixNano()
			fmt.Println("插入：", (emdTime-stTime)/1000000000)
			wg.Done()
		}()
	}
	wg.Wait()
	stTime = time.Now().UnixNano()

	for n := 0; n < 100; n++ {
		go func() {
			//for i := 0; i < len(strList); i++ {
			//	ok := rList.Exist(strList[i])
			//	if ok {
			//	}
			//}
			rList.Exist(str)
			emdTime := time.Now().UnixNano()
			fmt.Println("读出：", (emdTime-stTime)/100000)
		}()
	}
}
