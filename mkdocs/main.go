package main

import (
	"dagger/mkdocs/internal/dagger"
	"fmt"
)

type Mkdocs struct{}

// Build the Mkdocs site and return the site directory
func (m *Mkdocs) Build(
	// src is the directory containing the Mkdocs project
	src *dagger.Directory,
	// The directory which container the built site
	// +optional
	// +default="site"
	siteDir string,
	// version is the version of the mkdocs-material image to use
	// +optional
	// +default="latest"
	version string,
) *dagger.Directory {
	workingDir := "/src"
	mkDocsImage := fmt.Sprintf("squidfunk/mkdocs-material:%s", version)
	return dag.Container().
		From(mkDocsImage).
		WithMountedDirectory(workingDir, src).
		WithWorkdir(workingDir).
		WithExec([]string{"mkdocs", "build"}).
		Directory(fmt.Sprintf("%s/%s", workingDir, siteDir))
}

func (m *Mkdocs) Serve(
	// src is the directory containing the Mkdocs project
	src *dagger.Directory,
	// The directory which container the built site
	// +optional
	// +default="site"
	siteDir string,
	// version is the version of the mkdocs-material image to use
	// +optional
	// +default="latest"
	version string,
	// port is the port to serve the Mkdocs site on
	// +optional
	// +default=80
	port int,
) *dagger.Container {
	buildDir := m.Build(src, siteDir, version)
	return dag.Container().
		From("nginx:alpine").
		WithEnvVariable("NGINX_PORT", fmt.Sprintf("%d", port)).
		WithMountedDirectory("/usr/share/nginx/html", buildDir).
		WithExposedPort(port)
}