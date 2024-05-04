package keys

import (
	"fmt"
	"github.com/magiconair/properties"
)

type KeysProvider interface {
	GetProdKey(keyName string) (string, bool)
}

type KeysProviderImpl struct {
	keys map[string]string
}

func NewKeyProvider() *KeysProviderImpl {
	return &KeysProviderImpl{
		keys: make(map[string]string),
	}
}

func (p *KeysProviderImpl) LoadFromFile(paths []string) error {
	var prodKeys *properties.Properties

	for _, path := range paths {
		keyFile, err := properties.LoadFile(path, properties.UTF8)
		if err == nil {
			prodKeys = keyFile
			break
		}
	}
	if prodKeys == nil {
		return fmt.Errorf("could not open files in any of provided paths")
	}

	p.keys = make(map[string]string)
	for _, key := range prodKeys.Keys() {
		value, _ := prodKeys.Get(key)
		p.keys[key] = value
	}
	return nil
}

func (p *KeysProviderImpl) GetProdKey(keyName string) (string, bool) {
	k, ok := p.keys[keyName]
	return k, ok
}
