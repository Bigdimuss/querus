# Brave Search Engine

## Description

Ce paquet Go permet d'effectuer des recherches dans le moteur de recherche Brave. Il gère différents types de recherches tels que les résultats web, d'images et d'actualités, tout en respectant la vie privée de l'utilisateur.

## Fonctionnalités

- **Types de recherche** : Prise en charge des recherches web, d'images et d'actualités.
- **Pondération** : Possibilité de spécifier une pondération pour les résultats.
- **Options de fraîcheur** : Filtrage des résultats selon la fraîcheur (dernière heure, jour, semaine, etc.).
- **Langues** : Support pour plusieurs langues avec des préférences par défaut.

## Installation

Pour utiliser ce paquet, ajoutez-le à votre projet Go :


```bash
go get github.com/votre-repo/brave
```

## Utilisation
Créer une instance de Brave

```go
braveEngine := NewBraveEngine(ponderation)
```
## Effectuer une recherche web
```go

options := engine.RequestOptions{
    Query:       "Exemple de recherche",
    SafeSearch:  "moderate",
    Lang:       "fr-FR",
}

results, err := braveEngine.WebSearch(options)
if err != nil {
    log.Fatal(err)
}
```

## Recherche d'images

```go
imageResults, err := braveEngine.ImagesSearch(options)
if err != nil {
    log.Fatal(err)
}
```

## Recherche d'actualités

```go
newsResults, err := braveEngine.NewsSearch(options)
if err != nil {
    log.Fatal(err)
}
```
## Configuration des options
Vous pouvez configurer les options de recherche comme suit :

```go

options := engine.RequestOptions{
    MaxResults:  20,
    Period:      "at", // Options: "at", "ph", "pd", "pw", "pm", "py"
    SafeSearch:  "moderate",
    Lang:        "en-EN",
}
```

## Structure des résultats
Les résultats de recherche sont retournés sous forme de slice d'objets engine.Result_Search, comprenant des détails comme le titre, l'URL, la description, la source et les scores basés sur la position.

## Contributions
Si vous souhaitez contribuer, n'hésitez pas à créer une pull request sur le dépôt ou à déposer un problème dans l'onglet des issues.

## Auteurs
Votre Nom
License
Ce projet est sous licence MIT. Voir le fichier LICENSE pour plus de détails.