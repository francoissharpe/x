// TODO: Implement restore (with deps caching), format, build and test functions for dotnet
// TODO: Implement functionality to get the ProjectReference from the csproj file and only copy the required projects during the build/publish process, basically get all local path dependencies and copy them to the build container
// TODO: Somehow figure out if an app is a console or web app and set the entrypoint and environment variables accordingly
// TODO: Set OCI labels using the module for the container based on the source dir as well as Properties from the csproj file
// TODO: Allow configuring a private/other nuget feed for the restore process
// TODO: Add a function for installing dotnet tools

package main

import (
	"context"
	"dagger/dotnet/internal/dagger"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"encoding/xml"
)

type Project struct {
	Name string
	Path string
	Type string
}

type Solution struct {
	Projects []Project
}

type Dotnet struct{}

func (m *Dotnet) GetProjectFromSolutionByName(ctx context.Context, src *dagger.Directory, project string) (*Project, error) {
	solution, err := m.SolutionInfo(ctx, src)
	if err != nil {
		return nil, err
	}
	if len(solution.Projects) == 0 {
		return nil, fmt.Errorf("no projects found in solution")
	}
	for _, p := range solution.Projects {
		if p.Name == project {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("project %s not found in the solution", project)
}

func (m *Dotnet) CsProjectVersion(ctx context.Context, src *dagger.Directory, project string) (string, error) {
	projectObj, err := m.GetProjectFromSolutionByName(ctx, src, project)
	if projectObj.Path == "" {
		return "", fmt.Errorf("project %s not found in the solution", project)
	}

	csprojFileContents, err := src.File(projectObj.Path).Contents(ctx)
	if err != nil {
		return "", err
	}
	type PropertyGroup struct {
		TargetFramework string `xml:"TargetFramework"`
	}

	type ProjectFile struct {
		PropertyGroups []PropertyGroup `xml:"PropertyGroup"`
	}

	var projectFile ProjectFile
	err = xml.Unmarshal([]byte(csprojFileContents), &projectFile)
	if err != nil {
		return "", err
	}

	for _, pg := range projectFile.PropertyGroups {
		if pg.TargetFramework != "" {
			// Remove the net from the target framework version
			return strings.TrimPrefix(pg.TargetFramework, "net"), nil
		}
	}
	return "", fmt.Errorf("target framework not found in project file")
}

func (m *Dotnet) Publish(
	ctx context.Context,
	src *dagger.Directory,
	project string,
) (*dagger.Container, error) {
	projectObj, err := m.GetProjectFromSolutionByName(ctx, src, project)
	if projectObj.Path == "" {
		return nil, fmt.Errorf("project %s not found in the solution", project)
	}
	csProjectVersion, err := m.CsProjectVersion(ctx, src, project)
	if err != nil {
		return nil, err
	}
	publishDir := dag.Container().
		From(fmt.Sprintf("mcr.microsoft.com/dotnet/sdk:%s", csProjectVersion)).
		WithWorkdir("/src").
		WithMountedDirectory("/src", src).
		WithExec([]string{"dotnet", "publish", "-o", "./bin", projectObj.Path}).
		Directory("/src/bin")

	return dag.Container().
		From(fmt.Sprintf("mcr.microsoft.com/dotnet/aspnet:%s", csProjectVersion)).
		WithWorkdir("/app").
		WithDirectory("/app", publishDir).
		WithEntrypoint([]string{"dotnet"}).
		WithDefaultArgs([]string{fmt.Sprintf("%s.dll", project)}), nil
}

func (m *Dotnet) SolutionInfo(
	ctx context.Context,
	// +ignore=["*", "!**/*.sln"]
	src *dagger.Directory,
) (*Solution, error) {
	projects, err := dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:8.0").
		WithWorkdir("/src").
		WithMountedDirectory("/src", src).
		WithExec([]string{"dotnet", "sln", "list"}).
		Stdout(ctx)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(projects, "\n")
	var projectList []Project

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.Contains(line, "Project(s)") || strings.Contains(line, "---") {
			continue
		}

		name := path.Base(line)
		path := strings.TrimSpace(line)
		ext := filepath.Ext(name)[1:]

		projectList = append(projectList, Project{
			// trim the filename suffix
			Name: strings.TrimSuffix(name, fmt.Sprintf(".%s", ext)),
			Path: path,
			Type: ext,
		})
	}

	solution := Solution{
		Projects: projectList,
	}

	return &solution, nil
}
