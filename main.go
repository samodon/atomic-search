package main

import (
	"fmt"
	"samodon/search/indexing"
	"samodon/search/pkg"
	"samodon/search/searching"
	"strings"
)

func main() {
	directory := "/Users/samo/Library/Mobile Documents/com~apple~CloudDocs/Documents/Obsidian Vaults/Projects/Notes/Atomic/"
	documents := indexing.GetDocuments(directory)

	uniquewords, uniquetags := (pkg.Tokenize(documents))

	invertedindex, tagindex := indexing.CreateIndex(uniquewords, uniquetags, documents)
	indexing.Writeout(invertedindex, tagindex)

	fmt.Print("Index Created !")
	tmpterms := "Create an Array or Slice in Go"
	terms := strings.Fields(tmpterms)

	ranking, keywords := searching.GetRankingByFrequency(terms, invertedindex)
	results := pkg.ConvertMap(ranking)
	for i := range results {
		content, _, _ := pkg.ReadFile(results[i].NoteLocation)
		strarr := strings.Fields(content)
		results[i].Proximityscore = searching.Getproximityscore(keywords, strarr)
		// fmt.Println(results[i].searchscore)
		tags := searching.GetRankingbyTagInclusion(terms, tagindex, results[i].NoteLocation)
		if tags != nil {
			if len(tags) > 1 {
				for _, tag := range tags {
					results[i].Tags = append(results[i].Tags, tag)
				}
			} else {
				results[i].Tags = append(results[i].Tags, tags[0])
			}
		}
		// fmt.Println(results[i].Tags)

		results[i].Searchscore = float32(searching.CalculateSearchScore(results[i].Proximityscore, results[i].Frequency, 0.6, 0.3, len(results[i].Tags), 0.1))

	}
	sortedResults := searching.SortBySearchScore(results)

	for _, result := range sortedResults {
		content, _, _ := pkg.ReadFile(result.NoteLocation)
		fmt.Print(content)
		fmt.Print("")
		fmt.Println(result.Searchscore)
	}
}
