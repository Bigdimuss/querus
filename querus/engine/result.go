package engine

// Item représente un élément de recherche avec diverses propriétés.
type Item struct {
	Title    string   `json:"title"`    // Titre de l'élément.
	Url      string   `json:"url"`      // URL de l'élément.
	Favicon  string   `json:"Favicon"`  // URL de l'icône de l'élément.
	Body     string   `json:"body"`     // Corps du contenu.
	Source   string   `json:"source"`   // Source de l'élément.
	Img      string   `json:"img"`      // URL de l'image associée.
	Height   int32    `json:"height"`   // Hauteur de l'image.
	Width    int32    `json:"width"`    // Largeur de l'image.
	Duration string   `json:"duration"` // Durée (s'il s'agit d'un média).
	Date     int      `json:"date"`     // Date de publication au format timestamp.
	Embed    string   `json:"embed"`    // Code d'intégration pour le contenu.
	Author   string   `json:"author"`   // Auteur de l'élément.
	Engines  []string `json:"engines"`  // Moteurs où cet élément est référencé.
	Position int      `json:"position"` // Position de l'élément dans les résultats.
	Score    float64  `json:"score"`    // Score d'évaluation de l'élément.
}

// Result_Search représente un résultat de recherche incluant des résultats connexes.
type Result_Search struct {
	Item                   // Embedding de la structure Item.
	Related_Results []Item `json:"related_results"` // Liste des résultats connexes.
}

// Scoring calcule un score basé sur la position et un coefficient de pondération.
func Scoring(position int, ponderation float64) float64 {
	// Renvoie le score calculé.
	return float64(1) / float64(position) * ponderation
}

// RemoveResult supprime un résultat de la liste des résultats en fonction d'un index donné.
func RemoveResult(result []Result_Search, index int) []Result_Search {
	// Renvoie une nouvelle slice sans l'élément à l'index spécifié.
	return append(result[:index], result[index+1:]...)
}
