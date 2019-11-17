package main

import (
	"github.com/manuito/sicobo/api"
	"github.com/manuito/sicobo/application"
	"github.com/manuito/sicobo/store"
)

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
