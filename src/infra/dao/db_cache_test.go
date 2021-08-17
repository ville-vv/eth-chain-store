package dao

import (
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/common/utils"
	"testing"
	"time"
)

func TestNewDbCache(t *testing.T) {
	log.Init()
	chDb := NewDbCache(nil)
	_ = chDb.Start()
	go func() {
		for {
			time.Sleep(time.Millisecond)
			str := utils.RandStringBytesMask(10)
			chDb.Insert("abc", str)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Millisecond)
			str := utils.RandStringBytesMask(10)
			chDb.Insert("dddcc", str)
		}
	}()
	select {}
}
