package yep

import (
	"encoding/json"
	"fmt"
	"net/url"
	"search/querus/engine"
	"search/querus/useragents"
	/* jsoniter "github.com/json-iterator/go"*/)

/* var json = jsoniter.ConfigCompatibleWithStandardLibrary*/

var YepLink string = "https://api.yep.com/fs/2/search"

type Yep struct {
	user_agent  string
	headers     map[string]string
	Name        string
	ponderation float64
}

func NewYepEngine(ponderation float64) *Yep {
	y := &Yep{user_agent: useragents.Get_random_ua(), Name: "Yep", ponderation: ponderation}
	y.init()
	return y
}

func (y *Yep) GetName() string {
	return y.Name
}
func (y *Yep) init() {
	y.headers = map[string]string{
		"Accept":          "*/*",
		"Accept-encoding": "gzip, deflate, br, zstd",
		"Origin":          "https://yep.com",
		"Priority":        "u=1, i",
		"Referer":         "https://yep.com/",
		"Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
	}

	//y.headers["User-Agent"] = y.user_agent
}

func (y *Yep) set_lang(lang string) {
	y.headers["Accept-Language"] = fmt.Sprintf("%s,fr;q=0.9,en-US;q=0.8,en;q=0.7", lang)
}

func (y *Yep) requestSearch(options engine.RequestOptions, searchType string) (string, error) {

	y.init()
	y.set_lang(options.Lang)

	client := "web"
	var orderedKeys []string
	params := url.Values{}

	params.Set("client", client)
	params.Set("gl", options.Lang[3:])
	params.Set("no_correct", "false")
	params.Set("q", options.Query)
	params.Set("safeSearch", options.SafeSearch)
	params.Set("type", searchType)

	orderedKeys = append(orderedKeys, "client", "gl", "no_correct", "q", "safeSearch", "type")

	if (searchType == "web" || searchType == "news") && options.MaxResults > 0 {
		params.Set("limit", fmt.Sprint(options.MaxResults))
	}

	response, err := engine.DoGetRequest(YepLink, params, orderedKeys, y.headers)
	if err != nil {
		return "", fmt.Errorf("error : %v", err)
	}
	return string(response), nil
}

func (y *Yep) controlOptions(options *engine.RequestOptions) {
	if options.Lang == "" {
		options.Lang = "en-EN"
	}
	if options.SafeSearch == "" {
		options.SafeSearch = "moderate"
	}
	if options.MaxResults < 0 {
		options.MaxResults = 0
	}
}
func (y *Yep) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return y.search(options, "web")
}
func (y *Yep) NewsSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return y.search(options, "news")
}
func (y *Yep) ImagesSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return make([]engine.Result_Search, 0), nil
}
func (y *Yep) search(options engine.RequestOptions, searchType string) ([]engine.Result_Search, error) {
	y.controlOptions(&options)

	response, err := y.requestSearch(options, searchType)
	if err != nil {
		return make([]engine.Result_Search, 0), fmt.Errorf("request error : %v", err)
	}
	results, ok := y.processWebSearch(response)
	if ok != nil {
		return results, fmt.Errorf("can't process web : %v", err)
	}
	if len(results) > 0 {
		return results, nil
	}
	return results, fmt.Errorf("error")

}

func (y *Yep) processWebSearch(jsonStr string) ([]engine.Result_Search, error) {
	// Décodage du JSON en une structure de données dynamique

	var decodedData []interface{}
	finalResults := make([]engine.Result_Search, 0)

	err := json.Unmarshal([]byte(jsonStr), &decodedData)
	if err != nil {
		return finalResults, fmt.Errorf("erreur de décodage JSON : %v", err)
	}
	nb := len(decodedData) - 1
	// Récupération des résultats de type "web"
	results := decodedData[int(nb)].(map[string]interface{})["results"].([]interface{})
	//fmt.Println(results)
	// Parcours des résultats et affichage des résultats de type "web"
	for key, data := range results {
		Data, err := data.(map[string]interface{})
		if !err {
			continue
		}
		title, err := Data["title"].(string)
		if !err {
			continue
		}
		url, err := Data["url"].(string)
		if !err {
			continue
		}
		description := Data["snippet"].(string)
		if !err {
			continue
		}
		position := key + 1
		finalResults = append(finalResults,
			engine.Result_Search{
				Item: engine.Item{

					Title:    engine.Normalize(title),
					Body:     engine.Normalize(description),
					Url:      engine.NormalizeUrl(url),
					Source:   engine.ExtractDomain(url),
					Engines:  []string{"Yep"},
					Position: position,
					Score:    engine.Scoring(position, y.ponderation),
					Favicon:  fmt.Sprintf("https://%v/favicon.ico", engine.ExtractDomain(url)),
				}})
	}
	if len(finalResults) > 0 {
		return finalResults, nil
	}
	return finalResults, fmt.Errorf("error")
}
