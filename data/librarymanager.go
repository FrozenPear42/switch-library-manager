package data

import (
	"errors"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/keys"
	"github.com/FrozenPear42/switch-library-manager/storage"
	"github.com/FrozenPear42/switch-library-manager/switchfs"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedExtension      = errors.New("unsupported extension")
	ErrFailedToReadFileMetadata  = errors.New("failed to read file metadata")
	ErrFailedToReadTitleID       = errors.New("failed to read title ID")
	ErrFailedToReadTitleVersion  = errors.New("failed to read title version")
	ErrFailedToCalculateChecksum = errors.New("failed to calculate checksum")
)

var (
	versionRegex = regexp.MustCompile(`\[[vV]?(?P<version>[0-9]{1,10})]`)
	titleIdRegex = regexp.MustCompile(`\[(?P<titleId>[A-Z,a-z0-9]{16})]`)
)

type fileInfo struct {
	FullPath string
	Name     string
	Size     int
	Modified int
}

type SwitchFileGame struct {
	IDPrefix string
	ID       string
	Version  int

	Name            map[string]string
	ReadableVersion string
	ISBN            string
}

type SwitchFileDLC struct {
	ForIDPrefix string
	ID          string
	Version     int
}

type SwitchFileUpdate struct {
	ForIDPrefix     string
	ID              string
	Version         int
	ReadableVersion string
}

type LibraryGameFileMetadata struct {
	BaseGames      []SwitchFileGame
	DLCs           []SwitchFileDLC
	Updates        []SwitchFileUpdate
	ExtractionType ExtractionType
	IsMultiContent bool
}

type ExtractionType string

const (
	ExtractionTypeKey      ExtractionType = "key"
	ExtractionTypeFilename ExtractionType = "filename"
)

type LibraryFileEntry struct {
	FilePath string
	FileSize int
	// FileModified is modification date used with FileSize as a "checksum"
	FileModified int
	IsSplit      bool
	*LibraryGameFileMetadata
}

type LibraryManager interface {
	// Rescan performs a scan of library and reports back progress
	Rescan(hardRescan bool, progressCallback ProgressCallback) error
	GetEntries() ([]LibraryFileEntry, error)
	GetFilesForID(id string) ([]LibraryFileEntry, error)
	Clear() error
}

type LibraryManagerImpl struct {
	logger          *zap.SugaredLogger
	db              storage.SwitchDatabaseLibrary
	keysProvider    keys.KeysProvider
	allowedFormats  []string
	scanDirectories []string

	// TODO: replace with persistence
	entries []LibraryFileEntry
}

func NewLibraryManager(logger *zap.SugaredLogger, keysProvider keys.KeysProvider, scanDirectories []string) *LibraryManagerImpl {
	return &LibraryManagerImpl{
		logger:          logger,
		db:              nil,
		keysProvider:    keysProvider,
		allowedFormats:  []string{"xci", "nsp", "nsz", "xcz"},
		scanDirectories: scanDirectories,
		entries:         nil,
	}
}

func (l *LibraryManagerImpl) Rescan(hardRescan bool, progressCallback ProgressCallback) error {
	var recursive bool = true

	var files []fileInfo
	errs := make(map[string]error)

	for dirIdx, path := range l.scanDirectories {
		dirProgress := func(filePath, fileName string) {
			if progressCallback != nil {
				progressCallback(dirIdx, len(l.scanDirectories), filePath)
			}
		}
		l.traverseFolder(path, recursive, dirProgress, &files, errs)
	}
	if len(errs) > 0 {
		l.logger.Warnf("errors: %v", errs)
	}
	// TODO: handle errors

	errs = make(map[string]error)
	fileEntries := make([]LibraryFileEntry, 0, len(files))
	for idx, file := range files {
		if progressCallback != nil {
			progressCallback(idx, len(files), "processing file: "+file.Name)
		}

		fileEntry, err := l.processFile(file)
		if err != nil {
			errs[file.FullPath] = err
			continue
		}
		fileEntries = append(fileEntries, *fileEntry)
	}
	if len(errs) > 0 {
		l.logger.Warnf("errors: %v", errs)
	}

	//TODO:  store to DB instead
	l.entries = fileEntries
	return nil
}

func (l *LibraryManagerImpl) traverseFolder(directory string, recursive bool, progress undeterminedFileProgressCallback, files *[]fileInfo, errs map[string]error) {
	_ = filepath.WalkDir(directory, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			errs[path] = err
			return nil
		}
		if path == directory {
			return nil
		}
		if info.Name()[0] == '.' {
			return nil
		}
		if info.IsDir() {
			if !recursive {
				return filepath.SkipDir
			}
			return nil
		}
		fileExtension := strings.TrimPrefix(filepath.Ext(path), ".")
		if !slices.Contains(l.allowedFormats, fileExtension) {
			errs[path] = ErrUnsupportedExtension
			return nil
		}

		if progress != nil {
			progress(path, info.Name())
		}
		fullInfo, err := info.Info()
		if err != nil {
			return fmt.Errorf("could not get file details: %w", err)
		}
		fullPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("could not get file absolute path: %w", err)
		}
		*files = append(*files, fileInfo{
			FullPath: fullPath,
			Name:     info.Name(),
			Size:     int(fullInfo.Size()),
			Modified: int(fullInfo.ModTime().Unix()),
		})
		return nil
	})
}

