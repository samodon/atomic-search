package pkg

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

func ReadFile(path string) (string, []string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "Error occured", nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{
				Mode: frontmatter.SetMetadata,
			},
		),
	)

	root := md.Parser().Parse(text.NewReader(content))
	doc := root.OwnerDocument()
	meta := doc.Meta()
	var tags []string
	if meta["tags"] != nil {
		strtags := fmt.Sprint(meta["tags"])
		strtags = strings.Trim(strtags, "[]")
		tags = strings.Fields(strtags)
	}

	return stripPunctuation(string(content)), tags, err
}

/*
* strips punctuations from a string
 */
func stripPunctuation(text string) string {
	var sb strings.Builder
	for _, c := range text {
		if !unicode.IsPunct(c) && !unicode.IsDigit(c) && !unicode.IsNumber(c) && !unicode.IsSymbol(c) {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

/*
Returns an a list of tokenized words from an array a documents
*/
func Tokenize(documents []Document) ([]string, []string) {
	var checkedwords []string
	var uniquewords []string
	var uniquetags []string
	for _, document := range documents {
		if document.Tags != nil {
			for _, tag := range document.Tags {
				if !Contains(uniquetags, tag) {
					uniquetags = append(uniquetags, tag)
				}
			}
		}
		words := strings.Fields(document.Content)
		for _, word := range words {
			if !Contains(checkedwords, word) {
				uniquewords = append(uniquewords, word)
				checkedwords = append(checkedwords, word)
			}
		}
	}
	sort.Strings(uniquewords)
	return uniquewords, uniquetags
}

/*
Checks if there is an item in an an array of strings
*/
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func ConvertMap(m map[string]int) []Ranking {
	var kvSlice []Ranking
	for k, v := range m {
		kvSlice = append(kvSlice, Ranking{k, v, 0, nil, 0, nil})
	}
	// TODO:
	// - [] sort by overall score in the end

	return kvSlice
}
