package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/thegeeklab/wp-plugin-go/v2/file"
	"github.com/thegeeklab/wp-plugin-go/v2/types"
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

	batchCmd := make([]*types.Cmd, 0)

	batchCmd = append(batchCmd, p.Settings.Ansible.Version())

	if p.Settings.PrivateKey != "" {
		if p.Settings.Ansible.PrivateKeyFile, err = file.WriteTmpFile("privateKey", p.Settings.PrivateKey); err != nil {
			return err
		}

		defer os.Remove(p.Settings.Ansible.PrivateKeyFile)
	}

	if p.Settings.VaultPassword != "" {
		if p.Settings.Ansible.VaultPasswordFile, err = file.WriteTmpFile("vaultPass", p.Settings.VaultPassword); err != nil {
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
		cmd.Env = append(os.Environ(), "ANSIBLE_FORCE_COLOR=1")

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
