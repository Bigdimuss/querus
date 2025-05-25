package mojeek

import (
	"fmt"
	"net/url"
	"search/querus/engine"
	"search/querus/useragents"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
)

type Mojeek struct {
	user_agent  string            // User-agent pour les requêtes HTTP
	headers     map[string]string // En-têtes HTTP nécessaires pour les requêtes
	Name        string            // Nom du moteur de recherche
	ponderation float64           // Ponderation pour le scoring des résultats
}

func NewMojeek(ponderation float64) *Mojeek {
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language": "fr,en-US;q=0.5,en;q=0.3",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Host":            "www.mojeek.com",
	}

	m := Mojeek{user_agent: useragents.Get_random_ua(), Name: "Mojeek", headers: headers, ponderation: ponderation}
	return &m
}

func (mojeek *Mojeek) GetName() string {
	return mojeek.Name
}

func (mojeek *Mojeek) SetLang(lang string) {
	mojeek.headers["Accept-Language"] = fmt.Sprintf("%s,%s,en-US;q=0.9,en;q=0.8", lang[:2], lang)
}
func (mojeek *Mojeek) controlOptions(options *engine.RequestOptions, searchType string) {
	if options.MaxResults <= 0 {
		options.MaxResults = 10 // Valeur par défaut pour Maximum des résultats
	}
	if options.IndexPage <= 1 {
		options.IndexPage = 0 // Valeur minimale pour IndexPage
	}
}

func (mojeek *Mojeek) requestSearch(options engine.RequestOptions, offset int, searchType string) (string, error) {
	params := url.Values{}
	var orderedKeys []string

	params.Set("q", options.Query)
	orderedKeys = append(orderedKeys, "q")
	cookies := "tlen=128; dlen=250"
	if offset > 0 {
		params.Set("s", fmt.Sprint(offset))
		orderedKeys = append(orderedKeys, "s")
	}
	if options.Country != "" {
		params.Set("reg", strings.ToLower(options.Country[:2]))
		orderedKeys = append(orderedKeys, "reg")
	}
	if options.Lang != "" {
		cookies += fmt.Sprintf("; %s", strings.ToLower(options.Lang[:2]))
	}
	if (searchType == "images") || (searchType == "news") {
		params.Set("fmt", searchType)
		orderedKeys = append(orderedKeys, "fmt")
	}
	mojeek.headers["Cookie"] = cookies

	baseUrl := "https://www.mojeek.com/search"
	mojeek.headers["User-Agent"] = mojeek.user_agent
	data, err := engine.DoGetRequest(baseUrl, params, orderedKeys, mojeek.headers)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(data), nil
}

func (mojeek *Mojeek) processWebSearch(resultString string, page int, searchType string) ([]engine.Result_Search, error) {
	finalResults := make([]engine.Result_Search, 0)
	var position int = 1 * (page + 1)

	if resultString == "" {
		return finalResults, fmt.Errorf("resultString is empty")
	}
	// Parser le document HTML
	doc := soup.HTMLParse(resultString)
	if doc.Error != nil {
		return finalResults, fmt.Errorf("document empty")
	}

	// Sélectionner des éléments et récupérer des informations
	titles := doc.Find("ul", "class", "results-standard").FindAll("li")
	if len(titles) == 0 {
		return finalResults, fmt.Errorf("no results found")
	}
	for _, result := range titles {
		if searchType == "web" {
			title := result.Find("h2").FullText()
			content := result.Find("p", "class", "s").FullText()
			url := result.Find("a", "class", "title").Attrs()["href"]
			if (title != "") && (content != "") && (url != "") {
				resultItem := engine.Item{
					Title:    title,
					Body:     content,
					Url:      url,
					Source:   engine.ExtractDomain(url),
					Engines:  []string{mojeek.GetName()},
					Position: position,
					Score:    engine.Scoring(position, mojeek.ponderation)}

				finalResults = append(finalResults, engine.Result_Search{Item: resultItem})
				position += 1

			} else {
				continue
			}
		} else {
			continue
		}
	}
	if len(finalResults) < 1 {
		return finalResults, fmt.Errorf("result : %s empty", mojeek.GetName())
	}
	return finalResults, nil
}

func (mojeek *Mojeek) search(options engine.RequestOptions, searchType string, nb_results_by_request int) ([]engine.Result_Search, error) {

	mojeek.controlOptions(&options, searchType) // Contrôle des options

	var offset int
	var page_nb int = options.MaxResults / nb_results_by_request // Calcul du nombre de pages
	if options.MaxResults%nb_results_by_request != 0 {
		page_nb += 1 // Ajout d'une page si nécessaire
	}
	var results []engine.Result_Search = make([]engine.Result_Search, 0)
	var resultsMutex sync.Mutex // Mutex pour la protection des résultats
	var wg sync.WaitGroup       // WaitGroup pour la gestion des goroutines
	for i := options.IndexPage; i < options.IndexPage+page_nb; i++ {
		offset = i * nb_results_by_request // Calculer l'offset
		wg.Add(1)                          // Incrémente le compteur de waitgroup
		go func(offset int, i int, results *[]engine.Result_Search, wg *sync.WaitGroup) {
			defer wg.Done() // Indique que la goroutine est terminée

			r, err := mojeek.requestSearch(options, offset, searchType) // Requête de recherche
			if err != nil {
				fmt.Println(err) // Log de l'erreur
				return
			}
			v, err := mojeek.processWebSearch(r, page_nb, searchType) // Traitement des résultats
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

func (mojeek *Mojeek) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	var nb_result_page int = 10
	mojeek.controlOptions(&options, "web")
	return mojeek.search(options, "web", nb_result_page)
}
func (Mojeek *Mojeek) ImagesSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return make([]engine.Result_Search, 0), nil
}
func (Mojeek *Mojeek) NewsSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	return make([]engine.Result_Search, 0), nil
}

/*
func (m *Mojeek) WebSearch(options engine.RequestOptions) ([]engine.Result_Search, error) {
	results, err := m.requestSearch(options, offset, searchType)
	if err != nil {
		fmt.Println(err)
		return make([]engine.Result_Search, 0), err
	}
	fmt.Println(results)
	return make([]engine.Result_Search, 0), nil
}
*/
