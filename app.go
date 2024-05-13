package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/data"
	"github.com/FrozenPear42/switch-library-manager/keys"
	"github.com/FrozenPear42/switch-library-manager/nut"
	"github.com/FrozenPear42/switch-library-manager/settings"
	"github.com/FrozenPear42/switch-library-manager/storage"
	"github.com/FrozenPear42/switch-library-manager/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"path/filepath"
	"sync"
)

type ServerReporter struct {
	ctx    context.Context
	logger *zap.SugaredLogger
}

func (s *ServerReporter) ReportProgress(filePath string, downloaded, total int64) {
	s.logger.Infof("file download progress: %s %d/%d", filePath, downloaded, total)
	//runtime.EventsEmit(s.ctx, )
}

// App struct
type App struct {
	mutex              sync.Mutex
	ctx                context.Context
	sugarLogger        *zap.SugaredLogger
	fullDB             storage.SwitchDatabase
	configProvider     settings.ConfigurationProvider
	libraryManager     data.LibraryManager
	nutServer          *nut.Server
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

	reporter := &ServerReporter{
		ctx:    a.ctx,
		logger: logger.Sugar(),
	}
	a.nutServer = nut.NewServer(config.NUTSettings.Host, config.NUTSettings.Port, libraryManager, reporter)

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

	_, err = a.fullDB.GetCatalogEntries(nil, 1, 0)
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

	err = a.nutServer.Listen()
	if err != nil {
		sugar.Error("Failed to start NUT server\n", err)
		runtime.Quit(a.ctx)
	}

	a.sugarLogger.Infof("initialized")
}

