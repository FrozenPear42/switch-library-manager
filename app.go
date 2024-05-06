package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/data"
	"github.com/FrozenPear42/switch-library-manager/db"
	"github.com/FrozenPear42/switch-library-manager/keys"
	"github.com/FrozenPear42/switch-library-manager/settings"
	"github.com/FrozenPear42/switch-library-manager/storage"
	"github.com/FrozenPear42/switch-library-manager/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"path/filepath"
)

type SwitchTitle struct {
	Name        string               `json:"name"`
	TitleID     string               `json:"titleID"`
	Icon        string               `json:"icon"`
	Banner      string               `json:"banner"`
	Region      string               `json:"region"`
	ReleaseDate int                  `json:"releaseDate"`
	Version     string               `json:"version"`
	Description string               `json:"description"`
	Intro       string               `json:"intro"`
	Publisher   string               `json:"publisher"`
	InLibrary   bool                 `json:"inLibrary"`
	DLCs        []SwitchTitle        `json:"DLCs"`
	Versions    []SwitchTitleVersion `json:"versions"`
}

type SwitchTitleVersion struct {
	Version     int    `json:"version"`
	ReleaseDate string `json:"releaseDate"`
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
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	workingDirectory, err := utils.GetExecDir()
	if err != nil {
		fmt.Println("Failed to get working directory. Aborting.")
		runtime.Quit(a.ctx)
	}

	configurationProvider, err := settings.NewConfigurationProvider(filepath.Join(workingDirectory, "settings.yaml"))
	if err != nil {
		fmt.Printf("Failed to initialize config provider. Aborting. Reason: %v\n", err)
		runtime.Quit(a.ctx)
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

	logger := createLogger(workingDirectory, true)
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	sugar.Info("[SLM starts]")
	sugar.Infof("[Working directory: %v]", workingDirectory)

	database, err := storage.NewDatabase(filepath.Join(workingDirectory, "slm_full.db"))
	if err != nil {
		sugar.Error("Failed to initialize database\n", err)
		runtime.Quit(a.ctx)
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

	a.fullDB = database
	a.configProvider = configurationProvider
	a.sugarLogger = logger.Sugar()

	err = a.initializeSwitchDB()
	if err != nil {
		sugar.Error("Failed to initialize database\n", err)
		runtime.Quit(a.ctx)
	}

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
		dlcs := make([]SwitchTitle, 0, len(entry.DLCs))
		for _, dlc := range entry.DLCs {
			dlcs = append(dlcs, SwitchTitle{
				Name:        dlc.Name,
				TitleID:     dlc.ID,
				Icon:        dlc.IconURL,
				Banner:      dlc.BannerURL,
				Region:      dlc.Region,
				Version:     dlc.Version,
				Description: dlc.Description,
				Intro:       dlc.Intro,

				// TODO: update
				Publisher:   "",
				ReleaseDate: 0,
				InLibrary:   false,
			})
		}

		versions := make([]SwitchTitleVersion, 0, len(entry.Versions))
		for _, version := range entry.Versions {
			versions = append(versions, SwitchTitleVersion{
				Version:     version.Version,
				ReleaseDate: version.ReleaseDate,
			})
		}

		result = append(result, SwitchTitle{
			Name:        entry.Name,
			TitleID:     entry.ID,
			Icon:        entry.IconURL,
			Banner:      entry.BannerURL,
			Region:      entry.Region,
			Version:     entry.Version,
			Description: entry.Description,
			Intro:       entry.Intro,
			// TODO: update
			ReleaseDate: 0,
			Publisher:   "",
			InLibrary:   false,

			DLCs:     dlcs,
			Versions: versions,
		})
	}
	a.sugarLogger.Infof("loading catalog of %v items to frontend", len(result))
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
