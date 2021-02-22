package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("External Links", func(t *testing.T) {
		dirPath := "test-markdowns"
		content, err := readMarkdown("test-markdowns/external_links.md")
		require.NoError(t, err)

		expected := Links{
			Link{
				AbsPath: "https://twitter.com",
				TypeOf:  ExternalLink,
			},
			Link{
				AbsPath: "https://github.com",
				TypeOf:  ExternalLink,
			},
			Link{
				AbsPath: "http://dont.exist.link.com",
				TypeOf:  ExternalLink,
			},
		}

		parser := &Parser{}
		result := parser.Links(content, dirPath)

		assert.Equal(t, expected, result)
	})

	t.Run("Internal Links", func(t *testing.T) {
		dirPath := "test-markdowns/sub_path"
		content, err := readMarkdown("test-markdowns/sub_path/internal_links.md")
		require.NoError(t, err)

		expected := Links{
			Link{
				AbsPath: "test-markdowns/external_links.md",
				RelPath: "../external_links.md",
				TypeOf:  InternalLink,
			},
			Link{
				AbsPath: "test-markdowns/sub_path/sub_sub_path/without_links.md",
				RelPath: "sub_sub_path/without_links.md",
				TypeOf:  InternalLink,
			},
			Link{
				AbsPath: "test-markdowns/sub_path/absolute_path.md",
				RelPath: "absolute_path.md",
				TypeOf:  InternalLink,
			},
			Link{
				AbsPath: "test-markdowns/sub_path/invalid.md",
				RelPath: "invalid.md",
				TypeOf:  InternalLink,
			},
		}

		parser := &Parser{}
		result := parser.Links(content, dirPath)

		assert.Equal(t, expected, result)
	})

	t.Run("Hash Internal Links", func(t *testing.T) {
		dirPath := "test-markdowns"
		content, err := readMarkdown("test-markdowns/hash_internal_links.md")
		require.NoError(t, err)

		expected := Links{
			Link{
				AbsPath: "https://github.com",
				TypeOf:  ExternalLink,
			},
			Link{
				AbsPath: "https://github.com",
				TypeOf:  ExternalLink,
			},
			Link{
				RelPath: "#first-header",
				TypeOf:  HashInternalLink,
			},
			Link{
				RelPath: "#second-header",
				TypeOf:  HashInternalLink,
			},
			Link{
				RelPath: "#third-header",
				TypeOf:  HashInternalLink,
			},
			Link{
				RelPath: "#header",
				TypeOf:  HashInternalLink,
			},
			Link{
				RelPath: "#header-with-block",
				TypeOf:  HashInternalLink,
			},
			Link{
				RelPath: "#header-with-link",
				TypeOf:  HashInternalLink,
			},
			Link{
				RelPath: "#very-strange-header-really-people-create-headers-look-like-this",
				TypeOf:  HashInternalLink,
			},
		}

		parser := &Parser{}
		result := parser.Links(content, dirPath)

		assert.Equal(t, expected, result)
	})

	t.Run("Headers", func(t *testing.T) {
		content, err := readMarkdown("test-markdowns/hash_internal_links.md")
		require.NoError(t, err)

		expected := []string{
			"First Header",
			"Second Header",
			"Third Header",
			"Header with link",
			"Header with block",
			"Very strange header really people create headers look like this",
			"Links",
		}

		parser := &Parser{}
		result := parser.Headers(content)

		assert.Equal(t, expected, result)
	})

	t.Run("Absolute Internal Path", func(t *testing.T) {
		//TODO: we have to initialize this global variable, because this test don't pass when launched individually
		//This variable was initialized in `config_file_test.go`.
		// In future such dependency should be removed.
		SetBasePath("test-markdowns", false)
		dirPath := "test-markdowns/sub_path"
		content, err := readMarkdown("test-markdowns/sub_path/absolute_path.md")
		require.NoError(t, err)

		expected := Links{
			Link{
				AbsPath: "test-markdowns/external_links.md",
				RelPath: "/external_links.md",
				TypeOf:  InternalLink,
			},
		}

		parser := &Parser{}
		result := parser.Links(content, dirPath)

		assert.Equal(t, expected, result)
	})
}
