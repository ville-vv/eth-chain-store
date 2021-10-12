package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGenRandIntStr(t *testing.T) {
	strMap := make(map[string]int)
	strList := make([]string, 0, 100000)

	stTime := time.Now().UnixNano()
	for i := 0; i < 50000000; i++ {
		str := GenRandIntStr(64)
		strMap[str] = i
		if i < 100000 {
			strList = append(strList, str)
		}
	}
	emdTime := time.Now().UnixNano()
	fmt.Println((emdTime - stTime) / 100000000)

	stTime = time.Now().UnixNano()
	for i := 0; i < len(strList); i++ {
		a, ok := strMap[strList[i]]
		if ok {
			_ = a
			fmt.Println(a)
		}
	}
	emdTime = time.Now().UnixNano()
	fmt.Println((emdTime - stTime) / 100000)
}
