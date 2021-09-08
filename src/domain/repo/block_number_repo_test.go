package repo

import (
	"fmt"
	"os"
	"testing"
)

func TestNewBlockNumberRepoV2(t *testing.T) {
	f, err := os.OpenFile("sync_data.json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Write([]byte("{}")))

	NewBlockNumberRepoV2(f)
}
