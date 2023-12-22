package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/execabs"
)

var (
	ErrPluginPlaybookNotSet  = errors.New("playbook is required")
	ErrPluginInventoryNotSet = errors.New("inventory is required")
)

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
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
	if len(p.Settings.Playbooks.Value()) == 0 {
		return ErrPluginPlaybookNotSet
	}

	if len(p.Settings.Inventories.Value()) == 0 {
		return ErrPluginInventoryNotSet
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	if err := p.playbooks(); err != nil {
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

	commands := []*execabs.Cmd{
		p.versionCommand(),
	}

	if p.Settings.PythonRequirements != "" {
		commands = append(commands, p.pythonRequirementsCommand())
	}

	if p.Settings.GalaxyRequirements != "" {
		commands = append(commands, p.galaxyRequirementsCommand())
	}

	for _, inventory := range p.Settings.Inventories.Value() {
		commands = append(commands, p.ansibleCommand(inventory))
	}

	for _, cmd := range commands {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "ANSIBLE_FORCE_COLOR=1")

		trace(cmd)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
