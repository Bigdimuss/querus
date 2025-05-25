package you

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"search/querus/engine"
	"search/querus/useragents"
	/* jsoniter "github.com/json-iterator/go"*/)

// Utilisation de la bibliothèque json-iterator pour le traitement JSON
/* var json = jsoniter.ConfigCompatibleWithStandardLibrary*/

// GLOBAL VALUE //

// URL des différentes API utilisées pour les recherches
var URLS map[string]string = map[string]string{
	"web":    "https://youcare.world/api/v2/search/all",
	"images": "https://youcare.world/api/v2/search/images",
	"news":   "https://youcare.world/api/v2/search/news",
	"videos": "https://youcare.world/api/v2/search/videos",
}

// URL pour initialiser le client
var INIT_URL string = "https://youcare.world/api/v2/client/initialize"

// STRUCT
// Définition de la structure pour l'API You
type You struct {
	user_agent  string
	Name        string
	ponderation float64
	ua          string
}

// Fonction de création d'un nouvel objet You avec un user agent aléatoire
func NewYouEngine(ponderation float64) *You {
	y := &You{user_agent: useragents.Get_random_ua(), Name: "You", ponderation: ponderation}
	return y
}

// Retourne le nom de l'instance
func (y *You) GetName() string {
	return y.Name
}

// Méthode pour obtenir un UUID à partir de l'API YouCare
func (y *You) get_u() (string, error) {
	// Définition des en-têtes pour la requête
	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-encoding": "gzip, deflate, br",
		"Accept-language": "fr,en-US;q=0.9,en;q=0.8",
		"Client-version":  "1.38.4",
		"Content-length":  "72",
		"Content-type":    "application/json",
		"Origin":          "https://youcare.world",
		"Referer":         "https://youcare.world/",
		"User-Agent":      y.user_agent,
	}

	// Paramètres de la requête JSON
	params := map[string]interface{}{
		"uuid":               "",
		"activeGoodDeedType": "",
		"parameters": map[string]string{
			"language": "fr",
		},
	}

	// Convertit les paramètres en JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		panic(err) // Erreur fatale si la conversion échoue
	}

	// Crée une requête POST avec les données JSON
	request, err := http.NewRequest("POST", INIT_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err) // Erreur fatale lors de la création de la requête
	}

	// Ajoute les en-têtes à la requête
	engine.Creat_headers(request, headers)

	// Envoie la requête et gère la réponse
	response, err := engine.DoRequest(request)
	if err != nil {
		return "", fmt.Errorf("cant get response : %v", err)
	}
	responseString := string(response)

	// Décodage de la réponse JSON
	var result map[string]interface{}
	json.Unmarshal([]byte(responseString), &result)
	uuid, ok := result["user"].(map[string]interface{})["uuid"].(string)
	if !ok {
		return "", fmt.Errorf("can't extract uuid : %v", response)
	}
	return uuid, nil
}

// Méthode pour effectuer une recherche en fonction des options et du type de recherche
func (y *You) requestSearch(options engine.RequestOptions, offset int, searchType string) (string, error) {
	// Paramètres de la requête
	var orderedKeys []string
	params := url.Values{}
	params.Set("q", options.Query)
	params.Set("o", string(fmt.Sprint(offset)))
	orderedKeys = append(orderedKeys, "q", "o")
	// En-têtes de la requête
	var headers map[string]string = map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": fmt.Sprintf("%s,en-US;q=0.9,en;q=0.8", options.Lang[:2]),
		"Client-Version":  "1.38.4",
		"Referer": engine.Creat_url_with_params(
			URLS["search"], params, orderedKeys),
		"User-Agent": y.user_agent,
	}

	params.Set("m", options.Lang)
	params.Set("f", options.SafeSearch)
	params.Set("u", y.ua)
	params.Set("s", "")
	params.Set("g", "null")
	orderedKeys = append(orderedKeys, "m", "f", "u", "s", "g")
	// Envoie la requête GET et gère les erreurs
	response, err := engine.DoGetRequest(URLS[searchType], params, orderedKeys, headers)
	if err != nil {
		return "", fmt.Errorf("cant requestsearch : %v", err)
	}
	responseString := string(response)
	return responseString, nil
}

// Méthode pour contrôler et ajuster les options de recherche
func (y *You) controlOptions(options *engine.RequestOptions) {
	if options.Lang == "" {
		options.Lang = "en-EN" // Définit la langue par défaut
	}
	if options.SafeSearch == "" {
		options.SafeSearch = "moderate" // Définit le niveau de filtre par défaut
	}
	if options.IndexPage < 1 {
		options.IndexPage = 1 // Définit l'index de la page par défaut
	}
	if options.MaxResults == 0 || options.MaxResults < 10 {
		options.MaxResults = 10 // Définit le nombre maximal de résultats par défaut
	}
}

