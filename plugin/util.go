package plugin

import (
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
