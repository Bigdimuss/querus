package qwant

// Package qwant implémente un moteur de recherche utilisant l'API de Qwant.
// Il définit des constantes pour les liens d'API, les options de recherche et les langues supportées.

import (
	"fmt"
	"math"
	"net/url"
	"search/querus/engine"
	"search/querus/useragents"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// Variable json configurable pour la compatibilité avec la bibliothèque standard.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// QwantLINKS contient les liens vers les différentes APIs de recherche de Qwant.
var QwantLINKS map[string]map[string]string = map[string]map[string]string{
	// Regroupement des différentes APIs selon le type de session (normal ou junior).
	"normal": {
		"web":       "https://fdn.qwant.com/v3/search/web",
		"knowledge": "https://fdn.qwant.com/v3/ia/knowledge2",
		"news":      "https://fdn.qwant.com/v3/search/news",
		"images":    "https://fdn.qwant.com/v3/search/images",
		"videos":    "https://fdn.qwant.com/v3/search/videos",
		"shopping":  "https://fdn.qwant.com/v3/search/shopping"},
	"junior": {
		"web":    "https://fdn.qwant.com/v3/egp/search/web",
		"images": "https://fdn.qwant.com/v3/egp/search/images",
		"videos": "https://fdn.qwant.com/v3/egp/search/videos"},
}

// QwantOPTIONS définit les options disponibles pour la recherche, comme la fraîcheur des résultats et l'ordre.
var QwantOPTIONS map[string]map[string]string = map[string]map[string]string{
	"freshness": {
		"undefined":  "all",
		"last_hour":  "hour",
		"last_day":   "day",
		"last_week":  "week",
		"last_month": "month"},
	"order": {
		"by_relevance": "relevance",
		"by_date":      "date",
	},
}

// QwantLangs contient une liste de langues supportées par le moteur de recherche.
var QwantLangs []string = []string{
	"fr-FR", "en-GB", "de-DE", "it-IT", "es-AR", // Langues disponibles
	"en-AU", "en-US", "es-ES", "ca-ES", // ...
	"cs-CZ", "ro-RO", "el-GB", "zh-CN", // ...
	// Ajout d'autres langues
}

// Qwant représente le moteur de recherche Qwant avec ses paramètres de configuration.
type Qwant struct {
	user_agent   string            // User-agent pour les requêtes HTTP
	headers      map[string]string // En-têtes HTTP nécessaires pour les requêtes
	session_type string            // Type de session (normal ou junior)
	Name         string            // Nom du moteur de recherche
	ponderation  float64           // Ponderation pour le scoring des résultats
}

// NewQwantEngine crée une nouvelle instance du moteur Qwant avec les paramètres par défaut.
func NewQwantEngine(ponderation float64) *Qwant {
	// Initialisation d'une instance de Qwant avec l'agent utilisateur et les paramètres par défaut.
	q := &Qwant{user_agent: useragents.Get_random_ua(), session_type: "normal", Name: "Qwant", ponderation: ponderation}
	q.init() // Appel à la méthode init pour initialiser les en-têtes
	return q
}

// NewQwantJuniorEngine crée une instance du moteur Qwant pour un public junior.
func NewQwantJuniorEngine(ponderation float64) *Qwant {
	q := &Qwant{user_agent: useragents.Get_random_ua(), session_type: "junior", Name: "QwantJunior", ponderation: ponderation}
	q.init()
	return q
}

// GetName retourne le nom du moteur de recherche.
func (q *Qwant) GetName() string {
	return q.Name
}

// init initialise les en-têtes HTTP par défaut pour les requêtes.
func (q *Qwant) init() {
	q.headers = map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "fr,en-US;q=0.9,en;q=0.8",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Origin":          "https://www.qwant.com",
		"Referer":         "https://www.qwant.com/",
	}
	q.headers["User-Agent"] = q.user_agent // Ajout de l'agent utilisateur aux en-têtes
}

