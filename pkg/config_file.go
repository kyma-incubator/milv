package pkg

import (
	"fmt"
	"path"
	"strings"
	"time"
)

type FileConfig struct {
	BasePath              string
	Backoff               time.Duration `yaml:"backoff"`
	ExternalLinksToIgnore []string      `yaml:"external-links-to-ignore"`
	InternalLinksToIgnore []string      `yaml:"internal-links-to-ignore"`
	Timeout               *int          `yaml:"timeout"`
	RequestRepeats        *int          `yaml:"request-repeats"`
	AllowRedirect         *bool         `yaml:"allow-redirect"`
	AllowCodeBlocks       *bool         `yaml:"allow-code-blocks"`
	IgnoreExternal        *bool         `yaml:"ignore-external"`
	IgnoreInternal        *bool         `yaml:"ignore-internal"`
}

func NewFileConfig(filePath string, config *Config) FileConfig {
	if config == nil {
		return FileConfig{}
	}

	cfg := *config
	file := File{}

	if found, foundFile := findFile(filePath, cfg.Files); found {
		file = foundFile
	}

	fileCfg := FileConfig{}
	if file.Config != nil {
		fileCfg = *file.Config
	}

	timeout := getDefaultIntIfNil(cfg.Timeout, fileCfg.Timeout)
	requestRepeats := getDefaultIntIfNil(cfg.RequestRepeats, fileCfg.RequestRepeats)
	allowRedirect := getDefaultBoolIfNil(cfg.AllowRedirect, fileCfg.AllowRedirect)
	allowCodeBlocks := getDefaultBoolIfNil(cfg.AllowCodeBlocks, fileCfg.AllowCodeBlocks)

	backoff := getDefaultDurationIfNotProvided(cfg.Backoff, fileCfg.Backoff)
	ignoreInternal := getInternalIgnorePolicy(filePath, cfg, fileCfg)
	ignoreExternal := getDefaultBoolIfNil(cfg.IgnoreExternal, fileCfg.IgnoreExternal)

	externalLinksToIgnore := getExternalLinksToIgnore(cfg, file.Config)
	internalLinksToIgnore := getInternalLinksToIgnore(cfg, file.Config)

	return FileConfig{
		BasePath:              config.BasePath,
		Backoff:               backoff,
		ExternalLinksToIgnore: externalLinksToIgnore,
		InternalLinksToIgnore: internalLinksToIgnore,
		Timeout:               &timeout,
		RequestRepeats:        &requestRepeats,
		AllowRedirect:         &allowRedirect,
		AllowCodeBlocks:       &allowCodeBlocks,
		IgnoreExternal:        &ignoreExternal,
		IgnoreInternal:        &ignoreInternal,
	}
}

func findFile(filePath string, files []File) (bool, File) {
	for _, file := range files {
		if filePath == file.RelPath && file.Config != nil {
			return true, file
		}
	}
	return false, File{}
}

func getExternalLinksToIgnore(config Config, fileConfig *FileConfig) []string {
	externalLinksToIgnore := config.ExternalLinksToIgnore
	if fileConfig != nil {
		externalLinksToIgnore = unique(append(config.ExternalLinksToIgnore, fileConfig.ExternalLinksToIgnore...))
	}

	return externalLinksToIgnore
}

func getInternalLinksToIgnore(config Config, fileConfig *FileConfig) []string {
	internalLinksToIgnore := config.InternalLinksToIgnore
	if fileConfig != nil {
		internalLinksToIgnore = unique(append(config.InternalLinksToIgnore, fileConfig.InternalLinksToIgnore...))
	}

	return internalLinksToIgnore
}

func getInternalIgnorePolicy(filepath string, config Config, fileConfig FileConfig) bool {
	internalIgnore := config.IgnoreInternal

	if isFileIgnored(filepath, config.FilesToIgnoreInternalLinksIn) {
		internalIgnore = true
	}

	if fileConfig.IgnoreInternal != nil {
		internalIgnore = *fileConfig.IgnoreInternal
	}

	return internalIgnore
}

func isFileIgnored(filePath string, filesToIgnore []string) bool {
	for _, fileToIgnore := range filesToIgnore {
		if strings.HasPrefix(fileToIgnore, ".") {
			if checkIfFileIsInIgnorePath(fileToIgnore, filePath) {
				return true
			}
		} else {
			if checkIfFilePathContainsIgnoredDir(fileToIgnore, filePath) {
				return true
			}
		}
	}
	return false
}

func getDefaultBoolIfNil(defaultValue bool, value *bool) bool {
	if value == nil {
		return defaultValue
	}
	return *value
}

func getDefaultIntIfNil(defaultValue int, value *int) int {
	if value == nil {
		return defaultValue
	}
	return *value
}

func getDefaultDurationIfNotProvided(defaultValue, value time.Duration) time.Duration {
	if value != 0 {
		return value
	}
	return defaultValue
}

func checkIfFileIsInIgnorePath(fileToIgnore, filePath string) bool {
	startingPath := path.Clean(fileToIgnore)
	cleanFilePath := path.Clean(filePath)

	output := strings.HasPrefix(cleanFilePath, startingPath)
	return output
}

func checkIfFilePathContainsIgnoredDir(fileToIgnore, filePath string) bool {
	rootedFilePath := fmt.Sprintf(`/%s`, filePath)
	dirToIgnore := fmt.Sprintf(`/%s/`, fileToIgnore)
	return strings.Contains(rootedFilePath, dirToIgnore)
}
