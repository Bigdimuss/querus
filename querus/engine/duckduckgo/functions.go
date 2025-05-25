package duckduckgo

import (
	"fmt"
	"regexp"
)

// extractVqd extrait la valeur de vqd d'une chaîne HTML donnée.
// Prend en entrée une chaîne HTML (htmlStrings) et des mots-clés (keywords).
// Si la valeur de vqd est trouvée, elle est retournée avec nil pour l'erreur.
// Sinon, une erreur est retournée indiquant que l'extraction a échoué.
func extractVqd(htmlStrings string, keywords string) (string, error) {
	regexVqd := regexp.MustCompile(`vqd=([^&"]+)`)
	match := regexVqd.FindStringSubmatch(htmlStrings)
	if len(match) == 2 {
		return string(match[1]), nil
	}

	return "", fmt.Errorf("_extract_vqd() keywords=%s: Could not extract vqd", keywords)
}
