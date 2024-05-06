package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipInstall(t *testing.T) {
	tests := []struct {
		name         string
		requirements string
		want         []string
	}{
		{
			name:         "with valid requirements file",
			requirements: "requirements.txt",
			want:         []string{pipBin, "install", "--upgrade", "--requirement", "requirements.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := PipInstall(tt.requirements)
			assert.Equal(t, tt.want, cmd.Cmd.Args)
		})
	}
}
