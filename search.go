package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type document struct {
	documentloc string
	content     string
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

// this implementation requires the  content first
// this should return a unique list of words
func tokenize(documents []document) []string {
	var checkedwords []string
	var uniquewords []string
	for _, document := range documents {
		words := strings.Fields(document.content)
		for _, word := range words {
			if !contains(uniquewords, word) {
				uniquewords = append(uniquewords, word)
				checkedwords = append(checkedwords, word)
			}
		}
	}
	return uniquewords
}

func createindex(words []string, documents []document) map[string]string {
	invertedindex := make(map[string]string)
	for _, word := range words {
		for _, document := range documents {
			if contains(strings.Fields(document.content), word) {
				invertedindex[word] = document.documentloc
			}
		}
	}
	return invertedindex
}

func main() {
	directory := "."

	documents := getdocuments(directory)
	uniquewords := (tokenize(documents))
	invertedindex := createindex(uniquewords, documents)
	fmt.Print(invertedindex)
}
