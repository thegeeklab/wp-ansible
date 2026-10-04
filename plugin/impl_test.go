package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thegeeklab/wp-ansible/ansible"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		files     map[string]string
		playbooks []string
		wantErr   error
	}{
		{
			name:      "playbook found",
			files:     map[string]string{"playbook.yml": "---\n"},
			playbooks: []string{"playbook.yml"},
		},
		{
			name:      "no playbook matches glob",
			playbooks: []string{"*.yml"},
			wantErr:   ansible.ErrAnsiblePlaybookNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			for name, content := range tt.files {
				require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
			}

			playbooks := make([]string, 0, len(tt.playbooks))
			for _, playbook := range tt.playbooks {
				playbooks = append(playbooks, filepath.Join(dir, playbook))
			}

			p := &Plugin{Settings: &Settings{Ansible: ansible.Ansible{Playbooks: playbooks}}}

			err := p.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)

				return
			}

			assert.NoError(t, err)
		})
	}
}
