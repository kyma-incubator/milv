package pkg

type LinkConfig struct {
	Timeout        *int  `yaml:"timeout"`
	RequestRepeats *int  `yaml:"request-repeats"`
	AllowRedirect  *bool `yaml:"allow-redirect"`
}

func NewLinkConfig(link Link, file *File) *LinkConfig {
	if file.Config != nil {
		for _, linkFile := range file.Links {
			if (link.RelPath == linkFile.RelPath || link.AbsPath == linkFile.RelPath) && linkFile.Config != nil {
				var timeout *int
				if linkFile.Config.Timeout != nil {
					timeout = linkFile.Config.Timeout
				} else {
					timeout = file.Config.Timeout
				}

				var requestRepeats *int
				if linkFile.Config.RequestRepeats != nil {
					requestRepeats = linkFile.Config.RequestRepeats
				} else {
					requestRepeats = file.Config.RequestRepeats
				}

				var allowRedirect *bool
				if linkFile.Config.AllowRedirect != nil {
					allowRedirect = linkFile.Config.AllowRedirect
				} else {
					allowRedirect = file.Config.AllowRedirect
				}

				return &LinkConfig{
					Timeout:        timeout,
					RequestRepeats: requestRepeats,
					AllowRedirect:  allowRedirect,
				}
			}
		}
		return &LinkConfig{
			Timeout:        file.Config.Timeout,
			RequestRepeats: file.Config.RequestRepeats,
			AllowRedirect:  file.Config.AllowRedirect,
		}
	}
	return nil
}
