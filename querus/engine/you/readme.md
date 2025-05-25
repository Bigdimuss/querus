# YouCare Search Engine Package

## Description

Ce package Go permet d'interagir avec l'API de recherche YouCare. Il fournit des méthodes pour effectuer des recherches sur le web, des images, des nouvelles et plus encore. Il utilise la bibliothèque `json-iterator` pour le traitement des données JSON.

## Installation

Assurez-vous d'avoir Go installé sur votre système. Ajoutez le package à votre projet avec la commande suivante :

```bash
go get github.com/json-iterator/go
```

## Utilisation
Création d'un Nouvel Engine
Pour créer une nouvelle instance de You, utilisez la fonction suivante :

```go
youEngine := NewYouEngine(ponderation)
```
## Méthodes Disponibles

- GetName(): Retourne le nom de l'instance.

```go
name := youEngine.GetName()
```

- WebSearch(options engine.RequestOptions): Effectue une recherche web avec les options spécifiées.

```go
results, err := youEngine.WebSearch(options)
```

- ImagesSearch(options engine.RequestOptions): Effectue une recherche d'images.

```go
results, err := youEngine.ImagesSearch(options)
```
- NewsSearch(options engine.RequestOptions): Effectue une recherche de nouvelles.

```go
results, err := youEngine.NewsSearch(options)
```
## Options de Recherche
Les options de recherche doivent être spécifiées dans une struct engine.RequestOptions, qui doit contenir :

- Query: La requête de recherche.
- Lang: La langue de recherche (ex. "fr" pour le français).
- SafeSearch: Niveau de filtrage (ex. "moderate").

## Gestion des Erreurs
Chaque méthode retourne une erreur si une opération échoue. Assurez-vous de vérifier et de gérer les erreurs retournées par les méthodes.

### Exemples
Voici un exemple de recherche :

```go
options := engine.RequestOptions{
    Query: "dernières nouvelles",
    Lang:  "fr",
}

results, err := youEngine.WebSearch(options)
if err != nil {
    fmt.Println("Erreur lors de la recherche :", err)
} else {
    for _, result := range results {
        fmt.Println(result.Item.Title, result.Item.Url)
    }
}
```

## Contributeurs
Pour contribuer, ouvrez un problème ou soumettez une pull request sur le dépôt.

## License
Ce package est sous licence MIT. Consultez le fichier LICENSE pour plus de détails.