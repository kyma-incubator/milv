package pkg

import (
	"testing"
	"time"

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
			Files: []File{{
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
			Backoff:               2 * time.Second,
		}

		result, err := NewConfig(commands)

		require.NoError(t, err)
		assert.Equal(t, expected.Files, result.Files)
		assert.Equal(t, expected.Backoff, result.Backoff)
		assert.ElementsMatch(t, expected.ExternalLinksToIgnore, result.ExternalLinksToIgnore)
		assert.ElementsMatch(t, expected.InternalLinksToIgnore, result.InternalLinksToIgnore)
		assert.ElementsMatch(t, expected.FilesToIgnore, result.FilesToIgnore)
	})
}
