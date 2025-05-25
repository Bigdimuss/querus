package wiby

import (
	"encoding/json"
	"fmt"
	"net/url"
	"search/querus/engine"
	"search/querus/useragents"
)

var URL string = "https://wiby.me/json/"

type Wiby struct {
	user_agent  string
	Name        string
	ponderation float64
}

func NewWiby(ponderation float64) *Wiby {
	w := Wiby{user_agent: useragents.Get_random_ua(), Name: "Wiby", ponderation: ponderation}
	return &w
}

func (w *Wiby) requestSearch(options engine.RequestOptions, offset int, searchType string) (string, error) {
	params := url.Values{}
	var orderedKeys []string
	params.Set("q", options.Query)
	orderedKeys = append(orderedKeys, "q")
	headers := map[string]string{
		"User-Agent": w.user_agent,
	}
	nb := int(offset / 12)
	if nb > 1 {
		params.Set("p", fmt.Sprint(nb))
	}
	response, err := engine.DoGetRequest(URL, params, orderedKeys, headers)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la requête : %s", err)
	}
	return string(response), nil
}

func (w *Wiby) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return w.search(options, "web", 12)
}

// search effectue la recherche et retourne les résultats en fonction des options fournies.
func (w *Wiby) search(options engine.RequestOptions, searchType string, nbResultPage int) ([]engine.Result_Search, error) {
	w.controlOptions(&options, "web") // Vérification et ajustement des options
	return engine.SearchGeneric(options, searchType, nbResultPage, w.requestSearch, w.processWebSearch)
}

func (w *Wiby) ImagesSearch() ([]engine.Result_Search, error) {
	return make([]engine.Result_Search, 0), nil
}

func (w *Wiby) NewsSearch() ([]engine.Result_Search, error) {
	return make([]engine.Result_Search, 0), nil
}

func (w *Wiby) processWebSearch(jsonStr string, page int, searchType string) ([]engine.Result_Search, error) {
	finalResults := make([]engine.Result_Search, 0)
	var decodedData []interface{}
	var position int = 1 * (page + 1)
	// Décodage du JSON
	err := json.Unmarshal([]byte(jsonStr), &decodedData)
	if err != nil {
		return finalResults, fmt.Errorf("erreur de décodage JSON : %v", err)
	}
	for _, data := range decodedData {
		item, ok := data.(map[string]interface{})
		if !ok {
			return finalResults, fmt.Errorf("aucun résultat trouvé")
		}
		finalResults = append(finalResults,
			engine.Result_Search{
				Item: engine.Item{
					Title:    engine.Normalize(item["Title"].(string)),
					Body:     engine.Normalize(item["Snippet"].(string)),
					Url:      engine.NormalizeUrl(item["URL"].(string)),
					Source:   engine.ExtractDomain(item["URL"].(string)),
					Engines:  []string{w.GetName()},
					Position: position,
					Score:    engine.Scoring(position, w.ponderation),
				}})
	}
	return finalResults, nil
}
func (w *Wiby) GetName() string {
	return w.Name
}
func (wi *Wiby) controlOptions(options *engine.RequestOptions, searchtype string) {
	options.MaxResults = 12
	if options.IndexPage < 0 {
		options.IndexPage = 0
	}
}
