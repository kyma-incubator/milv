package pkg

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	codeBlockPattern = `(?m)^(.*\x60{3}).*\n(.*|\n)+?\n(.*\x60{3})$`
)

func fileExists(file string) error {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func headerExists(link string, headers []string) bool {
	link = strings.TrimPrefix(link, "#")
	for _, header := range headers {
		if link == strings.ToLower(strings.Replace(header, " ", "-", -1)) {
			return true
		}
	}
	return false
}

func unique(elements []string) []string {
	encountered := map[string]bool{}
	for v := range elements {
		encountered[elements[v]] = true
	}

	result := []string{}
	for key := range encountered {
		if string(key) != "" {
			result = append(result, key)
		}
	}
	return result
}

func removeIgnoredFiles(filePaths, filesToIgnore []string) []string {
	var newFilePaths []string
	for _, file := range filePaths {
		exists := false
		for _, fileToIgnore := range filesToIgnore {
			if strings.Contains(file, fileToIgnore) {
				exists = true
				break
			}
		}
		if !exists {
			newFilePaths = append(newFilePaths, file)
		}
	}
	return newFilePaths
}

func readMarkdown(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func removeCodeBlocks(markdown string) string {
	re := regexp.MustCompile(codeBlockPattern)
	return re.ReplaceAllString(markdown, "")
}

func contains(slice []string, value string) bool {
	for _, el := range slice {
		if value == el {
			return true
		}
	}
	return false
}
