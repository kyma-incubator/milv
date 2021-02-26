package pkg

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/kyma-incubator/milv/cli"
)

type Config struct {
	BasePath                     string
	Files                        []File   `yaml:"files"`
	ExternalLinksToIgnore        []string `yaml:"external-links-to-ignore"`
	InternalLinksToIgnore        []string `yaml:"internal-links-to-ignore"`
	FilesToIgnore                []string `yaml:"files-to-ignore"`
	FilesToIgnoreInternalLinksIn []string `yaml:"files-to-ignore-internal-links-in"`
	Timeout                      int      `yaml:"timeout"`
	RequestRepeats               int      `yaml:"request-repeats"`
	AllowRedirect                bool     `yaml:"allow-redirect"`
	AllowCodeBlocks              bool     `yaml:"allow-code-blocks"`
	IgnoreExternal               bool     `yaml:"ignore-external"`
	IgnoreInternal               bool     `yaml:"ignore-internal"`
}

func NewConfig(commands cli.Commands) (*Config, error) {
	config := &Config{}

	err := fileExists(commands.ConfigFile)
	if commands.ConfigFile != "milv.config.yaml" && err != nil {
		return nil, err
	}
	if err == nil {
		yamlFile, err := ioutil.ReadFile(commands.ConfigFile)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(yamlFile, config)
		if err != nil {
			return nil, err
		}
	}
	return config.combine(commands), nil
}

func (c *Config) combine(commands cli.Commands) *Config {
	var timeout int
	if commands.FlagsSet["timeout"] {
		timeout = commands.Timeout
	} else {
		timeout = c.Timeout
	}

	var requestRepeats int
	if commands.FlagsSet["request-repeats"] {
		requestRepeats = commands.RequestRepeats
	} else {
		requestRepeats = c.RequestRepeats
	}

	var allowRedirect, allowCodeBlocks, ignoreExternal, ignoreInternal bool
	if commands.FlagsSet["allow-redirect"] {
		allowRedirect = commands.AllowRedirect
	} else {
		allowRedirect = c.AllowRedirect
	}
	if commands.FlagsSet["allow-code-blocks"] {
		allowCodeBlocks = commands.AllowCodeBlocks
	} else {
		allowCodeBlocks = c.AllowCodeBlocks
	}
	if commands.FlagsSet["ignore-external"] {
		ignoreExternal = commands.IgnoreExternal
	} else {
		ignoreExternal = c.IgnoreExternal
	}
	if commands.FlagsSet["ignore-internal"] {
		ignoreInternal = commands.IgnoreInternal
	} else {
		ignoreInternal = c.IgnoreInternal
	}

	return &Config{
		BasePath:                     commands.BasePath,
		Files:                        c.Files,
		ExternalLinksToIgnore:        unique(append(c.ExternalLinksToIgnore, commands.ExternalLinksToIgnore...)),
		InternalLinksToIgnore:        unique(append(c.InternalLinksToIgnore, commands.InternalLinksToIgnore...)),
		FilesToIgnoreInternalLinksIn: unique(append(c.FilesToIgnoreInternalLinksIn, commands.FilesToIgnoreInternalLinksIn...)),
		FilesToIgnore:                unique(append(c.FilesToIgnore, commands.FilesToIgnore...)),
		Timeout:                      timeout,
		RequestRepeats:               requestRepeats,
		AllowRedirect:                allowRedirect,
		AllowCodeBlocks:              allowCodeBlocks,
		IgnoreExternal:               ignoreExternal,
		IgnoreInternal:               ignoreInternal,
	}
}
