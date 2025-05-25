package engine

import (
	"fmt"
	"sync"
)

// AvailableEngines est un type pour une liste d'engines.
type AvailableEngines []Engine

// Engine définit une interface pour les moteurs de recherche.
// Elle inclut des méthodes pour la recherche web, d'images, et de nouvelles, ainsi qu'une méthode pour obtenir le nom de l'engin.
type Engine interface {
	WebSearch(RequestOptions) ([]Result_Search, error)
	ImagesSearch(RequestOptions) ([]Result_Search, error)
	NewsSearch(RequestOptions) ([]Result_Search, error)
	GetName() string
}

// RequestOptions contient les paramètres pour effectuer une recherche.
// Elle inclut des informations telles que le type de recherche, la requête, la langue,
// la sécurité de recherche, le pays, la fraîcheur, la période, les sources, le numéro de page et le nombre maximum de résultats.
type RequestOptions struct {
	Type       string
	Query      string
	Lang       string
	SafeSearch string
	Country    string
	Freshness  string
	Period     string
	Sources    string
	IndexPage  int
	MaxResults int
	Related    bool
	Engines    []string
}

// SearchGeneric effectue une recherche générique basée sur les options fournies.
// Elle divise les résultats sur plusieurs pages en gérant la concurrence via goroutines.
func SearchGeneric(
	options RequestOptions,
	searchType string,
	maxResultsPerPage int,
	requestSearchFunc func(options RequestOptions, offset int, searchType string) (string, error),
	processWebSearchFunc func(result string, offset int, searchType string) ([]Result_Search, error)) ([]Result_Search, error) {

	// Calcule le nombre de pages nécessaires en fonction du nombre maximum de résultats.
	var page_nb int = (options.MaxResults + maxResultsPerPage - 1) / maxResultsPerPage
	if options.MaxResults > page_nb*maxResultsPerPage {
		page_nb += 1
	}

	// Initialise un slice pour stocker les résultats de la recherche.
	var results []Result_Search = make([]Result_Search, 0)

	// Vérification de la sécurité des accès aux résultats avec un mutex.
	var resultsMutex sync.Mutex
	var wg sync.WaitGroup

	// Lance des goroutines pour effectuer la recherche sur plusieurs pages.
	for i := options.IndexPage; i < options.IndexPage+page_nb; i++ {
		offset := i * maxResultsPerPage
		wg.Add(1) // Déclare qu'une goroutine sera ajoutée au groupe.
		go func(options RequestOptions, offset int, results *[]Result_Search, wg *sync.WaitGroup) {
			defer wg.Done() // Indique que la goroutine est terminée une fois qu'elle a réussi à faire son travail.

			// Effectue la requête de recherche.
			r, err := requestSearchFunc(options, offset, searchType)
			if err != nil {
				fmt.Printf("Error in requestSearch: %v\n", err)
				return
			}
			// Traite les résultats de la recherche.
			v, err := processWebSearchFunc(r, offset, searchType)
			if err != nil {
				fmt.Printf("Error in processWebSearch: %v\n", err)
				return
			}

			// Protéger l'accès aux résultats partagés avec un mutex.
			resultsMutex.Lock()
			(*results) = append((*results), v...)
			resultsMutex.Unlock()
		}(options, offset, &results, &wg)
	}

	// Attend que toutes les goroutines terminent leur exécution.
	wg.Wait()

	// Renvoyer les résultats ou une erreur si aucun résultat n'est trouvé.
	if len(results) > 0 {
		return results, nil
	}
	return results, fmt.Errorf("error during search")
}
