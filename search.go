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

/*
Accepts directoy name and returns an array of type documet
*/
func getDocuments(dirname string) []document {
	var documents []document
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "md") {
			var tmpdocument document
			tmpdocument.documentloc = path
			tmpdocument.content, tmpdocument.tags, err = readFile(tmpdocument.documentloc)
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

// TODO:
//- [ ] look into processing this in chunks for improved efficiency, if necessary
//- [x] return tags in addition to the document content

/*
readFile function returns a string containing the content of a document and a string slice containing the tags for that document
*/
func readFile(path string) (string, []string, error) {
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
Checks if there is an item in an an array of strings
*/
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

/*
Returns an a list of tokenized words from an array a documents
*/
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

/*
Returns both invereted indices as hashmaps, requires []string, []string, []document
*/
func createIndex(words []string, tags []string, documents []document) (map[string][]string, map[string][]string) {
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

func getResults(invertedindex map[string][]string) [][]string {
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

/*
Writes out both indices as json
*/
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

/* TODO:
- [x] after getting the complete rankings sort them
- [x] stem all the words before they are searched
*/

func getRankingByFrequency(terms []string, index map[string][]string) (map[string]int, []string) {
	ranking := make(map[string]int)
	var keywords []string
	for _, term := range terms {

		term, _ = snowball.Stem(term, "english", true)
		// fmt.Println(term)
		if index[term] != nil {
			//			fmt.Println(len(index[term]))
			for _, loc := range index[term] {
				ranking[loc] = ranking[loc] + 1
				keywords = append(keywords, term)
			}
		}
	}

	return ranking, keywords
}

func convertMap(m map[string]int) []Ranking {
	var kvSlice []Ranking
	for k, v := range m {
		kvSlice = append(kvSlice, Ranking{k, v, 0, nil, 0})
	}
	// TODO:
	// - [] sort by overall score in the end

	return kvSlice
}

func sortBySearchScore(rankings []Ranking) []Ranking {
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].searchscore > rankings[j].searchscore
	})
	return rankings
}

// TODO:
// Propably make the get the proximity score from the average proximity of the elements in that content
type Ranking struct {
	index          string
	frequency      int
	proximityscore float32
	keywords       []string
	searchscore    float32
}

func getDistance(keyword string, content []string) int {
	keywordStem := english.Stem(keyword, true)
	position := -1

	for i, word := range content {
		if keywordStem == english.Stem(word, true) {
			position = i
			break
		}
	}

	if position == -1 || position == len(content)-1 {
		return -1
	}

	return len(content) - position - 1
}

func getproximityscore(words []string, content []string) float32 {
	var distances []int
	totalDistance := 0
	for _, word := range words {
		tmpdist := getDistance(word, content)
		totalDistance = totalDistance + tmpdist
		distances = append(distances, tmpdist)
	}
	return float32(totalDistance) / float32(len(distances))
}

func calculateSearchScore(proximityScore float32, frequency int, weightProximity float32, weightFrequency float32) float32 {
	var normalizedProximity float32 = 0.0
	if proximityScore != 0 {
		normalizedProximity = 1.0 / proximityScore
	}
	searchScore := (weightProximity * normalizedProximity) + (weightFrequency * float32(frequency))
	return searchScore
}

func main() {
	directory := "/Users/samo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Obsidian Vaults/Projects/Notes/Atomic/"
	// directory := "."
	documents := getDocuments(directory)

	uniquewords, uniquetags := (tokenize(documents))

	invertedindex, tagindex := createIndex(uniquewords, uniquetags, documents)
	writeout(invertedindex, tagindex)

	tmpterms := "Altering a table in SQL"
	terms := strings.Fields(tmpterms)

	ranking, keywords := getRankingByFrequency(terms, invertedindex)
	results := convertMap(ranking)

	for i := range results {
		content, _, _ := readFile(results[i].index)
		strarr := strings.Fields(content)
		results[i].proximityscore = getproximityscore(keywords, strarr)
		results[i].searchscore = float32(calculateSearchScore(results[i].proximityscore, results[i].frequency, 0.7, 0.3))
		// fmt.Println(results[i].searchscore)
	}
	sortedResults := sortBySearchScore(results)

	for _, result := range sortedResults {
		content, _, _ := readFile(result.index)
		fmt.Print(content)
		fmt.Print("")
		fmt.Println(result.searchscore)
	}
}