// requestSearch effectue une recherche en fonction du type spécifié et des options fournies.
func (q *Qwant) requestSearch(options engine.RequestOptions, offset int, searchType string) (string, error) {
	q.init() // Réinitialiser les en-têtes avant chaque requête

	var orderedKeys []string
	// Paramètres de base communs à toutes les requêtes
	params := url.Values{}
	params.Set("q", options.Query)
	params.Set("locale", strings.Replace(options.Lang, "-", "_", -1))
	params.Set("offset", fmt.Sprint(offset))
	params.Set("device", "desktop")
	params.Set("tgp", fmt.Sprint(4))
	params.Set("safesearch", q.setSafeSearch(options.SafeSearch))
	orderedKeys = append(orderedKeys, "q", "locale", "offset", "device", "tgp", "safesearch")
	// Ajustement des paramètres selon le type de recherche
	switch searchType {
	case "web":
		params.Set("freshness", options.Freshness)
		params.Set("count", fmt.Sprint(10))
		params.Set("displayed", "true") // Spécifique à la recherche web
		orderedKeys = append(orderedKeys, "freshness", "count", "displayed")
	case "images":
		params.Set("t", "images")
		params.Set("count", fmt.Sprint(50))
		orderedKeys = append(orderedKeys, "t", "count")
	case "news":
		params.Set("t", "news")
		params.Set("count", fmt.Sprint(10))
		orderedKeys = append(orderedKeys, "t", "count")
	default:
		return "", fmt.Errorf("type de recherche non supporté : %s", searchType)
	}

	// Récupération du lien API correspondant au type de recherche
	link := QwantLINKS[q.session_type][searchType]

	// Envoi de la requête GET et gestion des erreurs éventuelles
	response, err := engine.DoGetRequest(link, params, orderedKeys, q.headers)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la requête : %v", err)
	}
	return string(response), nil
}

// controlOptions ajuste les options de recherche pour respecter les valeurs par défaut si nécessaire.
func (q *Qwant) controlOptions(options *engine.RequestOptions) {
	if options.MaxResults <= 0 {
		options.MaxResults = 20 // Valeur par défaut pour Maximum des résultats
	}
	if !engine.ContainsString(QwantLangs, options.Lang) {
		options.Lang = "en-GB" // Langue par défaut
	}
	if options.SafeSearch == "" {
		options.SafeSearch = "moderate" // Niveau de SafeSearch par défaut
	}
	if options.Freshness == "" {
		options.Freshness = QwantOPTIONS["freshness"]["undefined"] // Fraîcheur par défaut
	}
	if options.IndexPage < 0 {
		options.IndexPage = 0 // Valeur minimale pour IndexPage
	}
}

// setSafeSearch ajuste le niveau de SafeSearch selon le paramètre fourni.
func (q *Qwant) setSafeSearch(value string) string {
	switch value {
	case "off":
		return "0" // Désactiver SafeSearch
	case "moderate":
		return "1" // Niveau modéré
	case "strict", "on":
		return "2" // Niveau strict
	default:
		return "1" // Valeur par défaut : modéré
	}
}

// search effectue la recherche et retourne les résultats en fonction des options fournies.
func (q *Qwant) search(options engine.RequestOptions, searchType string, nbResultPage int) ([]engine.Result_Search, error) {
	q.controlOptions(&options) // Vérification et ajustement des options
	return engine.SearchGeneric(options, searchType, nbResultPage, q.requestSearch, q.processWebSearch)
}

