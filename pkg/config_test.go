package pkg

import (
	"testing"

	"github.com/kyma-incubator/milv/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("Config File", func(t *testing.T) {
		commands := cli.Commands{
			ConfigFile: "test-markdowns/milv-test.config.yaml",
		}

		expected := &Config{
			Files: []File{
				File{
					RelPath: "./src/foo.md",
					Config: &FileConfig{
						ExternalLinksToIgnore: []string{"github.com"},
						InternalLinksToIgnore: []string{"#contributing"},
					},
				},
			},
			ExternalLinksToIgnore: []string{"localhost", "abc.com"},
			InternalLinksToIgnore: []string{"LICENSE"},
			FilesToIgnore:         []string{"./README.md"},
		}

		result, err := NewConfig(commands)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