func (l *LibraryManagerImpl) processFile(file fileInfo) (*LibraryFileEntry, error) {
	//scan sub-folders if flag is present

	isSplit := false
	// check if it's a split file (ending with xx where x is a digit )
	if partNum, err := strconv.Atoi(file.Name[len(file.Name)-2:]); err == nil {
		if partNum == 0 {
			isSplit = true
		}
	}

	fileExtension := strings.TrimPrefix(filepath.Ext(file.Name), ".")

	if !isSplit && !slices.Contains(l.allowedFormats, fileExtension) {
		return nil, ErrUnsupportedExtension
	}

	gameFileMetadata, err := l.getGameMetadata(file.FullPath, file.Name, fileExtension)
	if err != nil {
		return nil, err
	}

	fileEntry := &LibraryFileEntry{
		FilePath:                file.FullPath,
		FileSize:                file.Size,
		FileModified:            file.Modified,
		IsSplit:                 isSplit,
		LibraryGameFileMetadata: gameFileMetadata,
	}

	return fileEntry, nil
}

func (l *LibraryManagerImpl) getGameMetadata(filePath, fileName, fileFormat string) (*LibraryGameFileMetadata, error) {
	var metadata map[string]*switchfs.ContentMetaAttributes
	var err error
	var extractionType ExtractionType

	_, isKeyAvailable := l.keysProvider.GetProdKey("header_key")
	if !isKeyAvailable {
		extractionType = ExtractionTypeFilename
		res := titleIdRegex.FindStringSubmatch(fileName)
		if len(res) != 2 {
			return nil, ErrFailedToReadTitleID
		}
		titleId := strings.ToLower(res[1])

		res = versionRegex.FindStringSubmatch(fileName)
		if len(res) != 2 {
			return nil, ErrFailedToReadTitleVersion
		}
		version, err := strconv.Atoi(res[1])
		if err != nil {
			return nil, ErrFailedToReadTitleVersion
		}

		metadata = map[string]*switchfs.ContentMetaAttributes{}
		metadata[titleId] = &switchfs.ContentMetaAttributes{TitleId: titleId, Version: version}
	} else {
		extractionType = ExtractionTypeKey
		switch fileFormat {
		case "nsp", "nsz":
			metadata, err = switchfs.ReadNspMetadata(l.keysProvider, filePath)
		case "xci":
			metadata, err = switchfs.ReadXciMetadata(l.keysProvider, filePath)
		case "00":
			metadata, err = switchfs.ReadSplitFileMetadata(l.keysProvider, filePath)
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToReadFileMetadata, err)
		}
	}

	result := &LibraryGameFileMetadata{
		BaseGames:      nil,
		DLCs:           nil,
		Updates:        nil,
		ExtractionType: extractionType,
		IsMultiContent: false,
	}

	entriesCount := 0
	for _, entry := range metadata {
		entryID := strings.ToUpper(entry.TitleId)
		entryPrefix := entryID[:len(entryID)-4]

		if strings.HasSuffix(entryID, "000") {
			// base game
			var name map[string]string
			var isbn string
			var readableVersion string

			if entry.Ncap != nil {
				name = make(map[string]string, len(entry.Ncap.TitleName))
				for lang, title := range entry.Ncap.TitleName {
					if title.Title != "" {
						name[lang] = title.Title
					}
				}
				isbn = entry.Ncap.Isbn
				readableVersion = entry.Ncap.DisplayVersion
			}

			result.BaseGames = append(result.BaseGames, SwitchFileGame{
				IDPrefix:        entryPrefix,
				ID:              entryID,
				Version:         entry.Version,
				Name:            name,
				ReadableVersion: readableVersion,
				ISBN:            isbn,
			})
			entriesCount += 1
		} else if strings.HasSuffix(entry.TitleId, "800") {
			// update
			var readableVersion string
			if entry.Ncap != nil {
				readableVersion = entry.Ncap.DisplayVersion
			}

			result.Updates = append(result.Updates, SwitchFileUpdate{
				ForIDPrefix:     entryPrefix,
				ID:              entryID,
				Version:         entry.Version,
				ReadableVersion: readableVersion,
			})
			entriesCount += 1
		} else {
			// DLC
			result.DLCs = append(result.DLCs, SwitchFileDLC{
				ForIDPrefix: entryPrefix,
				ID:          entryID,
				Version:     entry.Version,
			})
			entriesCount += 1
		}
		if entriesCount > 1 {
			result.IsMultiContent = true
		}
	}

	return result, nil
}

func (l *LibraryManagerImpl) Clear() error {
	return nil
}

func (l *LibraryManagerImpl) GetEntries() ([]LibraryFileEntry, error) {
	return l.entries, nil
}

func (l *LibraryManagerImpl) GetFilesForID(id string) ([]LibraryFileEntry, error) {
	var result []LibraryFileEntry
outer:
	for _, entry := range l.entries {
		for _, t := range entry.BaseGames {
			if t.ID == id {
				result = append(result, entry)
				continue outer
			}
		}
		for _, t := range entry.DLCs {
			if t.ID == id {
				result = append(result, entry)
				continue outer
			}
		}
		for _, t := range entry.Updates {
			if t.ID == id {
				result = append(result, entry)
				continue outer
			}
		}
	}
	return result, nil
}
