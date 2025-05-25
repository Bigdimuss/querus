package brave

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"search/querus/engine"
	"search/querus/useragents"
	"strings"

	"github.com/robertkrimen/otto"
)

// BraveLINKS contient les URLs pour différents types de recherche sur Brave
var BraveLINKS map[string]string = map[string]string{
	"web":     "https://search.brave.com/search",
	"images":  "https://search.brave.com/images",
	"news":    "https://search.brave.com/news",
	"videos":  "https://search.brave.com/videos",
	"goggles": "https://search.brave.com/goggles",
}

// BraveHEADERS contient les en-têtes HTTP pat défaut à utiliser pour les requêtes.
var BraveHEADERS map[string]string = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"Accept-Encoding": "gzip, deflate", // br BOGUE
	"Accept-Language": "fr,en-US;q=0.9,en;q=0.8",
}

// BraveOPTIONS définit des options de fraîcheur pour les requêtes.
var BraveOPTIONS map[string]map[string]string = map[string]map[string]string{
	"freshness": {
		"undefined":  "at",
		"last_hour":  "ph",
		"last_day":   "pd",
		"last_week":  "pw",
		"last_month": "pm",
		"last_years": "py"},
}

// Brave structure qui représente l'agent de recherche Brave
type Brave struct {
	user_agent  string
	Name        string
	ponderation float64
}

// NewBraveEngine crée une nouvelle instance de Brave avec une pondération spécifiée.
func NewBraveEngine(ponderation float64) *Brave {
	return &Brave{user_agent: useragents.Get_random_ua(), Name: "Brave", ponderation: ponderation}
}

// GetName retourne le nom de l'agent de recherche
func (b *Brave) GetName() string {
	return b.Name
}

// requestSearch effectue une recherche vers l'URL approprié avec les paramètres spécifiés.
func (b *Brave) requestSearch(options engine.RequestOptions, indexPage int, searchType string) (string, error) {
	/*
		// Définit des cookies pour la requête.
		cookies := map[string]string{
			"useLocation": "0",
			"safesearch":  options.SafeSearch,
			"country":     options.Lang[:2],
		}
	*/
	// Prépare les paramètres de recherche.
	params := url.Values{}
	var orderedKeys []string
	params.Set("q", options.Query)
	orderedKeys = append(orderedKeys, "q")
	if searchType == "web" {
		params.Set("source", "web")
		params.Set("tf", options.Freshness)
		orderedKeys = append(orderedKeys, "source", "tf")
	}

	// Gestion de la pagination
	if indexPage > 0 {
		params.Set("offset", fmt.Sprint(options.IndexPage))
		params.Set("spellcheck", "0")
		orderedKeys = append(orderedKeys, "offset", "spellcheck")
	}

	// Création de l'URL de requête.
	link := BraveLINKS[searchType]
	link = engine.Creat_url_with_params(link, params, orderedKeys)
	fmt.Println(link)
	fmt.Println("BONJOUR !")

	// Prépare l'en-tête de requête.
	header := BraveHEADERS
	header["User-Agent"] = b.user_agent
	header["Accept-Language"] = fmt.Sprintf("%s,%s;q=0.8,en-US;q=0.5,en;q=0.3", strings.ToLower(options.Lang[:2]), options.Lang)
	header["Referer"] = "https://search.brave.com/"

	// Création de la requête.
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println(err)
	}
	//engine.AddCookies(request, cookies)
	engine.Creat_headers(request, header)

	// Envoi de la requête et récuperation de la réponse.
	response, err := engine.DoRequest(request)
	if err != nil {
		return "", fmt.Errorf("error dorequest : %v", err)
	}
	return string(response), nil
}

