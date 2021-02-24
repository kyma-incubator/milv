package pkg_test

import (
	"fmt"
	"github.com/kyma-incubator/milv/cli"
	"github.com/kyma-incubator/milv/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilesToIgnoreInternalLinksIn(t *testing.T) {
	cliCommands := cli.Commands{ConfigFile: "test-markdowns/milv-test.config.yaml"}
	allMarkdowns := findAllMDFiles(t, "test-markdowns")
	cliCommands.Files = allMarkdowns
	cfg, err := pkg.NewConfig(cliCommands)
	require.NoError(t, err)

	files, err := pkg.NewFiles(cliCommands.Files, cfg)

	require.NoError(t, err)
	fmt.Printf("%+v", files)
	for _, file := range files {
		if strings.Contains(file.RelPath, "internal-ignore") {
			assert.NotNil(t, file.Config.IgnoreInternal)
			assert.True(t, *file.Config.IgnoreInternal)
		}
	}

}

func findAllMDFiles(t *testing.T, dir string) []string {
	var markdowns []string
	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			//location := path.Join(filePath, info.Name())
			markdowns = append(markdowns, filePath)
		}
		return nil
	})
	require.NoError(t, err)
	return markdowns
}
