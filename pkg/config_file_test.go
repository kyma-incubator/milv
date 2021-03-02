package pkg

import (
	"testing"
	"time"

	"github.com/kyma-incubator/milv/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileConfig(t *testing.T) {
	trueBool := true
	falseBool := false
	t.Run("Check if links to ignore are merged", func(t *testing.T) {
		//GIVEN
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

		//WHEN
		result := NewFileConfig("./src/foo.md", config)

		//THEN
		require.NoError(t, err)
		assert.ElementsMatch(t, expected.ExternalLinksToIgnore, result.ExternalLinksToIgnore)
		assert.ElementsMatch(t, expected.InternalLinksToIgnore, result.InternalLinksToIgnore)
	})

	t.Run("Check different scenario for ignoring internal links paths", func(t *testing.T) {
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

	t.Run("Config without file Configs", func(t *testing.T) {
		//GIVEN
		timeout := 5
		requestRepeats := 6
		cfg := &Config{
			BasePath:        "path",
			RequestRepeats:  requestRepeats,
			Timeout:         timeout,
			AllowRedirect:   true,
			AllowCodeBlocks: true,
			IgnoreExternal:  false,
			IgnoreInternal:  true,
			Backoff:         5 * time.Hour,
		}

		expectedCfg := &FileConfig{
			BasePath:        "path",
			Timeout:         &timeout,
			RequestRepeats:  &requestRepeats,
			AllowRedirect:   &trueBool,
			AllowCodeBlocks: &trueBool,
			IgnoreExternal:  &falseBool,
			IgnoreInternal:  &trueBool,
			Backoff:         5 * time.Hour,
		}
		//WHEN
		newConfig := NewFileConfig("any-path", cfg)

		//THEN
		require.NotNil(t, newConfig)
		assert.Equal(t, expectedCfg, newConfig)
	})

	t.Run("Config with matching File Configs", func(t *testing.T) {
		//GIVEN
		timeout := 5
		requestRepeats := 6
		filePath := "path"

		files := []File{
			{RelPath: "some-random/documentation.md"},
			{
				RelPath: filePath,
				Config: &FileConfig{
					Timeout:         &timeout,
					RequestRepeats:  &requestRepeats,
					AllowRedirect:   &trueBool,
					AllowCodeBlocks: &trueBool,
					IgnoreExternal:  &trueBool,
					IgnoreInternal:  &falseBool,
					Backoff:         10 * time.Second,
				}},
		}

		cfg := &Config{
			Files:           files,
			Timeout:         timeout,
			RequestRepeats:  requestRepeats,
			AllowRedirect:   false,
			AllowCodeBlocks: false,
			IgnoreExternal:  false,
			IgnoreInternal:  true,
			Backoff:         1 * time.Second,
		}

		expectedCfg := FileConfig{
			ExternalLinksToIgnore: []string{},
			InternalLinksToIgnore: []string{},
			Timeout:               &timeout,
			RequestRepeats:        &requestRepeats,
			AllowRedirect:         &trueBool,
			AllowCodeBlocks:       &trueBool,
			IgnoreExternal:        &trueBool,
			IgnoreInternal:        &falseBool,
			Backoff:               10 * time.Second,
		}

		//WHEN
		mergedFileConfig := NewFileConfig(filePath, cfg)

		//THEN
		require.NotNil(t, mergedFileConfig)
		assert.Equal(t, expectedCfg, *mergedFileConfig)
	})

	t.Run("Test with several ignore paths", func(t *testing.T) {
		//GIVEN
		firstPathToIgnore := "./website/src"
		secondPathToIgnore := "./website/content"
		filePath := "./website/content/news.md"
		cfg := Config{FilesToIgnoreInternalLinksIn: []string{firstPathToIgnore, secondPathToIgnore}}

		//WHEN
		output := NewFileConfig(filePath, &cfg)

		//THEN
		require.NotNil(t, output)
		require.NotNil(t, output.IgnoreInternal)
		require.True(t, *output.IgnoreInternal)
	})
}
