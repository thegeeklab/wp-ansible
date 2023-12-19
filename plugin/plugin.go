package plugin

import (
	wp "github.com/thegeeklab/wp-plugin-go/plugin"
	"github.com/urfave/cli/v2"
)

// Plugin implements provide the plugin.
type Plugin struct {
	*wp.Plugin
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	Requirements      string
	Galaxy            string
	Inventories       cli.StringSlice
	Playbooks         cli.StringSlice
	Limit             string
	SkipTags          string
	StartAtTask       string
	Tags              string
	ExtraVars         cli.StringSlice
	ModulePath        cli.StringSlice
	Check             bool
	Diff              bool
	FlushCache        bool
	ForceHandlers     bool
	ListHosts         bool
	ListTags          bool
	ListTasks         bool
	SyntaxCheck       bool
	Forks             int
	VaultID           string
	VaultPassword     string
	VaultPasswordFile string
	Verbose           int
	PrivateKey        string
	PrivateKeyFile    string
	User              string
	Connection        string
	Timeout           int
	SSHCommonArgs     string
	SFTPExtraArgs     string
	SCPExtraArgs      string
	SSHExtraArgs      string
	Become            bool
	BecomeMethod      string
	BecomeUser        string
}

func New(options wp.Options, settings *Settings) *Plugin {
	p := &Plugin{}

	if options.Execute == nil {
		options.Execute = p.run
	}

	p.Plugin = wp.New(options)
	p.Settings = settings

	return p
}