func (a *App) initializeSwitchDB() error {
	updateProgress := func(step, total int, message string) {
		a.sugarLogger.Infof("progress update: %v/%v %v", step, total, message)
		eventMessage := EventMessage{
			Type: string(EventTypeStartupProgress),
			Data: EventStartupProgressPayload{
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

	var result []CatalogSwitchGame
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
		dlcs := make([]CatalogDLCData, 0, len(entry.DLCs))
		for _, dlc := range entry.DLCs {
			dlcs = append(dlcs, CatalogDLCData{
				Name:        dlc.Name,
				TitleID:     dlc.ID,
				Banner:      dlc.BannerURL,
				Region:      dlc.Region,
				Version:     dlc.Version,
				Description: dlc.Description,
			})
		}

		versions := make([]CatalogVersionData, 0, len(entry.Versions))
		for _, version := range entry.Versions {
			versions = append(versions, CatalogVersionData{
				Version:     version.Version,
				ReleaseDate: version.ReleaseDate,
			})
		}

		result = append(result, CatalogSwitchGame{
			CatalogGameData: CatalogGameData{
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
			},
			DLCs:     dlcs,
			Versions: versions,
		})
	}
	a.sugarLogger.Infof("loading catalog of %v items to frontend", len(result))
	return CatalogPage{
		Games:      result,
		TotalGames: page.TotalCount,
		NextCursor: page.NextCursor,
		IsLastPage: page.IsLastPage,
	}, nil
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

func (a *App) LoadLibraryGames() ([]LibrarySwitchGame, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	fileEntries, err := a.libraryManager.GetEntries()
	if err != nil {
		return nil, fmt.Errorf("could not get file entries from library: %w", err)
	}

	// TODO: detect duplicates, still aggregate those but flag as duplicates, maybe even just put those in a list with versions?

	// map by title prefix
	titles := make(map[string]*LibrarySwitchGame)

	for _, fileEntry := range fileEntries {
		// base games
		for _, game := range fileEntry.BaseGames {
			title, ok := titles[game.IDPrefix]
			if !ok {
				title = &LibrarySwitchGame{
					DLCs:    make(map[string]LibraryDLCData),
					Updates: make(map[string]LibraryUpdateData),
				}
				titles[game.IDPrefix] = title
			}
			title.LibraryGameData.InLibrary = true
			title.LibraryGameData.Files = append(title.LibraryGameData.Files, LibraryGameDataFile{
				FileID:          fileEntry.FilePath,
				FilePath:        fileEntry.FilePath,
				ReadableVersion: game.ReadableVersion,
				ExtractionType:  fileEntry.ExtractionType,
			})
		}

		// dlcs
		for _, dlc := range fileEntry.DLCs {
			title, ok := titles[dlc.ForIDPrefix]
			if !ok {
				title = &LibrarySwitchGame{
					DLCs:    make(map[string]LibraryDLCData),
					Updates: make(map[string]LibraryUpdateData),
				}
				titles[dlc.ForIDPrefix] = title
			}
			titleDLC, ok := title.DLCs[dlc.ID]
			if !ok {
				titleDLC = LibraryDLCData{}
			}

			titleDLC.InLibrary = true
			titleDLC.Files = append(titleDLC.Files, LibraryDLCDataFile{
				FileID:         fileEntry.FilePath,
				FilePath:       fileEntry.FilePath,
				FileVersion:    dlc.Version,
				ExtractionType: fileEntry.ExtractionType,
			})
			title.DLCs[dlc.ID] = titleDLC
		}

		//updates
		for _, update := range fileEntry.Updates {
			title, ok := titles[update.ForIDPrefix]
			if !ok {
				title = &LibrarySwitchGame{
					DLCs:    make(map[string]LibraryDLCData),
					Updates: make(map[string]LibraryUpdateData),
				}
				titles[update.ForIDPrefix] = title
			}

			titleUpdate, ok := title.Updates[update.ID]
			if !ok {
				titleUpdate = LibraryUpdateData{}
			}
			titleUpdate.Files = append(titleUpdate.Files, LibraryUpdateDataFile{
				FileID:          fileEntry.FilePath,
				FilePath:        fileEntry.FilePath,
				FileVersion:     update.Version,
				ReadableVersion: update.ReadableVersion,
				ExtractionType:  fileEntry.ExtractionType,
			})
			title.Updates[update.ID] = titleUpdate
		}
	}

	// fill catalog details

	errs := make(map[string]error)
	for idPrefix, title := range titles {
		catalogData, err := a.fullDB.GetCatalogEntryByIDPrefix(idPrefix)
		if err != nil {
			errs[idPrefix] = err
			continue
		}
		title.CatalogGameData = CatalogGameData{
			Name:        catalogData.Name,
			TitleID:     catalogData.ID,
			Icon:        catalogData.IconURL,
			Banner:      catalogData.BannerURL,
			Region:      catalogData.Region,
			ReleaseDate: catalogData.ReleaseDate,
			Version:     catalogData.Version,
			Description: catalogData.Description,
			Intro:       catalogData.Intro,
			Publisher:   catalogData.Publisher,
			Screenshots: catalogData.Screenshots,
		}

		for _, dlc := range catalogData.DLCs {
			catalogDLCData := CatalogDLCData{
				Name:        dlc.Name,
				TitleID:     dlc.ID,
				Banner:      dlc.BannerURL,
				Region:      dlc.Region,
				Version:     dlc.Version,
				Description: dlc.Description,
			}

			dlcEntry, ok := title.DLCs[dlc.ID]
			if ok {
				dlcEntry.CatalogDLCData = catalogDLCData
			} else {
				dlcEntry = LibraryDLCData{
					CatalogDLCData: catalogDLCData,
					InLibrary:      false,
				}
			}
			title.DLCs[dlc.ID] = dlcEntry
		}

		// TODO: update versions
		//title.AllVersions =
		//title.IsRecentUpdateInLibrary =
	}

	// TODO: handle errors
	if len(errs) > 0 {
		return nil, fmt.Errorf("errors: %v", errs)
	}

	result := make([]LibrarySwitchGame, 0, len(titles))
	for _, title := range titles {
		result = append(result, *title)
	}
	return result, nil
}

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