// processWebSearch traite les résultats de recherche Web en décodant la réponse JSON.
func (q *Qwant) processWebSearch(jsonStr string, page int, searchType string) ([]engine.Result_Search, error) {
	finalResults := make([]engine.Result_Search, 0) // Résultats finaux à retourner
	var decodedData map[string]interface{}          // Structure de données pour le JSON décodé
	var position int = 1 * (page + 1)               // Position de départ pour le scoring

	// Décodage du JSON
	err := json.Unmarshal([]byte(jsonStr), &decodedData)
	if err != nil {
		return finalResults, fmt.Errorf("erreur de décodage JSON : %v", err)
	}

	// Récupération des résultats de la réponse
	results, ok := decodedData["data"].(map[string]interface{})["result"].(map[string]interface{})["items"]
	if !ok {
		return finalResults, fmt.Errorf("aucun résultat trouvé")
	}
	if searchType == "web" {
		results, ok = results.(map[string]interface{})["mainline"]
		if !ok {
			return finalResults, fmt.Errorf("aucun résultat trouvé")
		}
	}

	// Parcours des résultats et construction de la liste des résultats
	for _, data := range results.([]interface{}) {
		if searchType == "web" {
			Data, ok := data.(map[string]interface{})
			if !ok {
				continue
			}
			resultType, ok := Data["type"].(string)
			if !ok || resultType != "web" {
				continue
			}

			resultData, ok := Data["items"].([]interface{})
			if !ok {
				continue
			}
			for _, result := range resultData {
				var description string
				resultData := result.(map[string]interface{})
				title, ok := resultData["title"].(string)
				if !ok {
					continue
				}
				description, ok = resultData["desc"].(string)
				if !ok {
					continue
				}

				url, ok := resultData["url"].(string)
				if !ok {
					continue
				}
				finalResults = append(finalResults,
					engine.Result_Search{
						Item: engine.Item{

							Title:    engine.Normalize(title),                 // Normalisation du titre
							Body:     engine.Normalize(description),           // Normalisation de la description
							Url:      engine.NormalizeUrl(url),                // Normalisation de l'URL
							Source:   engine.ExtractDomain(url),               // Extraction du domaine source
							Engines:  []string{q.GetName()},                   // Marquer la source comme Qwant
							Position: position,                                // Position dans les résultats
							Score:    engine.Scoring(position, q.ponderation), // Calcul du score
						}})
				position += 1 // Incrémentation de la position
			}
		}
		if searchType == "images" || searchType == "news" {
			resultData, ok := data.(map[string]interface{})
			if !ok {
				continue
			}

			title, ok := resultData["title"].(string)
			if !ok {
				continue
			}
			url, ok := resultData["url"].(string)
			if !ok {
				continue
			}

			item := engine.Item{
				Title:    engine.Normalize(title),
				Url:      engine.NormalizeUrl(url),
				Source:   engine.ExtractDomain(url),
				Engines:  []string{"Qwant"},
				Position: position,
				Score:    engine.Scoring(position, q.ponderation)}

			if searchType == "news" {
				desc, ok := resultData["desc"].(string)
				if !ok {
					continue
				}
				item.Body = desc
				date, ok := resultData["date"].(int)
				if !ok {
					date = 0
				}
				item.Date = date // Enregistrement de la date pour les résultats de news
			}
			if searchType == "images" {
				img, ok := resultData["media"].(string)
				if !ok {
					continue
				}
				item.Img = img // Enregistrement de l'image

				height, ok := resultData["height"].(float64)
				if !ok {
					height = 0
				}

				width, ok := resultData["width"].(float64)
				if !ok {
					width = 0
				}
				item.Height = int32(math.Round(height)) // Assignation de la hauteur
				item.Width = int32(math.Round(width))   // Assignation de la largeur
			}

			finalResults = append(finalResults,
				engine.Result_Search{Item: item}) // Ajout de l'item à la liste des résultats
			position += 1
		}
	}

	return finalResults, nil // Retour des résultats finaux
}

// WebSearch effectue une recherche sur le web et retourne les résultats.
func (q *Qwant) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var max_result_per_page int = 10                     // Maximum de résultats par page
	return q.search(options, "web", max_result_per_page) // Appel à la méthode search
}

// ImagesSearch effectue une recherche d'images et retourne les résultats.
func (q *Qwant) ImagesSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var max_result_per_page int = 50 // Maximum de résultats par page pour les images
	return q.search(options, "images", max_result_per_page)
}

// NewsSearch effectue une recherche d'actualités et retourne les résultats.
func (q *Qwant) NewsSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var max_result_per_page int = 10 // Maximum de résultats par page pour les news
	return q.search(options, "news", max_result_per_page)
}
