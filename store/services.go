package store

import (
	"biblio/application"
	"biblio/clients"
	"errors"
)

// Monitors : all the monitoring elements in a channel
var booksToComplete = make(chan *Book)

// SelectActiveDatabase Select the active collection in mongo. Must be called before any mongo use
func SelectActiveDatabase(name string) {
	database := loadMongoDatabase(name)
	spec := loadOrInitDatabaseSpec(database, name)
	switchToDatabase(spec, database)
}

// StartPictureProcess start chan processing for book pictures completion
func StartPictureProcess() {
	go func() {
		for {
			book := <-booksToComplete
			application.Debug("Begin completion of picture for book", book.Isbn)
			CompleteBookPicture(book)
		}
	}()
	application.Info("PictureProcess linked to book updates")
}

// GetBook read details for a book
func GetBook(isbn string) (Book, error) {

	book := loadBook(isbn)

	if book.Isbn != "" {

		// Complete missing picture
		if book.PictureURL == "" {
			booksToComplete <- &book
		}

		return book, nil
	}

	return book, errors.New("No book found for isbn " + isbn)
}

// AddNewBook Init and add new book to Store
func AddNewBook(isbn string) (Book, error) {

	book := loadBook(isbn)

	if book.Isbn != "" {

		// Complete missing picture
		if book.PictureURL == "" {
			booksToComplete <- &book
		}

		// And stop with loaded book
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

	// Add to completion process
	booksToComplete <- &book

	return book, nil
}

// CompleteBookPicture From isbn, search for picture links, and download 1st picture
func CompleteBookPicture(book *Book) {

	pictureURLs := clients.SearchBingImage(book.Isbn)

	// Result found
	if len(pictureURLs.Value) > 0 {
		book.PictureURL = pictureURLs.Value[0].ContentUrl

		book.CandidateDetails.PictureURLs = make([]string, len(pictureURLs.Value))
		for i := range pictureURLs.Value {
			book.CandidateDetails.PictureURLs[i] = pictureURLs.Value[i].ContentUrl
		}

		// Download picture
		file, err := clients.DownloadFile(book.Isbn, book.PictureURL)

		if err == nil {
			book.Picture = file
		} else {
			application.Error("Failed to download picture from ", book.PictureURL, err)
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
