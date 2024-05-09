package utils

import (
	"hash/crc32"
	"io"
	"os"
)

func FileCRC32(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := crc32.New(crc32.IEEETable)
	_, err = io.Copy(hasher, f)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}
