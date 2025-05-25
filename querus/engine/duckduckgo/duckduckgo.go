package duckduckgo

import (
	"bytes" // Pour manipuler des byte slices
	"encoding/json"
	"fmt"                  // Pour formater des chaînes de caractères
	"io"                   // Pour les opérations d'entrée/sortie
	"math"                 // Pour les calculs mathématiques
	"net/http"             // Pour faire des requêtes HTTP
	"net/url"              // Pour manipuler les URL
	"search/querus/engine" // Module de moteur de recherche
	"search/querus/useragents"
	"strings" // Pour manipuler des chaînes de caractères
	"sync"    // Pour la gestion de la concurrence
	/*jsoniter "github.com/json-iterator/go"*/) // Librairie pour le traitement JSON

// Utilisation d'une configuration compatible avec la bibliothèque standard pour JSON
/*var json = jsoniter.ConfigCompatibleWithStandardLibrary*/

// Liens de DuckDuckGo pour différents types de recherches
var DuckDuckGOLink map[string]string = map[string]string{
	"web":    "https://links.duckduckgo.com/d.js",
	"images": "https://duckduckgo.com/i.js",
	"news":   "https://duckduckgo.com/news.js",
}

// Structure principale pour interagir avec DuckDuckGo
type DuckDuckGO struct {
	user_agent  string            // Agent utilisateur pour les requêtes HTTP
	headers     map[string]string // En-têtes HTTP
	vqd         string            // Valeur de requête variable
	ponderation float64           // Ponderation pour le score des résultats
	Name        string            // Nom du moteur de recherche
}

// Fonction de création d'une nouvelle instance de DuckDuckGO
func NewDuckDuckGoEngine(ponderation float64) *DuckDuckGO {
	// Création et initialisation d'un objet DuckDuckGO
	d := &DuckDuckGO{user_agent: useragents.Get_random_ua(), headers: map[string]string{"Referer": "https://duckduckgo.com/"}, Name: "DuckDuckGo", ponderation: ponderation}
	return d
}

// Retourne le nom du moteur de recherche
func (d *DuckDuckGO) GetName() string {
	return d.Name
}

// Fonction pour obtenir la valeur VQD nécessaire pour effectuer des recherches
func (d DuckDuckGO) getVqd(keywords string) (string, error) {
	data := url.Values{}
	data.Set("q", keywords)
	data.Add("ia", "web") // Indication que la recherche est pour le web

	// Création d'une requête POST
	request, err := http.NewRequest("POST", "https://duckduckgo.com", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err) // Gestion des erreurs
	}

	// Ajout de l'en-tête à la requête
	d.headers["Content-Type"] = "application/x-www-form-urlencoded"
	d.headers["User-Agent"] = d.user_agent

	engine.Creat_headers(request, d.headers) // Fonction pour créer des en-têtes

	// Exécution de la requête et traitement de la réponse
	resp, err := engine.DoPostData(request, data)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %v", err)
	}
	defer resp.Body.Close() // S'assure que le corps de la réponse est fermé

	body, err := io.ReadAll(resp.Body) // Lecture du corps de la réponse
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	stringBody := string(body)                   // Conversion en chaîne
	vqd, err := extractVqd(stringBody, keywords) // Extraction de la valeur VQD
	if err != nil {
		return "", fmt.Errorf("failed to extract VQD : %v", err)
	}
	return vqd, nil
}

// Fonction principale pour effectuer une recherche
func (d *DuckDuckGO) requestSearch(query string, lang string, safeSearch string, offset int, vqd string, searchType string) ([]byte, error) {
	params := url.Values{} // Paramètres pour la requête
	var orderedKeys []string

	// Configuration des paramètres selon le type de recherche (web, images, news)
	if searchType == "web" {

		// Gérer les paramètres de recherche sécurisée
		switch safeSearch {
		case "moderate":
			params.Set("ex", "-1")
			orderedKeys = append(orderedKeys, "ex")
		case "off":
			params.Set("ex", "-2")
			orderedKeys = append(orderedKeys, "ex")
		case "on":
			params.Set("p", "-1")
			orderedKeys = append(orderedKeys, "p")

		default:
			params.Set("ex", "-1")
			orderedKeys = append(orderedKeys, "ex")
		}
		params.Set("l", lang)
		params.Set("q", query)
		params.Set("vqd", vqd)
		orderedKeys = append(orderedKeys, "l", "q", "vqd")
	}
	// Le cas spécifique pour les nouvelles
	if searchType == "news" {
		params.Set("l", lang)
		params.Set("o", "json")
		params.Set("noamp", "1")
		params.Set("q", query)
		params.Set("vqd", vqd)
		orderedKeys = append(orderedKeys, "l", "o", "noamp", "q", "vqd", "p")
		switch safeSearch {
		case "moderate":
			params.Set("p", "-1")
		case "off":
			params.Set("p", "-2")
		case "on":
			params.Set("p", "1")
		default:
			params.Set("p", "-1")
		}
	}

	if searchType == "images" {
		// Configuration des paramètres pour la recherche d'images
		params.Set("l", lang)
		params.Set("o", "json")
		params.Set("q", query)
		params.Set("vqd", vqd)
		params.Set("f", ",,,,")
		orderedKeys = append(orderedKeys, "l", "o", "q", "vqd", "f", "p")
		switch safeSearch {
		case "moderate":
			params.Set("p", "1")
		case "off":
			params.Set("p", "-1")
		case "on":
			params.Set("p", "1")
		default:
			params.Set("p", "1")
		}
	}

	// Gestion du décalage pour la pagination
	if offset > 0 {
		params.Set("s", fmt.Sprint(offset))
		orderedKeys = append(orderedKeys, "s")
	}

	// Mise à jour des en-têtes d'acceptation des langues
	d.headers["Accept-Languages"] = fmt.Sprintf("%v,en-US;q=0.9,en;q=0.8", strings.ToLower(lang[:2]))

	var link string = DuckDuckGOLink[searchType]                               // Récupération du lien approprié
	response, err := engine.DoGetRequest(link, params, orderedKeys, d.headers) // Envoi de la requête GET
	if err != nil {
		return nil, fmt.Errorf("request error : %v", err) // Gestion des erreurs
	}
	return response, nil // Retourne la réponse
}

