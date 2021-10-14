package file

import (
	"testing"
)

func TestAutoSplitFile_WriteString(t *testing.T) {
	fileSplit, err := NewAutoSplitFile("data/testxsdg", 0)
	if err != nil {
		t.Error(err)
		return
	}
	fileSplit.SetMaxSizeMb(20)

	for i := 0; i < 1000000; i++ {
		fileSplit.WriteString("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
	fileSplit.Close()
}
