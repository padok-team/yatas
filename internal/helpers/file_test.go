package helpers

import (
	"fmt"
	"testing"
)

func TestReadFile(t *testing.T) {
	data, err := ReadFile("./file_test.go")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", data)
}
