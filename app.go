package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/data"
	"github.com/FrozenPear42/switch-library-manager/keys"
	"github.com/FrozenPear42/switch-library-manager/settings"
	"github.com/FrozenPear42/switch-library-manager/storage"
	"github.com/FrozenPear42/switch-library-manager/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"path/filepath"
	"sync"
)

type SwitchTitle struct {
	Name        string               `json:"name"`
	TitleID     string               `json:"titleID"`
	Icon        string               `json:"icon"`
	Banner      string               `json:"banner"`
	Region      string               `json:"region"`
	ReleaseDate string               `json:"releaseDate"`
	Version     string               `json:"version"`
	Description string               `json:"description"`
	Intro       string               `json:"intro"`
	Publisher   string               `json:"publisher"`
	InLibrary   bool                 `json:"inLibrary"`
	Screenshots []string             `json:"screenshots"`
	DLCs        []SwitchTitle        `json:"dlcs"`
	Versions    []SwitchTitleVersion `json:"versions"`
}

type SwitchTitleVersion struct {
	Version     int    `json:"version"`
	ReleaseDate string `json:"releaseDate"`
}

type CatalogFiltersSortBy string

type CatalogFilters struct {
	SortBy storage.CatalogFiltersSortBy `json:"sortBy"`
	Name   *string                      `json:"name"`
	ID     *string                      `json:"id"`
	Region []string                     `json:"region"`
	Cursor int                          `json:"cursor"`
	Limit  int                          `json:"limit"`
}

type CatalogPage struct {
	Titles      []SwitchTitle `json:"titles"`
	TotalTitles int           `json:"totalTitles"`
	NextCursor  int           `json:"nextCursor"`
	IsLastPage  bool          `json:"isLastPage"`
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
	mutex              sync.Mutex
	ctx                context.Context
	sugarLogger        *zap.SugaredLogger
	fullDB             storage.SwitchDatabase
	configProvider     settings.ConfigurationProvider
	libraryManager     data.LibraryManager
	recentStartupEvent EventMessage
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		mutex: sync.Mutex{},
	}
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
				runtime.Quit(a.ctx)
			}
		} else {
			fmt.Printf("Failed to load configuration file. Aborting. (%v)\n", err)
			runtime.Quit(a.ctx)
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
		runtime.Quit(a.ctx)
	}

	libraryManager := data.NewLibraryManager(logger.Sugar(), keyProvider, config.ScanDirectories)

	a.fullDB = database
	a.configProvider = configurationProvider
	a.sugarLogger = logger.Sugar()
	a.libraryManager = libraryManager

	a.sugarLogger.Infof("startup")

	a.mutex.Lock()
	defer a.mutex.Unlock()
	err = a.initializeSwitchDB()
	if err != nil {
		sugar.Error("Failed to initialize database\n", err)
		runtime.Quit(a.ctx)
	}

	err = a.libraryManager.Rescan(false, func(current, total int, message string) {
		sugar.Debugf("processing: %v/%v, %v", current, total, message)
	})
	if err != nil {
		sugar.Error("Failed to init scan files\n", err)
		runtime.Quit(a.ctx)
	}

	a.sugarLogger.Infof("initialized")
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

func (a *App) LoadCatalog(filters CatalogFilters) (CatalogPage, error) {
	a.sugarLogger.Infof("load catalog request %v", filters)

	var result []SwitchTitle
	catalogFilters := &storage.CatalogFilters{
		SortBy: filters.SortBy,
		Name:   filters.Name,
		ID:     filters.ID,
		Region: filters.Region,
	}
	page, err := a.fullDB.GetCatalogEntries(catalogFilters, filters.Limit, filters.Cursor)
	if err != nil {
		return CatalogPage{}, err
	}

	for _, entry := range page.Data {
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
				Publisher:   dlc.Publisher,
				ReleaseDate: dlc.ReleaseDate,
				Screenshots: dlc.Screenshots,
				// TODO: update
				InLibrary: false,
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
			ReleaseDate: entry.ReleaseDate,
			Publisher:   entry.Publisher,
			Screenshots: entry.Screenshots,
			// TODO: update
			InLibrary: false,

			DLCs:     dlcs,
			Versions: versions,
		})
	}
	a.sugarLogger.Infof("loading catalog of %v items to frontend", len(result))
	return CatalogPage{
		Titles:      result,
		TotalTitles: page.TotalCount,
		NextCursor:  page.NextCursor,
		IsLastPage:  page.IsLastPage,
	}, nil
}

type LibraryFileEntry struct {
	FilePath string `json:"filePath"`
	FileSize int    `json:"fileSize"`
}

func (a *App) LoadLibraryFiles() ([]LibraryFileEntry, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	fmt.Println(a)
	fmt.Println(a.libraryManager)
	fmt.Println(a.sugarLogger)
	a.sugarLogger.Debugf("manager: %v", a.libraryManager)
	entries, err := a.libraryManager.GetEntries()
	if err != nil {
		return nil, err
	}
	var res []LibraryFileEntry
	for _, e := range entries {
		res = append(res, LibraryFileEntry{
			FilePath: e.FilePath,
			FileSize: e.FileSize,
		})
	}
	return res, nil
}

//type LibraryGameEntry struct {
//	TitleID string
//	Name string
//
//}

//func (a *App) LoadLibraryGames() ([]data.LibraryFileEntry, error) {
//	files, err :=  a.libraryManager.GetEntries()
//
//	for _, fileEntry := range files {
//
//	}
//}

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
