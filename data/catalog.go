package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/storage"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type titlesJsonEntry struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Version     json.Number `json:"version"`
	Region      string      `json:"region"`
	ReleaseDate int         `json:"releaseDate"`
	Developer   string      `json:"developer"`
	Publisher   string      `json:"publisher"`
	IconURL     string      `json:"iconUrl"`
	BannerURL   string      `json:"bannerUrl"`
	Screenshots []string    `json:"screenshots"`
	Description string      `json:"description"`
	Size        int         `json:"size"`
	Category    []string    `json:"category"`
	FrontBoxArt string      `json:"frontBoxArt"`
	Intro       string      `json:"intro"`
	Key         string      `json:"key"`
	IsDemo      bool        `json:"isDemo"`
	Language    string      `json:"language"`
	Languages   []string    `json:"languages"`
}

type titlesJson map[string]titlesJsonEntry

type versionsJsonEntry map[string]string

type versionsJson map[string]versionsJsonEntry

type ProgressCallback func(current, total int, message string)

const (
	dialTimeout = 3 * time.Second
)

var (
	ErrNoUpdateAvailable = errors.New("no update available")
)

func BuildCatalog(db storage.SwitchDatabaseCatalog, titlesURL, versionsURL string, callback ProgressCallback) error {
	// TODO: all or nothing - if only one download fails redownload the other one

	totalSteps := 7

	callback(1, totalSteps, "Preparing")
	metadata, err := db.GetCatalogMetadata()
	if err != nil {
		return fmt.Errorf("failed to fetch catalog metadata: %w", err)
	}

	// TODO: temp measure
	if metadata.TitlesETag != "" || metadata.VersionsETag != "" {
		callback(2, 2, "Done!")
		return nil
	}

	tmpDir, err := os.MkdirTemp("", "slm")
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}

	zap.S().Infof("TMP dir: %v", tmpDir)

	callback(2, totalSteps, "Downloading titles data...")
	titlesPath := filepath.Join(tmpDir, "titles.json")
	titlesEtag := metadata.TitlesETag

	titlesFile, newTitlesEtag, err := downloadFileWithEtag(titlesURL, titlesPath, titlesEtag)
	if err != nil && !errors.Is(err, ErrNoUpdateAvailable) {
		return fmt.Errorf("failed to download switch titles: %w", err)
	}
	defer func() {
		if titlesFile != nil {
			titlesFile.Close()
			os.Remove(titlesPath)
		}
	}()

	callback(3, totalSteps, "Downloading versions data...")
	versionsPath := filepath.Join(tmpDir, "versions.json")
	versionsEtag := metadata.VersionsETag

	versionsFile, newVersionsEtag, err := downloadFileWithEtag(versionsURL, versionsPath, versionsEtag)
	if err != nil && !errors.Is(err, ErrNoUpdateAvailable) {
		return fmt.Errorf("failed to download switch titles: %w", err)
	}
	defer func() {
		if versionsFile != nil {
			versionsFile.Close()
			os.Remove(versionsPath)
		}
	}()

	callback(4, totalSteps, "Processing data...")
	entries, err := parseCatalogFiles(titlesFile, versionsFile)
	if err != nil {
		return fmt.Errorf("failed to parse catalog data: %w", err)
	}
	zap.S().Infof("Processed data: %v", tmpDir)

	callback(5, totalSteps, "Updating local DB...")
	err = db.ClearCatalog()
	if err != nil {
		return fmt.Errorf("failed to clear catalog: %w", err)
	}
	err = db.AddCatalogEntries(entries)
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}
	callback(6, totalSteps, "Finishing up...")
	err = db.UpdateCatalogMetadata(storage.CatalogMetadata{
		VersionsETag: newVersionsEtag,
		TitlesETag:   newTitlesEtag,
	})
	if err != nil {
		return fmt.Errorf("failed to update catalog metadata: %w", err)
	}
	callback(7, totalSteps, "Done...")
	return nil
}

