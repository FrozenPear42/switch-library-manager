package settings

import (
	_ "embed"
	"github.com/hashicorp/go-version"
	"io"
	"net/http"
	"strings"
)

const (
	updateURL = "https://raw.githubusercontent.com/FrozenPear42/switch-library-manager/master/settings/.version"
)

//go:embed .version
var versionBytes []byte

// GetCurrentVersion returns current version tag
func GetCurrentVersion() string {
	return strings.TrimSpace(string(versionBytes))
}

// CheckForUpdates checks for updates and returns (new version, is an update, error)
func CheckForUpdates() (string, bool, error) {
	currentVersionString := GetCurrentVersion()

	res, err := http.Get(updateURL)
	if err != nil {
		return "", false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", false, err
	}
	remoteVersionString := strings.TrimSpace(string(body))

	currentVersion, err := version.NewVersion(currentVersionString)
	if err != nil {
		return "", false, err
	}
	remoteVersion, err := version.NewVersion(remoteVersionString)
	if err != nil {
		return "", false, err
	}
	if remoteVersion.GreaterThan(currentVersion) {
		return remoteVersionString, true, nil
	}
	return remoteVersionString, false, nil
}
