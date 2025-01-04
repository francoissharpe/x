package main

import (
	"context"
	"dagger/procfile/internal/dagger"
	"fmt"
	"html/template"
	"strings"
)

// createDockerEntryPointSh returns a shell script that given an argument will start the process with the given name.
func createDockerEntryPointSh(procs []Process) string {
	tmpl := `#!/bin/sh
set -e

# Function to handle SIGTERM and SIGINT
_term() {
  echo "Caught SIGTERM signal!"
  kill -TERM "$child" 2>/dev/null
}

# Trap SIGTERM and SIGINT
trap _term SIGTERM SIGINT

# Start the process based on the argument
case "$1" in
{{- range . }}
  {{ .Name }})
	echo "Starting {{ .Name }}..."
	exec {{ .Command }} &
	child=$!
	wait "$child"
	;;
{{- end }}
  *)
	echo "Invalid process name: $1"
	exit 1
	;;
esac
`
	t, err := template.New("entrypoint").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	var script strings.Builder
	err = t.Execute(&script, procs)
	if err != nil {
		panic(err)
	}
	return script.String()
}

type Process struct {
	Name    string
	Command string
}

type Procfile struct{}

func (m *Procfile) Parse(ctx context.Context, procfile *dagger.File) ([]Process, error) {
	var procs []Process
	procfileContents, err := procfile.Contents(ctx)
	if err != nil {
		return nil, err
	}
	// Parse the procfile and populate the procs slice.
	for _, line := range strings.Split(procfileContents, "\n") {
		// Split the line into process name and command.
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid procfile line: %s", line)
		}
		processName, command := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		procs = append(procs, Process{Name: processName, Command: command})
	}
	return procs, nil
}

// Given a Procfile and container, return the container with the entrypoint configured to start with the given process.
func (m *Procfile) WithProcfile(
	ctx context.Context,
	// The container to configure.
	container *dagger.Container,
	// The Procfile to use.
	procfile *dagger.File,
	// Path to the entrypoint script.
	// +optional
	// +default="/docker-entrypoint.sh"
	entryPointPath string,
) *dagger.Container {
	procs, err := m.Parse(ctx, procfile)
	if err != nil {
		fmt.Errorf("failed to parse Procfile: %v", err)
	}
	return container.WithNewFile(entryPointPath, createDockerEntryPointSh(procs)).WithEntrypoint([]string{entryPointPath}).WithDefaultArgs([]string{procs[0].Name})
}
