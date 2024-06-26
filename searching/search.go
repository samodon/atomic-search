package searching

import (
	"samodon/search/pkg"
	"sort"

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
	var normalizedProximity float32 = 0.0
	if proximityScore > 0 {
		normalizedProximity = 1.0 / proximityScore
	}
	searchScore := (weightProximity * normalizedProximity) + (weightFrequency * float32(frequency)) + float32(numoftags)*weightTag
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