// Contrôle et ajuste les options de la recherche
func (d *DuckDuckGO) controlOptions(options *engine.RequestOptions) {
	if options.Lang == "" {
		options.Lang = "fr-FR" // Langue par défaut
	}
	options.Lang = strings.ToLower(options.Lang) // Normalisation de la langue

	// Validation des index de page et des résultats maximum
	if options.IndexPage < 0 {
		options.IndexPage = 0
	}
	if options.MaxResults < 50 {
		options.MaxResults = 50
	}

	// Vérification de l'option de recherche sécurisée
	switch options.SafeSearch {
	case "moderate", "on", "off": // Options valides
	default:
		options.SafeSearch = "moderate" // Option par défaut
	}
}

// Fonction principale pour la recherche
func (d *DuckDuckGO) search(options engine.RequestOptions, searchType string, nb_results_by_request int) ([]engine.Result_Search, error) {
	d.controlOptions(&options) // Contrôle des options

	var offset int
	var page_nb int = options.MaxResults / nb_results_by_request // Calcul du nombre de pages
	if options.MaxResults%nb_results_by_request != 0 {
		page_nb += 1 // Ajout d'une page si nécessaire
	}
	var results []engine.Result_Search = make([]engine.Result_Search, 0) // Initialiser le tableau des résultats
	vqd, err := d.getVqd(options.Query)                                  // Obtenir la valeur VQD
	if err != nil {
		fmt.Println(err) // Log de l'erreur
		return results, fmt.Errorf("error vqd : %v", err)
	}
	d.vqd = vqd // Stocker la valeur VQD

	var resultsMutex sync.Mutex // Mutex pour la protection des résultats
	var wg sync.WaitGroup       // WaitGroup pour la gestion des goroutines
	for i := options.IndexPage; i < options.IndexPage+page_nb; i++ {
		offset = i * nb_results_by_request // Calculer l'offset
		wg.Add(1)                          // Incrémente le compteur de waitgroup
		go func(offset int, i int, results *[]engine.Result_Search, wg *sync.WaitGroup) {
			defer wg.Done() // Indique que la goroutine est terminée

			r, err := d.requestSearch(options.Query, options.Lang, options.SafeSearch, offset, vqd, searchType) // Requête de recherche
			if err != nil {
				fmt.Println(err) // Log de l'erreur
				return
			}
			v, err := d.processWebSearch(r, i, options.Query, searchType) // Traitement des résultats
			if err != nil {
				fmt.Printf("quelque chose se passe mal : %v", err) // Log de l'erreur
				return
			}
			resultsMutex.Lock()                   // Verrouiller l'accès aux résultats
			(*results) = append((*results), v...) // Ajouter les résultats
			resultsMutex.Unlock()                 // Déverrouiller l'accès
		}(offset, i, &results, &wg) // Appel de la goroutine
	}
	wg.Wait() // Attendre que toutes les goroutines soient terminées

	if len(results) > 0 {
		return results, nil // Retourne les résultats trouvés
	}
	return results, fmt.Errorf("error") // Retourne une erreur si aucun résultat
}

