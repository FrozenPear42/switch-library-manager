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
	Cover       string `json:"cover"`
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
	switchDB           *db.SwitchTitlesDB
	localDB            *db.LocalSwitchFilesDB
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
		switchDB:       nil,
		localDB:        nil,
		localDbManager: nil,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = a.initializeSwitchDB()
}

func (a *App) RequestStartupProgress() {
	runtime.EventsEmit(a.ctx, string(EventTypeStartupProgress), a.recentStartupEvent)
}

func (a *App) Rescan(hard bool) {

}

func (a *App) OrganizeLibrary() {
	//folderToScan := settings.ReadSettings(g.baseFolder).AppDataDirectory
	//options := settings.ReadSettings(g.baseFolder).OrganizeOptions
	//if !process.IsOptionsValid(options) {
	//	zap.S().Error("the organize options in settings.json are not valid, please check that the template contains file/folder name")
	//	g.state.window.SendMessage(Message{Name: "error", Payload: "the organize options in settings.json are not valid, please check that the template contains file/folder name"}, func(m *astilectron.EventMessage) {})
	//	return
	//}
	//process.OrganizeByFolders(folderToScan, g.state.localDB, g.state.switchDB, g)
	//if settings.ReadSettings(g.baseFolder).OrganizeOptions.DeleteOldUpdateFiles {
	//	process.DeleteOldUpdates(g.baseFolder, g.state.localDB, g)
	//}
}

func (a *App) LoadCatalogue() []SwitchTitle {
	var result []SwitchTitle
	for k, v := range a.switchDB.TitlesMap {
		_, isInLibrary := a.localDB.TitlesMap[k]
		if v.Attributes.Name == "" || v.Attributes.Id == "" {
			continue
		}
		result = append(result, SwitchTitle{
			TitleId:     v.Attributes.Id,
			Name:        v.Attributes.Name,
			Cover:       v.Attributes.BannerUrl,
			Icon:        v.Attributes.IconUrl,
			Region:      v.Attributes.Region,
			ReleaseDate: v.Attributes.ReleaseDate,
			Version:     v.Attributes.Version.String(),
			Description: v.Attributes.Description,
			Publisher:   v.Attributes.Publisher,
			InLibrary:   isInLibrary,
		})
	}
	return result
}

func (a *App) LoadLocalGames() {
	//folderToScan := settings.ReadSettings(g.baseFolder).AppDataDirectory
	//recursiveMode := settings.ReadSettings(g.baseFolder).ScanRecursively
	//
	//scanFolders := settings.ReadSettings(g.baseFolder).ScanDirectories
	//scanFolders = append(scanFolders, folderToScan)
	//localDB, err := a.localDbManager.CreateLocalSwitchFilesDB(scanFolders, g, recursiveMode, ignoreCache)
	//a.localDB = localDB
	//
	//// get ignore ids
	//settingsObj := settings.ReadSettings(g.baseFolder)
	//ignoreIds := map[string]struct{}{}
	//for _, id := range settingsObj.IgnoreDLCTitleIDs {
	//	ignoreIds[strings.ToLower(id)] = struct{}{}
	//}
	//
	//missingDLC := process.ScanForMissingDLC(a.localDB.TitlesMap, a.switchDB.TitlesMap, ignoreIds)
	//missingUpdates := process.ScanForMissingUpdates(a.localDB.TitlesMap, a.switchDB.TitlesMap)
	//
	//missingDLCTitles := make([]process.IncompleteTitle, 0, len(missingDLC))
	//for _, missingUpdate := range missingDLC {
	//	missingDLCTitles = append(missingDLCTitles, missingUpdate)
	//}
	//
	//missingUpdatesTitles := make([]process.IncompleteTitle, len(missingUpdates))
	//for _, missingUpdate := range missingUpdates {
	//	missingUpdatesTitles = append(missingUpdatesTitles, missingUpdate)
	//}
}

func (a *App) CheckUpdate() (string, bool, error) {
	recentVersion, isUpdateAvailable, err := settings.CheckForUpdates()
	if err != nil {
		a.sugarLogger.Error(err)
		return "", false, err
	}

	return recentVersion, isUpdateAvailable, nil
}

func (a *App) updateProgress(curr int, total int, message string) {
	//progressMessage := ProgressUpdate{curr, total, message}
	//a.sugarLogger.Debugf("%v (%v/%v)", message, curr, total)
	//msg, err := json.Marshal(progressMessage)
	//if err != nil {
	//	a.sugarLogger.Error(err)
	//	return
	//}
	//// TODO: send event
	//a.state.window.SendMessage(Message{Name: "updateProgress", Payload: string(msg)}, func(m *astilectron.EventMessage) {})
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
