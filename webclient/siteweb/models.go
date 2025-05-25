package siteweb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

func SetupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{}, &Page{}, &Article{}, &Category{}, &Tag{}, &Comment{})
	return db, nil
}

// User structure with basic user information.
type User struct {
	gorm.Model
	Username string    `json:"username" gorm:"unique;not null"`
	Email    string    `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string    `json:"password" gorm:"not null"`
	Articles []Article `json:"articles" gorm:"foreignKey:AuthorID"` // One-to-many relationship
	Comments []Comment `json:"comments" gorm:"foreignKey:AuthorID"` // One-to-many relationship
}
type Staff struct {
	User
	Autorisations []string `json:"autorisations"`
}

// Page structure for static pages.
type Page struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	Likes    int    `json:"likes"`
	Draft    bool   `json:"is_draft"`
	AuthorID uint   `json:"author_id" gorm:"index"` // Foreign key for User
}

// Article structure with its content and metadata.
type Article struct {
	gorm.Model
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Likes      int         `json:"likes"`
	Draft      bool        `json:"is_draft"`
	AuthorID   uint        `json:"author_id" gorm:"index"` // Foreign key for User
	Comments   []Comment   `json:"comments"`               // One-to-many relationship
	Categories []*Category `gorm:"many2many:article_categories;"`
	Tags       []*Tag      `gorm:"many2many:article_tags;"`
}

// Category structure for categorizing articles.
type Category struct {
	gorm.Model
	Name     string     `json:"name"`
	Articles []*Article `gorm:"many2many:article_categories;"`
}

// Tag structure for tagging articles.
type Tag struct {
	gorm.Model
	Name     string     `json:"name"`
	Articles []*Article `gorm:"many2many:article_tags;"`
}

// Comment structure for article comments.
type Comment struct {
	gorm.Model
	AuthorID  uint   `json:"author_id" gorm:"index"` // Foreign key for User
	Content   string `json:"content"`
	ArticleID uint   `json:"article_id" gorm:"index"` // Foreign key for Article
}
