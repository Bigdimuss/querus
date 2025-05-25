package alltheinternet

import "regexp"

// extract_json traite la reponse de l'api en extrayant la partie JSON
func extract_json(response *string) {
	texte := *response
	apiDebut := "/*O_o*/\ngoogle.search.cse.api("
	apiFin := "\n);"
	// Modifie la réponse pour remplacer le format illisible par un JSON valide
	*response = texte[len(apiDebut):len(texte)-len(apiFin)] + "}"
}

// do_extract utilise une expression régulière pour trouver la valeur d'un clé spécifique dans la réponse JSON
func do_extract(response string, value string) string {

	// compile une rexpression régulière pour capturer les valeurs associées à 'value'.
	regex := regexp.MustCompile(`"` + value + `"\s*:\s*"([^"]+)"`)
	// Cherche les correspondances dans la réponse.
	match := regex.FindStringSubmatch(string(response))
	var cxValue string
	// Si correspondance, récupère la valeur.
	if len(match) > 1 {
		cxValue = match[1]
	}
	return cxValue // Renvoie la valeur trouver ou une chaine vide.
}

// extract_key extrait les valeurs des clés spécifiques de la réponse JSON et les renvoie sous forme de map.
func extract_keys(response string) map[string]string {

	// Initialise un map pour associer des clés à leurs noms dans la réponse.
	values := map[string]string{"cx": "cx", "cse_tok": "cse_token", "cselibv": "cselibVersion"}

	// Parcourt chaque clé pour extraire ses valeurs à l'aide de do_extract.
	for key, v := range values {
		value := do_extract(response, v)
		values[key] = value // Met à jour le map avec la valeur trouvée.
	}
	return values // Renvoie le map contenant les clés et leurs valeurs associée.
}
