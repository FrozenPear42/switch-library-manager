package settings

import (
	"encoding/json"
	"os"
)

type etagData struct {
	TitlesETag   string `json:"titlesETag"`
	VersionsETag string `json:"versionsETag"`
}

type EtagStore struct {
	data   etagData
	loaded bool
}

func (s *EtagStore) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	err = json.NewDecoder(f).Decode(&s.data)
	if err != nil {
		return err
	}
	return nil
}

func (s *EtagStore) Save(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(s.data)
	if err != nil {
		return err
	}
	return nil
}

func (s *EtagStore) GetTitlesETag() (string, bool) {
	return s.data.TitlesETag, s.loaded
}

func (s *EtagStore) GetVersionsETag() (string, bool) {
	return s.data.VersionsETag, s.loaded
}

func (s *EtagStore) UpdateTitlesETag(etag string) {
	s.data.TitlesETag = etag
}

func (s *EtagStore) UpdateVersionsETag(etag string) {
	s.data.VersionsETag = etag
}
