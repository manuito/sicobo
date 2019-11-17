package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manuito/sicobo/application"
	"github.com/manuito/sicobo/store"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func selectDb(c *gin.Context) (db string) {
	db = c.Param("database")
	store.SelectActiveDatabase(db)
	return
}

// GET http://1.2.3.4:8080/api/v1/library/{database}
// @Summary Get detail specs for one specified library, identified by its name
// @Param database	path	string	true	"library code"
// @Success 200 {object} store.BookDatabaseSpec	"Spec for the requested library"
// @Router /api/v1/library/{database} [get]
func getLibrarySpec(c *gin.Context) {
	selectDb(c)
	application.Debug("Get library spec")
	c.JSON(http.StatusOK, store.ActiveSpec())
}

// GET http://1.2.3.4:8080/api/v1/library/{database}/books
// @Summary Get all managed books for one specified library
// @Param database	path	string	true	"library code"
// @Success 200 {array} store.Book	"All books"
// @Router /api/v1/library/{database}/books [get]
func getLibraryBooks(c *gin.Context) {
	selectDb(c)
	application.Debug("Get library books")
	c.JSON(http.StatusOK, store.ListExistingBooks())
}

// GET http://1.2.3.4:8080/api/v1/library/{database}/books/{isbn}
// @Summary Get details for one book, from a library, identified by its isbn code
// @Param database	path	string	true	"library code"
// @Param isbn		path	string	true	"book ISBN code"
// @Success 200 {object} store.Book	"Available book for specified ISBN code"
// @Router /api/v1/library/{database}/books/{isbn} [get]
func getLibraryBook(c *gin.Context) {
	selectDb(c)
	isbn := c.Param("isbn")
	application.Debug("Get library book with isbn", isbn)
	book, err := store.GetBook(isbn)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	} else {
		c.JSON(http.StatusOK, book)
	}
}

// POST http://1.2.3.4:8080/api/v1/library/{database}/books/{isbn}
// @Summary Add a new book to a specified library. The book is only identified by its ISBN code
// @Param database	path	string	true	"library code"
// @Param isbn		path	string	true	"book ISBN code. App will search for details using various service providers"
// @Success 200 {object} store.Book	"Created book info for specified ISBN code"
// @Router /api/v1/library/{database}/books/{isbn} [post]
func addBookToLibrary(c *gin.Context) {
	db := selectDb(c)
	isbn := c.Param("isbn")
	application.Debug("Add book with isbn", isbn, " to library books", db)
	book, _ := store.AddNewBook(isbn)
	c.JSON(http.StatusOK, book)
}

// StartServer start the gin server route on all identified API entry points
// @title SICOBO - Simple personal comic book collection management
// @version 0.1
// @description SICOBO backend services
// @contact.name manuito
// @contact.url https://github.com/manuito
// @license.name Do What The Fuck You Want To Public License (WTFPL)
// @license.url http://www.wtfpl.net/about/
// @BasePath /api/v1
func StartServer() {
	r := gin.Default()

	// Simple group: v2
	v1 := r.Group("/api/v1")
	{
		lib := v1.Group("/library")
		{
			lib.GET("/:database", getLibrarySpec)
			lib.GET("/:database/books", getLibraryBooks)
			lib.GET("/:database/books/:isbn", getLibraryBook)
			lib.POST("/:database/books/:isbn", addBookToLibrary)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
