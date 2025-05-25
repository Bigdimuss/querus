package main

import (
	"runtime"
	"search/webclient"
)

/*
func main() {
	y := you.NewYouEngine(0.5)
	r, err := y.WebSearch(engine.RequestOptions{Query: "site:manjaro"})
	if err != nil {
		wrappedErr := fmt.Errorf("ERREUR ! : %w", err)
		fmt.Println(wrappedErr)
	}
	for key, v := range r {
		fmt.Println(key, " - ", v)
	}
	fmt.Println(len(r))
}
*/
/*
func main() {

	a := alltheinternet.NewAllTheInternetEngine(0.5)
	t, err := a.WebSearch(engine.RequestOptions{Query: "Chat", Lang: "fr-FR", SafeSearch: "off", IndexPage: 0, MaxResults: 15})
	if err != nil {
		fmt.Printf("error websearch function : %v\n", err)
	} else {
		for key, v := range t {
			fmt.Println(key, "-", v)
		}
		fmt.Println(len(t))
	}
}
*/
/*
func main() {
	y := qwant.NewQwantEngine(0.5)
	r, err := y.NewsSearch(engine.RequestOptions{Query: "Bonjour", Lang: "fr-FR", SafeSearch: "off", IndexPage: 0, MaxResults: 10})
	if err != nil {
		fmt.Println("ERREUR")
	}
	for key, v := range r {
		fmt.Println(key, " - ", v)
	}
	fmt.Println(len(r))
}*/

/*
func main() {
	y := brave.NewBraveEngine(0.5)
	r, err := y.WebSearch(engine.RequestOptions{Query: "France", Period: "undefined", Lang: "fr-FR", SafeSearch: "moderate", IndexPage: 0, MaxResults: 20})
	if err != nil {
		fmt.Println("ERREUR")
	}
	for key, v := range r {
		fmt.Println(key, " - ", v)
	}
	fmt.Println("nombre de resultat : ", len(r))
}*/

/*
func main() {
	y := carrot2.NewCarrot2Engine(0.5)
	r, err := y.WebSearch(engine.RequestOptions{Query: "Chien"})
	if err != nil {
		fmt.Println("Error")
	}
	for key, v := range r {
		fmt.Println(key, " - ", v)
	}
	fmt.Println("nombre de resultat : ", len(r))
}*/
/*
func main() {
	y := yep.NewYepEngine(0.5)
	r, err := y.WebSearch(engine.RequestOptions{Query: "manjaro", Lang: "fr-FR"})
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
	}
	for key, v := range r {
		fmt.Println(key, " - ", v)
	}
	fmt.Println("nombre de resultat : ", len(r))
}
*/
/*
func main() {
	d := duckduckgo.NewDuckDuckGoEngine(1.5)
	data, err := d.WebSearch(engine.RequestOptions{Query: "abricot", MaxResults: 150, Lang: "fr-FR"})
	if err != nil {
		fmt.Println("Error")
	}
	for key, d := range data {
		fmt.Println(key, " - ", fmt.Sprintf("%+v", d))
	}
}
*/
/*
func main() {
	d := mojeek.NewMojeek(1.5)
	data, err := d.WebSearch(engine.RequestOptions{Query: "Macron", MaxResults: 20, Lang: "fr-FR", IndexPage: 0, Country: "fr"})
	if err != nil {
		fmt.Println("Error")
	}
	for key, d := range data {
		fmt.Println(key, " - ", fmt.Sprintf("%+v", d))
	}
}*/

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	//client()
	webclient.RunServer()

}
