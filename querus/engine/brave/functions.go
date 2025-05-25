package brave

import (
	"fmt"
	"regexp"

	"github.com/anaskhan96/soup"
)

// extract_json extrait un tableau JSON d'une réponse HTML.
// Il attend en entrée une chaîne de caractères `response`,
// qui contient le HTML à analyser, et retourne le contenu trouvé
// ou une erreur si la recherche échoue.

func extract_json(response string) (string, error) {
	// Analyser la réponse HTML en utilisant la bibliothèque Soup
	doc := soup.HTMLParse(response)

	// Trouver tous les éléments <script> dans le document HTML
	scripts := doc.FindAllStrict("script")

	// Prendre le dernier élément <script>
	script := scripts[len(scripts)-1].Text()
	fmt.Println(script)
	// Définir une expression régulière pour extraire le tableau JSON
	re := regexp.MustCompile(`const data = (\[.*?\]);`)

	// Exécuter l'expression régulière sur le contenu du <script>
	parsedJs := re.FindStringSubmatch(script)

	// Vérifier si un résultat a été trouvé
	if len(parsedJs) > 1 {
		// Retourner le tableau JSON extrait
		return parsedJs[1], nil
	}

	// Si aucune correspondance n'a été trouvée, retourner une erreur
	return "", fmt.Errorf("error extract_json")
}
