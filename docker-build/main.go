package main

import (
	"context"
	"dagger/docker-build/internal/dagger"
	"fmt"
	"strings"
)

type DockerBuild struct{}

// Returns a container that echoes whatever string argument is provided
func (m *DockerBuild) Build(
	// Directory to build
	src *dagger.Directory,
	// Path to the Dockerfile to use (e.g., "frontend.Dockerfile").
	// +optional
	// +default="Dockerfile"
	dockerFile string,
	// Comma-separated list of build arguments in the format key=value
	// +optional
	// +default=[]
	buildArgs []string,
) (*dagger.Container, error) {
	// Create Build args from slice
	buildArgsSlice := make([]dagger.BuildArg, 0, len(buildArgs))
	for _, arg := range buildArgs {
		split := strings.Split(arg, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid build arg format: %s", arg)
		}
		buildArgsSlice = append(buildArgsSlice, dagger.BuildArg{
			Name:  split[0],
			Value: split[1],
		})
	}

	dockerFileBuildOpts := dagger.DirectoryDockerBuildOpts{
		Dockerfile: dockerFile,
		BuildArgs:  buildArgsSlice,
	}
	return src.DockerBuild(dockerFileBuildOpts), nil
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *DockerBuild) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}
