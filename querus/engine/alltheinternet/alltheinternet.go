package alltheinternet

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"search/querus/engine"
	"search/querus/useragents"

	"github.com/anaskhan96/soup"
	jsoniter "github.com/json-iterator/go"
)

// fonction de la bibliothèque JSON
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// STRUCT
// Définition de la structure AllTheInternet
type AllTheInternet struct {
	user_agent  string
	Name        string
	ponderation float64
}

// Fonction pour créer une nouvelle instance de AllTheInternet
func NewAllTheInternetEngine(ponderation float64) *AllTheInternet {
	return &AllTheInternet{user_agent: useragents.Get_random_ua(), Name: "AllTheInternet", ponderation: ponderation}
}

// Renvoie le Nom du moteur de recherche
func (ati *AllTheInternet) GetName() string {
	return ati.Name
}

// Récupère les certificats de recherche nécessaire
func (ati *AllTheInternet) get_certif(query string, searchType string) (url.Values, []string, error) {

	var baseUrl string = "https://www.alltheinternet.com/"
	params := url.Values{}
	var orderedKeys []string
	params.Set("q", query)

	orderedKeys = append(orderedKeys, "q")

	// En-tête de requête
	headers := map[string]string{
		"User-Agent": ati.user_agent,
		"Referer":    baseUrl,
	}

	// Envoi la requète GET
	body, err := engine.DoGetRequest(baseUrl, params, orderedKeys, headers)
	if err != nil {
		return url.Values{}, make([]string, 0), fmt.Errorf("error DoGetRequest")
	}

	var orderedFinalKeys []string
	orderedFinalKeys = append(orderedFinalKeys, "cx")
	// Analyse du HTML obtenu
	doc := soup.HTMLParse(string(body))
	scripts := doc.FindAllStrict("script")
	for _, v := range scripts {
		src, err := v.Attrs()["src"]
		if !err {
			continue
		}
		re := regexp.MustCompile(`https://cse.google.com/cse.js`)
		if re.MatchString(src) {
			u, err := url.Parse(src)
			if err != nil {
				return url.Values{}, make([]string, 0), fmt.Errorf("error with parse : %v", err)
			}
			queryValues, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				return url.Values{}, make([]string, 0), fmt.Errorf("error with parsequery : %v", err)
			}
			keysValues := url.Values{}
			keysValues.Set("cx", string(queryValues.Get("cx")))
			return keysValues, orderedFinalKeys, nil
		}
	}
	return url.Values{}, make([]string, 0), fmt.Errorf("error get_certif function ")
}

// Récupère les clés nécessaire pour la recherche
func (ati *AllTheInternet) get_keys(query string, searchType string) (map[string]string, error) {
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3",
		"accept-encoding": "gzip, deflate",
		"user-agent":      ati.user_agent}

	// Obtention des certificats
	params, orderedKeys, err := ati.get_certif(query, searchType)

	if err != nil {
		return make(map[string]string), fmt.Errorf("error : %v", err)
	}
	rbyte, err := engine.DoGetRequest("https://cse.google.com/cse.js", params, orderedKeys, headers)
	if err != nil {
		return make(map[string]string), err
	}
	response := string(rbyte)
	return extract_keys(string(response)), nil
}

// Fonction principale pour effectuer une recherche
func (ati *AllTheInternet) requestSearch(options engine.RequestOptions, offset int, searchType string) (string, error) {

	link := "https://cse.google.com/cse/element/v1"

	headers := map[string]string{
		"user-agent":      ati.user_agent,
		"accept-language": fmt.Sprintf("%s,%s;q=0.8,en-US;q=0.5,en;q=0.3", options.Lang[:2], options.Lang)}

	params := url.Values{}
	var orderedKeys []string
	rurl := fmt.Sprintf("https://www.alltheinternet.com/?q=%s", options.Query)
	// Si le type de recherche n'est pas "web", il est ajouté aux paramètres
	if searchType != "web" && (searchType == "news" || searchType == "image" || searchType == "video") {
		prefix := fmt.Sprintf("&area=%s", searchType)
		rurl += prefix
	}
	params.Set("rsz", fmt.Sprint(options.MaxResults))
	params.Set("num", fmt.Sprint(options.MaxResults))
	params.Set("hl", options.Lang[:2])
	params.Set("source", "gcsc")
	params.Set("start", fmt.Sprint(offset))
	params.Set("gss", ".com")
	params.Set("q", options.Query)
	params.Set("safe", options.SafeSearch)
	params.Set("lr", options.Lang)
	params.Set("cr", "")
	params.Set("gl", "")
	params.Set("filter", "0")
	//fmt.Sprintf("%d", int(indexPage*int(maxResults))),
	params.Set("sort", "")
	params.Set("as_oq", "")
	params.Set("as_sitesearch", "")
	params.Set("as_filetype", "")
	params.Set("exp", "cc")
	params.Set("callback", "google.search.cse.api")
	params.Set("rurl", rurl)

	orderedKeys = append(orderedKeys,
		"rsz",
		"num",
		"hl",
		"source",
		"gss",
		"start",
		"cselibv",
		"cx",
		"q",
		"safe",
		"cse_tok",
		"lr",
		"cr",
		"gl",
		"filter",
		"sort",
		"as_oq",
		"as_sitesearch",
		"as_filetype",
		"exp",
		"callback",
		"rurl")

	if searchType == "image" {
		params.Set("searchtype", searchType)
	}

	keysData, err := ati.get_keys(options.Query, searchType)

	if err != nil {
		return "", fmt.Errorf("error : %v", err)
	}

	// Ajout des clés aux paramètres
	for k, v := range keysData {
		params.Set(k, v)
	}

	rbyte, err := engine.DoGetRequest(link, params, orderedKeys, headers)
	if err != nil {
		return "", fmt.Errorf("error : %v", err)
	}

	response := string(rbyte)
	extract_json(&response)
	return response, nil
}

