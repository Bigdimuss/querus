package webclient

import (
	"html/template"
	"net/http"
	"search/webclient/siteweb"
	"search/webclient/siteweb/admin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	gin.SetMode(gin.ReleaseMode)

	db, err := siteweb.SetupDatabase()
	if err != nil {
		panic("failed to connect to database")
	}
	app := siteweb.App{DB: db}
	admin := admin.Admin{DB: db}

	router := gin.Default()
	admin.SetRout(router)
	app.SetRout(router)
	// Configuration du middleware de sessions
	store := cookie.NewStore([]byte("ceciestsecret")) // Changez "secret" par une clé sécurisée
	router.Use(sessions.Sessions("session", store))   // Appliquez le middleware ici

	// Configuration des templates et des statiques
	router.SetFuncMap(template.FuncMap{
		"RenderEngines": RenderEnginesfunc,
		"truncate":      truncate,
	})
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Routes publiques
	router.GET("/", getIndex)
	router.GET("/search", getSearchWeb)
	router.GET("/search/:page", getSearchWeb)

	// Routes pour l'API
	api := router.Group("/api")
	api.GET("/search", getSearchAPI)
	api.GET("/articles", func(c *gin.Context) {
		var articles []siteweb.Article
		db.Preload("Comments").Preload("Categories").Find(&articles)
		c.JSON(http.StatusOK, articles)
	})

	// Lancement du serveur
	router.Run(":8080")
}
