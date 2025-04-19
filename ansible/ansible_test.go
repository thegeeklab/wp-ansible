package ansible

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string
		ansible *Ansible
		want    []string
	}{
		{
			name:    "test version command",
			ansible: &Ansible{},
			want:    []string{ansibleBin, "--version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.ansible.Version()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestGalaxyInstall(t *testing.T) {
	tests := []struct {
		name    string
		ansible *Ansible
		want    []string
	}{
		{
			name: "with valid requirements file and no verbosity",
			ansible: &Ansible{
				GalaxyRequirements: "requirements.yml",
			},
			want: []string{ansibleGalaxyBin, "install", "--force", "--role-file", "requirements.yml"},
		},
		{
			name: "with valid requirements file and verbosity level 1",
			ansible: &Ansible{
				GalaxyRequirements: "requirements.yml",
				Verbose:            1,
			},
			want: []string{ansibleGalaxyBin, "install", "--force", "--role-file", "requirements.yml", "-v"},
		},
		{
			name: "with valid requirements file and verbosity level 3",
			ansible: &Ansible{
				GalaxyRequirements: "requirements.yml",
				Verbose:            3,
			},
			want: []string{ansibleGalaxyBin, "install", "--force", "--role-file", "requirements.yml", "-vvv"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.ansible.GalaxyInstall()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestAnsibleCommand(t *testing.T) {
	tests := []struct {
		name    string
		ansible *Ansible
		want    []string
	}{
		{
			name: "with inventory and no other settings",
			ansible: &Ansible{
				Inventories: []string{"inventory.yml"},
			},
			want: []string{ansiblePlaybookBin, "--inventory", "inventory.yml", "--forks", "0"},
		},
		{
			name: "with inventory and module path",
			ansible: &Ansible{
				Inventories: []string{"inventory.yml"},
				ModulePath:  []string{"/path/to/modules"},
			},
			want: []string{
				ansiblePlaybookBin, "--inventory", "inventory.yml", "--module-path",
				"/path/to/modules", "--forks", "0",
			},
		},
		{
			name: "with inventory, module path, and vault ID",
			ansible: &Ansible{
				Inventories: []string{"inventory.yml"},
				ModulePath:  []string{"/path/to/modules"},
				VaultID:     "my_vault_id",
			},
			want: []string{
				ansiblePlaybookBin, "--inventory", "inventory.yml", "--module-path", "/path/to/modules",
				"--vault-id", "my_vault_id", "--forks", "0",
			},
		},
		{
			name: "with inventory, module path, vault ID, and vault password file",
			ansible: &Ansible{
				Inventories:       []string{"inventory.yml"},
				ModulePath:        []string{"/path/to/modules"},
				VaultID:           "my_vault_id",
				VaultPasswordFile: "/path/to/vault/password/file",
			},
			want: []string{
				ansiblePlaybookBin, "--inventory", "inventory.yml", "--module-path", "/path/to/modules",
				"--vault-id", "my_vault_id", "--vault-password-file", "/path/to/vault/password/file",
				"--forks", "0",
			},
		},
		{
			name: "with inventory, module path, vault ID, vault password file, and extra vars",
			ansible: &Ansible{
				Inventories:       []string{"inventory.yml"},
				ModulePath:        []string{"/path/to/modules"},
				VaultID:           "my_vault_id",
				VaultPasswordFile: "/path/to/vault/password/file",
				ExtraVars:         []string{"var1=value1", "var2=value2"},
			},
			want: []string{
				ansiblePlaybookBin, "--inventory", "inventory.yml", "--module-path", "/path/to/modules",
				"--vault-id", "my_vault_id", "--vault-password-file", "/path/to/vault/password/file",
				"--extra-vars", "var1=value1", "--extra-vars", "var2=value2", "--forks", "0",
			},
		},
		{
			name: "with inventory and list hosts",
			ansible: &Ansible{
				ListHosts:   true,
				Inventories: []string{"inventory.yml"},
				Playbooks:   []string{"playbook1.yml", "playbook2.yml"},
			},
			want: []string{
				ansiblePlaybookBin, "--inventory", "inventory.yml", "--list-hosts",
				"playbook1.yml", "playbook2.yml",
			},
		},
		{
			name: "with inventory and syntax check",
			ansible: &Ansible{
				SyntaxCheck: true,
				Inventories: []string{"inventory.yml"},
				Playbooks:   []string{"playbook1.yml", "playbook2.yml"},
			},
			want: []string{
				ansiblePlaybookBin, "--inventory", "inventory.yml", "--syntax-check",
				"playbook1.yml", "playbook2.yml",
			},
		},
		{
			name: "with all options",
			ansible: &Ansible{
				Check:          true,
				Diff:           true,
				FlushCache:     true,
				ForceHandlers:  true,
				Forks:          10,
				Limit:          "host1,host2",
				ListTags:       true,
				ListTasks:      true,
				SkipTags:       "tag1,tag2",
				StartAtTask:    "task_name",
				Tags:           "tag3,tag4",
				PrivateKeyFile: "/path/to/private/key",
				User:           "remote_user",
				Connection:     "ssh",
				Timeout:        60,
				SSHCommonArgs:  "-o StrictHostKeyChecking=no",
				SFTPExtraArgs:  "-o IdentitiesOnly=yes",
				SCPExtraArgs:   "-r",
				SSHExtraArgs:   "-o ForwardAgent=yes",
				Become:         true,
				BecomeMethod:   "sudo",
				BecomeUser:     "root",
				Verbose:        2,
				Inventories:    []string{"inventory.yml"},
				Playbooks:      []string{"playbook1.yml", "playbook2.yml"},
			},
			want: []string{
				ansiblePlaybookBin, "--inventory", "inventory.yml", "--check", "--diff", "--flush-cache",
				"--force-handlers", "--forks", "10", "--limit", "host1,host2", "--list-tags", "--list-tasks",
				"--skip-tags", "tag1,tag2", "--start-at-task", "task_name", "--tags", "tag3,tag4",
				"--private-key", "/path/to/private/key", "--user", "remote_user", "--connection", "ssh",
				"--timeout", "60", "--ssh-common-args", "-o StrictHostKeyChecking=no", "--sftp-extra-args",
				"-o IdentitiesOnly=yes", "--scp-extra-args", "-r", "--ssh-extra-args", "-o ForwardAgent=yes",
				"--become", "--become-method", "sudo", "--become-user", "root", "-vv", "playbook1.yml", "playbook2.yml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.ansible.Play()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}
