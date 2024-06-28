package searching

import (
	"samodon/search/pkg"
	"sort"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
)

func SortBySearchScore(rankings []pkg.Ranking) []pkg.Ranking {
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Searchscore > rankings[j].Searchscore
	})
	return rankings
}

/* TODO:
* If a file contains a tag that is in the search term return a bool stating it contains
* the tag, this should be signifant in determining the rankings
 */
func GetRankingbyTagInclusion(terms []string, tagindex map[string][]string, location string) []string {
	var tags []string
	for _, term := range terms {
		// fmt.Println(tagindex[term])
		if tagindex[term] == nil {
			continue
		} else {
			for _, loc := range tagindex[term] {
				if location == loc {
					tags = append(tags, term)
				}
			}
		}
	}
	return tags
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

func Getproximityscore(words []string, content []string) float32 {
	var distances []int
	totalDistance := 0
	for _, word := range words {
		tmpdist := getDistance(word, content)
		totalDistance = totalDistance + tmpdist
		distances = append(distances, tmpdist)
	}
	return float32(totalDistance) / float32(len(distances))
}

func CalculateSearchScore(proximityScore float32, frequency int, weightProximity float32, weightFrequency float32, numoftags int, weightTag float32) float32 {
	// var normalizedProximity float32 = 0.0
	// if proximityScore != 0 {
	// 	normalizedProximity = 1.0 / proximityScore
	// }
	searchScore := (weightProximity * proximityScore) + (weightFrequency * float32(frequency)) + float32(numoftags)*weightTag
	return searchScore
}

func GetRankingByFrequency(terms []string, index map[string][]string) (map[string]int, []string) {
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

// Stemmmed files

func stemFiles(searchTerms []string) []string {
	for _, term := range searchTerms {
		term = english.Stem(term, true)
		// fmt.Println(term)
	}
	return searchTerms
}

//export GetSearchRanking
func GetSearchRanking(searchTerms string, mapdir string, tagdir string) []pkg.Ranking {
	// fmt.Println(searchTerms)
	termsarr := strings.Fields(searchTerms)
	// fmt.Println(termsarr)
	termsarr = stemFiles(termsarr)
	contentIndex, _ := pkg.JSONToMap(mapdir)
	tagindex, _ := pkg.JSONToMap(tagdir)
	ranking, keywords := GetRankingByFrequency(termsarr, contentIndex)
	results := pkg.ConvertMap(ranking)
	for i := range results {
		content, _, _ := pkg.ReadFile(results[i].NoteLocation)
		strarr := strings.Fields(content)
		results[i].Proximityscore = Getproximityscore(keywords, strarr)
		// fmt.Println(results[i].searchscore)
		tags := GetRankingbyTagInclusion(termsarr, tagindex, results[i].NoteLocation)
		if tags != nil {
			if len(tags) > 1 {
				for _, tag := range tags {
					results[i].Tags = append(results[i].Tags, tag)
				}
			} else {
				results[i].Tags = append(results[i].Tags, tags[0])
			}
		}
		results[i].Searchscore = float32(CalculateSearchScore(results[i].Proximityscore, results[i].Frequency, 0.6, 0.3, len(results[i].Tags), 0.1))
	}

	return SortBySearchScore(results)
}
