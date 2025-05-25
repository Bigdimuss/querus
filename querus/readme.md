# Querus Engine

## Description

**Querus** est un package Go conçu pour effectuer des recherches sur différents types de contenus (web, images, nouvelles) en utilisant plusieurs moteurs de recherche. Ce package se distingue par sa capacité de fusionner les résultats, d'éliminer les doublons, et de fournir des statistiques précises sur les performances de chaque moteur de recherche.

## Sommaire :

[Fonctionnalités](#fonctionnalités)

## Fonctionnalités

- **Recherche Multimoteur**: Effectuez des recherches sur plusieurs moteurs simultanément.
- **Fusion des Résultats**: Élimination des doublons et mise à jour des scores.
- **Analyse des Performances**: Suivez le délai de réponse pour chaque moteur de recherche.
- **Support de Wikipedia**: Intégration d'informations pertinentes à partir de Wikipedia.
- **Personnalisation**: Options de recherche flexibles pour adapter les requêtes selon les besoins.

## Installation

Pour installer Querus, clonez le dépôt et ajoutez-le à votre projet Go :

```bash
go get github.com/yourusername/querus
```

Assurez-vous également que les dépendances, telles que search/querus/engine et search/wikipedia, sont correctement configurées dans votre projet.

# Utilisation

## Importation

Commencez par importer le package dans votre code Go :

```go
import (
    "github.com/yourusername/querus"
    "search/querus/engine"
)
```

## Exemple de Recherche

Voici comment effectuer une recherche en utilisant le package Querus :

```go
package main

import (
    "fmt"
    "github.com/yourusername/querus"
    "search/querus/engine"
)

func main() {
    options := engine.RequestOptions{
        Query: "dernières nouvelles",
        Type:  "news", // Types valides : "web", "images", "news"
    }

    engines := []engine.Engine{ /* initialisez vos moteurs ici */ }
  
    results := querus.Search(options, engines, "asc") // "asc" ou "desc" pour trier
  
    // Afficher les résultats
    for _, result := range results.Results {
        fmt.Printf("Titre: %s, URL: %s, Score: %.2f\n", result.Item.Title, result.Item.Url, result.Score)
    }
  
    fmt.Printf("Nombre total de résultats: %d\n", results.Nb)
    fmt.Printf("Délai total: %.2f secondes\n", results.Delay)
}
```

## Structure des Résultats

La structure Result contient les informations suivantes :

```go
type Result struct {
    Type         string                 `json:"type"`
    Query        string                 `json:"query"`
    Engines      []string               `json:"engines"`
    EngineDelays []float64              `json:"enginedelays"`
    Delay        float64                `json:"delay"`
    Nb           int                    `json:"nb"`
    Wikipedia    wikipedia.Wiki_Result  `json:"wikipedia"`
    Results      []engine.Result_Search `json:"results"`
}
```

## Méthodes Utiles

- CompareAndMergeResults: Compare et fusionne les résultats d'une nouvelle recherche avec ceux déjà existants.
- processRelatedResults: Traite les résultats liés pour économiser de l'espace et ajouter des informations supplémentaires.
- searchEngine: Effectue la recherche sur un moteur spécifique selon le type de recherche souhaité.

## Gestion des Erreurs

Chaque méthode renvoie une erreur en cas d'échec. Il est recommandé de toujours vérifier et gérer ces erreurs pour garantir la robustesse de votre application.

## Contribution

Les contributions sont bienvenues ! Pour contribuer :

- Forkez le dépôt.
- Créez une branche (git checkout -b feature-branch).
- Commitez vos modifications (git commit -m 'Add some feature').
- Poussez la branche (git push origin feature-branch).
- Ouvrez une Pull Request.

## License

Ce projet est sous licence MIT. Consultez le fichier LICENSE pour plus de détails.

## Auteurs

Votre Nom

## Remerciements

Merci à tous les contributeurs et à la communauté Go pour leur soutien et leurs contributions.
