package wikipedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	gowiki "github.com/trietmn/go-wiki"
)

type Wiki_Result struct {
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Url       string   `json:"url"`
	Thumbnail string   `json:"thumbnail"`
	Related   []string `json:"related"`
	Suggest   string   `json:"suggest"`
	Content   string   `json:"content"`
	Empty     bool     `json:"empty"`
}

func GetImageLink(pageTitle string, lang string, size int) (string, error) {
	// Construire l'URL de la requête
	apiURL := fmt.Sprintf("https://%v.wikipedia.org/w/api.php", lang[:2])
	params := url.Values{}
	params.Set("action", "query")
	params.Set("titles", pageTitle)
	params.Set("prop", "pageimages")
	params.Set("format", "json")
	params.Set("pithumbsize", fmt.Sprint(size))

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// Envoyer la requête GET
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Analyser la réponse JSON
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Extraire le lien de l'image
	query := result["query"].(map[string]interface{})
	pages := query["pages"].(map[string]interface{})

	var imageLink string
	for _, page := range pages {
		pageData := page.(map[string]interface{})
		thumbnail, ok := pageData["thumbnail"].(map[string]interface{})
		if ok {
			imageLink = thumbnail["source"].(string)
			break
		}
	}
	return imageLink, nil
}

func Search(query string, max_results int, lang string, suggestate bool) (Wiki_Result, error) {
	gowiki.SetLanguage(strings.ToLower(lang[:2]))
	search_result, suggest, err := gowiki.Search(query, max_results, suggestate)
	WikiPediaResult := Wiki_Result{Empty: true}
	if err != nil {
		return WikiPediaResult, err
	} else {
		if len(search_result) > 1 {
			result := search_result[0]
			// Get the page
			page, err := gowiki.GetPage(result, -1, false, true)
			if err != nil {
				return WikiPediaResult, err
			}

			imageURL, err := GetImageLink(page.Title, lang, 250)
			if err != nil {
				imageURL = ""
				fmt.Println(err)
			}
			summary, err := page.GetSummary()
			if err != nil {
				fmt.Println(err)
				WikiPediaResult.Summary = ""
			} else {
				WikiPediaResult.Summary = summary
			}
			WikiPediaResult.Title = page.OriginalTitle
			WikiPediaResult.Url = page.URL
			WikiPediaResult.Thumbnail = imageURL
			if len(search_result) > 1 {
				WikiPediaResult.Related = search_result[1:]
			}
			WikiPediaResult.Suggest = suggest
			WikiPediaResult.Empty = false
			return WikiPediaResult, nil
		} else {
			WikiPediaResult.Suggest = suggest
			return WikiPediaResult, err
		}

	}
}

/*
// Search for the Wikipedia page title

// Get the page
page, err := gowiki.GetPage("Rafael Nadal", -1, false, true)
if err != nil {
	fmt.Println(err)
}

// Get the content of the page
content, err := page.GetContent()
if err != nil {
	fmt.Println(err)
}
fmt.Printf("This is the page content: %v\n", content)*/
