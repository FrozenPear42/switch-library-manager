package nut

import (
	"encoding/json"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func NewRouter(db data.LibraryManager, reporter ProgressReporter) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.NotFound(HandleNotFound())
	router.Route("/api", func(r chi.Router) {
		r.Get("/search", HandleGetSearch(db))
		r.Get("/download/{titleId}/{fileName}", HandleGetDownload(db, reporter, false))
		r.Get("/download/{titleId}/{fileName}/{start}", HandleGetDownload(db, reporter, false))
		r.Get("/download/{titleId}/{fileName}/{start}/{stop}", HandleGetDownload(db, reporter, false))
		r.Head("/download/{titleId}/{fileName}", HandleGetDownload(db, reporter, true))
		r.Head("/download/{titleId}/{fileName}/{start}", HandleGetDownload(db, reporter, true))
		r.Head("/download/{titleId}/{fileName}/{start}/{stop}", HandleGetDownload(db, reporter, true))

		// those are not used by tinfoil so we skip implementation for now
		//r.Get("/user", HandleGetUser)
		//r.Get("/scan", HandleGetUser)
		//r.Get("/titles", HandleGetUser)
		//r.Get("/titleImage", HandleGetUser)
		//r.Get("/bannerImage", HandleGetUser)
		//r.Get("/frontArtBoxImage", HandleGetUser)
		//r.Get("/screenshotImage", HandleGetUser)
		//r.Get("/preload", HandleGetUser)
		//r.Get("/install", HandleGetUser)
		//r.Get("/offsetAndSize", HandleGetUser)
		//r.Get("/directoryList", HandleGetUser)
		//r.Get("/file", HandleGetUser)
		//r.Head("/file", HandleGetUser)
		//r.Get("/fileSize", HandleGetUser)
		//r.Get("/titleUpdates", HandleGetUser)
		//r.Get("/organize", HandleGetUser)
		//r.Get("/updateDb", HandleGetUser)
		//r.Get("/export", HandleGetUser)
		//r.Get("/importRegions", HandleGetUser)
		//r.Get("/regions", HandleGetUser)
		//r.Get("/updateLatest", HandleGetUser)
		//r.Get("/updateAllVersions", HandleGetUser)
		//r.Get("/scrapeShogun", HandleGetUser)
		//r.Get("/submitKey", HandleGetUser)
		//r.Post("/tinfoilSetInstalledApps", HandleGetUser)
		//r.Get("/switchList", HandleGetUser)
		//r.Get("/switchInstalled", HandleGetUser)
	})

	return router
}

type searchResultDTO []searchResultEntryDTO

type searchResultEntryDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Size    int    `json:"size"`
	Version int    `json:"version"`
}

func HandleNotFound() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.NotFound(writer, request)
	}
}

