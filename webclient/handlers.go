package webclient

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func getSearchWeb(c *gin.Context) {
	data := getSearch(c.Copy())

	queryParams := c.Request.URL.Query()
	log.Println("Query Params:", queryParams)

	// Construction manuelle de la chaîne de requête
	var queryString string
	for key, values := range queryParams {
		for _, value := range values {
			if queryString != "" {
				queryString += "&"
			}
			queryString += key + "=" + value // Construction de la chaîne sans encodage
		}
	}

	log.Println("Query String:", queryString)
	nextPageURL := fmt.Sprintf("/search/%d?%s", data.IndexPage+1, queryString)
	previousPageURL := fmt.Sprintf("/search/%d?%s", data.IndexPage-1, queryString)

	c.HTML(http.StatusOK, "search.html", gin.H{
		"type":              data.Type,
		"query":             data.Query,
		"number":            len(data.Results),
		"time":              data.Delay,
		"wiki":              data.Wikipedia,
		"results":           data.Results,
		"engines":           data.Engines,
		"enginesdelays":     data.EngineDelays,
		"previous_page_nb":  data.IndexPage - 1,
		"previous_page_url": previousPageURL,
		"index_page":        data.IndexPage,
		"next_page_nb":      data.IndexPage + 1,
		"next_page_url":     nextPageURL,
	})
}

func getSearchAPI(c *gin.Context) {
	result := getSearch(c.Copy())
	c.JSON(200, result)

}
