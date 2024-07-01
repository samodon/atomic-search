package indexing

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
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
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	os.Mkdir(filepath.Join(usr.HomeDir, "/index/"), os.FileMode(0755))
	file, err := os.Create(filepath.Join(usr.HomeDir, "/index/wordindex.json"))
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
	file, err = os.Create(filepath.Join(usr.HomeDir, "/index/tagindex.json"))
	if err != nil {
		log.Panic(err)
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
//export CreateIndex
func CreateIndex(directory string) /* (map[string][]string, map[string][]string)*/ {
	documents := GetDocuments(directory)
	words, tags := (pkg.Tokenize(documents))
	wordindex := make(map[string][]string)
	fmt.Println("Indexing...")
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
	Writeout(wordindex, tagsindex)
	// return wordindex, tagsindex
}
