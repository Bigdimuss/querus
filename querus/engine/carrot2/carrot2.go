package carrot2

import (
	"fmt"
	"net/url"
	"search/querus/engine"
	"search/querus/useragents"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var Carrot2URL string = "https://www.etools.ch/partnerSearch.do"

var Carrot2ChoicesOptions map[string]map[string]string = map[string]map[string]string{
	"language": {
		"English": "en-EN",
		"French":  "fr-FR",
		"German":  "de-DE",
		"Italian": "it-IT",
		"Spanish": "es-ES",
	},
	"sources": {
		"all":     "all",
		"fastest": "fastest",
	},
	"country": {
		"all":           "all",
		"Austria":       "AT",
		"France":        "FR",
		"Germany":       "DE",
		"Great Britain": "GB",
		"Italy":         "IT",
		"Lichtenstein":  "LI",
		"Spain":         "ES",
		"Switzerland":   "CH",
	},
	"Safesearch": {
		"on":       "true",
		"off":      "false",
		"moderate": "true",
	},
}

type Carrot2 struct {
	user_agent  string
	headers     map[string]string
	ponderation float64
	Name        string
}

func NewCarrot2Engine(ponderation float64) *Carrot2 {
	q := &Carrot2{user_agent: useragents.Get_random_ua(), Name: "Carrot2", ponderation: ponderation}
	q.init()
	return q
}
func (c *Carrot2) GetName() string {
	return c.Name
}
func (c *Carrot2) init() {
	c.headers = map[string]string{
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate, br",
		"Referer":         "https://search.carrot2.org/",
		"Origin":          "https://search.carrot2.org",
		"User-Agent":      c.user_agent,
	}
}

func (c *Carrot2) requestSearch(query string, lang string, safeSearch string, maxResults int, sources string, country string) (string, error) {
	var orderedKeys []string
	params := url.Values{}
	params.Set("partner", "Carrot2Json")
	params.Set("query", query)
	params.Set("dataSourceResults", "40")
	params.Set("maxRecords", fmt.Sprintln(maxResults))
	params.Set("safeSearch", safeSearch)
	params.Set("dataSources", sources)
	params.Set("language", lang[:2])
	params.Set("country", country)
	orderedKeys = append(orderedKeys,
		"partner",
		"query",
		"dataSourceResults",
		"maxRecords",
		"safeSearch",
		"dataSources",
		"language",
		"country",
	)
	c.headers["Accept-Language"] = fmt.Sprintf("%s,%s;q=0.8,en-US;q=0.5,en;q=0.3", lang[:2], lang)
	response, err := engine.DoGetRequest(Carrot2URL, params, orderedKeys, c.headers)
	if err != nil {
		return "", fmt.Errorf("request search error : %v", err)
	}

	return string(response), nil
}

func (c *Carrot2) controlOptions(options *engine.RequestOptions) {
	if options.MaxResults != 0 {
		if options.MaxResults > 200 {
			options.MaxResults = 200
		}
	} else {
		options.MaxResults = 20
	}

	if options.Lang == "" {
		options.Lang = "en-EN"
	}
	if options.SafeSearch != "true" && options.SafeSearch != "false" {
		options.SafeSearch = "true"
	}
	if options.Country == "" {
		options.Country = "web"
	}
	if options.Sources == "" {
		options.Sources = "all"
	}
	if options.Country == "" {
		options.Country = strings.ToUpper(options.Lang[2:])
	}

}

func (c *Carrot2) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) { // query string, lang string, safeSearch string, maxResults int, country string) []engine.Result_Search {
	c.controlOptions(&options)
	response, err := c.requestSearch(options.Query, options.Lang, options.SafeSearch, options.MaxResults, options.Sources, options.Country)
	if err != nil {
		return make([]engine.Result_Search, 0), fmt.Errorf("error : %v", err)
	}
	data, ok := c.processWebSearch(response)
	if ok != nil {
		return data, fmt.Errorf("error : %v", err)
	}
	if len(data) > 0 {
		return data, nil
	}
	return data, fmt.Errorf("error")
}

func (c *Carrot2) processWebSearch(jsonStr string) ([]engine.Result_Search, error) {
	var decodedData map[string]interface{}

	finalResults := make([]engine.Result_Search, 0)

	err := json.Unmarshal([]byte(jsonStr), &decodedData)
	if err != nil {
		return finalResults, fmt.Errorf("erreur de décodage JSON : %v", err)
	}
	results, ok := decodedData["response"].(map[string]interface{})["mergedRecords"].([]interface{})
	if !ok {
		return finalResults, fmt.Errorf("decodedata error")
	}

	for k, v := range results {
		result := v.(map[string]interface{})
		title := result["title"].(string)
		url := result["url"].(string)
		description := result["text"].(string)
		position := k + 1
		finalResults = append(finalResults,
			engine.Result_Search{
				Item: engine.Item{
					Title:    engine.Normalize(title),
					Body:     engine.Normalize(description),
					Url:      engine.Normalize(url),
					Source:   engine.ExtractDomain(url),
					Engines:  []string{"Carrot2"},
					Position: position,
					Score:    engine.Scoring(position, c.ponderation),
					Favicon:  fmt.Sprintf("https://%v/favicon.ico", engine.ExtractDomain(url)),
				}})
	}
	if len(finalResults) > 0 {
		return finalResults, nil
	}
	return finalResults, fmt.Errorf("error")
}
