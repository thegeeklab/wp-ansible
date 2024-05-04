package main

import (
	"github.com/thegeeklab/wp-ansible/plugin"
)

//nolint:gochecknoglobals
var (
	BuildVersion = "devel"
	BuildDate    = "00000000"
)

func main() {
	plugin.New(nil, BuildVersion, BuildDate).Run()
}
