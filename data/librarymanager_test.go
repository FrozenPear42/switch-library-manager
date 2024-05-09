package data

import (
	"github.com/FrozenPear42/switch-library-manager/keys"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestProcessFile(t *testing.T) {
	var zapConfig zap.Config
	zapConfig = zap.NewDevelopmentConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}
	logger, err := zapConfig.Build()

	assert.Nil(t, err)

	keysProvider := keys.NewKeyProvider()
	err = keysProvider.LoadFromFile([]string{"../fixtures/prod.keys"})
	assert.Nil(t, err)

	manager := LibraryManagerImpl{
		logger:          logger.Sugar(),
		db:              nil,
		keysProvider:    keysProvider,
		allowedFormats:  []string{"nsp"},
		scanDirectories: []string{"../fixtures"},
	}

	err = manager.Rescan(true, func(current, total int, message string) {
		logger.Sugar().Debugf("progress: %v/%v: %v", current, total, message)
	})
	assert.Nil(t, err)
}
