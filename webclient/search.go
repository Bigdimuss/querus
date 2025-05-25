package webclient

import (
	"fmt"
	"search/querus"
	"search/querus/engine"
	"search/querus/engine/alltheinternet"
	"search/querus/engine/duckduckgo"
	"search/querus/engine/mojeek"
	"search/querus/engine/qwant"
	"search/querus/engine/yep"
	"search/querus/engine/you"
	"search/wikipedia"
	"sync"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
)

type NamedEngine struct {
	Name   string
	Engine engine.Engine
}

func generateEnginesList(engines *[]string) engine.AvailableEngines {

	var PonderationsConfig engine.Ponderations = engine.Ponderations{
		"AllTheInternet": 0.52,
		//"Brave":          0.55,
		"Carrot2":    0.3,
		"DuckDuckGo": 0.6,
		"Qwant":      0.56,
		"Yep":        0.42,
		"You":        0.42,
		"Mojeek":     0.4,
	}

	defEngine := []NamedEngine{
		{Name: "AllTheInternet", Engine: alltheinternet.NewAllTheInternetEngine(PonderationsConfig["AllTheInternet"])},
		{Name: "DuckDuckGo", Engine: duckduckgo.NewDuckDuckGoEngine(PonderationsConfig["DuckDuckGo"])},
		{Name: "Qwant", Engine: qwant.NewQwantEngine(PonderationsConfig["Qwant"])},
		{Name: "Yep", Engine: yep.NewYepEngine(PonderationsConfig["Yep"])},
		//{Name: "Brave", Engine: brave.NewBraveEngine(PonderationsConfig["Brave"])},
		{Name: "You", Engine: you.NewYouEngine(PonderationsConfig["You"])},
		{Name: "Mojeek", Engine: mojeek.NewMojeek(PonderationsConfig["Mojeek"])},
	}

	availableengines := engine.AvailableEngines{}
	if len(*engines) > 0 {
		for _, v := range defEngine {
			for _, n := range *engines {
				if v.Name == n {
					availableengines = append(availableengines, v.Engine)
				}
			}
		}
	} else {
		for _, v := range defEngine {
			availableengines = append(availableengines, v.Engine)
		}
	}

	fmt.Println(availableengines)
	return availableengines
}

func getSearch(c *gin.Context) querus.Result {
	start := time.Now()
	//Récuperer les parametre de recherche
	var reqOptions engine.RequestOptions

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 1
	}
	reqOptions.IndexPage = page

	engines := c.QueryArray("engines")
	reqOptions.Engines = engines

	fmt.Println("Engines ", engines)

	reqOptions.Type = c.Query("type")
	if reqOptions.Type == "" {
		reqOptions.Type = "web"
	}

	reqOptions.Query = c.Query("q")

	r := c.Query("related")
	if r == "false" {
		reqOptions.Related = false
	} else {
		reqOptions.Related = true
	}

	maxResult, err := strconv.Atoi(c.Query("max_result"))
	if err != nil {
		maxResult = 10
	}
	reqOptions.MaxResults = maxResult

	reqOptions.Lang = c.Query("lang")
	if reqOptions.Lang == "" {
		reqOptions.Lang = "en-EN"
	}

	reqOptions.SafeSearch = c.Query("safe")

	if reqOptions.SafeSearch != "moderate" &&
		reqOptions.SafeSearch != "on" &&
		reqOptions.SafeSearch != "off" {

		reqOptions.SafeSearch = "moderate"
	}

	fmt.Println("type : ", reqOptions.Type)
	fmt.Println("index : ", reqOptions.IndexPage)
	fmt.Println("query : ", reqOptions.Query)
	fmt.Println("max result : ", reqOptions.MaxResults)
	fmt.Println("safe_search : ", reqOptions.SafeSearch)

	data := client(reqOptions)
	delay := time.Since(start)
	data.Delay = delay.Seconds()
	data.Type = reqOptions.Type
	data.Engines = engines
	return data

}
func client(reqOptions engine.RequestOptions) querus.Result {
	var wg sync.WaitGroup
	var resultsMutex sync.Mutex
	var result querus.Result
	var wiki wikipedia.Wiki_Result

	if reqOptions.Type == "web" || reqOptions.Type == "news" {
		wg.Add(1)
		go func(query string, wiki *wikipedia.Wiki_Result, wg *sync.WaitGroup) {
			defer wg.Done()
			w, err := wikipedia.Search(query, 6, reqOptions.Lang, true)
			if err != nil {
				fmt.Printf("%v - %v\n", err, w)
				return
			}
			if w.Empty {
				w, err = wikipedia.Search(query, 6, "en", true)
				if err != nil {
					fmt.Printf("%v - %v\n", err, w)
					return
				}
				if !w.Empty {
					resultsMutex.Lock()
					(*wiki) = w
					resultsMutex.Unlock()
				}
			} else {
				resultsMutex.Lock()
				(*wiki) = w
				resultsMutex.Unlock()
			}

		}(reqOptions.Query, &wiki, &wg)
	}
	engines := generateEnginesList(&reqOptions.Engines)
	wg.Add(1)

	go func(options engine.RequestOptions, engines engine.AvailableEngines, result *querus.Result, wg *sync.WaitGroup) {
		defer wg.Done()
		var r querus.Result = querus.Result{Query: options.Query}
		if options.Type == "web" || options.Type == "images" || options.Type == "news" {
			r = querus.Search(options, engines, "ascending")
		}
		resultsMutex.Lock()
		(*result) = r
		resultsMutex.Unlock()

	}(reqOptions, engines, &result, &wg)

	wg.Wait()
	result.Wikipedia = wiki
	result.IndexPage = reqOptions.IndexPage
	// Récupérer les données de l'utilisateur à partir de la base de données ou d'une autre source
	// Retourner les données de l'utilisateur en tant que réponse JSON

	return result

}
