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

	t.Run("File has IgnoreInternal config set to True", func(t *testing.T) {
		tcs := []struct {
			Name                         string
			FilePath                     string
			Expected                     bool
			FilesToIgnoreInternalLinksIn []string
		}{
			{
				Name:                         "File has IgnoreInternal config set to True",
				FilePath:                     "ignore-me-internally/my-markdown.md",
				Expected:                     true,
				FilesToIgnoreInternalLinksIn: []string{"ignore-me-internally"},
			}, {
				Name:                         "File has IgnoreInternal config set to True, relative path",
				FilePath:                     "./ignore-me-internally/my-markdown.md",
				Expected:                     true,
				FilesToIgnoreInternalLinksIn: []string{"ignore-me-internally"},
			}, {
				Name:                         "Ignore concrete directory",
				FilePath:                     "./ignore-me-internally/my-markdown.md",
				Expected:                     true,
				FilesToIgnoreInternalLinksIn: []string{"./ignore-me-internally"},
			}, {
				Name:                         "",
				FilePath:                     "ignore-me-internally/my-markdown.md",
				Expected:                     true,
				FilesToIgnoreInternalLinksIn: []string{"./ignore-me-internally"},
			}, {
				Name:                         "File should be ignored even contains ignored substring ",
				FilePath:                     "not-ignore-me/my-markdown.md",
				Expected:                     false,
				FilesToIgnoreInternalLinksIn: []string{"ignore"},
			},
			{
				Name:                         "File should be ignored",
				FilePath:                     "not-ignore-me/my-markdown.md",
				Expected:                     false,
				FilesToIgnoreInternalLinksIn: []string{"./ignore-me-internally"},
			}}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				//GIVEN
				cfg := &Config{
					FilesToIgnoreInternalLinksIn: tc.FilesToIgnoreInternalLinksIn,
				}

				//WHEN
				fileCfg := NewFileConfig(tc.FilePath, cfg)

				//THEN
				require.NotNil(t, fileCfg)
				require.NotNil(t, fileCfg.IgnoreInternal)
				require.Equal(t, tc.Expected, *fileCfg.IgnoreInternal)
			})
		}
	})
}
