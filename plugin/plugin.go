package plugin

import (
	"fmt"

	"github.com/thegeeklab/wp-ansible/ansible"
	plugin_base "github.com/thegeeklab/wp-plugin-go/v6/plugin"
	"github.com/urfave/cli/v3"
)

//go:generate go run ../internal/docs/main.go -output=../docs/data/data-raw.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*plugin_base.Plugin
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	PythonRequirements string
	PrivateKey         string
	VaultPassword      string
	Ansible            ansible.Ansible
}

func New(e plugin_base.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := plugin_base.Options{
		Name:                "wp-ansible",
		Description:         "Manage infrastructure with Ansible",
		Flags:               Flags(p.Settings, plugin_base.FlagsPluginCategory),
		Execute:             p.run,
		HideWoodpeckerFlags: true,
	}

	if len(build) > 0 {
		options.Version = build[0]
	}

	if len(build) > 1 {
		options.VersionMetadata = fmt.Sprintf("date=%s", build[1])
	}

	if e != nil {
		options.Execute = e
	}

	p.Plugin = plugin_base.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "python-requirements",
			Usage:       "path to python requirements file",
			Sources:     cli.EnvVars("PLUGIN_PYTHON_REQUIREMENTS"),
			Destination: &settings.PythonRequirements,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "galaxy-requirements",
			Usage:       "path to galaxy requirements file",
			Sources:     cli.EnvVars("PLUGIN_GALAXY_REQUIREMENTS"),
			Destination: &settings.Ansible.GalaxyRequirements,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "inventory",
			Usage:       "path to inventory file",
			Sources:     cli.EnvVars("PLUGIN_INVENTORY", "PLUGIN_INVENTORIES"),
			Required:    true,
			Destination: &settings.Ansible.Inventories,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "playbook",
			Usage:       "list of playbooks to apply",
			Sources:     cli.EnvVars("PLUGIN_PLAYBOOK", "PLUGIN_PLAYBOOKS"),
			Required:    true,
			Destination: &settings.Ansible.Playbooks,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "limit",
			Usage:       "limit selected hosts to an additional pattern",
			Sources:     cli.EnvVars("PLUGIN_LIMIT"),
			Destination: &settings.Ansible.Limit,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "skip-tags",
			Usage:       "only run plays and tasks whose tags do not match",
			Sources:     cli.EnvVars("PLUGIN_SKIP_TAGS"),
			Destination: &settings.Ansible.SkipTags,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "start-at-task",
			Usage:       "start the playbook at the task matching this name",
			Sources:     cli.EnvVars("PLUGIN_START_AT_TASK"),
			Destination: &settings.Ansible.StartAtTask,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "tags",
			Usage:       "only run plays and tasks tagged with these values",
			Sources:     cli.EnvVars("PLUGIN_TAGS"),
			Destination: &settings.Ansible.Tags,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "extra-vars",
			Usage:       "set additional variables as `key=value`",
			Sources:     cli.EnvVars("PLUGIN_EXTRA_VARS", "ANSIBLE_EXTRA_VARS"),
			Destination: &settings.Ansible.ExtraVars,
			Category:    category,
		},
		&cli.StringSliceFlag{
			Name:        "module-path",
			Usage:       "prepend paths to module library",
			Sources:     cli.EnvVars("PLUGIN_MODULE_PATH"),
			Destination: &settings.Ansible.ModulePath,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "check",
			Usage:       "run a check, do not apply any changes",
			Sources:     cli.EnvVars("PLUGIN_CHECK"),
			Destination: &settings.Ansible.Check,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "diff",
			Usage:       "show the differences, may print secrets",
			Sources:     cli.EnvVars("PLUGIN_DIFF"),
			Destination: &settings.Ansible.Diff,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "flush-cache",
			Usage:       "clear the fact cache for every host in inventory",
			Sources:     cli.EnvVars("PLUGIN_FLUSH_CACHE"),
			Destination: &settings.Ansible.FlushCache,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "force-handlers",
			Usage:       "run handlers even if a task fails",
			Sources:     cli.EnvVars("PLUGIN_FORCE_HANDLERS"),
			Destination: &settings.Ansible.ForceHandlers,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "list-hosts",
			Usage:       "outputs a list of matching hosts",
			Sources:     cli.EnvVars("PLUGIN_LIST_HOSTS"),
			Destination: &settings.Ansible.ListHosts,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "list-tags",
			Usage:       "list all available tags",
			Sources:     cli.EnvVars("PLUGIN_LIST_TAGS"),
			Destination: &settings.Ansible.ListTags,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "list-tasks",
			Usage:       "list all tasks that would be executed",
			Sources:     cli.EnvVars("PLUGIN_LIST_TASKS"),
			Destination: &settings.Ansible.ListTasks,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "syntax-check",
			Usage:       "perform a syntax check on the playbook",
			Sources:     cli.EnvVars("PLUGIN_SYNTAX_CHECK"),
			Destination: &settings.Ansible.SyntaxCheck,
			Category:    category,
		},
		&cli.IntFlag{
			Name:        "forks",
			Usage:       "specify number of parallel processes to use",
			Sources:     cli.EnvVars("PLUGIN_FORKS"),
			Value:       ansible.AnsibleForksDefault,
			Destination: &settings.Ansible.Forks,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "vault-id",
			Usage:       "the vault identity to use",
			Sources:     cli.EnvVars("PLUGIN_VAULT_ID", "ANSIBLE_VAULT_ID"),
			Destination: &settings.Ansible.VaultID,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "vault-password",
			Usage:       "the vault password to use",
			Sources:     cli.EnvVars("PLUGIN_VAULT_PASSWORD", "ANSIBLE_VAULT_PASSWORD"),
			Destination: &settings.VaultPassword,
			Category:    category,
		},
		&cli.IntFlag{
			Name:        "verbose",
			Usage:       "level of verbosity, 0 up to 4",
			Sources:     cli.EnvVars("PLUGIN_VERBOSE"),
			Destination: &settings.Ansible.Verbose,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "private-key",
			Usage:       "SSH private key used to authenticate the connection",
			Sources:     cli.EnvVars("PLUGIN_PRIVATE_KEY", "ANSIBLE_PRIVATE_KEY"),
			Destination: &settings.PrivateKey,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "user",
			Usage:       "connect as this user",
			Sources:     cli.EnvVars("PLUGIN_USER", "ANSIBLE_USER"),
			Destination: &settings.Ansible.User,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "connection",
			Usage:       "connection type to use",
			Sources:     cli.EnvVars("PLUGIN_CONNECTION"),
			Destination: &settings.Ansible.Connection,
			Category:    category,
		},
		&cli.IntFlag{
			Name:        "timeout",
			Usage:       "override the connection timeout in seconds",
			Sources:     cli.EnvVars("PLUGIN_TIMEOUT"),
			Destination: &settings.Ansible.Timeout,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "ssh-common-args",
			Usage:       "specify common arguments to pass to SFTP, SCP and SSH connections",
			Sources:     cli.EnvVars("PLUGIN_SSH_COMMON_ARGS"),
			Destination: &settings.Ansible.SSHCommonArgs,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "sftp-extra-args",
			Usage:       "specify extra arguments to pass to SFTP connections only",
			Sources:     cli.EnvVars("PLUGIN_SFTP_EXTRA_ARGS"),
			Destination: &settings.Ansible.SFTPExtraArgs,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "scp-extra-args",
			Usage:       "specify extra arguments to pass to SCP connections only",
			Sources:     cli.EnvVars("PLUGIN_SCP_EXTRA_ARGS"),
			Destination: &settings.Ansible.SCPExtraArgs,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "ssh-extra-args",
			Usage:       "specify extra arguments to pass to SSH connections only",
			Sources:     cli.EnvVars("PLUGIN_SSH_EXTRA_ARGS"),
			Destination: &settings.Ansible.SSHExtraArgs,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "become",
			Usage:       "enable privilege escalation",
			Sources:     cli.EnvVars("PLUGIN_BECOME"),
			Destination: &settings.Ansible.Become,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "become-method",
			Usage:       "privilege escalation method to use",
			Sources:     cli.EnvVars("PLUGIN_BECOME_METHOD", "ANSIBLE_BECOME_METHOD"),
			Destination: &settings.Ansible.BecomeMethod,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "become-user",
			Usage:       "privilege escalation user to use",
			Sources:     cli.EnvVars("PLUGIN_BECOME_USER", "ANSIBLE_BECOME_USER"),
			Destination: &settings.Ansible.BecomeUser,
			Category:    category,
		},
	}
}
