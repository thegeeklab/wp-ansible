package plugin

import (
	"context"
	"fmt"
	"os"

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
	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	batchCmd := make([]*types.Cmd, 0)
	batchCmd = append(batchCmd, p.versionCommand())

	if err := p.getPlaybooks(); err != nil {
		return err
	}

	if err := p.ansibleConfig(); err != nil {
		return err
	}

	if p.Settings.PrivateKey != "" {
		if err := p.privateKey(); err != nil {
			return err
		}

		defer os.Remove(p.Settings.PrivateKeyFile)
	}

	if p.Settings.VaultPassword != "" {
		if err := p.vaultPass(); err != nil {
			return err
		}

		defer os.Remove(p.Settings.VaultPasswordFile)
	}

	if p.Settings.PythonRequirements != "" {
		batchCmd = append(batchCmd, p.pythonRequirementsCommand())
	}

	if p.Settings.GalaxyRequirements != "" {
		batchCmd = append(batchCmd, p.galaxyRequirementsCommand())
	}

	for _, inventory := range p.Settings.Inventories.Value() {
		batchCmd = append(batchCmd, p.ansibleCommand(inventory))
	}

	for _, cmd := range batchCmd {
		cmd.Env = append(cmd.Env, "ANSIBLE_FORCE_COLOR=1")

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
