package plugin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"
)

const (
	AnsibleForksDefault = 5

	ansibleFolder = "/etc/ansible"
	ansibleConfig = "/etc/ansible/ansible.cfg"

	pipBin             = "/usr/local/bin/pip"
	ansibleBin         = "/usr/local/bin/ansible"
	ansibleGalaxyBin   = "/usr/local/bin/ansible-galaxy"
	ansiblePlaybookBin = "/usr/local/bin/ansible-playbook"

	strictFilePerm = 0o600
)

const ansibleContent = `
[defaults]
host_key_checking = False
`

var ErrAnsiblePlaybookNotFound = errors.New("no playbook found")

// ansibleConfig creates the Ansible configuration directory and file.
// It ensures the directory exists and writes the Ansible configuration
// content to the config file with strict file permissions.
func (p *Plugin) ansibleConfig() error {
	if err := os.MkdirAll(ansibleFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create ansible directory: %w", err)
	}

	if err := os.WriteFile(ansibleConfig, []byte(ansibleContent), strictFilePerm); err != nil {
		return fmt.Errorf("failed to create ansible config: %w", err)
	}

	return nil
}

// privateKey creates a temporary file containing the private key specified in the plugin settings,
// and sets the PrivateKeyFile field in the plugin settings to the name of the temporary file.
// This is used to pass the private key to the Ansible command.
func (p *Plugin) privateKey() error {
	tmpfile, err := os.CreateTemp("", "privateKey")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}

	if _, err := tmpfile.Write([]byte(p.Settings.PrivateKey)); err != nil {
		return fmt.Errorf("failed to write private key file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close private key file: %w", err)
	}

	p.Settings.PrivateKeyFile = tmpfile.Name()

	return nil
}

// vaultPass creates a temporary file containing the vault password and sets the VaultPasswordFile
// field in the Plugin's Settings. This allows the vault password to be used when running
// Ansible commands that require it.
func (p *Plugin) vaultPass() error {
	tmpfile, err := os.CreateTemp("", "vaultPass")
	if err != nil {
		return fmt.Errorf("failed to create vault password file: %w", err)
	}

	if _, err := tmpfile.Write([]byte(p.Settings.VaultPassword)); err != nil {
		return fmt.Errorf("failed to write vault password file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close vault password file: %w", err)
	}

	p.Settings.VaultPasswordFile = tmpfile.Name()

	return nil
}

// getPlaybooks retrieves a list of playbook files based on the configured playbook patterns.
// If any of the patterns fail to match any files, the original pattern is included in the list.
// If no playbooks are found, ErrAnsiblePlaybookNotFound is returned.
func (p *Plugin) getPlaybooks() error {
	var playbooks []string

	for _, pb := range p.Settings.Playbooks.Value() {
		files, err := filepath.Glob(pb)
		if err != nil {
			playbooks = append(playbooks, pb)

			continue
		}

		playbooks = append(playbooks, files...)
	}

	if len(playbooks) == 0 {
		log.Debug().Strs("patterns", p.Settings.Playbooks.Value()).Msg("no playbooks found")

		return ErrAnsiblePlaybookNotFound
	}

	p.Settings.Playbooks = *cli.NewStringSlice(playbooks...)

	return nil
}

func (p *Plugin) versionCommand() *Cmd {
	args := []string{
		"--version",
	}

	return &Cmd{
		Cmd: execabs.Command(ansibleBin, args...),
	}
}

// pythonRequirementsCommand returns an execabs.Cmd that runs the pip install
// command with the specified Python requirements file and upgrades any existing
// packages.
func (p *Plugin) pythonRequirementsCommand() *Cmd {
	args := []string{
		"install",
		"--upgrade",
		"--requirement",
		p.Settings.PythonRequirements,
	}

	return &Cmd{
		Cmd: execabs.Command(pipBin, args...),
	}
}

// galaxyRequirementsCommand returns an execabs.Cmd that runs the Ansible Galaxy
// install command with the specified role file and verbosity level.
func (p *Plugin) galaxyRequirementsCommand() *Cmd {
	args := []string{
		"install",
		"--force",
		"--role-file",
		p.Settings.GalaxyRequirements,
	}

	if p.Settings.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", p.Settings.Verbose)))
	}

	return &Cmd{
		Cmd: execabs.Command(ansibleGalaxyBin, args...),
	}
}

