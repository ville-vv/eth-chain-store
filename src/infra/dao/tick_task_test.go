package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestTickTask_Exit(t *testing.T) {
	exec := func() {
		fmt.Println(time.Now())
		time.Sleep(time.Second * 4)
	}
	tsk := NewTickTask("aaa", time.Second*1, exec)
	tsk.Start()
	select {}
}
