package clients

import (
	"biblio/application"
	"encoding/json"
	"net/http"
)

type GBookList struct {
	Kind       string  `json:"Kind"`
	TotalItems int     `json:"totalItems"`
	Items      []GBook `json:"items"`
}

type GBook struct {
	Kind       string   `json:"kind"`
	Id         string   `json:"id"`
	VolumeInfo GVolume  `json:"volumeInfo"`
	SearchInfo GSnippet `json:"searchInfo"`
}

type GVolume struct {
	Title               string        `json:"title"`
	Authors             []string      `json:"authors"`
	PublishedDate       string        `json:"publishedDate"`
	IndustryIdentifiers []GIdentifier `json:"industryIdentifiers"`
	PageCount           int           `json:"pageCount"`
	Language            string        `json:"language"`
}

type GIdentifier struct {
	Identifier string `json:"identifier"`
}

type GSnippet struct {
	TextSnippet string `json:"textSnippet"`
}

func SearchGoogleBooks(isbn string) GBookList {
	res, err := http.Get("https://www.googleapis.com/books/v1/volumes?q=+isbn:" + isbn + "&key=" + application.State.Config.GoogleBookAPIKey)
	if err != nil {
		application.Error("Cannot connect to google book", err)
	}

	decoder := json.NewDecoder(res.Body)
	result := GBookList{}
	err = decoder.Decode(&result)
	if err != nil {
		application.Error("Cannot read result from google book", err)
	}
	return result
}