func parseCatalogFiles(titlesFile, versionsFile io.Reader) (map[string]storage.CatalogEntry, error) {
	// assuming titles file is sorted
	var versionsData versionsJson
	err := json.NewDecoder(versionsFile).Decode(&versionsData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	var titlesData titlesJson
	err = json.NewDecoder(titlesFile).Decode(&titlesData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode titles: %w", err)
	}

	entries := make(map[string]storage.CatalogEntry)

	for id, data := range titlesData {
		id = strings.ToUpper(id)
		// parse id
		// Main       ends with 000
		// Updates    ends with 800
		// Dlc        a running counter (starting with 001) in the 4 last chars

		mainTitleId := id[:len(id)-4] + "xxxx"

		if _, ok := entries[mainTitleId]; !ok {
			entries[mainTitleId] = storage.CatalogEntry{}
		}

		entry := entries[mainTitleId]

		if strings.HasSuffix(id, "000") {
			// main title
			entry.CatalogEntryData = storage.CatalogEntryData{
				ID:          data.ID,
				Name:        data.Name,
				Version:     data.Version.String(),
				BannerURL:   data.BannerURL,
				IconURL:     data.IconURL,
				Description: data.Description,
				Intro:       data.Intro,
				Region:      data.Region,
				Key:         data.Key,
				ReleaseDate: parseReleaseDate(data.ReleaseDate),
				Publisher:   data.Publisher,
				IsDemo:      data.IsDemo,
				Screenshots: data.Screenshots,
			}
			var vs []storage.CatalogEntryVersion
			for versionNumberStr, releaseDate := range versionsData[id[:len(id)-3]+"000"] {
				versionNumber, err := strconv.Atoi(versionNumberStr)
				if err != nil {
					continue
				}
				vs = append(vs, storage.CatalogEntryVersion{
					Version:     versionNumber,
					ReleaseDate: releaseDate,
				})
			}
			slices.SortStableFunc(vs, func(a, b storage.CatalogEntryVersion) int {
				return a.Version - b.Version
			})
			entry.Versions = vs
		} else if strings.HasSuffix(id, "800") {
			// update
			v, err := data.Version.Int64()
			if err != nil {
				continue
			}
			entry.RecentUpdate = storage.CatalogEntryRecentUpdate{
				ID:      id,
				Version: int(v),
				Key:     data.Key,
			}
		} else {
			// dlc
			entry.DLCs = append(entry.DLCs, storage.CatalogEntryDLC{
				CatalogEntryData: storage.CatalogEntryData{
					ID:          data.ID,
					Name:        data.Name,
					Version:     data.Version.String(),
					BannerURL:   data.BannerURL,
					IconURL:     data.IconURL,
					Description: data.Description,
					Intro:       data.Intro,
					Region:      data.Region,
					Key:         data.Key,
					ReleaseDate: parseReleaseDate(data.ReleaseDate),
					Publisher:   data.Publisher,
					Screenshots: data.Screenshots,
				},
			})
		}
		entries[mainTitleId] = entry
	}
	return entries, nil
}

// downloadFileWithEtag  downloads a file from a given url using etag header, returns (file, new etag, error). Returned file has to ble closed by caller.
func downloadFileWithEtag(url string, path string, etag string) (*os.File, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("If-None-Match", etag)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: dialTimeout,
			}).DialContext,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	newEtag := resp.Header.Get("Etag")

	if resp.StatusCode != http.StatusOK {
		// FIXME: magic, fix it
		if resp.StatusCode < 400 {
			return nil, "", ErrNoUpdateAvailable
		}
		return nil, "", fmt.Errorf("got a non 200 response - %v", resp.Status)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return nil, "", fmt.Errorf("could not open file: %w", err)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("could not download (copy): %w", err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, "", fmt.Errorf("could not download (seek): %w", err)
	}

	return f, newEtag, nil
}

func parseReleaseDate(date int) string {
	year := date / 1_00_00
	yearR := date % 1_00_00
	month := yearR / 1_00
	day := yearR % 1_00

	if year == 0 || month == 0 || day == 0 {
		return ""
	}
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
