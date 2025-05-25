package siteweb

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

func (app *App) GetArticles(c *gin.Context) {
	var articles []Article
	if err := app.DB.Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "blog-index.html", gin.H{"articles": articles})
}

func (app *App) GetArticle(c *gin.Context) {
	c.HTML(http.StatusOK, "blog-article.html", gin.H{})
}

func (app *App) GetPage(c *gin.Context) {
	c.HTML(http.StatusOK, "page.html", gin.H{})
}
func (app *App) GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}
func (app *App) GetConnect(c *gin.Context) {
	c.HTML(http.StatusOK, "connect.html", gin.H{})
}

var currentUser *User // Variable pour stocker l'utilisateur connecté

func (app *App) Login(c *gin.Context) {
	var credentials User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	var user User
	if err := app.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiants incorrects"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiants incorrects"})
		return
	}

	// Stocker l'utilisateur dans la session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la sauvegarde de la session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie"})
}
func (app *App) Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hacher le mot de passe
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hachage du mot de passe"})
		return
	}
	user.Password = hashedPassword

	// Enregistrer l'utilisateur dans la base de données
	if err := app.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement de l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Utilisateur enregistré avec succès"})
}

// Méthode pour ajouter une page (protégée)
func (app *App) CreatePage(c *gin.Context) {
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
		return
	}
	// Logic pour créer une page...
}

func (app *App) SetRout(router *gin.Engine) {
	// Routes protégées
	blog := router.Group("/admin")
	blog.GET("/blog", app.GetArticles)
	blog.GET("/blog/:slug", app.GetArticle)

	site := router.Group("/page")
	site.GET("/:slug", app.GetPage)
	router.POST("/login", app.Login) // Assurez-vous que cette route est définie après le middleware de session
	router.POST("/register", app.Register)
	router.GET("/register", app.GetRegister)
	router.GET("/admin", app.GetConnect)
}
