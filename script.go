package asometer

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/mailgun/mailgun-go"
	"sort"
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

func MeasureTraffic(title string, keywords string) {
	queries := getQueries(title, keywords)

	var bestQueries []TrafficQuery
	for i, query := range queries {
		words := strings.Split(query, " ")
		fmt.Println(query, fmt.Sprintf("(%d / %d)", i, len(queries)))
		term := strings.Join(words, "+")
		url := fmt.Sprintf("https://itunes.apple.com/search?term=%s&country=us&entity=software", term)

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
			if i > 30 {
				break
			}
			if entry.TrackCensoredName == title {
				bestQueries = append(bestQueries, TrafficQuery{query, i + 1, len(pageResp.Results)})
				break
			}
		}
	}
	//bestQueries := []TrafficQuery{{"asdsad", 32, 12}, {"Sdas", 0, 12}, {"rjrer", 2332, 9012}}
	sort.Sort(ByPosition(bestQueries))
	fmt.Println("RESULTS:")
	report := ""
	report = report + title + "\n\n"
	for _, query := range bestQueries {
		line := fmt.Sprintf("%2v / %2v  %s", query.Position, query.TotalAppsCount, query.Query)
		fmt.Println(line)
		report = report + line + "\n"
	}
	sendEmail(title, report)
}

func sendEmail(title string, body string) {
	mg := mailgun.NewMailgun("sandbox8b923124ef234dfdb45b74f1ac03503a.mailgun.org", "key-4a4100c30e521251c772ca9e80fba232", "")

	m := mg.NewMessage(  
    	"Asometer <postmaster@sandbox8b923124ef234dfdb45b74f1ac03503a.mailgun.org>",
    	title,
    	body,
    	"Goga <gogosapiens@gmail.com>")

	_, _, err := mg.Send(m)

	if err != nil {  
  		fmt.Println(err)
	} else {
		fmt.Println("Mail Sent")
	}
}

func main() {
	// title := "GetSpace Free - Delete Duplicate Photo from Device Storage"
	// keywords := "delete, clean, disk, master, manager, cleaner, memory, cache, camera, boost, aid, battery, saver, system, gallery"
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
	anagrams = append(anagrams, getAnagrams(allKeywords, 2, 0)...)
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