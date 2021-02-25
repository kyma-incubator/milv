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

func NewFileConfig(filePath string, config *Config) *FileConfig {
	if config != nil {
		for _, file := range config.Files {
			if filePath == file.RelPath && file.Config != nil {
				var timeout *int
				if file.Config.Timeout != nil {
					timeout = file.Config.Timeout
				} else {
					timeout = &config.Timeout
				}

				var requestRepeats *int8
				if file.Config.Timeout != nil {
					requestRepeats = file.Config.RequestRepeats
				} else {
					requestRepeats = &config.RequestRepeats
				}

				var allowRedirect, allowCodeBlocks, ignoreExternal *bool

				if file.Config.AllowCodeBlocks != nil {
					allowCodeBlocks = file.Config.AllowCodeBlocks
				} else {
					allowCodeBlocks = &config.AllowCodeBlocks
				}
				if file.Config.AllowRedirect != nil {
					allowRedirect = file.Config.AllowRedirect
				} else {
					allowRedirect = &config.AllowRedirect
				}

				ignoreInternal := applyInternalIgnorePolicy(filePath, *config, file.Config)

				if file.Config.IgnoreExternal != nil {
					ignoreExternal = file.Config.IgnoreExternal
				} else {
					ignoreExternal = &config.IgnoreExternal
				}

				return &FileConfig{
					BasePath:              config.BasePath,
					ExternalLinksToIgnore: unique(append(config.ExternalLinksToIgnore, file.Config.ExternalLinksToIgnore...)),
					InternalLinksToIgnore: unique(append(config.InternalLinksToIgnore, file.Config.InternalLinksToIgnore...)),
					Timeout:               timeout,
					RequestRepeats:        requestRepeats,
					AllowRedirect:         allowRedirect,
					AllowCodeBlocks:       allowCodeBlocks,
					IgnoreExternal:        ignoreExternal,
					IgnoreInternal:        &ignoreInternal,
				}
			}
		}

		ignoreInternal := applyInternalIgnorePolicy(filePath, *config, nil)

		return &FileConfig{
			BasePath:              config.BasePath,
			ExternalLinksToIgnore: config.ExternalLinksToIgnore,
			InternalLinksToIgnore: config.InternalLinksToIgnore,
			Timeout:               &config.Timeout,
			RequestRepeats:        &config.RequestRepeats,
			AllowRedirect:         &config.AllowRedirect,
			AllowCodeBlocks:       &config.AllowCodeBlocks,
			IgnoreExternal:        &config.IgnoreExternal,
			IgnoreInternal:        &ignoreInternal,
		}
	}
	return nil
}

func findFileConfig(filePath string, files Files) (bool, *File) {
	for _, file := range files {
		if filePath == file.RelPath && file.Config != nil {
			return true, file
		}
	}
	return false, nil
}

func applyInternalIgnorePolicy(filepath string, config Config, fileConfig *FileConfig) bool {
	var internalIgnore = false

	if config.IgnoreInternal {
		internalIgnore = true
	}

	//check if file is covered by ignore internal links in policy
	if isFileIgnored(filepath, config.FilesToIgnoreInternalLinksIn) {
		internalIgnore = true
	}

	//this is JavaGoScript
	//apply the most important file specific policy
	if fileConfig != nil && fileConfig.IgnoreInternal != nil && *fileConfig.IgnoreInternal {
		internalIgnore = true
	}

	return internalIgnore
}

func isFileIgnored(filePath string, filesToIgnore []string) bool {
	for _, fileToIgnore := range filesToIgnore {
		if strings.HasPrefix(fileToIgnore, ".") {
			//fileToIgnore := path.Join(basePath, fileToIgnore)
			//block concrete path or file
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
