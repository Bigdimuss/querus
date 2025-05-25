package engine

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Définit un délai d'attente de 5 secondes pour les requêtes HTTP.
var delay = time.Duration(3 * time.Second)

// DoRequest exécute une requête HTTP et retourne le corps de la réponse sous forme de tableau d'octets.
func DoRequest(request *http.Request) ([]byte, error) {
	// Création d'un client HTTP avec un délai d'attente spécifié.
	var client http.Client = http.Client{Timeout: delay}

	// Exécute la requête HTTP.
	response, err := client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("error : %v", err) // Retourne une erreur si la requête échoue.
	}

	// Vérifie si le statut de la réponse est OK (200).
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", response.StatusCode)
	} else {
		fmt.Println(response.StatusCode)
	}

	defer response.Body.Close() // Ferme le corps de la réponse à la fin de la fonction.

	var reader io.ReadCloser
	// Gère la compression gzip si elle est spécifiée dans l'en-tête de la réponse.
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error : %v", err)
		}
		defer reader.Close() // Ferme le lecteur gzip à la fin de la fonction.
	default:
		reader = response.Body // Utilise le corps de la réponse non compressé.
	}

	// Lit le corps de la réponse.
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}
	return body, nil // Retourne le corps sous forme de tableau d'octets.
}

// DoGetRequest exécute une requête GET avec des paramètres et des en-têtes supplémentaires.
func DoGetRequest(link string, params url.Values, orderedKeys []string, headers map[string]string) ([]byte, error) {
	var client http.Client = http.Client{Timeout: delay} // Création d'un client HTTP.

	// Crée l'URL complète avec les paramètres.
	link = Create_url_with_params(link, params, orderedKeys)

	// Crée une nouvelle requête GET.
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println(err)
	}

	Creat_headers(request, headers) // Ajoute les en-têtes à la requête.
	// Exécute la requête HTTP.
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making the GET request:", err)
		return nil, fmt.Errorf("error : %v", err)
	}

	// Vérifie le statut de la réponse.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", response.StatusCode)
	} else {
		fmt.Println(response.StatusCode)
	}

	defer response.Body.Close() // Ferme le corps de la réponse.

	var reader io.ReadCloser
	// Gère la compression gzip si nécessaire.
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error : %v", err)
		}
		defer reader.Close() // Ferme le lecteur gzip.
	default:
		reader = response.Body
	}

	// Lit le corps de la réponse.
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}
	return body, nil // Retourne le corps sous forme de tableau d'octets.
}

// DoPostData exécute une requête POST en ajoutant des données sous forme de paramètres URL.
func DoPostData(request *http.Request, data url.Values) (*http.Response, error) {
	var client http.Client = http.Client{Timeout: delay} // Création d'un client HTTP.

	// Encode les données en tant que paramètres de requête.
	request.URL.RawQuery = data.Encode()

	// Exécute la requête HTTP.
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}

	// Vérifie le statut de la réponse.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", response.StatusCode)
	}

	return response, nil // Retourne la réponse si tout est correct.
}
