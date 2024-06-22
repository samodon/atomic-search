package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

type document struct {
	documentloc string
	content     string
	tags        []string
}

func getdocuments(dirname string) []document {
	var documents []document
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "md") {
			var tmpdocument document
			tmpdocument.documentloc = path
			tmpdocument.content, tmpdocument.tags, err = readfile(tmpdocument.documentloc)
			if err != nil {
				fmt.Print(err)
				return nil
			}
			documents = append(documents, tmpdocument)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return documents
}

// TODO
// modify this to return date, tags, topic, etc
// look into processing this in chunks for improved efficiency
func readfile(path string) (string, []string, error) {
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

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// TODO
// Sort the words in alphabetical order
// only return the root words
func tokenize(documents []document) ([]string, []string) {
	var checkedwords []string
	var uniquewords []string
	var uniquetags []string
	for _, document := range documents {
		if document.tags != nil {
			for _, tag := range document.tags {
				if !contains(uniquetags, tag) {
					uniquetags = append(uniquetags, tag)
				}
			}
		}
		words := strings.Fields(document.content)
		for _, word := range words {
			if !contains(checkedwords, word) {
				uniquewords = append(uniquewords, word)
				checkedwords = append(checkedwords, word)
			}
		}
	}
	sort.Strings(uniquewords)
	return uniquewords, uniquetags
}

func createindex(words []string, tags []string, documents []document) (map[string][]string, map[string][]string) {
	wordindex := make(map[string][]string)
	fmt.Printf("Indexing...")
	for _, word := range words {
		for _, document := range documents {
			if contains(strings.Fields(document.content), word) {
				stemmed, _ := snowball.Stem(word, "english", true)
				wordindex[stemmed] = append(wordindex[stemmed], document.documentloc)
			}
		}
	}
	tagsindex := make(map[string][]string)
	for _, tag := range tags {
		for _, document := range documents {
			if contains(document.tags, tag) {
				// fmt.Println(tag)
				// fmt.Println(document.documentloc)
				tagsindex[tag] = append(tagsindex[tag], document.documentloc)
			}
		}
	}
	return wordindex, tagsindex
}

func stripPunctuation(text string) string {
	var sb strings.Builder
	for _, c := range text {
		if !unicode.IsPunct(c) && !unicode.IsDigit(c) && !unicode.IsNumber(c) && !unicode.IsSymbol(c) {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

func getresults(invertedindex map[string][]string) [][]string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter search term:")
	input, _ := reader.ReadString('\n')

	var results [][]string
	wordlist := strings.Fields(input)

	for _, word := range wordlist {
		word = english.Stem(word, true)

		if invertedindex[word] != nil {
			path := invertedindex[word]
			results = append(results, path)
		}

	}
	return results
}

func writeout(invertedindex map[string][]string, tagindex map[string][]string) {
	file, err := os.Create("wordindex")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(invertedindex)
	if err != nil {
		panic(err)
	}
	file, err = os.Create("tagindex")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder = json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(tagindex)
	if err != nil {
		panic(err)
	}
}

func main() {
	directory := "/Users/samo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Obsidian Vaults/Projects/Notes/Atomic/"

	documents := getdocuments(directory)

	uniquewords, uniquetags := (tokenize(documents))

	// fmt.Println(uniquetags)
	invertedindex, tagindex := createindex(uniquewords, uniquetags, documents)
	// fmt.Print(tagindex)
	writeout(invertedindex, tagindex)
	// tags inverted index
}
