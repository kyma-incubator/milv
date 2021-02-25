package pkg

import (
	"testing"

	"github.com/kyma-incubator/milv/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFile(t *testing.T) {
	t.Run("Config File", func(t *testing.T) {
		commands := cli.Commands{
			ConfigFile: "test-markdowns/milv-test.config.yaml",
			BasePath:   "test-markdowns",
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
			ShouldBeIgnored              bool
			FilesToIgnoreInternalLinksIn []string
		}{
			{
				Name:                         "File and Ignore has relative path",
				FilePath:                     "ignore-me-internally/my-markdown.md",
				ShouldBeIgnored:              true,
				FilesToIgnoreInternalLinksIn: []string{"ignore-me-internally"},
			}, {
				Name:                         "File has relative path with ./ path and Ignore has relative path",
				FilePath:                     "./ignore-me-internally/my-markdown.md",
				ShouldBeIgnored:              true,
				FilesToIgnoreInternalLinksIn: []string{"ignore-me-internally"},
			}, {
				Name:                         "File and Ignore has relative path with ./",
				FilePath:                     "./ignore-me-internally/my-markdown.md",
				ShouldBeIgnored:              true,
				FilesToIgnoreInternalLinksIn: []string{"./ignore-me-internally"},
			}, {
				Name:                         "File has relative path path and Ignore has relative path with ./",
				FilePath:                     "ignore-me-internally/my-markdown.md",
				ShouldBeIgnored:              true,
				FilesToIgnoreInternalLinksIn: []string{"./ignore-me-internally"},
			}, {
				Name:                         "File should not be ignored when contains ignored substring in path",
				FilePath:                     "not-ignore-me/my-markdown.md",
				ShouldBeIgnored:              false,
				FilesToIgnoreInternalLinksIn: []string{"ignore"},
			},
			{
				Name:                         "File should be ignored",
				FilePath:                     "./ignore-me-internally/not-ignore-me/my-markdown.md",
				ShouldBeIgnored:              true,
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
				require.Equal(t, tc.ShouldBeIgnored, *fileCfg.IgnoreInternal)
			})
		}
	})
}
