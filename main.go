package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/db"
	"github.com/FrozenPear42/switch-library-manager/keys"
	"github.com/FrozenPear42/switch-library-manager/settings"
	"github.com/FrozenPear42/switch-library-manager/utils"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"go.uber.org/zap"
	"path/filepath"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	workingDirectory, err := utils.GetExecDir()
	if err != nil {
		fmt.Println("Failed to get working directory. Aborting.")
		return
	}

	configurationProvider, err := settings.NewConfigurationProvider(filepath.Join(workingDirectory, "settings.yaml"))
	if err != nil {
		fmt.Printf("Failed to initialize config provider. Aborting. Reason: %v\n", err)
		return
	}

	err = configurationProvider.LoadFromFile()
	if err != nil {
		if errors.Is(err, settings.ErrConfigurationFileNotFound) {
			fmt.Println("Configuration file not found. Using default configuration. Creating new configuration file.")
			err = configurationProvider.SaveToFile()
			if err != nil {
				fmt.Println("Could not save new configuration file. Aborting.")
				return
			}
		} else {
			fmt.Printf("Failed to load configuration file. Aborting. (%v)\n", err)
			return
		}
	}

	config := configurationProvider.GetCurrentConfig()

	logger := createLogger(workingDirectory, config.Debug)
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	sugar.Info("[SLM starts]")
	sugar.Infof("[Working directory: %v]", workingDirectory)

	localDbManager, err := db.NewLocalSwitchDBManager(workingDirectory)
	if err != nil {
		sugar.Error("Failed to create local files db\n", err)
		return
	}

	keyProvider := keys.NewKeyProvider()
	keyPaths := []string{
		config.ProdKeysPath,
		filepath.Join(workingDirectory, "prod.keys"),
		"${HOME}/.switch/prod.keys",
	}

	err = keyProvider.LoadFromFile(keyPaths)
	if err != nil {
		sugar.Error("Failed to initialize keys\n", err)
		return
	}

	// Create an instance of the app structure
	app := NewApp(sugar, localDbManager)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Switch Library Manager",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func createLogger(workingFolder string, debug bool) *zap.Logger {
	logPath := filepath.Join(workingFolder, "slm.log")

	// TODO: swap to zapcore and force file
	//f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	_ = logPath

	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("failed to create logger - %v", err)
		panic(1)
	}
	zap.ReplaceGlobals(logger)
	return logger
}
