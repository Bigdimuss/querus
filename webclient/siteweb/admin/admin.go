package admin

import (
	"net/http"
	"search/webclient/siteweb"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Admin struct {
	DB *gorm.DB
}

func (admin *Admin) GetPagesAdmin(c *gin.Context) {

}

func (admin *Admin) NewPageForm(c *gin.Context) {
	c.HTML(http.StatusOK, "page-form.html", nil)
}

func (admin *Admin) CreatePage(c *gin.Context) {
	var page siteweb.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Définir AuthorID si nécessaire
	page.AuthorID = 1 // Remplacez ceci par la logique pour obtenir l'ID de l'utilisateur connecté
	if err := admin.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, page)
}
func (admin *Admin) EditPageAdmin(c *gin.Context) {

}

func (admin *Admin) SetRout(router *gin.Engine) {
	// Routes protégées
	adm := router.Group("/admin")
	adm.Use(siteweb.AuthMiddleware())
	{
		adm.POST("/pages", admin.CreatePage)
		adm.GET("/pages", admin.GetPagesAdmin)
		adm.GET("/pages/edit/:id", admin.EditPageAdmin)
		adm.GET("/pages/new", admin.NewPageForm)
	}
}
