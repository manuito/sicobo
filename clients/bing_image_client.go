package clients

import (
	"encoding/json"
	"net/http"
	"sicobo/application"
	"strconv"
)

type BSearchResult struct {
	Value []BPicture `json:"value"`
}

type BPicture struct {
	Name         string `json:"name"`
	ContentUrl   string `json:"contentUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

func SearchBingImage(isbn string) BSearchResult {

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.cognitive.microsoft.com/bing/v7.0/images/search?q="+isbn+"&count="+strconv.Itoa(bingImageMaxItems)+"&offset=0", nil)
	req.Header.Add(bingHeaderName, application.State.Config.BingAPIkey)
	res, err := client.Do(req)
	if err != nil {
		application.Error("Cannot request bing img", err)
	}

	decoder := json.NewDecoder(res.Body)
	result := BSearchResult{}
	err = decoder.Decode(&result)
	if err != nil {
		application.Error("Cannot read result from bing img", err)
	}
	return result
}
