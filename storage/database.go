package storage

import (
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

type SwitchDatabaseLocal interface {
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
	path string
	db   *bolthold.Store
}

func NewDatabase(path string) (*Database, error) {
	db, err := bolthold.Open(path, 0644, &bolthold.Options{Options: &bbolt.Options{
		NoSync: true,
	}})
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	return &Database{
		db: db,
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
	q := &bolthold.Query{}
	q = q.Skip(cursor).Limit(pageSize)

	if filters != nil {
		switch filters.SortBy {
		case CatalogFiltersSortByID:
			q.SortBy("ID")
		case CatalogFiltersSortByName:
			q.SortBy("Name")
		}
	}

	var entries []CatalogEntry
	err := d.db.Find(&entries, q)
	if err != nil {
		return Page[CatalogEntry]{}, err
	}

	count, err := d.db.Count(&CatalogEntry{}, nil)
	if err != nil {
		return Page[CatalogEntry]{}, err
	}

	nextCursor := cursor + pageSize + 1
	return Page[CatalogEntry]{
		Data:       entries,
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
