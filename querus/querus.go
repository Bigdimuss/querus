package querus

import (
	"fmt"
	"search/querus/engine"
	"search/wikipedia"
	"sort"
	"sync"
	"time"
)

// Result représente les résultats d'une recherche.
type Result struct {
	Type         string                 `json:"type"`         // Type de recherche (web, images, news)
	Query        string                 `json:"query"`        // La requête utilisée pour la recherche
	Engines      []string               `json:"engines"`      // Liste des moteurs de recherche utilisés
	EngineDelays map[string]float64     `json:"enginedelays"` // Délai de réponse pour chaque moteur
	Delay        float64                `json:"delay"`        // Délai total de la recherche
	Nb           int                    `json:"nb"`           // Nombre total de résultats
	Wikipedia    wikipedia.Wiki_Result  `json:"wikipedia"`    // Résultat retourné par Wikipedia
	Results      []engine.Result_Search `json:"results"`      // Liste des résultats de recherche
	IndexPage    int                    `json:"index_page"`
}

// CompareAndMergeResults compare et fusionne les résultats cibles avec la liste existante.
func CompareAndMergeResults(cibles []engine.Result_Search, results *[]engine.Result_Search) {
	for _, cible := range cibles {
		foundIndex := -1
		for i, existing := range *results {
			if existing.Url == cible.Url {
				foundIndex = i
				break
			}
		}

		if foundIndex != -1 {
			// Mise à jour des moteurs de recherche et du score
			for _, engine := range cible.Engines {
				if !contains((*results)[foundIndex].Engines, engine) {
					(*results)[foundIndex].Engines = append((*results)[foundIndex].Engines, engine)
					(*results)[foundIndex].Score += cible.Score
				}
			}
		} else {
			*results = append(*results, cible)
		}
	}
}

// contains vérifie si une valeur existe dans un tableau de chaînes.
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Search effectue une recherche en utilisant les options et les moteurs fournis.
func Search(options engine.RequestOptions, engines []engine.Engine, order string) Result {
	results := Result{Query: options.Query, EngineDelays: make(map[string]float64)}
	start := time.Now()

	var wg sync.WaitGroup
	var resultsMutex sync.Mutex
	if options.Query != "" {
		for _, e := range engines {
			wg.Add(1)
			go func(en engine.Engine) {
				defer wg.Done()
				start := time.Now()
				engineResults, err := searchEngine(en, options)
				if err != nil {
					fmt.Printf("Erreur avec : %v\n", en.GetName())
					return
				}

				resultsMutex.Lock()
				CompareAndMergeResults(engineResults, &results.Results)
				results.Engines = append(results.Engines, en.GetName())
				results.EngineDelays[en.GetName()] = time.Since(start).Seconds()
				resultsMutex.Unlock()
			}(e)
		}
		fmt.Println("QUERUS WAIT")
		wg.Wait()
		if options.Related {
			// Traitement des résultats liés
			processRelatedResults(&results)
		}
		// Tri des résultats par score
		sort.Sort(&engine.Sort{Results: results.Results, Order: order})

		// Attribution des positions et calcul du délai total
		for i := range results.Results {
			results.Results[i].Position = i + 1
		}
		results.Nb = len(results.Results)
		results.Delay = time.Since(start).Seconds()
	}
	return results
}

// searchEngine effectue une recherche sur un moteur donné.
func searchEngine(en engine.Engine, options engine.RequestOptions) ([]engine.Result_Search, error) {
	switch options.Type {
	case "web":
		return en.WebSearch(options)
	case "images":
		return en.ImagesSearch(options)
	case "news":
		return en.NewsSearch(options)
	default:
		return nil, fmt.Errorf("type de recherche inconnu: %s", options.Type)
	}
}

// processRelatedResults traite les résultats liés et les fusionne en cas de doublons.
// processRelatedResults traite les résultats liés et les fusionne en cas de doublons.
func processRelatedResults(results *Result) {
	for v := 0; v < len(results.Results); v++ {
		for i := 0; i < len(results.Results); i++ {
			if i != v && results.Results[i].Source == results.Results[v].Source {
				if results.Results[v].Score-0.1*float64(len(results.Results[v].Related_Results)) > results.Results[i].Score {
					results.Results[v].Related_Results = append(results.Results[v].Related_Results, results.Results[i].Item)
				} else {
					// Fusion des résultats
					results.Results[i].Related_Results = append(results.Results[i].Related_Results, results.Results[v].Item)
					results.Results[i].Related_Results = append(results.Results[i].Related_Results, results.Results[v].Related_Results...)
					results.Results[i].Score += 0.1 * float64(len(results.Results[i].Related_Results))
					results.Results[v] = results.Results[i]
				}
				results.Results = append(results.Results[:i], results.Results[i+1:]...) // Suppression de l'élément
				i--                                                                     // Décrémente i pour rester au même index après la suppression
			}
		}
	}
}
