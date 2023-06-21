package pkg

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	file, err := Parse("../go.mod")
	if err != nil {
		t.Errorf("it should not return error: %s", err.Error())
		return
	}

	bts, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		t.Errorf("MarshalIndent error: %s", err.Error())
		return
	}
	fmt.Println(string(bts))
}