// ansibleCommand returns an execabs.Cmd that runs the Ansible playbook with the
// specified inventory file and various configuration options set on the Plugin struct.
func (p *Plugin) ansibleCommand(inventory string) *Cmd {
	args := []string{
		"--inventory",
		inventory,
	}

	if len(p.Settings.ModulePath.Value()) > 0 {
		args = append(args, "--module-path", strings.Join(p.Settings.ModulePath.Value(), ":"))
	}

	if p.Settings.VaultID != "" {
		args = append(args, "--vault-id", p.Settings.VaultID)
	}

	if p.Settings.VaultPasswordFile != "" {
		args = append(args, "--vault-password-file", p.Settings.VaultPasswordFile)
	}

	for _, v := range p.Settings.ExtraVars.Value() {
		args = append(args, "--extra-vars", v)
	}

	if p.Settings.ListHosts {
		args = append(args, "--list-hosts")
		args = append(args, p.Settings.Playbooks.Value()...)

		return &Cmd{
			Cmd: execabs.Command(ansiblePlaybookBin, args...),
		}
	}

	if p.Settings.SyntaxCheck {
		args = append(args, "--syntax-check")
		args = append(args, p.Settings.Playbooks.Value()...)

		return &Cmd{
			Cmd: execabs.Command(ansiblePlaybookBin, args...),
		}
	}

	if p.Settings.Check {
		args = append(args, "--check")
	}

	if p.Settings.Diff {
		args = append(args, "--diff")
	}

	if p.Settings.FlushCache {
		args = append(args, "--flush-cache")
	}

	if p.Settings.ForceHandlers {
		args = append(args, "--force-handlers")
	}

	if p.Settings.Forks != AnsibleForksDefault {
		args = append(args, "--forks", strconv.Itoa(p.Settings.Forks))
	}

	if p.Settings.Limit != "" {
		args = append(args, "--limit", p.Settings.Limit)
	}

	if p.Settings.ListTags {
		args = append(args, "--list-tags")
	}

	if p.Settings.ListTasks {
		args = append(args, "--list-tasks")
	}

	if p.Settings.SkipTags != "" {
		args = append(args, "--skip-tags", p.Settings.SkipTags)
	}

	if p.Settings.StartAtTask != "" {
		args = append(args, "--start-at-task", p.Settings.StartAtTask)
	}

	if p.Settings.Tags != "" {
		args = append(args, "--tags", p.Settings.Tags)
	}

	if p.Settings.PrivateKeyFile != "" {
		args = append(args, "--private-key", p.Settings.PrivateKeyFile)
	}

	if p.Settings.User != "" {
		args = append(args, "--user", p.Settings.User)
	}

	if p.Settings.Connection != "" {
		args = append(args, "--connection", p.Settings.Connection)
	}

	if p.Settings.Timeout != 0 {
		args = append(args, "--timeout", strconv.Itoa(p.Settings.Timeout))
	}

	if p.Settings.SSHCommonArgs != "" {
		args = append(args, "--ssh-common-args", p.Settings.SSHCommonArgs)
	}

	if p.Settings.SFTPExtraArgs != "" {
		args = append(args, "--sftp-extra-args", p.Settings.SFTPExtraArgs)
	}

	if p.Settings.SCPExtraArgs != "" {
		args = append(args, "--scp-extra-args", p.Settings.SCPExtraArgs)
	}

	if p.Settings.SSHExtraArgs != "" {
		args = append(args, "--ssh-extra-args", p.Settings.SSHExtraArgs)
	}

	if p.Settings.Become {
		args = append(args, "--become")
	}

	if p.Settings.BecomeMethod != "" {
		args = append(args, "--become-method", p.Settings.BecomeMethod)
	}

	if p.Settings.BecomeUser != "" {
		args = append(args, "--become-user", p.Settings.BecomeUser)
	}

	if p.Settings.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", p.Settings.Verbose)))
	}

	args = append(args, p.Settings.Playbooks.Value()...)

	return &Cmd{
		Cmd:     execabs.Command(ansiblePlaybookBin, args...),
		Private: false,
	}
}
