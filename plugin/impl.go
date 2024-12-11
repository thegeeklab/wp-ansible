package plugin

import (
	"context"
	"fmt"
	"os"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v4/exec"
	plugin_file "github.com/thegeeklab/wp-plugin-go/v4/file"
)

func (p *Plugin) run(_ context.Context) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if err := p.Settings.Ansible.GetPlaybooks(); err != nil {
		return err
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	var err error

	batchCmd := make([]*plugin_exec.Cmd, 0)

	batchCmd = append(batchCmd, p.Settings.Ansible.Version())

	if p.Settings.PrivateKey != "" {
		p.Settings.Ansible.PrivateKeyFile, err = plugin_file.WriteTmpFile("privateKey", p.Settings.PrivateKey)
		if err != nil {
			return err
		}

		defer os.Remove(p.Settings.Ansible.PrivateKeyFile)
	}

	if p.Settings.VaultPassword != "" {
		p.Settings.Ansible.VaultPasswordFile, err = plugin_file.WriteTmpFile("vaultPass", p.Settings.VaultPassword)
		if err != nil {
			return err
		}

		defer os.Remove(p.Settings.Ansible.VaultPasswordFile)
	}

	if p.Settings.PythonRequirements != "" {
		batchCmd = append(batchCmd, PipInstall(p.Settings.PythonRequirements))
	}

	if p.Settings.Ansible.GalaxyRequirements != "" {
		batchCmd = append(batchCmd, p.Settings.Ansible.GalaxyInstall())
	}

	batchCmd = append(batchCmd, p.Settings.Ansible.Play())

	for _, cmd := range batchCmd {
		if cmd == nil {
			continue
		}

		cmd.Env = append(os.Environ(), "ANSIBLE_FORCE_COLOR=1")

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
