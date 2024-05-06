package plugin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
			require.Equal(t, tt.want, cmd.Cmd.Args)
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		wantErr  bool
	}{
		{
			name:     "successful write",
			filename: "test.txt",
			content:  "test content",
			wantErr:  false,
		},
		{
			name:     "empty content",
			filename: "test.txt",
			content:  "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := WriteFile(tt.filename, tt.content)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			defer os.Remove(path)

			content, err := os.ReadFile(path)
			require.NoError(t, err)
			require.Equal(t, tt.content, string(content))
		})
	}
}
