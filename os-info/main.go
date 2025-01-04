package main

import (
	"context"
	"dagger/os-info/internal/dagger"
	"encoding/json"
	"strings"
)

type Summary struct {
	// Name of the operating system
	Name string `json:"name"`
	// Pretty name of the operating system
	PrettyName string `json:"prettyName"`
	// Numeric version of the operating system
	Version string `json:"version"`
	// Codename of the operating system version
	VersionCodename string `json:"versionCodename"`
}

// Returns a JSON representation of the Summary
func (s *Summary) Json() string {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

type OsInfo struct{}

// Returns information about the operating system of the given container
func (m *OsInfo) Info(ctx context.Context, container *dagger.Container) *Summary {
	osRelease, err := container.WithExec([]string{"cat", "/etc/os-release"}).Stdout(ctx)
	if err != nil {
		return nil
	}
	osReleaseLines := strings.Split(osRelease, "\n")
	osReleaseMap := make(map[string]string)
	for _, line := range osReleaseLines {
		split := strings.Split(line, "=")
		if len(split) == 2 {
			osReleaseMap[split[0]] = strings.Trim(split[1], "\"")
		}
	}
	return &Summary{
		Name:            osReleaseMap["ID"],
		PrettyName:      osReleaseMap["PRETTY_NAME"],
		Version:         osReleaseMap["VERSION_ID"],
		VersionCodename: osReleaseMap["VERSION_CODENAME"],
	}
}