// Ajuste les options de recherche
func (ati *AllTheInternet) controlOptions(options *engine.RequestOptions, searchtype string) {
	if options.MaxResults <= 0 {
		if searchtype == "web" {
			options.MaxResults = 15
		}
		if searchtype == "images" {
			options.MaxResults = 18
		}
	}
	if options.Lang == "" {
		options.Lang = "en-EN"
	}
	if options.SafeSearch == "" {
		options.SafeSearch = "moderate"
	}
	if options.IndexPage < 0 {
		options.IndexPage = 0
	}
}

// Fonction de recherche principale
func (ati *AllTheInternet) search(options engine.RequestOptions, searchType string, nbResultPage int) ([]engine.Result_Search, error) {
	ati.controlOptions(&options, searchType)
	return engine.SearchGeneric(options, searchType, nbResultPage, ati.requestSearch, ati.processWebSearch)
}

// Traitement des résultats de recherche Web
func (ati *AllTheInternet) processWebSearch(jsonStr string, page int, searchType string) ([]engine.Result_Search, error) {

	var decodedData map[string]interface{}
	finalResults := make([]engine.Result_Search, 0)
	var position int = 1 * (page + 1)

	err := json.Unmarshal([]byte(jsonStr), &decodedData)

	if err != nil {
		return finalResults, fmt.Errorf("erreur de décodage JSON : %v", err)
	}

	// Récupération des résultats
	results, ok := decodedData["results"].([]interface{})
	if !ok {
		return finalResults, fmt.Errorf("aucun résultat trouvé dans decodeData")
	}

	for _, result := range results {
		resultData, ok := result.(map[string]interface{})
		if !ok || resultData == nil {
			continue
		}

		title, ok := resultData["titleNoFormatting"].(string)
		if !ok || title == "" {
			continue
		}

		description, ok := resultData["contentNoFormatting"].(string)
		if !ok || description == "" {
			continue
		}
		result := engine.Item{
			Title: title,
			Body:  description,
		}
		var img string

		if searchType == "image" {
			url, ok := resultData["originalContextUrl"].(string)
			if !ok {
				continue
			}
			result.Url = url
			img, ok := resultData["url"].(string)
			if !ok {
				continue
			}
			result.Img = img
			height, ok := resultData["height"].(float64)
			if !ok {
				height = 0
			}
			result.Height = int32(math.Round(height))
			width, ok := resultData["width"].(float64)
			if !ok {
				width = 0
			}
			result.Width = int32(math.Round(width))
		} else {
			url, ok := resultData["url"].(string)
			if !ok {
				continue
			}
			result.Url = url
			rich, ok := resultData["richSnippet"].(map[string]interface{})
			if !ok || rich == nil {
				// Le champ "richSnippet" est inexistant ou vide
				img = ""
			} else {
				cseImage, ok := rich["cseImage"].(map[string]interface{})
				if !ok || cseImage == nil {
					// Le champ "cseImage" est inexistant ou vide
					img = ""
				} else {
					src, ok := cseImage["src"].(string)
					if !ok {
						// Le champ "src" n'est pas une chaîne de caractères
						img = ""
					} else {
						// Tout est en ordre, affectez la valeur de "src" à "img"
						img = src
					}
				}
			}
			result.Img = img
		}
		result.Source = engine.ExtractDomain(result.Url)
		result.Engines = []string{ati.GetName()}
		result.Position = position
		result.Score = engine.Scoring(position, ati.ponderation)
		finalResults = append(finalResults, engine.Result_Search{Item: result})
		position += 1
	}

	return finalResults, nil
}

// Fonction pour rechercher sur le Web
func (ati *AllTheInternet) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var nb_result_page int = 15
	return ati.search(options, "web", nb_result_page)
}

// Fonction pour rechercher des images
func (ati *AllTheInternet) ImagesSearch(options engine.RequestOptions) ([]engine.Result_Search, error) { //MOdify processSearch to Image
	var nb_result_page int = 18
	return ati.search(options, "images", nb_result_page)
}

// Fonction pour rechercher des nouvelles
func (ati *AllTheInternet) NewsSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var nb_result_page int = 10
	return ati.search(options, "news", nb_result_page)
}
