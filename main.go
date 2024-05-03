package main

import (
	"embed"
	"fmt"
	"github.com/giwty/switch-library-manager/db"
	"github.com/giwty/switch-library-manager/settings"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get executable directory, please ensure app has sufficient permissions. Aborting.")
		return
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get working directory. Aborting.")
		return
	}

	appSettings := settings.ReadSettings(workingDirectory)

	logger := createLogger(workingDirectory, appSettings.Debug)
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	sugar.Info("[SLM starts]")
	sugar.Infof("[Executable: %v]", exePath)
	sugar.Infof("[Working directory: %v]", workingDirectory)

	localDbManager, err := db.NewLocalSwitchDBManager(workingDirectory)
	if err != nil {
		sugar.Error("Failed to create local files db\n", err)
		return
	}

	_, err = settings.InitSwitchKeys(workingDirectory)
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
