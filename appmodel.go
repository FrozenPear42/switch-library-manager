package main

import (
	"github.com/FrozenPear42/switch-library-manager/data"
	"github.com/FrozenPear42/switch-library-manager/storage"
)

type OrganizeSettings struct {
}

// Catalog

type CatalogGameData struct {
	Name        string   `json:"name"`
	TitleID     string   `json:"titleID"`
	Icon        string   `json:"icon"`
	Banner      string   `json:"banner"`
	Region      string   `json:"region"`
	ReleaseDate string   `json:"releaseDate"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Intro       string   `json:"intro"`
	Publisher   string   `json:"publisher"`
	Screenshots []string `json:"screenshots"`
}

type CatalogDLCData struct {
	Name        string `json:"name"`
	TitleID     string `json:"titleID"`
	Banner      string `json:"banner"`
	Region      string `json:"region"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type CatalogVersionData struct {
	Version     int    `json:"version"`
	ReleaseDate string `json:"releaseDate"`
}

type CatalogSwitchGame struct {
	CatalogGameData
	DLCs     []CatalogDLCData     `json:"dlcs"`
	Versions []CatalogVersionData `json:"versions"`
}

type CatalogPage struct {
	Games      []CatalogSwitchGame `json:"games"`
	TotalGames int                 `json:"totalTitles"`
	NextCursor int                 `json:"nextCursor"`
	IsLastPage bool                `json:"isLastPage"`
}

type CatalogFilters struct {
	SortBy storage.CatalogFiltersSortBy `json:"sortBy"`
	Name   *string                      `json:"name"`
	ID     *string                      `json:"id"`
	Region []string                     `json:"region"`
	Cursor int                          `json:"cursor"`
	Limit  int                          `json:"limit"`
}

// Library

type LibraryFileEntry struct {
	FileID   string `json:"fileID"`
	FilePath string `json:"filePath"`
	FileSize int    `json:"fileSize"`
}

type LibraryGameData struct {
	CatalogGameData
	InLibrary bool                  `json:"inLibrary"`
	Files     []LibraryGameDataFile `json:"files"`
}

type LibraryGameDataFile struct {
	FileID          string              `json:"fileID"`
	FilePath        string              `json:"filePath"`
	ReadableVersion string              `json:"readableVersion"`
	ExtractionType  data.ExtractionType `json:"extractionType"`
}

type LibraryDLCData struct {
	CatalogDLCData
	InLibrary bool                 `json:"inLibrary"`
	Files     []LibraryDLCDataFile `json:"files"`
}

type LibraryDLCDataFile struct {
	FileID         string              `json:"fileID"`
	FilePath       string              `json:"filePath"`
	FileVersion    int                 `json:"fileVersion"`
	ExtractionType data.ExtractionType `json:"extractionType"`
}

type LibraryUpdateData struct {
	Files []LibraryUpdateDataFile `json:"files"`
}

type LibraryUpdateDataFile struct {
	FileID          string              `json:"fileID"`
	FilePath        string              `json:"filePath"`
	FileVersion     int                 `json:"fileVersion"`
	ReadableVersion string              `json:"readableVersion"`
	ExtractionType  data.ExtractionType `json:"extractionType"`
}

type LibrarySwitchGame struct {
	LibraryGameData
	DLCs                    map[string]LibraryDLCData    `json:"dlcs"`
	Updates                 map[string]LibraryUpdateData `json:"updates"`
	AllVersions             []CatalogVersionData         `json:"allVersions"`
	IsRecentUpdateInLibrary bool                         `json:"isRecentUpdateInLibrary"`
}

// Events

type EventType string

const (
	EventTypeStartupProgress EventType = "startupProgress"
)

type EventMessagePayload interface {
	__eventMessagePayload()
}
type _eventMessagePayload struct{}

func (_ _eventMessagePayload) __eventMessagePayload() {}

type EventMessage struct {
	Type string              `json:"type"`
	Data EventMessagePayload `json:"data"`
}

type EventStartupProgressPayload struct {
	_eventMessagePayload
	Completed bool   `json:"completed"`
	Running   bool   `json:"running"`
	Message   string `json:"message"`
	Current   int    `json:"current"`
	Total     int    `json:"total"`
}
