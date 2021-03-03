package pkg

import (
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

type Headers []string

type File struct {
	RelPath string `yaml:"path"`
	AbsPath string
	DirPath string
	Content string
	Links   Links `yaml:"links"`
	Headers Headers
	Status  bool
	Config  *FileConfig `yaml:"config"`
	Stats   *FileStats
	parser  *Parser
	valid   *Validator
}

func NewFile(filePath string, fileLinks Links, config FileConfig) (*File, error) {
	if match, _ := regexp.MatchString(`.md$`, filePath); !match {
		return nil, errors.New("The specified file isn't a markdown file")
	}

	absPath, _ := filepath.Abs(filePath)
	if err := fileExists(absPath); err != nil {
		return nil, err
	}
	content, err := readMarkdown(absPath)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	waiter := NewWaiter(config.Backoff)

	return &File{
		RelPath: filePath,
		AbsPath: absPath,
		DirPath: filepath.Dir(filePath),
		Content: content,
		Links:   fileLinks,
		Config:  &config,
		parser:  &Parser{},
		valid:   NewValidator(client, waiter),
	}, nil
}

func (f *File) Run() {
	f.ExtractLinks().
		ExtractHeaders().
		ValidateLinks().
		ExtractStats()
}

func (f *File) ExtractLinks() *File {
	externalLinksToIgnore, internalLinksToIgnore := []string{}, []string{}
	if f.Config != nil {
		externalLinksToIgnore = f.Config.ExternalLinksToIgnore
		internalLinksToIgnore = f.Config.InternalLinksToIgnore
	}

	content := f.Content
	if f.Config != nil && f.Config.AllowCodeBlocks != nil && !*f.Config.AllowCodeBlocks {
		content = removeCodeBlocks(content)
	}

	basePath := ""
	if f.Config != nil {
		basePath = f.Config.BasePath
	}
	f.Links = f.parser.
		Links(basePath, content, f.DirPath).
		AppendConfig(f).
		RemoveIgnoredLinks(externalLinksToIgnore, internalLinksToIgnore).
		Filter(func(link Link) bool {
			if f.Config != nil && f.Config.IgnoreInternal != nil && *f.Config.IgnoreInternal && (link.TypeOf == HashInternalLink || link.TypeOf == InternalLink) {
				return false
			}

			if f.Config != nil && f.Config.IgnoreExternal != nil && *f.Config.IgnoreExternal && link.TypeOf == ExternalLink {
				return false
			}

			return true
		})
	return f
}

func (f *File) ExtractHeaders() *File {
	f.Headers = f.parser.Headers(f.Content)
	return f
}

func (f *File) ValidateLinks() *File {
	f.Links = f.valid.Links(f.Links, f.Headers)
	f.Status = f.Links.CheckStatus()
	return f
}

func (f *File) ExtractStats() *File {
	f.Stats = NewFileStats(f)
	return f
}

func (f *File) WriteStats() *File {
	writeStats(f)
	return f
}

func (f *File) Summary() *File {
	summaryOfFile(f)
	return f
}
