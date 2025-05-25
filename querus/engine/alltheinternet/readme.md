# AllTheInternet Search Engine

## Description

Le package `alltheinternet` permet de réaliser des recherches via le moteur de recherche AllTheInternet. Il est conçu pour effectuer des recherches sur le web, des images et des nouvelles, tout en respectant les préférences de l'utilisateur.

## Fonctionnalités

- **Recherche Web** : Effectue des recherches standard sur le web.
- **Recherche d'Images** : Spécifiquement pour les résultats d'images.
- **Recherche d'Actualités** : Accède aux dernières nouvelles.
- **Personnalisation** : Prend en charge la pondération et le filtrage des résultats.

## Installation

Pour intégrer ce package dans votre projet Go, utilisez la commande suivante :

```bash
go get github.com/votre-repo/alltheinternet
```

## Utilisation

Créer une instance de AllTheInternet

```go
ati := NewAllTheInternetEngine(ponderation)
```

## Recherche sur le Web
```go
options := engine.RequestOptions{
    Query:       "Exemple de recherche",
    SafeSearch:  "moderate",
    Lang:        "fr-FR",
}

results, err := ati.WebSearch(options)
if err != nil {
    log.Fatal(err)
}
```
## Recherche d'Images
```go
imageResults, err := ati.ImagesSearch(options)
if err != nil {
    log.Fatal(err)
}
```
## Recherche d'Actualités

```go
newsResults, err := ati.NewsSearch(options)
if err != nil {
    log.Fatal(err)
}
```

## Structure des Résultats
Les résultats des recherches sont retournés sous forme de slice d'objets engine.Result_Search, contenant les informations suivantes :

- Title : Titre du résultat
- Body : Description du résultat
- Url : Lien vers le résultat
- Img : (pour les images) URL de l'image
- Position : Position du résultat dans la liste
- Score : Score basé sur la pondération

## Contribution
Si vous souhaitez contribuer au projet, ouvrez une pull request ou signalez un problème dans la section des issues sur le dépôt.

## Auteurs
Votre Nom
## License
Ce projet est sous licence MIT. Consultez le fichier LICENSE pour plus de détails.