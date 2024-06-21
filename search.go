package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kljensen/snowball"
)

type document struct {
	documentloc string
	content     string
}

func test(doc document) {
}

func getdocuments(dirname string) []document {
	var documents []document
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "txt") {
			var tmpdocument document
			tmpdocument.documentloc = path
			tmpdocument.content, err = readfile(tmpdocument.documentloc)
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
func readfile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "Error occured", err
	}
	return string(content), err
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
func tokenize(documents []document) []string {
	var checkedwords []string
	var uniquewords []string
	for _, document := range documents {
		words := strings.Fields(document.content)
		for _, word := range words {
			if !contains(checkedwords, word) {
				uniquewords = append(uniquewords, word)
				checkedwords = append(checkedwords, word)
			}
		}
	}
	sort.Strings(uniquewords)
	return uniquewords
}

func createindex(words []string, documents []document) map[string][]string {
	invertedindex := make(map[string][]string)
	for _, word := range words {
		for _, document := range documents {
			if contains(strings.Fields(document.content), word) {
				stemmed, _ := snowball.Stem(word, "english", true)
				invertedindex[stemmed] = append(invertedindex[stemmed], document.documentloc)
			}
		}
	}
	return invertedindex
}

func getresults(words []string) []document {
	words = append(words, "test")
	return nil
}

func main() {
	directory := "."
	reader := bufio.NewReader(os.Stdin)

	documents := getdocuments(directory)

	uniquewords := (tokenize(documents))
	invertedindex := createindex(uniquewords, documents)
	// fmt.Println(invertedindex)

	fmt.Print("Enter search term:")
	input, _ := reader.ReadString('\n')
	input, _ = snowball.Stem(input, "english", true)
	input = strings.TrimSpace(input)
	print(input)
	fmt.Println(invertedindex[input])
}
