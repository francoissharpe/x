package main

import (
	"context"
	"dagger/gitty/internal/dagger"
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"strings"
)

func setupGitContainer(src *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("alpine/git").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		// Pull tags
		WithExec([]string{"git", "fetch", "--tags"})
}

func convertNewlineStringToSlice(str string) []string {
	return strings.Split(strings.TrimSpace(str), "\n")
}

type Gitty struct{}

func (m *Gitty) Debug(ctx context.Context, src *dagger.Directory) (string, error) {
	return setupGitContainer(src).WithExec([]string{"ls", "-1", ".git/refs/heads"}).Stdout(ctx)
}

func (m *Gitty) GitInfo(ctx context.Context, src *dagger.Directory) (*GitInfo, error) {
	eg, gctx := errgroup.WithContext(ctx)

	results := make(map[string]*string)
	commands := map[string][]string{
		"author":    {"git", "log", "-1", "--pretty=format:%an"},
		"hash":      {"git", "log", "-1", "--pretty=format:%H"},
		"date":      {"git", "log", "-1", "--pretty=format:%ai"},
		"shortHash": {"git", "log", "-1", "--pretty=format:%h"},
		"message":   {"git", "log", "-1", "--pretty=format:%s"},
		"remoteUrl": {"git", "config", "--get", "remote.origin.url"},
		"branch":    {"ls", "-1", ".git/refs/heads"},
		"tags":      {"git", "tag", "-l", "--contains", "HEAD"},
	}

	for key, cmd := range commands {
		results[key] = new(string)
		eg.Go(func() error {
			var err error
			result, err := setupGitContainer(src).WithExec(cmd).Stdout(gctx)
			*results[key] = strings.TrimSpace(result)
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &GitInfo{
		RemoteUrl:       *results["remoteUrl"],
		CommitHash:      *results["hash"],
		CommitHashShort: *results["shortHash"],
		CommitAuthor:    *results["author"],
		CommitDate:      *results["date"],
		CommitMessage:   *results["message"],
		Branch:          *results["branch"],
		Tags:            convertNewlineStringToSlice(*results["tags"]),
	}, nil
}

type GitInfo struct {
	RemoteUrl       string   `json:"remoteUrl"`
	Branch          string   `json:"branch"`
	CommitHash      string   `json:"commitHash"`
	CommitHashShort string   `json:"commitHashShort"`
	CommitAuthor    string   `json:"commitAuthor"`
	CommitDate      string   `json:"commitDate"`
	CommitMessage   string   `json:"commitMessage"`
	Tags            []string `json:"tags"`
}

func (g *GitInfo) Json() (string, error) {
	data, err := json.Marshal(g)
	return string(data), err
}
