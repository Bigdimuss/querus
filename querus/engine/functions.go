package engine

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// REGEX_STRIP_TAGS est une expression régulière utilisée pour supprimer les balises HTML.
var REGEX_STRIP_TAGS = regexp.MustCompile("<.*?>")

// Normalize prend une chaîne de caractères contenant du HTML brut et renvoie une version sans balises HTML.
func Normalize(rawHTML string) string {
	// Vérifie si la chaîne est vide.
	if rawHTML == "" {
		return ""
	}
	// Supprime les balises HTML et effectue une désévasion des entités HTML.
	return html.UnescapeString(REGEX_STRIP_TAGS.ReplaceAllString(rawHTML, ""))
}

// NormalizeUrl remplace les espaces dans une URL par des signes plus.
func NormalizeUrl(url string) string {
	return strings.ReplaceAll(url, " ", "+")
}

func Creat_url_with_params(baseURL string, values url.Values, orderKeys []string) string {
	// Crée une nouvelle instance de url.Values
	// Construit l'URL finale
	finalURL, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	finalURL.RawQuery = values.Encode()
	fmt.Println(finalURL.String())
	return finalURL.String()
}

func Create_url_with_params(baseURL string, params url.Values, orderedKeys []string) string {
	var builder strings.Builder

	// Ajoute chaque paramètre dans l'ordre spécifié.
	for i, key := range orderedKeys {
		if values, exists := params[key]; exists {
			for _, value := range values {
				if i > 0 || builder.Len() > 0 {
					builder.WriteString("&")
				}
				builder.WriteString(fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
			}
		} else {
			if i > 0 || builder.Len() > 0 {
				builder.WriteString("&")
			}
			builder.WriteString(fmt.Sprintf("%s=", url.QueryEscape(key)))
		}
	}

	// Construit l'URL finale.
	if builder.Len() > 0 {
		fmt.Printf("%s?%s", baseURL, builder.String())
		return fmt.Sprintf("%s?%s", baseURL, builder.String())
	}
	return baseURL // Retourne simplement le baseURL si aucun paramètre n'est présent.
}

// Creat_headers ajoute des en-têtes personnalisés à une requête HTTP.
func Creat_headers(request *http.Request, headers map[string]string) {
	if request == nil {
		return
	}
	if headers == nil {
		return
	}
	// Parcourt chaque paire clé-valeur d'en-têtes et les ajoute à la requête.
	for key, value := range headers {
		request.Header.Set(key, value)
	}
}

// ExtractDomain extrait le nom de domaine d'une chaîne ressemblant à une URL.
func ExtractDomain(urlLikeString string) string {
	// Supprime les espaces inutiles de la chaîne.
	urlLikeString = strings.TrimSpace(urlLikeString)

	// Vérifie si la chaîne commence par http ou https.
	if regexp.MustCompile(`^https?`).MatchString(urlLikeString) {
		read, _ := url.Parse(urlLikeString)
		urlLikeString = read.Host // Récupère le domaine.
	}

	// Supprime le préfixe www. si présent.
	if regexp.MustCompile(`^www\.`).MatchString(urlLikeString) {
		urlLikeString = regexp.MustCompile(`^www\.`).ReplaceAllString(urlLikeString, "")
	}

	// Renvoie le nom de domaine trouvé.
	return regexp.MustCompile(`([a-z0-9\-]+\.)+[a-z0-9\-]+`).FindString(urlLikeString)
}

// AddCookies ajoute des cookies à une requête HTTP.
func AddCookies(request *http.Request, cookies map[string]string) {
	for k, v := range cookies {
		// Crée un cookie et l'ajoute à la requête.
		request.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
}

// GetType renvoie le nom du type d'une variable passée en paramètre.
func GetType(myvar interface{}) string {
	// Récupère le type de la variable et gère les pointeurs.
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name() // Renvoie le nom du type pointeur.
	} else {
		return t.Name() // Renvoie le nom du type normal.
	}
}

// ContainsString vérifie si une chaîne de caractères est présente dans un slice de chaînes.
func ContainsString(slice []string, searchString string) bool {
	for _, s := range slice {
		if s == searchString {
			return true // Retourne vrai si la chaîne est trouvée.
		}
	}
	return false // Retourne faux si la chaîne n'est pas trouvée.
}

type Ponderations map[string]float64

func SetPonderation(key string, ponderation float64, base *Ponderations) {
	if 0.0 > ponderation && ponderation > 1.0 {
		ponderation = 0.5
	}
	(*base)[key] = ponderation
}
