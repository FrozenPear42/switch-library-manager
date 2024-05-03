package main

import (
	"context"
	"encoding/json"
	"github.com/asticode/go-astilectron"
	"github.com/giwty/switch-library-manager/db"
	"github.com/giwty/switch-library-manager/process"
	"github.com/giwty/switch-library-manager/settings"
	"go.uber.org/zap"
	"strings"
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

// App struct
type App struct {
	ctx            context.Context
	sugarLogger    *zap.SugaredLogger
	switchDB       *db.SwitchTitlesDB
	localDB        *db.LocalSwitchFilesDB
	localDbManager *db.LocalSwitchDBManager
}

// NewApp creates a new App application struct
func NewApp(sugarLogger *zap.SugaredLogger, localDbManager *db.LocalSwitchDBManager) *App {
	return &App{
		sugarLogger:    sugarLogger,
		switchDB:       nil,
		localDB:        nil,
		localDbManager: nil,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Rescan(hard bool) {

}

func (a *App) OrganizeLibrary() {
	folderToScan := settings.ReadSettings(g.baseFolder).Folder
	options := settings.ReadSettings(g.baseFolder).OrganizeOptions
	if !process.IsOptionsValid(options) {
		zap.S().Error("the organize options in settings.json are not valid, please check that the template contains file/folder name")
		g.state.window.SendMessage(Message{Name: "error", Payload: "the organize options in settings.json are not valid, please check that the template contains file/folder name"}, func(m *astilectron.EventMessage) {})
		return
	}
	process.OrganizeByFolders(folderToScan, g.state.localDB, g.state.switchDB, g)
	if settings.ReadSettings(g.baseFolder).OrganizeOptions.DeleteOldUpdateFiles {
		process.DeleteOldUpdates(g.baseFolder, g.state.localDB, g)
	}
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
	folderToScan := settings.ReadSettings(g.baseFolder).Folder
	recursiveMode := settings.ReadSettings(g.baseFolder).ScanRecursively

	scanFolders := settings.ReadSettings(g.baseFolder).ScanFolders
	scanFolders = append(scanFolders, folderToScan)
	localDB, err := a.localDbManager.CreateLocalSwitchFilesDB(scanFolders, g, recursiveMode, ignoreCache)
	a.localDB = localDB

	// get ignore ids
	settingsObj := settings.ReadSettings(g.baseFolder)
	ignoreIds := map[string]struct{}{}
	for _, id := range settingsObj.IgnoreDLCTitleIds {
		ignoreIds[strings.ToLower(id)] = struct{}{}
	}

	missingDLC := process.ScanForMissingDLC(a.localDB.TitlesMap, a.switchDB.TitlesMap, ignoreIds)
	missingUpdates := process.ScanForMissingUpdates(a.localDB.TitlesMap, a.switchDB.TitlesMap)

	missingDLCTitles := make([]process.IncompleteTitle, 0, len(missingDLC))
	for _, missingUpdate := range missingDLC {
		missingDLCTitles = append(missingDLCTitles, missingUpdate)
	}

	missingUpdatesTitles := make([]process.IncompleteTitle, len(missingUpdates))
	for _, missingUpdate := range missingUpdates {
		missingUpdatesTitles = append(missingUpdatesTitles, missingUpdate)
	}
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
	progressMessage := ProgressUpdate{curr, total, message}
	a.sugarLogger.Debugf("%v (%v/%v)", message, curr, total)
	msg, err := json.Marshal(progressMessage)
	if err != nil {
		a.sugarLogger.Error(err)
		return
	}
	// TODO: send event
	a.state.window.SendMessage(Message{Name: "updateProgress", Payload: string(msg)}, func(m *astilectron.EventMessage) {})
}
