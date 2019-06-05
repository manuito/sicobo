package store

import (
	"biblio/application"
	"biblio/clients"
	"errors"
)

// SelectActiveDatabase Select the active collection in mongo. Must be called before any mongo use
func SelectActiveDatabase(name string) {
	database := loadMongoDatabase(name)
	spec := loadOrInitDatabaseSpec(database, name)
	switchToDatabase(spec, database)
}

// AddNewBook Init and add new book to Store
func AddNewBook(isbn string) (Book, error) {

	book := loadBook(isbn)

	if book.Isbn != "" {
		return book, errors.New("already exist")
	}

	book = Book{Isbn: isbn}

	loadBookFromGoogleBooks(&book)

	// Try again from bing
	if book.Title == "" {
		loadBookFromBing(&book)
	}

	// Identify collection
	detectCollectionCandidates(&book)

	insertBook(book)

	return book, nil
}

// CompleteBookPicture From isbn, search for picture links, and download 1st picture
func CompleteBookPicture(book *Book) {

	pictures := clients.SearchBingImage(book.Isbn)

	// Result found
	if len(pictures.Value) > 0 {
		book.Picture = pictures.Value[0].ContentUrl

		book.CandidateDetails.Pictures = make([]string, len(pictures.Value))
		for i := range pictures.Value {
			book.CandidateDetails.Pictures[i] = pictures.Value[i].ContentUrl
		}
	}

	updateBook(book)
}

func ListExistingBooks() []Book {
	return loadBooks()
}

// CompleteBookCollection Search for existing collection, link to specified Book
func CompleteBookCollection(book *Book) {

}

func loadBookFromGoogleBooks(book *Book) {
	result := clients.SearchGoogleBooks(book.Isbn)

	// Result found
	if len(result.Items) > 0 {

		vol := result.Items[0].VolumeInfo

		application.Debug("Found in Google book : " + vol.Title)

		book.Title = vol.Title

		// All candidates
		book.CandidateDetails.Titles = make([]string, len(result.Items))
		for i := range result.Items {
			book.CandidateDetails.Titles[i] = result.Items[i].VolumeInfo.Title
		}

		// All authors
		book.Authors = make([]Author, len(vol.Authors))
		for i := range book.Authors {
			book.Authors[i].Name = vol.Authors[i]
		}

		book.PageCount = vol.PageCount
		book.Snippet = result.Items[0].SearchInfo.TextSnippet
		book.PublishedDate = vol.PublishedDate
	} else {
		application.Info(book.Isbn + " not found in Google book !")
	}
}

// Less available data
func loadBookFromBing(book *Book) {
	bresult := clients.SearchBingWeb(book.Isbn)

	book.CandidateDetails.DegradedSource = true

	// Result found
	if len(bresult) > 0 {
		book.Title = bresult[0]
		book.CandidateDetails.Titles = bresult
	}
}

func detectCollectionCandidates(book *Book) {

}
