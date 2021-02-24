package pkg

type Files []*File

func NewFiles(filePaths []string, config *Config) (Files, error) {
	var files Files

	filePaths = removeIgnoredFiles(filePaths, config.FilesToIgnore)
	for _, filePath := range filePaths {
		file, err := NewFile(filePath, NewLinks(filePath, config), NewFileConfig(filePath, config))
		if err != nil {
			return Files{}, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (f Files) Run(verbose bool) {
	for _, file := range f {
		file.Run()
		if verbose {
			file.WriteStats()
		}
	}
}

func (f Files) Summary() bool {
	return summaryOfFiles(f)
}
