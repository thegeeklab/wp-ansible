package plugin

import (
	"fmt"

	wp "github.com/thegeeklab/wp-plugin-go/v2/plugin"
	"github.com/urfave/cli/v2"
)

//go:generate go run ../internal/docs/main.go -output=../docs/data/data-raw.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*wp.Plugin
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	PythonRequirements string
	GalaxyRequirements string
	Inventories        cli.StringSlice
	Playbooks          cli.StringSlice
	Limit              string
	SkipTags           string
	StartAtTask        string
	Tags               string
	ExtraVars          cli.StringSlice
	ModulePath         cli.StringSlice
	Check              bool
	Diff               bool
	FlushCache         bool
	ForceHandlers      bool
	ListHosts          bool
	ListTags           bool
	ListTasks          bool
	SyntaxCheck        bool
	Forks              int
	VaultID            string
	VaultPassword      string
	VaultPasswordFile  string
	Verbose            int
	PrivateKey         string
	PrivateKeyFile     string
	User               string
	Connection         string
	Timeout            int
	SSHCommonArgs      string
	SFTPExtraArgs      string
	SCPExtraArgs       string
	SSHExtraArgs       string
	Become             bool
	BecomeMethod       string
	BecomeUser         string
}

func New(e wp.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := wp.Options{
		Name:                "wp-ansible",
		Description:         "Manage infrastructure with Ansible",
		Flags:               Flags(p.Settings, wp.FlagsPluginCategory),
		HideWoodpeckerFlags: true,
	}

	if len(build) > 0 {
		options.Version = build[0]
	}

	if len(build) > 1 {
		options.VersionMetadata = fmt.Sprintf("date=%s", build[1])
	}

	if e == nil {
		options.Execute = p.run
	}

	p.Plugin = wp.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "python-requirements",
			Usage:       "path to python requirements file",
			EnvVars:     []string{"PLUGIN_PYTHON_REQUIREMENTS"},
			Destination: &settings.PythonRequirements,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "galaxy-requirements",
			Usage:       "path to galaxy requirements file",
			EnvVars:     []string{"PLUGIN_GALAXY_REQUIREMENTS"},
			Destination: &settings.GalaxyRequirements,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "inventory",
			Usage:       "path to inventory file",
			EnvVars:     []string{"PLUGIN_INVENTORY", "PLUGIN_INVENTORIES"},
			Required:    true,
			Destination: &settings.Inventories,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "playbook",
			Usage:       "list of playbooks to apply",
			EnvVars:     []string{"PLUGIN_PLAYBOOK", "PLUGIN_PLAYBOOKS"},
			Required:    true,
			Destination: &settings.Playbooks,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "limit",
			Usage:       "limit selected hosts to an additional pattern",
			EnvVars:     []string{"PLUGIN_LIMIT"},
			Destination: &settings.Limit,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "skip-tags",
			Usage:       "only run plays and tasks whose tags do not match",
			EnvVars:     []string{"PLUGIN_SKIP_TAGS"},
			Destination: &settings.SkipTags,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "start-at-task",
			Usage:       "start the playbook at the task matching this name",
			EnvVars:     []string{"PLUGIN_START_AT_TASK"},
			Destination: &settings.StartAtTask,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "tags",
			Usage:       "only run plays and tasks tagged with these values",
			EnvVars:     []string{"PLUGIN_TAGS"},
			Destination: &settings.Tags,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "extra-vars",
			Usage:       "set additional variables as `key=value`",
			EnvVars:     []string{"PLUGIN_EXTRA_VARS", "ANSIBLE_EXTRA_VARS"},
			Destination: &settings.ExtraVars,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "module-path",
			Usage:       "prepend paths to module library",
			EnvVars:     []string{"PLUGIN_MODULE_PATH"},
			Destination: &settings.ModulePath,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "check",
			Usage:       "run a check, do not apply any changes",
			EnvVars:     []string{"PLUGIN_CHECK"},
			Destination: &settings.Check,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "diff",
			Usage:       "show the differences, may print secrets",
			EnvVars:     []string{"PLUGIN_DIFF"},
			Destination: &settings.Diff,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "flush-cache",
			Usage:       "clear the fact cache for every host in inventory",
			EnvVars:     []string{"PLUGIN_FLUSH_CACHE"},
			Destination: &settings.FlushCache,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "force-handlers",
			Usage:       "run handlers even if a task fails",
			EnvVars:     []string{"PLUGIN_FORCE_HANDLERS"},
			Destination: &settings.ForceHandlers,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "list-hosts",
			Usage:       "outputs a list of matching hosts",
			EnvVars:     []string{"PLUGIN_LIST_HOSTS"},
			Destination: &settings.ListHosts,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "list-tags",
			Usage:       "list all available tags",
			EnvVars:     []string{"PLUGIN_LIST_TAGS"},
			Destination: &settings.ListTags,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "list-tasks",
			Usage:       "list all tasks that would be executed",
			EnvVars:     []string{"PLUGIN_LIST_TASKS"},
			Destination: &settings.ListTasks,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "syntax-check",
			Usage:       "perform a syntax check on the playbook",
			EnvVars:     []string{"PLUGIN_SYNTAX_CHECK"},
			Destination: &settings.SyntaxCheck,
			Category:    category,
		},
		&cli.IntFlag{
			Name:        "forks",
			Usage:       "specify number of parallel processes to use",
			EnvVars:     []string{"PLUGIN_FORKS"},
			Value:       AnsibleForksDefault,
			Destination: &settings.Forks,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "vault-id",
			Usage:       "the vault identity to use",
			EnvVars:     []string{"PLUGIN_VAULT_ID", "ANSIBLE_VAULT_ID"},
			Destination: &settings.VaultID,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "vault-password",
			Usage:       "the vault password to use",
			EnvVars:     []string{"PLUGIN_VAULT_PASSWORD", "ANSIBLE_VAULT_PASSWORD"},
			Destination: &settings.VaultPassword,
			Category:    category,
		},
		&cli.IntFlag{
			Name:        "verbose",
			Usage:       "level of verbosity, 0 up to 4",
			EnvVars:     []string{"PLUGIN_VERBOSE"},
			Destination: &settings.Verbose,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "private-key",
			Usage:       "SSH private key used to authenticate the connection",
			EnvVars:     []string{"PLUGIN_PRIVATE_KEY", "ANSIBLE_PRIVATE_KEY"},
			Destination: &settings.PrivateKey,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "user",
			Usage:       "connect as this user",
			EnvVars:     []string{"PLUGIN_USER", "ANSIBLE_USER"},
			Destination: &settings.User,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "connection",
			Usage:       "connection type to use",
			EnvVars:     []string{"PLUGIN_CONNECTION"},
			Destination: &settings.Connection,
			Category:    category,
		},
		&cli.IntFlag{
			Name:        "timeout",
			Usage:       "override the connection timeout in seconds",
			EnvVars:     []string{"PLUGIN_TIMEOUT"},
			Destination: &settings.Timeout,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "ssh-common-args",
			Usage:       "specify common arguments to pass to SFTP, SCP and SSH connections",
			EnvVars:     []string{"PLUGIN_SSH_COMMON_ARGS"},
			Destination: &settings.SSHCommonArgs,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "sftp-extra-args",
			Usage:       "specify extra arguments to pass to SFTP connections only",
			EnvVars:     []string{"PLUGIN_SFTP_EXTRA_ARGS"},
			Destination: &settings.SFTPExtraArgs,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "scp-extra-args",
			Usage:       "specify extra arguments to pass to SCP connections only",
			EnvVars:     []string{"PLUGIN_SCP_EXTRA_ARGS"},
			Destination: &settings.SCPExtraArgs,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "ssh-extra-args",
			Usage:       "specify extra arguments to pass to SSH connections only",
			EnvVars:     []string{"PLUGIN_SSH_EXTRA_ARGS"},
			Destination: &settings.SSHExtraArgs,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "become",
			Usage:       "enable privilege escalation",
			EnvVars:     []string{"PLUGIN_BECOME"},
			Destination: &settings.Become,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "become-method",
			Usage:       "privilege escalation method to use",
			EnvVars:     []string{"PLUGIN_BECOME_METHOD", "ANSIBLE_BECOME_METHOD"},
			Destination: &settings.BecomeMethod,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "become-user",
			Usage:       "privilege escalation user to use",
			EnvVars:     []string{"PLUGIN_BECOME_USER", "ANSIBLE_BECOME_USER"},
			Destination: &settings.BecomeUser,
			Category:    category,
		},
	}
}