// Fonction pour traiter les résultats de recherche web
func (d *DuckDuckGO) processWebSearch(htmlBytes []byte, page_nb int, query string, searchType string) ([]engine.Result_Search, error) {
	var position int = 1 * (page_nb + 1) // Position de départ des résultats
	var start int
	var end int
	var data []byte

	if searchType == "web" {
		// Extraction des données spécifiques à la recherche web
		start = bytes.Index(htmlBytes, []byte("DDG.pageLayout.load('d',")) + 24
		end = bytes.Index(htmlBytes[start:], []byte(");DDG.duckbar.load(")) + start
		data = htmlBytes[start:end]
	} else {
		data = htmlBytes // Utiliser les données originales
	}

	var results []engine.Result_Search = make([]engine.Result_Search, 0) // Initialiser le tableau des résultats
	var cache []string                                                   // Cache pour éviter les doublons
	if searchType == "web" {
		var decodedData []interface{}             // Décodage des données JSON
		err := json.Unmarshal(data, &decodedData) // Décodage des données JSON

		if err != nil {
			return make([]engine.Result_Search, 0), fmt.Errorf("process result error : %v", err)
		}
		for _, data := range decodedData {
			Data, ok := data.(map[string]interface{}) // Conversion des données
			if !ok {
				continue
			}
			href, ok := Data["u"].(string) // Récupération de l'URL
			if !ok {
				continue
			}
			if href == "" {
				continue
			}
			if href != fmt.Sprintf("http://www.google.com/search?q=%v", query) { // Exclure les résultats de Google
				is_it := false
				for _, v := range cache { // Vérifier si l'URL est déjà en cache
					if v == href {
						is_it = true
						break
					}
				}
				if !is_it {
					// Normalisation et ajout des résultats
					title := engine.Normalize(Data["t"].(string))
					body := engine.Normalize(Data["a"].(string))
					if body != "" {
						results = append(results,
							engine.Result_Search{
								Item: engine.Item{
									// Remplir le résultat
									Title:    engine.Normalize(title),
									Url:      engine.NormalizeUrl(href),
									Body:     engine.Normalize(body),
									Source:   engine.ExtractDomain(href),
									Engines:  []string{"DuckDuckGo"},
									Position: position,
									Score:    engine.Scoring(position, d.ponderation),
								}})
					}
					cache = append(cache, href) // Ajouter à la cache
					position += 1               // Incrémenter la position
				}
			}
		}
	}

	if searchType == "images" {
		var decodedData map[string]interface{} // Décodage des données pour les images
		err := json.Unmarshal(data, &decodedData)

		if err != nil {
			return make([]engine.Result_Search, 0), fmt.Errorf("process result error : %v", err)
		}

		resultsData := decodedData["results"].([]interface{}) // Récupération des résultats
		for _, data := range resultsData {
			items := data.(map[string]interface{}) // Conversion
			title, ok := items["title"].(string)
			if !ok {
				continue
			}
			url, ok := items["url"].(string)
			if !ok {
				continue
			}
			img, ok := items["image"].(string)
			if !ok {
				continue
			}
			height, ok := items["height"].(float64) // Récupération de la hauteur
			if !ok {
				height = 0
			}
			width, ok := items["width"].(float64) // Récupération de la largeur
			if !ok {
				width = 0
			}
			results = append(results,
				engine.Result_Search{
					Item: engine.Item{
						Title:    engine.Normalize(title),
						Url:      engine.NormalizeUrl(url),
						Img:      engine.Normalize(img),
						Body:     "",
						Height:   int32(math.Round(height)),
						Width:    int32(math.Round(width)),
						Source:   engine.ExtractDomain(url),
						Engines:  []string{"DuckDuckGo"},
						Position: position,
						Score:    engine.Scoring(position, d.ponderation),
					}})
		}
	}
	if searchType == "news" {
		// Traitement spécifique pour les nouvelles
		var decodedData map[string]interface{}
		err := json.Unmarshal(data, &decodedData)

		if err != nil {
			return make([]engine.Result_Search, 0), fmt.Errorf("process result error : %v", err)
		}

		resultsData := decodedData["results"].([]interface{})
		for _, data := range resultsData {
			items := data.(map[string]interface{})
			title, ok := items["title"].(string)
			if !ok {
				continue
			}
			url, ok := items["url"].(string)
			if !ok {
				continue
			}
			body, ok := items["excerpt"].(string)
			if !ok {
				continue
			}
			image, ok := items["image"].(string)
			if !ok {
				image = ""
			}
			date, ok := items["date"].(int)
			if !ok {
				date = 0
			}
			results = append(results,
				engine.Result_Search{
					Item: engine.Item{
						Title:    engine.Normalize(title),
						Url:      engine.NormalizeUrl(url),
						Body:     body,
						Img:      image,
						Source:   engine.ExtractDomain(url),
						Engines:  []string{"DuckDuckGo"},
						Date:     date,
						Position: position,
						Score:    engine.Scoring(position, d.ponderation),
					}})
		}
	}
	return results, nil // Retourne les résultats traités
}

// Méthode de recherche pour le web
func (d *DuckDuckGO) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var nb_results_by_request int = 50                     // Nombre de résultats par requête
	return d.search(options, "web", nb_results_by_request) // Appel de la méthode de recherche
}

// Méthode de recherche pour les images
func (d *DuckDuckGO) ImagesSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var nb_results_by_request int = 100                       // Nombre de résultats par requête
	return d.search(options, "images", nb_results_by_request) // Appel de la méthode de recherche
}

// Méthode de recherche pour les nouvelles
func (d *DuckDuckGO) NewsSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var nb_results_by_request int = 100                    // Nombre de résultats par requête
	return d.search(options, "web", nb_results_by_request) // Appel de la méthode de recherche
}
