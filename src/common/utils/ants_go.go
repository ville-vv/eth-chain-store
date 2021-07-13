package utils

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

type WorkGo interface {
	Submit(task func()) error
	Release()
}

var antGo WorkGo
var antGoOnce sync.Once

func RunGo() {
	antGoOnce.Do(func() {
		antsPool, err := ants.NewPool(10000000)
		if err != nil {
			panic(err)
		}
		antGo = antsPool
	})
}
func DisGo() {
	antGo.Release()
}
func Go(f func()) {
	_ = antGo.Submit(f)
}
