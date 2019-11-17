package main

import (
	"github.com/manuito/sicobo/api"
	"github.com/manuito/sicobo/application"
	_ "github.com/manuito/sicobo/docs"
	"github.com/manuito/sicobo/store"
)

// @title SICOBO - Simple personal comic book collection management
// @version 0.1
// @description SICOBO backend services
// @contact.name manuito
// @contact.url https://github.com/manuito
// @license.name Do What The Fuck You Want To Public License (WTFPL)
// @license.url http://www.wtfpl.net/about/
// @BasePath /api/v1
func main() {

	application.Info("Starting sicobo MGNT app !")

	isbn := "9782377540037"

	store.StartConnect()
	store.SelectActiveDatabase("other")

	store.StartPictureProcess()

	book, err := store.AddNewBook(isbn)

	if err == nil {
		application.Info("add new book :", book.Title)
		store.CompleteBookPicture(&book)
	} else {
		application.Info("book already exists :", book.Title)
	}

	application.Info("Current store size :", store.ActiveSpec().TotalBooks)

	api.StartServer()
}
