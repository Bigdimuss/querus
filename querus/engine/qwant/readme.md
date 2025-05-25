# Qwant Go Package

Le package Go pour interagir avec l'API de Qwant est conçu pour fournir un moteur de recherche respectueux de la vie privée. Ce projet vous permet de faire des recherches Web, d'images et d'actualités facilement et efficacement depuis vos applications Go.

## Table des Matières

- [Introduction](#introduction)
- [Installation](#installation)
- [Utilisation](#utilisation)
  - [Recherche Web](#recherche-web)
  - [Recherche d'Images](#recherche-dimages)
  - [Recherche d'Actualités](#recherche-dactualités)
- [Fonctionnalités](#fonctionnalités)
- [Configuration](#configuration)
- [Contributions](#contributions)
- [License](#license)
- [Contact](#contact)

## Introduction

Qwant est un moteur de recherche qui protège la vie privée de ses utilisateurs en ne collectant ni ne traçant leurs données personnelles. Ce package vous permet d'accéder à l'API de Qwant pour effectuer des recherches de manière simple et efficace.

## Installation

Pour installer ce package, ouvrez votre terminal et exécutez la commande suivante :

```bash
go get github.com/votre-utilisateur/qwant
```

Assurez-vous que votre environnement Go est configuré correctement avant d'exécuter cette commande.

# Utilisation
## Recherche Web
Voici comment réaliser une recherche sur le web :

```go
package main

import (
    "fmt"
    "github.com/votre-utilisateur/qwant"
)

func main() {
    engine := qwant.NewQwantEngine(1.0)
    options := engine.RequestOptions{
        Query:        "Votre recherche",
        Lang:         "fr-FR",
        SafeSearch:   "moderate",
        MaxResults:   20,
    }

    results, err := engine.WebSearch(options)
    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }

    for _, result := range results {
        fmt.Println(result.Item.Title, result.Item.Url)
    }
}
```
## Recherche d'Images
Pour effectuer une recherche d'images, utilisez le code suivant :

```go
imageResults, err := engine.ImageSearch(options)
if err != nil {
    fmt.Println("Erreur :", err)
    return
}

for _, img := range imageResults {
    fmt.Println(img.Item.Title, img.Item.Url)
}
```
## Recherche d'Actualités
Pour rechercher des actualités, vous pouvez utiliser :

```go
newsResults, err := engine.NewsSearch(options)
if err != nil {
    fmt.Println("Erreur :", err)
    return
}

for _, news := range newsResults {
    fmt.Println(news.Item.Title, news.Item.Url)
}
```
## Fonctionnalités
- Respect de la vie privée : Aucune collecte de données personnelles.
- Recherche polyvalente : Web, images et actualités.
- Filtrage des résultats : Options pour SafeSearch et langue.
- Pagination des résultats : Gestion facile de plusieurs pages de résultats.

## Contributions
Les contributions sont les bienvenues ! Pour contribuer :

## Forkez le projet.
- Créez une branche pour votre fonctionnalité (git checkout -b feature/votre-fonctionnalité).
- Committez vos changements (git commit -m 'Ajout d\'une nouvelle fonctionnalité').
- Poussez votre branche (git push origin feature/votre-fonctionnalité).
- Ouvrez une pull request.

## License
Ce projet est sous licence MIT. Consultez le fichier LICENSE pour plus de détails.

## Contact
Pour toute question ou suggestion, vous pouvez contacter le mainteneur à l'adresse suivante : votre-email@exemple.com.

Merci d'avoir choisi le package Qwant Go ! Happy coding!