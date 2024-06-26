package indexing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"samodon/search/pkg"
	"strings"

	"github.com/kljensen/snowball"
)

func GetDocuments(dirname string) []pkg.Document {
	var documents []pkg.Document
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "md") {
			var tmpdocument pkg.Document
			tmpdocument.DocumentLocation = path
			tmpdocument.Content, tmpdocument.Tags, err = pkg.ReadFile(tmpdocument.DocumentLocation)
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

/*
Writes out both indices as json
*/
func Writeout(invertedindex map[string][]string, tagindex map[string][]string) {
	file, err := os.Create("index/wordindex.json")
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
	file, err = os.Create("index/tagindex.json")
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

/*
Returns both invereted indices as hashmaps, requires []string, []string, []document
*/
func CreateIndex(words []string, tags []string, documents []pkg.Document) (map[string][]string, map[string][]string) {
	wordindex := make(map[string][]string)
	fmt.Printf("Indexing...")
	for _, word := range words {
		for _, document := range documents {
			if pkg.Contains(strings.Fields(document.Content), word) {
				stemmed, _ := snowball.Stem(word, "english", true)
				wordindex[stemmed] = append(wordindex[stemmed], document.DocumentLocation)
			}
		}
	}
	tagsindex := make(map[string][]string)
	for _, tag := range tags {
		for _, document := range documents {
			if pkg.Contains(document.Tags, tag) {
				// fmt.Println(tag)
				// fmt.Println(document.documentloc)
				tagsindex[tag] = append(tagsindex[tag], document.DocumentLocation)
			}
		}
	}
	return wordindex, tagsindex
}
