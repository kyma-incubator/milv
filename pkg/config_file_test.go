package pkg

import (
	"testing"

	"github.com/kyma-incubator/milv/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFile(t *testing.T) {
	t.Run("Config File", func(t *testing.T) {
		SetBasePath("test-markdowns", false)

		commands := cli.Commands{
			ConfigFile: "test-markdowns/milv-test.config.yaml",
		}

		expected := &FileConfig{
			ExternalLinksToIgnore: []string{"localhost", "abc.com", "github.com"},
			InternalLinksToIgnore: []string{"LICENSE", "#contributing"},
		}

		config, err := NewConfig(commands)
		require.NoError(t, err)

		result := NewFileConfig("./src/foo.md", config)

		require.NoError(t, err)
		assert.ElementsMatch(t, expected.ExternalLinksToIgnore, result.ExternalLinksToIgnore)
		assert.ElementsMatch(t, expected.InternalLinksToIgnore, result.InternalLinksToIgnore)
	})
}
