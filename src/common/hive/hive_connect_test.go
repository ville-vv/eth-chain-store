package hive

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	fmt.Println(Connect("localhost:10000"))
}
