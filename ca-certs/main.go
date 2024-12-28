package main

import (
	"dagger/ca-certs/internal/dagger"
	"path"
)

func ShDashC(command string) []string {
	return []string{"sh", "-c", command}
}

func WithCommonEnvVars(container *dagger.Container) *dagger.Container {
	return container.WithEnvVariable("CURL_CA_BUNDLE", "/etc/ssl/certs/ca-certificates.crt").
		WithEnvVariable("NODE_EXTRA_CA_CERTS", "/etc/ssl/certs/ca-certificates.crt").
		WithEnvVariable("REQUESTS_CA_BUNDLE", "/etc/ssl/certs/ca-certificates.crt").
		WithEnvVariable("SSL_CERT_FILE", "/etc/ssl/certs/ca-certificates.crt").
		WithEnvVariable("GRPC_DEFAULT_SSL_ROOTS_FILE_PATH", "/etc/ssl/certs/ca-certificates.crt")
}

type CaCerts struct{}

// adds the specified CA certificates to the container and updates the CA certificates file in the given container.
func (m *CaCerts) WithCaCerts(
	// The container to add the CA certificates to.
	container *dagger.Container,
	// The list of HTTP urls of CA certificates to add.
	caCerts []string,
	// The directory to store the CA certificates in.
	// +optional
	caCertsDir string,
) *dagger.Container {
	if caCertsDir == "" {
		caCertsDir = "/usr/local/share/ca-certificates"
	}
	caCertsCtr := dag.Container().From("buildpack-deps:bookworm-curl")
	for _, caCert := range caCerts {
		certFilename := path.Base(caCert) + ".crt"
		outputPath := path.Join(caCertsDir, certFilename)
		caCertsCtr = caCertsCtr.WithExec(ShDashC("curl -sSL " + caCert + " -o " + outputPath))
	}

	caCertsCtr = caCertsCtr.WithExec(ShDashC("update-ca-certificates"))
	caCertificatesCrtFile := caCertsCtr.File("/etc/ssl/certs/ca-certificates.crt")

	return container.With(WithCommonEnvVars).
		WithDirectory(caCertsDir, caCertsCtr.Directory(caCertsDir)).
		WithFile("/etc/ssl/certs/ca-certificates.crt", caCertificatesCrtFile)
}
