package main

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func main() {
	title := "GetSpace Free - Delete Duplicate Photo from Device Storage"
	keywords := "delete, clean, disk, master, manager, cleaner, memory, cache, camera, boost, aid, battery, saver, system, gallery"
	queries := getQueries(title, keywords)
	var bestQueries []string
	for i, query := range queries {
		words := strings.Split(query, " ")
		fmt.Println(query, fmt.Sprintf("(%d / %d)", i, len(queries)))
		term := strings.Join(words, "+")
		url := fmt.Sprintf("https://itunes.apple.com/search?term=%s&country=us&entity=software", term)

		data, err := readBody(url)
		if err != nil {
			panic(err)
		}
		var pageResp PageResp
		err = json.Unmarshal(data, &pageResp)
		if err != nil {
			panic(err)
		}
		if len(pageResp.Results) < 10 {
			continue
		}
		for i, entry := range pageResp.Results {
			if i > 30 {
				break
			}
			if entry.TrackCensoredName == title {
				bestQueries = append(bestQueries, fmt.Sprintf("%2v / %2v  %s", i, len(pageResp.Results), query))
				break
			}
		}
	}
	fmt.Println("RESULTS:")
	for _, query := range bestQueries {
		fmt.Println(query)	
	}
}

type PageRespEntry struct {
	TrackCensoredName string `json:"trackCensoredName"`
}

type PageResp struct {
	Results []PageRespEntry `json:"results"`
}

func readBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil	
}

func getQueries(title string, keywords string) []string {
	title = strings.ToLower(title)
	keywords = strings.ToLower(keywords)
	for _, str := range []string{".", ",", ":", " -", "(", ")"} {
		title = strings.Replace(title, str, "", -1)	
		keywords = strings.Replace(keywords, str, "", -1)	
	}
	allKeywords := append(strings.Split(title, " "), strings.Split(keywords, " ")...)
	var anagrams []string
	anagrams = append(anagrams, getAnagrams(allKeywords, 1, 0)...)
	// anagrams = append(anagrams, getAnagrams(allKeywords, 2, 0)...)
	return anagrams
}

func getAnagrams(keywords []string, limit int, reqDepth int) []string {
	if reqDepth == limit - 1 || len(keywords) == 1 {
		return keywords
	}
	anagrams := make([]string, 0)
	for i, word := range keywords {
		var words []string
		words = append(words, keywords[:i]...)
		words = append(words, keywords[i + 1:]...)
		for _, queryRest := range getAnagrams(words, limit, reqDepth + 1) {
			anagrams = append(anagrams, []string{word + " " + queryRest}...)	
		}
	}
	return anagrams
}