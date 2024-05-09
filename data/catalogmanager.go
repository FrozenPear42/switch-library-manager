package data

type CatalogManager interface {
	RebuildCatalog(hard bool) error
}