// Méthode pour effectuer une recherche web
func (y *You) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return y.search(options, "web", 10) // Recherche par défaut de 10 résultats
}

// Méthode pour effectuer une recherche d'images
func (y *You) ImagesSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return y.search(options, "images", 45) // Recherche d'images avec une limite de 45 résultats
}

// Méthode pour effectuer une recherche de nouvelles
func (y *You) NewsSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return y.search(options, "news", 10) // Recherche de nouvelles avec 10 résultats par défaut
}

// Méthode générale de recherche
func (y *You) search(options engine.RequestOptions, searchType string, nbResultPage int) ([]engine.Result_Search, error) {
	y.controlOptions(&options) // Ajuste les options si nécessaire
	ua, err := y.get_u()       // Obtient l'UUID
	if err != nil {
		return make([]engine.Result_Search, 0), fmt.Errorf("can't get ua : %v", err)
	}
	y.ua = ua // Assigne l'UUID à l'instance
	return engine.SearchGeneric(options, searchType, nbResultPage, y.requestSearch, y.processWebSearch)
}

// Méthode pour traiter les résultats de recherche web
func (y *You) processWebSearch(jsonStr string, page int, searchType string) ([]engine.Result_Search, error) {
	// Décodage du JSON en une structure de données dynamique
	var decodedData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &decodedData)

	finalResults := make([]engine.Result_Search, 0) // Liste pour stocker les résultats finaux
	if err != nil {
		return finalResults, fmt.Errorf("erreur de décodage JSON : %v", err)
	}
	fmt.Println("RESULT :")
	fmt.Println(decodedData)
	// Récupération des résultats de type "web"
	results, ok := decodedData["results"].([]interface{})
	if !ok {
		return finalResults, fmt.Errorf("aucun résultat trouvé") // Gestion des résultats manquants
	}

	// Parcours des résultats et affichage des résultats de type "web"
	var position int = 1 * (page + 1) // Position de départ pour les résultats

	for _, result := range results {
		resultData, ok := result.(map[string]interface{})
		fmt.Println(resultData)
		if !ok {
			continue // Ignore les résultats invalides
		}

		if searchType == "web" || searchType == "news" {
			resultType, ok := resultData["type"].(string)
			if !ok || resultType != searchType {
				continue // Ignore les types de résultat non pertinents
			}

			// Récupération des informations des résultats
			title, ok := resultData["title"].(string)
			if !ok {
				continue
			}

			description, ok := resultData["description"].(string)
			if !ok {
				continue
			}

			url, ok := resultData["url"].(string)
			if !ok {
				continue
			}

			// Ajoute le résultat final à la liste
			finalResults = append(finalResults,
				engine.Result_Search{
					Item: engine.Item{
						Title:    engine.Normalize(title),
						Body:     engine.Normalize(description),
						Url:      engine.NormalizeUrl(url),
						Source:   engine.ExtractDomain(url),
						Engines:  []string{"You"},
						Position: position,
						Score:    engine.Scoring(position, y.ponderation),
					}})
			position += 1 // Incrémente la position pour le prochain résultat
		}

		// Traitement spécifique pour les images
		if searchType == "images" {
			resultType, ok := resultData["type"].(string)
			if !ok || resultType != "image" {
				continue // Ignore les résultats non-image
			}
			url, ok := resultData["hostPageUrl"].(string)
			if !ok {
				continue
			}
			img, ok := resultData["url"].(string)
			if !ok {
				continue
			}
			title, ok := resultData["name"].(string)
			if !ok {
				continue
			}
			width, ok := resultData["width"].(float64)
			if !ok {
				width = 0 // Défaut à 0 si la largeur n'est pas spécifiée
			}
			height, ok := resultData["height"].(float64)
			if !ok {
				height = 0 // Défaut à 0 si la hauteur n'est pas spécifiée
			}

			// Ajoute le résultat de l'image à la liste
			finalResults = append(finalResults,
				engine.Result_Search{
					Item: engine.Item{
						Title:    engine.Normalize(title),
						Img:      engine.NormalizeUrl(img),
						Url:      engine.NormalizeUrl(url),
						Source:   engine.ExtractDomain(url),
						Height:   int32(math.Round(height)),
						Width:    int32(math.Round(width)),
						Engines:  []string{"You"},
						Position: position,
						Score:    engine.Scoring(position, y.ponderation),
					}})
			position += 1 // Mise à jour de la position pour l'image
		}
	}

	// Retourne les résultats finaux ou une erreur s'il n'y a aucun résultat
	if len(finalResults) > 0 {
		return finalResults, nil
	}
	return finalResults, fmt.Errorf("error") // Erreur si aucun résultat
}
