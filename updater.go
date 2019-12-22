package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// IsLatestVersion checks for updates and returns true if no update is available
func IsLatestVersion() bool {
	// For now, we always just check the latest commit
	resp, err := http.Get("https://kraxarn.com/openrq/updater/commit")
	if err != nil {
		fmt.Println("error: failed to check for updates:", err)
		return false
	}
	// Read body
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return versionCommitHash == strings.TrimSpace(string(body))
}

func Update() error {
	// Check if running from 'go run'
	if strings.HasPrefix(os.Args[0], "/tmp") {
		//return fmt.Errorf("cannot update from temporary folder")
	}
	// Download to buffer
	resp, err := http.Get(fmt.Sprintf("https://kraxarn.com/openrq/updater/%v", runtime.GOOS))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// Get permission for current file
	fileInfo, _ := os.Stat(os.Args[0])
	// Try to remove running file first
	if err := os.Remove(os.Args[0]); err != nil {
		return err
	}
	// Write to current executable with same permissions
	return ioutil.WriteFile(os.Args[0], body, fileInfo.Mode())
}
