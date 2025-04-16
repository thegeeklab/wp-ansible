package plugin

import (
	"os"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v5/exec"
)

const pipBin = "/usr/local/bin/pip"

// PipInstall returns a command to install Python packages from a requirements file.
// The command will upgrade any existing packages and install the packages specified in the given requirements file.
func PipInstall(req string) *plugin_exec.Cmd {
	args := []string{
		"install",
		"--upgrade",
		"--requirement",
		req,
	}

	cmd := plugin_exec.Command(pipBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
