package helpers

import (
	"io"
	"os"
)

// Read a file and return its content as a byte array
func ReadFile(configPath string) ([]byte, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
