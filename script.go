package main

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/mailgun/mailgun-go"
	"sort"
	"sync"
)

type TrafficQuery struct {
	Query string
	Position int
	TotalAppsCount int
}

type ByPosition []TrafficQuery
func (a ByPosition) Len() int { return len(a) }
func (a ByPosition) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPosition) Less(i, j int) bool { return a[i].Position < a[j].Position }

type Progress struct {
    Mu sync.Mutex
    States map[string]float64
}

func (c *Progress) Set(title string, progress float64) {
    c.Mu.Lock()
    c.States[title] = progress
    c.Mu.Unlock()
}

func (c *Progress) GetStates() (states map[string]float64) {
    c.Mu.Lock()
    states = c.States
    c.Mu.Unlock()
    return
}

var progress = Progress{States: make(map[string]float64)}

func measureTraffic(title string, keywords string, country string, email string) {
	queries := getQueries(keywords)
	fmt.Println(len(queries))
	progress.Set(title, 0)

	var bestQueries []TrafficQuery
	for i, query := range queries {
		words := strings.Split(query, " ")
		fmt.Println(query, fmt.Sprintf("(%d / %d)", i, len(queries)))
		term := strings.Join(words, "+")
		url := fmt.Sprintf("https://itunes.apple.com/search?term=%s&country=%s&entity=software&limit=100", term, country)

		data, err := readBody(url)
		if err != nil {
			fmt.Println(err)
		}
		var pageResp PageResp
		err = json.Unmarshal(data, &pageResp)
		if err != nil {
			fmt.Println(err)
		}
		if len(pageResp.Results) < 10 {
			continue
		}
		for i, entry := range pageResp.Results {
			if entry.TrackCensoredName == title {
				bestQueries = append(bestQueries, TrafficQuery{query, i + 1, len(pageResp.Results)})
				break
			}
		}
		progress.Set(title, float64(i + 1) / float64(len(queries)))
	}
	sort.Sort(ByPosition(bestQueries))
	fmt.Println("RESULTS: ", len(bestQueries))
	report := "Results:\n"
	for _, query := range bestQueries {
		line := fmt.Sprintf("%2v / %2v  %s", query.Position, query.TotalAppsCount, query.Query)
		fmt.Println(line)
		report = report + line + "\n"
	}
	progress.Set(title, 1)
	sendEmail(title, report, email)
}

func sendEmail(title string, body string, address string) {
	mg := mailgun.NewMailgun("sandbox8b923124ef234dfdb45b74f1ac03503a.mailgun.org", "key-4a4100c30e521251c772ca9e80fba232", "")

	m := mg.NewMessage(  
    	"Asometer <postmaster@sandbox8b923124ef234dfdb45b74f1ac03503a.mailgun.org>",
    	title,
    	body,
    	"Goga <" + address + ">")

	_, _, err := mg.Send(m)

	if err != nil {  
		fmt.Println(body)
  		fmt.Println(err)
	} else {
		fmt.Println("Mail Sent")
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

func getQueries(keywords string) []string {
	allKeywords := strings.Split(keywords, ",")
	var anagrams []string
	anagrams = append(anagrams, getAnagrams(allKeywords, 1, 0)...)
	// anagrams = append(anagrams, getAnagrams(allKeywords, 2, 0)...)
	// anagrams = append(anagrams, getAnagrams(allKeywords, 3, 0)...)
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