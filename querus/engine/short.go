package engine

import "fmt"

// SortInterface définit les méthodes nécessaires pour un tri personnalisé.
type SortInterface interface {
	Len() int           // Retourne le nombre d'éléments à trier.
	Less(i, j int) bool // Détermine si l'élément à l'index i est "moins" que celui à l'index j.
	Swap(i, j int)      // Échange les éléments aux index i et j.
}

// Sort est une structure qui contient une liste de résultats à trier et l'ordre de tri.
type Sort struct {
	Results []Result_Search // Liste des résultats à trier.
	Order   string          // Indique l'ordre de tri, soit "ascending" soit "descending".
}

// Len retourne le nombre total de résultats dans la structure Sort.
func (a *Sort) Len() int {
	return len(a.Results) // Retourne la longueur de la slice Results.
}

// Less compare deux éléments de la slice Results en fonction de l'ordre spécifié.
func (a *Sort) Less(i, j int) bool {
	// Utilise un switch pour déterminer l'ordre de tri.
	switch a.Order {
	case "ascending":
		// Pour un ordre croissant, compare les scores.
		return a.Results[i].Score > a.Results[j].Score
	case "descending":
		// Pour un ordre décroissant, compare les scores.
		return a.Results[i].Score < a.Results[j].Score
	default:
		// Pour tout autre cas, utilise l'ordre par défaut (croissant).
		fmt.Println("Use default order (ascending)")
		return a.Results[i].Score > a.Results[j].Score
	}
}

// Swap échange deux résultats à des index spécifiés dans la slice Results.
func (a *Sort) Swap(i, j int) {
	// Échange les éléments aux index i et j.
	a.Results[i], a.Results[j] = a.Results[j], a.Results[i]
}
