package webclient

import (
	"fmt"
	"html/template"
	"strings"
)

// Fonction pour tronquer le texte sans couper un mot
func truncate(text string, length int) string {
	if len(text) <= length {
		return text
	}

	// Chercher le dernier espace avant la limite
	lastSpace := text[:length]
	if lastSpaceIndex := len(lastSpace) - 1; lastSpaceIndex >= 0 {
		if lastSpace[lastSpaceIndex] != ' ' {
			if index := lastSpaceIndex; index > 0 {
				// Trouver le dernier espace dans le texte
				for i := lastSpaceIndex; i >= 0; i-- {
					if lastSpace[i] == ' ' {
						return text[:i] + "..."
					}
				}
			}
		}
	}

	// Si aucun espace n'est trouvé, retourner le texte tronqué à la longueur spécifiée
	return text[:length] + "..."
}

func RenderEnginesfunc(engines []string) template.HTML {
	var result strings.Builder
	nbengines := len(engines)
	if nbengines > 0 {
		result.WriteString(`<div class="engines"><i>`)
		engineString := ""
		for key, engine := range engines {
			engineString += fmt.Sprintf("%v", engine)
			if key < nbengines-1 {
				engineString += " - "
			}
		}
		result.WriteString(template.HTMLEscapeString(engineString))
		result.WriteString(`</i></div>`)

	}
	return template.HTML(result.String())
}
