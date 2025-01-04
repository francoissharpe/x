package main

import (
	"dagger/opencontainers/internal/dagger"
)

type Opencontainers struct{}

func generateImageInfo(
	created, authors, url, documentation, source, version, revision, vendor, license, refName, title, description, baseDigest, baseName string,
) map[string]string {
	return map[string]string{
		"org.opencontainers.image.created":       created,
		"org.opencontainers.image.authors":       authors,
		"org.opencontainers.image.url":           url,
		"org.opencontainers.image.documentation": documentation,
		"org.opencontainers.image.source":        source,
		"org.opencontainers.image.version":       version,
		"org.opencontainers.image.revision":      revision,
		"org.opencontainers.image.vendor":        vendor,
		"org.opencontainers.image.license":       license,
		"org.opencontainers.image.ref.name":      refName,
		"org.opencontainers.image.title":         title,
		"org.opencontainers.image.description":   description,
		"org.opencontainers.image.base.digest":   baseDigest,
		"org.opencontainers.image.base.name":     baseName,
	}
}

// Returns a container that echoes whatever string argument is provided
func (m *Opencontainers) WithAnnotations(
	// Container to be annotated
	container *dagger.Container,
	// Date and time on which the image was built, conforming to RFC 3339.
	// +optional
	created string,
	// Contact details of the people or organization responsible for the image
	// +optional
	authors string,
	// URL to find more information on the image
	// +optional
	url string,
	// URL to get documentation on the image
	// +optional
	documentation string,
	// URL to get source code for building the image
	// +optional
	source string,
	// Version of the packaged software
	// +optional
	version string,
	// Source control revision identifier for the packaged software
	// +optional
	revison string,
	// Name of the distributing entity, organization or individual
	// +optional
	vendor string,
	// License(s) under which contained software is distributed as an SPDX License Expression
	// +optional
	license string,
	// Name of the reference for a target
	// +optional
	refName string,
	// Human-readable title of the image
	// +optional
	title string,
	// Human-readable description of the software packaged in the image
	// +optional
	description string,
	// Digest of the image this image is based on
	// +optional
	baseDigest string,
	// Image reference of the image this image is based on
	// +optional
	baseName string,
) *dagger.Container {
	// Add the annotations to the container
	for key, value := range generateImageInfo(
		created,
		authors,
		url,
		documentation,
		source,
		version,
		revison,
		vendor,
		license,
		refName,
		title,
		description,
		baseDigest,
		baseName,
	) {
		// Only add the annotation if the value is not empty
		if value == "" {
			continue
		}
		container = container.WithAnnotation(key, value)
	}
	return container
}

func (m *Opencontainers) WithLabels(
	// Container to be annotated
	container *dagger.Container,
	// Date and time on which the image was built, conforming to RFC 3339.
	// +optional
	created string,
	// Contact details of the people or organization responsible for the image
	// +optional
	authors string,
	// URL to find more information on the image
	// +optional
	url string,
	// URL to get documentation on the image
	// +optional
	documentation string,
	// URL to get source code for building the image
	// +optional
	source string,
	// Version of the packaged software
	// +optional
	version string,
	// Source control revision identifier for the packaged software
	// +optional
	revison string,
	// Name of the distributing entity, organization or individual
	// +optional
	vendor string,
	// License(s) under which contained software is distributed as an SPDX License Expression
	// +optional
	license string,
	// Name of the reference for a target
	// +optional
	refName string,
	// Human-readable title of the image
	// +optional
	title string,
	// Human-readable description of the software packaged in the image
	// +optional
	description string,
	// Digest of the image this image is based on
	// +optional
	baseDigest string,
	// Image reference of the image this image is based on
	// +optional
	baseName string,
) *dagger.Container {
	// Add the annotations to the container
	for key, value := range generateImageInfo(
		created,
		authors,
		url,
		documentation,
		source,
		version,
		revison,
		vendor,
		license,
		refName,
		title,
		description,
		baseDigest,
		baseName,
	) {
		// Only add the label if the value is not empty
		if value == "" {
			continue
		}
		container = container.WithLabel(key, value)
	}
	return container
}
