package storage

import (
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"golang.org/x/exp/slices"
	"strings"
	"sync"
)

type SwitchDatabaseLibrary interface {
}

type SwitchDatabaseCatalog interface {
	GetCatalogMetadata() (CatalogMetadata, error)
	UpdateCatalogMetadata(data CatalogMetadata) error

	AddCatalogEntries(entries map[string]CatalogEntry) error
	GetCatalogEntryByID(id string) (CatalogEntry, bool, error)
	GetCatalogEntries(filters *CatalogFilters, pageSize int, cursor int) (Page[CatalogEntry], error)
	ClearCatalog() error
}

type SwitchDatabase interface {
	SwitchDatabaseCatalog
}

type Database struct {
	path   string
	db     *bolthold.Store
	data   []CatalogEntry
	loaded bool
	mutex  sync.Mutex
}

func NewDatabase(path string) (*Database, error) {
	db, err := bolthold.Open(path, 0644, &bolthold.Options{Options: &bbolt.Options{
		NoSync: true,
	}})
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	return &Database{
		path:   path,
		db:     db,
		data:   make([]CatalogEntry, 0),
		loaded: false,
		mutex:  sync.Mutex{},
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetCatalogMetadata() (CatalogMetadata, error) {
	var metadata CatalogMetadata
	err := d.db.Get("metadata", &metadata)
	if err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return CatalogMetadata{
				VersionsETag: "",
				TitlesETag:   "",
			}, nil
		}

		return metadata, err
	}
	return metadata, nil
}

func (d *Database) UpdateCatalogMetadata(metadata CatalogMetadata) error {
	err := d.db.Upsert("metadata", metadata)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) AddCatalogEntries(entries map[string]CatalogEntry) error {
	errs := make(map[string]error)
	for key, entry := range entries {
		err := d.db.Upsert(key, entry)
		if err != nil {
			errs[key] = err
			continue
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("could not add some entries: %v", errs)
	}
	return nil
}

func (d *Database) GetCatalogEntryByID(id string) (CatalogEntry, bool, error) {
	panic("implement me")
}

func (d *Database) GetCatalogEntries(filters *CatalogFilters, pageSize int, cursor int) (Page[CatalogEntry], error) {
	// FIXME: bolthold is terrible, it basically loads entire database into memory then filters and limits it.
	// for now let's keep it as in memory DB
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.loaded {
		err := d.db.Find(&d.data, nil)
		if err != nil {
			return Page[CatalogEntry]{}, err
		}
		d.loaded = true
	}

	if filters != nil {
		switch filters.SortBy {
		case CatalogFiltersSortByID:
			slices.SortFunc(d.data, func(a, b CatalogEntry) int {
				return strings.Compare(a.ID, b.ID)
			})
		case CatalogFiltersSortByName:
			slices.SortFunc(d.data, func(a, b CatalogEntry) int {
				return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
			})
		}
	}

	var entries []CatalogEntry

	if filters != nil {
		for _, entry := range d.data {
			valid := true
			if filters.Name != nil {
				valid = valid && strings.Contains(strings.ToLower(entry.Name), strings.ToLower(*filters.Name))
			}
			if filters.ID != nil {
				valid = valid && strings.HasPrefix(strings.ToLower(entry.ID), strings.ToLower(*filters.ID))
			}
			if len(filters.Region) > 0 {
				valid = valid && slices.ContainsFunc(filters.Region, func(s string) bool {
					return strings.ToLower(s) == strings.ToLower(entry.Region)
				})
			}

			if valid {
				entries = append(entries, entry)
			}
		}
	} else {
		entries = d.data
	}

	count := len(entries)

	var data []CatalogEntry
	if cursor > len(entries) {
		data = []CatalogEntry{}
	} else {
		data = entries[max(0, cursor):min(cursor+pageSize, len(entries))]
	}
	nextCursor := cursor + pageSize + 1

	return Page[CatalogEntry]{
		Data:       data,
		NextCursor: nextCursor,
		TotalCount: count,
		IsLastPage: nextCursor > count,
	}, nil
}

func (d *Database) ClearCatalog() error {
	err := d.db.DeleteMatching(&CatalogEntry{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
