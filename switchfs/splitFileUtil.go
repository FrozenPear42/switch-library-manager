package switchfs

import (
	"errors"
	"github.com/FrozenPear42/switch-library-manager/keys"
	"os"
)

func ReadSplitFileMetadata(keyProvider keys.KeysProvider, filePath string) (map[string]*ContentMetaAttributes, error) {
	//check if this is a NS* or XC* file
	_, err := ReadPfs0File(filePath)
	isXCI := false
	if err != nil {
		_, err = readXciHeader(filePath)
		if err != nil {
			return nil, errors.New("split file is not an XCI/XCZ or NSP/NSZ")
		}
		isXCI = true
	}

	if isXCI {
		return ReadXciMetadata(keyProvider, filePath)
	} else {
		return ReadNspMetadata(keyProvider, filePath)
	}
}

func readXciHeader(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	header := make([]byte, 0x200)
	_, err = file.Read(header)
	if err != nil {
		return nil, err
	}

	if string(header[0x100:0x104]) != "HEAD" {
		return nil, errors.New("not an XCI/XCZ file")
	}
	return header, nil
}
