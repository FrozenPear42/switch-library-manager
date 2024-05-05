package main

import (
	"context"
	"github.com/FrozenPear42/switch-library-manager/data"
	"github.com/FrozenPear42/switch-library-manager/db"
	"github.com/FrozenPear42/switch-library-manager/settings"
	"github.com/FrozenPear42/switch-library-manager/storage"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

type SwitchTitle struct {
	Name        string `json:"name"`
	TitleId     string `json:"titleId"`
	Icon        string `json:"icon"`
	Banner      string `json:"banner"`
	Region      string `json:"region"`
	ReleaseDate int    `json:"releaseDate"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Publisher   string `json:"publisher"`
	InLibrary   bool   `json:"inLibrary"`
}

type OrganizeSettings struct {
}

type StartupProgressPayload struct {
	Completed bool   `json:"completed"`
	Running   bool   `json:"running"`
	Message   string `json:"message"`
	Current   int    `json:"current"`
	Total     int    `json:"total"`
}

type EventMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type EventType string

const (
	EventTypeStartupProgress EventType = "startupProgress"
)

// App struct
type App struct {
	ctx                context.Context
	sugarLogger        *zap.SugaredLogger
	fullDB             storage.SwitchDatabase
	localDbManager     *db.LocalSwitchDBManager
	configProvider     settings.ConfigurationProvider
	recentStartupEvent EventMessage
}

// NewApp creates a new App application struct
func NewApp(sugarLogger *zap.SugaredLogger, configProvider settings.ConfigurationProvider, database storage.SwitchDatabase) *App {
	return &App{
		sugarLogger:    sugarLogger,
		fullDB:         database,
		configProvider: configProvider,
		localDbManager: nil,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = a.initializeSwitchDB()
}

func (a *App) initializeSwitchDB() error {
	updateProgress := func(step, total int, message string) {
		a.sugarLogger.Infof("progress update: %v/%v %v", step, total, message)
		eventMessage := EventMessage{
			Type: string(EventTypeStartupProgress),
			Data: StartupProgressPayload{
				Completed: step == total,
				Running:   step > 0 && step != total,
				Message:   message,
				Current:   step,
				Total:     total,
			},
		}
		a.recentStartupEvent = eventMessage
		runtime.EventsEmit(a.ctx, string(EventTypeStartupProgress), eventMessage)
	}

	config := a.configProvider.GetCurrentConfig()

	err := data.BuildCatalog(a.fullDB, config.TitlesEndpoint, config.VersionsEndpoint, updateProgress)
	if err != nil {
		a.sugarLogger.Errorf("could not build title catalog: %v", err)
		runtime.Quit(a.ctx)

	}
	return nil
}

func (a *App) RequestStartupProgress() {
	runtime.EventsEmit(a.ctx, string(EventTypeStartupProgress), a.recentStartupEvent)
}

func (a *App) LoadCatalog() ([]SwitchTitle, error) {
	var result []SwitchTitle
	entries, err := a.fullDB.GetCatalogEntries(nil, 0, 0)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		result = append(result, SwitchTitle{
			TitleId:     entry.ID,
			Name:        entry.Name,
			Banner:      entry.BannerURL,
			Icon:        entry.IconURL,
			Region:      entry.Region,
			Version:     entry.Version,
			Description: entry.Description,
			// TODO: update
			Publisher:   "",
			ReleaseDate: 0,
			InLibrary:   false,
		})
	}
	return result, nil
}

//func (a *App) OrganizeLibrary() {
//	//folderToScan := settings.ReadSettings(g.baseFolder).AppDataDirectory
//	//options := settings.ReadSettings(g.baseFolder).OrganizeOptions
//	//if !process.IsOptionsValid(options) {
//	//	zap.S().Error("the organize options in settings.json are not valid, please check that the template contains file/folder name")
//	//	g.state.window.SendMessage(Message{Name: "error", Payload: "the organize options in settings.json are not valid, please check that the template contains file/folder name"}, func(m *astilectron.EventMessage) {})
//	//	return
//	//}
//	//process.OrganizeByFolders(folderToScan, g.state.localDB, g.state.switchDB, g)
//	//if settings.ReadSettings(g.baseFolder).OrganizeOptions.DeleteOldUpdateFiles {
//	//	process.DeleteOldUpdates(g.baseFolder, g.state.localDB, g)
//	//}
//}

//func (a *App) CheckUpdate() (string, error) {
//	//recentVersion, isUpdateAvailable, err := settings.CheckForUpdates()
//	//if err != nil {
//	//	a.sugarLogger.Error(err)
//	//	return "", false, err
//	//}
//	//
//	//return recentVersion, isUpdateAvailable, nil
//	return "", nil
//}
