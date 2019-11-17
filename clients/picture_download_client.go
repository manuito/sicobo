package clients

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sicobo/application"
	"strings"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(isbn string, url string) (string, error) {

	// Init dest folder if not done yet

	destFolder := prepareDownloadFolder()
	name := getLocalName(isbn, url)
	filepath := filepath.Join(destFolder, name)

	application.Info("Store picture for "+isbn+" from", url, " to", filepath)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return name, err
}

func prepareDownloadFolder() string {

	newpath := filepath.Join(".", application.State.Config.FileStore)
	os.MkdirAll(newpath, os.ModePerm)
	return newpath
}

func getLocalName(isbn string, url string) string {

	idx := strings.LastIndex(url, ".")

	return isbn + url[idx:]
}