// processWebSearch traite les résultats de recherche et les convertit en un format exploitable
func (b *Brave) processWebSearch(jsonStr string, page int, searchType string) ([]engine.Result_Search, error) {

	finalResults := make([]engine.Result_Search, 0)

	//Extraction du JSON de la réponse
	jsonStr, err := extract_json(jsonStr)
	if err != nil {
		return finalResults, fmt.Errorf("error processwebsearch : %v", err)
	}

	vm := otto.New()              // Initialise un interpréteur de JavaScript
	value, err := vm.Run(jsonStr) //Exécute le JSON comme un code JS
	if err != nil {
		return finalResults, fmt.Errorf("error javascript evaluation : %v", err)
	}

	// Exportation des données.
	decodedData, err := value.Export()
	if err != nil {
		fmt.Println("Erreur lors de l'exportation des données en JSON:", err)
		return finalResults, fmt.Errorf("error export data from json : %v", err)
	}

	datas, ok := decodedData.([]map[string]interface{}) // Conversion en map

	if !ok {
		return finalResults, fmt.Errorf("decodeData error")
	}

	// Extraction des résultats en fonction du type de recherche
	results, ok := datas[len(datas)-1]["data"]
	if !ok {
		return finalResults, fmt.Errorf("can't read data")
	}
	results, ok = results.(map[string]interface{})["body"]
	if !ok {
		return finalResults, fmt.Errorf("can't read data")
	}
	results, ok = results.(map[string]interface{})["response"]
	if !ok {
		return finalResults, fmt.Errorf("can't read data")
	}

	if searchType == "web" {
		results, ok = results.(map[string]interface{})["web"]
		if !ok {
			return finalResults, fmt.Errorf("can't read data")
		}
	}
	if searchType == "news" {
		results, ok = results.(map[string]interface{})["news"]
		fmt.Println("NEWS")
		if !ok {
			return finalResults, fmt.Errorf("can't read data")
		}
	}

	// Récupération des résultats spécifiques.
	results, ok = results.(map[string]interface{})["results"]
	if !ok {
		return finalResults, fmt.Errorf("can't read data")
	}
	resultsData := results.([]map[string]interface{})

	// Traitement des résultats pour création d'items.
	var position int = 1 * (page + 1)
	for _, Data := range resultsData {
		title, ok := Data["title"].(string)
		if !ok {
			continue // Ignore les resultats sans titre.
		}
		url, ok := Data["url"].(string)
		if !ok {
			continue // Ignore les resultats sans URL
		}
		body, ok := Data["description"].(string)
		if !ok {
			continue // Ignore les résultats sans description.
		}
		source := engine.ExtractDomain(url) // Extrait le domaine de l'URL
		var item engine.Item = engine.Item{
			Title:    engine.Normalize(title),                 // Normalise le titre
			Url:      engine.NormalizeUrl(url),                // Normalise l'URL
			Body:     engine.Normalize(body),                  // Normalise la description
			Source:   source,                                  // Source de l'item
			Engines:  []string{"Brave"},                       // indique l'agent de recherche
			Position: position,                                // position dans les resultats
			Score:    engine.Scoring(position, b.ponderation), // Calcule du score basé sur la position
		}

		// Traitement spécifiques pour les recherche d'images.
		if searchType == "images" {
			img, ok := Data["properties"].(map[string]interface{})["url"]
			if !ok {
				continue
			}
			item.Img = img.(string)
			width, ok := Data["properties"].(map[string]interface{})["width"].(float64)
			if !ok {
				width = 0
			}
			item.Width = int32(math.Round(width))
			height, ok := Data["properties"].(map[string]interface{})["height"].(float64)
			if !ok {
				height = 0
			}
			item.Height = int32(math.Round(height))
		}

		finalResults = append(finalResults, engine.Result_Search{Item: item}) // Ajout de l'item à la liste finale
		position += 1                                                         // Incrémente la position
	}
	return finalResults, nil
}

// controlOptions valide et ajuste les options de requête avant l'exécution.
func (b *Brave) controlOptions(options *engine.RequestOptions) {
	switch options.Period {
	case "at", "ph", "pd", "pw", "pm", "py":
	default:
		options.Period = "at"
	}
	if options.Lang == "" {
		options.Lang = "en-EN" // Langue par défaut.
	}
	if options.SafeSearch == "" {
		options.SafeSearch = "moderate" // Recherche sécurisé par défaut
	}
	if options.IndexPage < 0 {
		options.IndexPage = 0 // Indice de page minimum
	}
	if options.Freshness == "" {
		options.Freshness = BraveOPTIONS["freshness"]["undefined"] // Gestion de la fraîcheur.
	}
}

// search effectue une recherche générale et retourne les résultats formatés.
func (b *Brave) search(options engine.RequestOptions, searchType string, maxResultPerPage int) ([]engine.Result_Search, error) {
	b.controlOptions(&options) // Valide les options de recherche.
	if options.MaxResults < maxResultPerPage {
		options.MaxResults = maxResultPerPage // Ajuste le nombre de résultats maximum par page.
	}
	return engine.SearchGeneric(options, searchType, maxResultPerPage, b.requestSearch, b.processWebSearch) // Appelle la recherche générique.
}

// WebSearch effectue une recherche Web via Brave.
func (b *Brave) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var max_result_per_page int = 20 // Limite le nombre de résultats à 20.
	return b.search(options, "web", max_result_per_page)
}

// NewsSearch effectue une recherche d'actualités via Brave.
func (b *Brave) NewsSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {

	var max_result_per_page int = 50 // Limite le nombre de résultats à 50.
	return b.search(options, "news", max_result_per_page)
}

// ImagesSearch effectue une recherche d'images via Brave.
func (b *Brave) ImagesSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {

	b.controlOptions(&options) // Valide les options de recherche.
	var results []engine.Result_Search = make([]engine.Result_Search, 0)
	r, err := b.requestSearch(options, 0, "images") // Effectue la recherche d'images.

	if err != nil {
		fmt.Println("ERREUR REQUEST BRAVE")
		return results, fmt.Errorf("error")
	}

	v, err := b.processWebSearch(r, 0, "images") // Traite les résultats.
	if err != nil {
		fmt.Printf("quelque chose se passe mal ... \n !Error : %v", err)
		return results, fmt.Errorf("error")
	}

	results = append(results, v...) // Ajoute les résultats traités.
	return results, nil
}
