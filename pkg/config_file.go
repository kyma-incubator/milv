package pkg

import (
	"fmt"
	"path"
	"strings"
)

type FileConfig struct {
	BasePath              string
	ExternalLinksToIgnore []string `yaml:"external-links-to-ignore"`
	InternalLinksToIgnore []string `yaml:"internal-links-to-ignore"`
	Timeout               *int     `yaml:"timeout"`
	RequestRepeats        *int8    `yaml:"request-repeats"`
	AllowRedirect         *bool    `yaml:"allow-redirect"`
	AllowCodeBlocks       *bool    `yaml:"allow-code-blocks"`
	IgnoreExternal        *bool    `yaml:"ignore-external"`
	IgnoreInternal        *bool    `yaml:"ignore-internal"`
}

func CombineConfigsForFile(filePath string, config *Config) *FileConfig {
	if config != nil {
		cfg := *config
		file := File{}

		if found, foundFile := findFile(filePath, cfg.Files); found {
			file = foundFile
		}

		timeout := getTimeoutPolicy(cfg, file.Config)
		requestRepeats := getRequestRepeatPolicy(cfg, file.Config)
		allowRedirect := getAllowRedirectPolicy(cfg, file.Config)
		allowCodeBlocks := getAllowCodeBlocksPolicy(cfg, file.Config)

		ignoreInternal := getInternalIgnorePolicy(filePath, cfg, file.Config)
		ignoreExternal := getIgnoreExternalPolicy(cfg, file.Config)

		externalLinksToIgnore := getExternalLinksToIgnore(cfg, file.Config)
		internalLinksToIgnore := getInternalLinksToIgnore(cfg, file.Config)

		return &FileConfig{
			BasePath:              config.BasePath,
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
	return nil
}

func findFile(filePath string, files []File) (bool, File) {
	for _, file := range files {
		if filePath == file.RelPath && file.Config != nil {
			return true, file
		}
	}
	return false, File{}
}

func getRequestRepeatPolicy(config Config, fileConfig *FileConfig) int8 {
	requestRepeats := config.RequestRepeats

	if fileConfig != nil && fileConfig.RequestRepeats != nil {
		requestRepeats = *fileConfig.RequestRepeats
	}

	return requestRepeats
}

func getAllowRedirectPolicy(config Config, fileConfig *FileConfig) bool {
	allowRedirect := config.AllowRedirect
	if fileConfig != nil && fileConfig.AllowRedirect != nil {
		allowRedirect = *fileConfig.AllowRedirect
	}

	return allowRedirect
}

func getAllowCodeBlocksPolicy(config Config, fileConfig *FileConfig) bool {
	allowCodeBlocks := config.AllowCodeBlocks
	if fileConfig != nil && fileConfig.AllowCodeBlocks != nil {
		allowCodeBlocks = *fileConfig.AllowCodeBlocks
	}
	return allowCodeBlocks
}

func getTimeoutPolicy(config Config, fileConfig *FileConfig) int {
	timeout := config.Timeout

	if fileConfig != nil && fileConfig.Timeout != nil {
		timeout = *fileConfig.Timeout
	}

	return timeout
}

func getIgnoreExternalPolicy(config Config, fileConfig *FileConfig) bool {
	ignoreExternal := config.IgnoreExternal

	if fileConfig != nil && fileConfig.IgnoreExternal != nil {
		ignoreExternal = *fileConfig.IgnoreExternal
	}
	return ignoreExternal
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

func getInternalIgnorePolicy(filepath string, config Config, fileConfig *FileConfig) bool {
	internalIgnore := config.IgnoreExternal

	if isFileIgnored(filepath, config.FilesToIgnoreInternalLinksIn) {
		internalIgnore = true
	}

	if fileConfig != nil && fileConfig.IgnoreInternal != nil && *fileConfig.IgnoreInternal {
		internalIgnore = true
	}

	return internalIgnore
}

func isFileIgnored(filePath string, filesToIgnore []string) bool {
	for _, fileToIgnore := range filesToIgnore {
		if strings.HasPrefix(fileToIgnore, ".") {
			return checkIfFileIsInIgnorePath(fileToIgnore, filePath)
		} else {
			return checkIfFilePathContainsIgnoredDir(fileToIgnore, filePath)
		}
	}
	return false
}

func checkIfFileIsInIgnorePath(fileToIgnore, filePath string) bool {
	startingPath := path.Clean(fileToIgnore)
	cleanFilePath := path.Clean(filePath)

	return strings.HasPrefix(cleanFilePath, startingPath)
}

func checkIfFilePathContainsIgnoredDir(fileToIgnore, filePath string) bool {
	rootedFilePath := fmt.Sprintf(`/%s`, filePath)
	dirToIgnore := fmt.Sprintf(`/%s/`, fileToIgnore)
	return strings.Contains(rootedFilePath, dirToIgnore)
}
