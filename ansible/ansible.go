package ansible

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	plugin_exec "github.com/thegeeklab/wp-plugin-go/v4/exec"
	"github.com/urfave/cli/v2"
)

const (
	AnsibleForksDefault = 5

	ansibleBin         = "/usr/local/bin/ansible"
	ansibleGalaxyBin   = "/usr/local/bin/ansible-galaxy"
	ansiblePlaybookBin = "/usr/local/bin/ansible-playbook"
)

var ErrAnsiblePlaybookNotFound = errors.New("no playbook found")

type Ansible struct {
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
	VaultPasswordFile  string
	Verbose            int
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

// Version runs the Ansible binary with the --version flag to retrieve the current version.
func (a *Ansible) Version() *plugin_exec.Cmd {
	args := []string{
		"--version",
	}

	cmd := plugin_exec.Command(ansibleBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// GetPlaybooks retrieves the list of Ansible playbook files based on the configured playbook patterns.
func (a *Ansible) GetPlaybooks() error {
	var playbooks []string

	for _, pb := range a.Playbooks.Value() {
		files, err := filepath.Glob(pb)
		if err != nil {
			playbooks = append(playbooks, pb)

			continue
		}

		playbooks = append(playbooks, files...)
	}

	if len(playbooks) == 0 {
		log.Debug().Strs("patterns", a.Playbooks.Value()).Msg("no playbooks found")

		return ErrAnsiblePlaybookNotFound
	}

	a.Playbooks = *cli.NewStringSlice(playbooks...)

	return nil
}

// GalaxyInstall runs the ansible-galaxy install command with the configured options.
func (a *Ansible) GalaxyInstall() *plugin_exec.Cmd {
	args := []string{
		"install",
		"--force",
		"--role-file",
		a.GalaxyRequirements,
	}

	if a.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", a.Verbose)))
	}

	cmd := plugin_exec.Command(ansibleGalaxyBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// Play runs the Ansible playbook with the configured options.
//
//nolint:gocyclo
func (a *Ansible) Play() *plugin_exec.Cmd {
	args := make([]string, 0)

	for _, inventory := range a.Inventories.Value() {
		args = append(args, "--inventory", inventory)
	}

	if len(a.ModulePath.Value()) > 0 {
		args = append(args, "--module-path", strings.Join(a.ModulePath.Value(), ":"))
	}

	if a.VaultID != "" {
		args = append(args, "--vault-id", a.VaultID)
	}

	if a.VaultPasswordFile != "" {
		args = append(args, "--vault-password-file", a.VaultPasswordFile)
	}

	for _, v := range a.ExtraVars.Value() {
		args = append(args, "--extra-vars", v)
	}

	if a.ListHosts {
		args = append(args, "--list-hosts")
		args = append(args, a.Playbooks.Value()...)

		cmd := plugin_exec.Command(ansiblePlaybookBin, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd
	}

	if a.SyntaxCheck {
		args = append(args, "--syntax-check")
		args = append(args, a.Playbooks.Value()...)

		cmd := plugin_exec.Command(ansiblePlaybookBin, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd
	}

	if a.Check {
		args = append(args, "--check")
	}

	if a.Diff {
		args = append(args, "--diff")
	}

	if a.FlushCache {
		args = append(args, "--flush-cache")
	}

	if a.ForceHandlers {
		args = append(args, "--force-handlers")
	}

	if a.Forks != AnsibleForksDefault {
		args = append(args, "--forks", strconv.Itoa(a.Forks))
	}

	if a.Limit != "" {
		args = append(args, "--limit", a.Limit)
	}

	if a.ListTags {
		args = append(args, "--list-tags")
	}

	if a.ListTasks {
		args = append(args, "--list-tasks")
	}

	if a.SkipTags != "" {
		args = append(args, "--skip-tags", a.SkipTags)
	}

	if a.StartAtTask != "" {
		args = append(args, "--start-at-task", a.StartAtTask)
	}

	if a.Tags != "" {
		args = append(args, "--tags", a.Tags)
	}

	if a.PrivateKeyFile != "" {
		args = append(args, "--private-key", a.PrivateKeyFile)
	}

	if a.User != "" {
		args = append(args, "--user", a.User)
	}

	if a.Connection != "" {
		args = append(args, "--connection", a.Connection)
	}

	if a.Timeout != 0 {
		args = append(args, "--timeout", strconv.Itoa(a.Timeout))
	}

	if a.SSHCommonArgs != "" {
		args = append(args, "--ssh-common-args", a.SSHCommonArgs)
	}

	if a.SFTPExtraArgs != "" {
		args = append(args, "--sftp-extra-args", a.SFTPExtraArgs)
	}

	if a.SCPExtraArgs != "" {
		args = append(args, "--scp-extra-args", a.SCPExtraArgs)
	}

	if a.SSHExtraArgs != "" {
		args = append(args, "--ssh-extra-args", a.SSHExtraArgs)
	}

	if a.Become {
		args = append(args, "--become")
	}

	if a.BecomeMethod != "" {
		args = append(args, "--become-method", a.BecomeMethod)
	}

	if a.BecomeUser != "" {
		args = append(args, "--become-user", a.BecomeUser)
	}

	if a.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", a.Verbose)))
	}

	args = append(args, a.Playbooks.Value()...)

	cmd := plugin_exec.Command(ansiblePlaybookBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
