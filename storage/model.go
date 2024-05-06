package storage

type CatalogMetadata struct {
	VersionsETag string
	TitlesETag   string
}

type CatalogEntryData struct {
	ID          string
	Name        string
	Version     string
	BannerURL   string
	IconURL     string
	Description string
	Intro       string
	Region      string
	Key         string
	ReleaseDate string
	Publisher   string
	IsDemo      bool
	Screenshots []string
}

type CatalogEntryRecentUpdate struct {
	ID      string
	Version int
	Key     string
}

type CatalogEntryVersion struct {
	Version int
	// ReleaseDate is release date in ISO format
	ReleaseDate string
}

type CatalogEntryDLC struct {
	CatalogEntryData
}

type CatalogEntry struct {
	CatalogEntryData
	RecentUpdate CatalogEntryRecentUpdate
	Versions     []CatalogEntryVersion
	DLCs         []CatalogEntryDLC
}
