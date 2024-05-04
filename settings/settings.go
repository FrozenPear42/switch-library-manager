package settings

import (
	"errors"
	"fmt"
	"github.com/FrozenPear42/switch-library-manager/process"
	"github.com/FrozenPear42/switch-library-manager/utils"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type OrganizeOptions struct {
	CreateFolderPerGame  bool   `yaml:"createFolderPerGame" default:"false"`
	RenameFiles          bool   `yaml:"renameFiles" default:"false"`
	DeleteEmptyFolders   bool   `yaml:"deleteEmptyFolders" default:"false"`
	DeleteOldUpdateFiles bool   `yaml:"deleteOldUpdateFiles" default:"false"`
	SwitchSafeFileNames  bool   `yaml:"switchSafeFileNames" default:"true"`
	FolderNameTemplate   string `yaml:"folderNameTemplate" default:"-"`
	FileNameTemplate     string `yaml:"fileNameTemplate" default:"-"`
}

func (o *OrganizeOptions) SetDefaults() {
	if defaults.CanUpdate(o.FileNameTemplate) {
		o.FileNameTemplate = fmt.Sprintf("{%v} ({%v})[{%v}][v{%v}]",
			process.TemplateTokenTitleName, process.TemplateTokenDLCName,
			process.TemplateTokenTitleID, process.TemplateTokenVersion)
	}
	if defaults.CanUpdate(o.FolderNameTemplate) {
		o.FolderNameTemplate = fmt.Sprintf("{%v}", process.TemplateTokenTitleName)
	}
}

type AppSettings struct {
	Debug             bool            `yaml:"debug" default:"false"`
	IgnoreDLCTitleIDs []string        `yaml:"ignoreDLCTitleIDs" default:"[\"test\"]"`
	ProdKeysPath      string          `yaml:"prodKeysPath" default:"-"`
	AppDataDirectory  string          `yaml:"appDataDirectory" default:"-"`
	ScanDirectories   []string        `yaml:"scanDirectories" default:"[]"`
	ScanRecursive     bool            `yaml:"scanRecursive" default:"true"`
	TitlesFileName    string          `yaml:"titlesFileName" default:"titles.json"`
	VersionsFileName  string          `yaml:"versionsFileName" default:"versions.json"`
	TitlesEndpoint    string          `yaml:"titlesEndpoint" default:"https://tinfoil.media/repo/db/titles.json"`
	VersionsEndpoint  string          `yaml:"versionsEndpoint" default:"https://tinfoil.media/repo/db/versions.json"`
	OrganizeOptions   OrganizeOptions `yaml:"organizeOptions"`
}

func (o *AppSettings) SetDefaults() {
	if defaults.CanUpdate(o.AppDataDirectory) {
		dir, err := utils.GetExecDir()
		if err == nil {
			o.AppDataDirectory = dir
		}
	}
}

var (
	ErrConfigurationFileNotFound = errors.New("configuration file not found")
)

type ConfigurationChangedCallback func(old AppSettings, new AppSettings)
type UnsubscribeFunction func()

type ConfigurationProvider interface {
	// GetCurrentConfig retrieves current config or return s default if no config is present.
	GetCurrentConfig() AppSettings
	// UpdateConfig updates configuration and persists it in file if possible.
	UpdateConfig(settings AppSettings) error
	// OnConfigurationChanged registers callback on configuration changes. Returns a function that has to be called to unsubscribe.
	OnConfigurationChanged(callback ConfigurationChangedCallback) UnsubscribeFunction
}

type ConfigurationProviderImpl struct {
	mutex            sync.RWMutex
	settingsInstance AppSettings
	listeners        map[string]ConfigurationChangedCallback
	configFilePath   string
}

func NewConfigurationProvider(configFilePath string) (*ConfigurationProviderImpl, error) {
	var settings AppSettings
	err := defaults.Set(&settings)
	if err != nil {
		return nil, err
	}

	return &ConfigurationProviderImpl{
		mutex:            sync.RWMutex{},
		settingsInstance: settings,
		listeners:        make(map[string]ConfigurationChangedCallback),
		configFilePath:   configFilePath,
	}, nil

}

func (c *ConfigurationProviderImpl) LoadFromFile() error {
	f, err := os.Open(c.configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %w", ErrConfigurationFileNotFound, err)
		}
		return err
	}
	defer f.Close()

	var settings AppSettings
	err = yaml.NewDecoder(f).Decode(&settings)
	if err != nil {
		return err
	}
	err = defaults.Set(&settings)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	c.settingsInstance = settings
	defer c.mutex.Unlock()
	return nil
}

func (c *ConfigurationProviderImpl) SaveToFile() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	f, err := os.OpenFile(c.configFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = yaml.NewEncoder(f).Encode(&c.settingsInstance)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigurationProviderImpl) GetCurrentConfig() AppSettings {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.settingsInstance
}

func (c *ConfigurationProviderImpl) UpdateConfig(settings AppSettings) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.settingsInstance = settings
	// TODO: save to file
	return nil
}

func (c *ConfigurationProviderImpl) OnConfigurationChanged(callback ConfigurationChangedCallback) UnsubscribeFunction {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	callbackID := uuid.New().String()
	c.listeners[callbackID] = callback

	return func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		delete(c.listeners, callbackID)
	}
}
