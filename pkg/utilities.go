package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

func ParseMdData(path string) (string, map[string]interface{}, error) {
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
	// newmeta := fmt.Sprint(meta)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(content, &buf, parser.WithContext(context)); err != nil {
		return "Error occurred", nil, err
	}
	return buf.String(), meta, err
}

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
	// stringyy := root.Lines()
	// fmt.Println(stringyy)
	var tags []string
	if meta["tags"] != nil {
		strtags := fmt.Sprint(meta["tags"])
		strtags = strings.Trim(strtags, "[]")
		tags = strings.Fields(strtags)
	}

	return stripPunctuation(string(content)), tags, err
}

func JSONToMap(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result map[string][]string

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
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
		document.Content = RemoveWords(document.Content)
		words := strings.Fields(document.Content)
		words = stemFiles(words)
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

func stemFiles(searchTerms []string) []string {
	for _, term := range searchTerms {
		term = english.Stem(term, true)
		// fmt.Println(term)
	}
	return searchTerms
}

// Removes commond words in the string, helper to tokenizer
func RemoveWords(input string) string {
	wordsToRemove := []string{"and", "my", "the", "in", "as", "but", "like", "do", "a", "has", "it"}
	lowerInput := strings.ToLower(input)

	for _, word := range wordsToRemove {
		wordLower := strings.ToLower(word)
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(wordLower) + `\b`)
		lowerInput = re.ReplaceAllString(lowerInput, "")
	}
	lowerInput = strings.Join(strings.Fields(lowerInput), " ")
	return lowerInput
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
