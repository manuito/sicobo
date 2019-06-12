package main

import (
	"biblio/api"
	"biblio/application"
	"biblio/store"
)

func main() {

	application.Info("Starting biblio MGNT app !")

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
