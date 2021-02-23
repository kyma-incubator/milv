package pkg

type FileConfig struct {
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

				var allowRedirect, allowCodeBlocks, ignoreExternal, ignoreInternal *bool
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
				if file.Config.IgnoreExternal != nil {
					ignoreExternal = file.Config.IgnoreExternal
				} else {
					ignoreExternal = &config.IgnoreExternal
				}
				if file.Config.IgnoreInternal != nil {
					ignoreInternal = file.Config.IgnoreInternal
				} else {
					ignoreInternal = &config.IgnoreInternal
				}

				return &FileConfig{
					ExternalLinksToIgnore: unique(append(config.ExternalLinksToIgnore, file.Config.ExternalLinksToIgnore...)),
					InternalLinksToIgnore: unique(append(config.InternalLinksToIgnore, file.Config.InternalLinksToIgnore...)),
					Timeout:               timeout,
					RequestRepeats:        requestRepeats,
					AllowRedirect:         allowRedirect,
					AllowCodeBlocks:       allowCodeBlocks,
					IgnoreExternal:        ignoreExternal,
					IgnoreInternal:        ignoreInternal,
				}
			}
		}
		return &FileConfig{
			ExternalLinksToIgnore: config.ExternalLinksToIgnore,
			InternalLinksToIgnore: config.InternalLinksToIgnore,
			Timeout:               &config.Timeout,
			RequestRepeats:        &config.RequestRepeats,
			AllowRedirect:         &config.AllowRedirect,
			AllowCodeBlocks:       &config.AllowCodeBlocks,
			IgnoreExternal:        &config.IgnoreExternal,
			IgnoreInternal:        &config.IgnoreInternal,
		}
	}
	return nil
}
