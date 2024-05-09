package old

//
//import (
//	"encoding/json"
//	"github.com/FrozenPear42/switch-library-manager/db"
//	"path/filepath"
//	"strconv"
//)
//
//type Pair struct {
//	Key   string `json:"key"`
//	Value string `json:"value"`
//}
//
//type LocalLibraryData struct {
//	LibraryData []LibraryTemplateData `json:"library_data"`
//	Issues      []Pair                `json:"issues"`
//	NumFiles    int                   `json:"num_files"`
//}
//
//type SwitchTitle struct {
//	Name        string `json:"name"`
//	TitleId     string `json:"titleId"`
//	Icon        string `json:"icon"`
//	Region      string `json:"region"`
//	ReleaseDate int    `json:"release_date"`
//}
//
//type LibraryTemplateData struct {
//	Id      int    `json:"id"`
//	Name    string `json:"name"`
//	Version string `json:"version"`
//	Dlc     string `json:"dlc"`
//	TitleId string `json:"titleId"`
//	Path    string `json:"path"`
//	Icon    string `json:"icon"`
//	Update  int    `json:"update"`
//	Region  string `json:"region"`
//	Type    string `json:"type"`
//}
//
//type ProgressUpdate struct {
//	Curr    int    `json:"curr"`
//	Total   int    `json:"total"`
//	Message string `json:"message"`
//}
//
//func (g *GUI) handleMessage(m *astilectron.EventMessage) interface{} {
//	var retValue string
//	g.state.Lock()
//	defer g.state.Unlock()
//	msg := Message{}
//	err := m.Unmarshal(&msg)
//
//	if err != nil {
//		g.sugarLogger.Error("Failed to parse client message", err)
//		return ""
//	}
//
//	g.sugarLogger.Debugf("Received message from client [%v]", msg)
//
//	switch msg.Name {
//	case "updateLocalLibrary":
//		ignoreCache, _ := strconv.ParseBool(msg.Payload)
//		localDB, err := g.buildLocalDB(g.localDbManager, ignoreCache)
//		if err != nil {
//			g.sugarLogger.Error(err)
//			g.state.window.SendMessage(Message{Name: "error", Payload: err.Error()}, func(m *astilectron.EventMessage) {})
//			return ""
//		}
//		response := LocalLibraryData{}
//		libraryData := []LibraryTemplateData{}
//		issues := []Pair{}
//		for k, v := range localDB.TitlesMap {
//			if v.BaseExist {
//				version := ""
//				name := ""
//				if v.File.Metadata.Ncap != nil {
//					version = v.File.Metadata.Ncap.DisplayVersion
//					name = v.File.Metadata.Ncap.TitleName["AmericanEnglish"].Title
//				}
//
//				if v.Updates != nil && len(v.Updates) != 0 {
//					if v.Updates[v.LatestUpdate].Metadata.Ncap != nil {
//						version = v.Updates[v.LatestUpdate].Metadata.Ncap.DisplayVersion
//					} else {
//						version = ""
//					}
//				}
//				if title, ok := g.state.switchDB.TitlesMap[k]; ok {
//					if title.Attributes.Name != "" {
//						name = title.Attributes.Name
//					}
//					libraryData = append(libraryData,
//						LibraryTemplateData{
//							Icon:    title.Attributes.IconUrl,
//							Name:    name,
//							TitleId: v.File.Metadata.TitleId,
//							Update:  v.LatestUpdate,
//							Version: version,
//							Region:  title.Attributes.Region,
//							Type:    getType(v),
//							Path:    filepath.Join(v.File.ExtendedInfo.BaseFolder, v.File.ExtendedInfo.FileName),
//						})
//				} else {
//					if name == "" {
//						name = db.ParseTitleNameFromFileName(v.File.ExtendedInfo.FileName)
//					}
//					libraryData = append(libraryData,
//						LibraryTemplateData{
//							Name:    name,
//							Update:  v.LatestUpdate,
//							Version: version,
//							Type:    getType(v),
//							TitleId: v.File.Metadata.TitleId,
//							Path:    v.File.ExtendedInfo.FileName,
//						})
//				}
//
//			} else {
//				for _, update := range v.Updates {
//					issues = append(issues, Pair{Key: filepath.Join(update.ExtendedInfo.BaseFolder, update.ExtendedInfo.FileName), Value: "base file is missing"})
//				}
//				for _, dlc := range v.Dlc {
//					issues = append(issues, Pair{Key: filepath.Join(dlc.ExtendedInfo.BaseFolder, dlc.ExtendedInfo.FileName), Value: "base file is missing"})
//				}
//			}
//		}
//		for k, v := range localDB.Skipped {
//			issues = append(issues, Pair{Key: filepath.Join(k.BaseFolder, k.FileName), Value: v.ReasonText})
//		}
//
//		response.LibraryData = libraryData
//		response.NumFiles = localDB.NumFiles
//		response.Issues = issues
//		msg, _ := json.Marshal(response)
//		g.state.window.SendMessage(Message{Name: "libraryLoaded", Payload: string(msg)}, func(m *astilectron.EventMessage) {})
//	case "updateDB":
//		if g.state.switchDB == nil {
//			switchDb, err := g.buildSwitchDb()
//			if err != nil {
//				g.sugarLogger.Error(err)
//				g.state.window.SendMessage(Message{Name: "error", Payload: err.Error()}, func(m *astilectron.EventMessage) {})
//				return ""
//			}
//			g.state.switchDB = switchDb
//		}
//	}
//	g.sugarLogger.Debugf("Server response [%v]", retValue)
//	return retValue
//}
//
//func getType(gameFile *db.SwitchGameFiles) string {
//	if gameFile.IsSplit {
//		return "split"
//	}
//	if gameFile.MultiContent {
//		return "multi-content"
//	}
//	ext := filepath.Ext(gameFile.File.ExtendedInfo.FileName)
//	if len(ext) > 1 {
//		return ext[1:]
//	}
//	return ""
//}
