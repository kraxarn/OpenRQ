package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IsLatestVersion checks for updates and returns true if no update is available
func IsLatestVersion() (bool, error) {
	// Check the latest release on GitHub
	resp, err := http.Get("https://api.github.com/repos/kraxarn/OpenRQ/releases")
	if err != nil {
		return false, err
	}
	// Read body
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// Parse as JSON
	releases := make([]interface{}, 0)
	if err = json.Unmarshal(body, &releases); err != nil {
		return false, err
	}
	// Get the latest version
	latest, ok := releases[0].(map[string]interface{})["tag_name"]
	if !ok {
		return false, fmt.Errorf("failed to parse tag name")
	}
	return latest.(string) == versionTagName, nil
}
