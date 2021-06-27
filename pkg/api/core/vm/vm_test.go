package vm

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var testStruct VirtualMachine
	js, err := json.Marshal(testStruct)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(string(js))
}