func HandleGetSearch(db data.LibraryManager) http.HandlerFunc {
	logger := zap.S()

	return func(writer http.ResponseWriter, request *http.Request) {

		entries, err := db.GetEntries()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// avoid duplicates - maps id to a file entry
		fileMap := make(map[string]searchResultEntryDTO)

		for _, entry := range entries {
			// TODO: check if its supported, its just assumption it's not looking at NUT source code
			if entry.IsMultiContent {
				logger.Warnf("skipping file, multicontent not supported by NUT: %v", entry.FilePath)
				continue
			}
			if entry.IsSplit {
				logger.Warnf("skipping file, split files not supported by NUT: %v", entry.FilePath)
				continue
			}

			var id string
			var version int
			if len(entry.DLCs) > 0 {
				id = entry.DLCs[0].ID
				version = entry.DLCs[0].Version
			} else if len(entry.Updates) > 0 {
				id = entry.Updates[0].ID
				version = entry.Updates[0].Version
			} else if len(entry.BaseGames) > 0 {
				id = entry.BaseGames[0].ID
				version = entry.BaseGames[0].Version
			}
			id = strings.ToUpper(id)
			fileName := filepath.Base(entry.FilePath)
			if oldEntry, ok := fileMap[id]; !ok || oldEntry.Version < version {
				fileMap[id] = searchResultEntryDTO{
					ID:      id,
					Name:    fileName,
					Size:    entry.FileSize,
					Version: version,
				}
			}
		}

		dto := make(searchResultDTO, 0, len(fileMap))
		for _, v := range fileMap {
			dto = append(dto, v)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err = json.NewEncoder(writer).Encode(dto)
		if err != nil {
			logger.Errorf("failed to write json response: %v", err)
		}
	}
}

func HandleGetDownload(db data.LibraryManager, progress ProgressReporter, isHead bool) http.HandlerFunc {
	logger := zap.S()

	return func(writer http.ResponseWriter, request *http.Request) {
		var err error
		titleID := chi.URLParam(request, "titleId")
		fileName := chi.URLParam(request, "fileName")

		fileName, err = url.PathUnescape(fileName)
		if err != nil {
			http.Error(writer, fmt.Sprintf("invalid file format: %v", err), http.StatusBadRequest)
			return
		}

		files, err := db.GetFilesForID(titleID)
		if err != nil {
			http.Error(writer, fmt.Sprintf("could not get file list: %v", err), http.StatusInternalServerError)
			return
		}
		if len(files) == 0 {
			http.Error(writer, fmt.Sprintf("could not find specified file %v", fileName), http.StatusBadRequest)
			return
		}

		var filePath string
		idx := slices.IndexFunc(files, func(entry data.LibraryFileEntry) bool {
			fn := filepath.Base(entry.FilePath)
			return fn == fileName
		})
		if idx != -1 {
			filePath = files[idx].FilePath
		} else {
			http.Error(writer, fmt.Sprintf("could not find specified file %v", fileName), http.StatusBadRequest)
			return
			// TODO: validate if its expected behaviour
			//logger.Errorf("could not find specified file (%v), using first available (%v)", fileName, filePath)
			//filePath = files[0].FilePath
		}

		f, err := os.Open(filePath)
		if err != nil {
			http.Error(writer, fmt.Sprintf("could not open the file %v", fileName), http.StatusBadRequest)
			return
		}

		stat, err := f.Stat()
		if err != nil {
			http.Error(writer, fmt.Sprintf("could not stat the file %v", fileName), http.StatusBadRequest)
			return
		}
		fileSize := stat.Size()

		chunkSize := 0x400000
		start := int64(0)
		stop := fileSize
		if rangeHeader := request.Header.Get("Range"); rangeHeader != "" {
			ranges, err := ParseRange(rangeHeader, fileSize)
			if err != nil || len(ranges) > 1 {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(ranges) == 1 {
				start = ranges[0].Start
				stop = ranges[0].Start + ranges[0].Length
			}
		} else {
			if s := chi.URLParam(request, "start"); s != "" {
				start, err = strconv.ParseInt(s, 10, 64)
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			if s := chi.URLParam(request, "stop"); s != "" {
				stop, err = strconv.ParseInt(s, 10, 64)
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}

		logger.Debugf("serving file %v, in chunk %v-%v")
		_, err = f.Seek(int64(start), io.SeekStart)

		contentFileName := titleID + filepath.Ext(filePath)

		writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", contentFileName))
		writer.Header().Set("Content-type", "application/octet-stream")
		writer.Header().Set("Accept-Ranges", "bytes")
		writer.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, stop-1, fileSize))
		writer.Header().Set("Content-Length", strconv.FormatInt(stop-start, 10))
		writer.WriteHeader(http.StatusPartialContent)

		if isHead {
			return
		}

		totalWritten := int64(0)
		toWrite := stop - start
		for {
			if toWrite-totalWritten == 0 {
				break
			}
			n, err := io.CopyN(writer, f, slices.Min([]int64{int64(chunkSize), toWrite - totalWritten}))
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Errorf("error while writing file: %v", err)
				break
			}
			totalWritten += n
			progress.ReportProgress(filePath, totalWritten, toWrite)
		}

	}
}
