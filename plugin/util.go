package plugin

import (
	"fmt"
	"os"

	"github.com/thegeeklab/wp-plugin-go/v2/types"
	"golang.org/x/sys/execabs"
)

const pipBin = "/usr/local/bin/pip"

// PipInstall returns a command to install Python packages from a requirements file.
// The command will upgrade any existing packages and install the packages specified in the given requirements file.
func PipInstall(req string) *types.Cmd {
	args := []string{
		"install",
		"--upgrade",
		"--requirement",
		req,
	}

	return &types.Cmd{
		Cmd: execabs.Command(pipBin, args...),
	}
}

// WriteFile creates a temporary file with the given name and content, and returns the path to the created file.
func WriteFile(name, content string) (string, error) {
	tmpfile, err := os.CreateTemp("", name)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		return "", fmt.Errorf("failed to close file: %w", err)
	}

	return tmpfile.Name(), nil
}
