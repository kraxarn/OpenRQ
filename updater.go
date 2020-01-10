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
	// Download to buffer
	resp, err := http.Get(fmt.Sprintf("https://kraxarn.com/openrq/updater/%v", runtime.GOOS))
	}
	defer resp.Body.Close()
}
