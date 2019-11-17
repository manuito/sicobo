package clients

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"sicobo/application"
	"strconv"
)

func cleanUp(s string) string {
	re := regexp.MustCompile(`\b(\\\d\d\d)`)
	return re.ReplaceAllStringFunc(s, func(s string) string {
		return `\u0` + s[1:]
	})
}

func SearchBingWeb(isbn string) []string {

	re := regexp.MustCompile("\"name\": \"([^\"]*)\"")
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.cognitive.microsoft.com/bing/v7.0/search?q="+isbn+"&count="+strconv.Itoa(bingWebMaxItems)+"&offset=0", nil)
	req.Header.Add(bingHeaderName, application.State.Config.BingAPIkey)
	req.Close = true
	res, err := client.Do(req)
	if err != nil {
		application.Error("Cannot connect to bing web", err)
	}

	b := res.Body
	bodyBytes, err := ioutil.ReadAll(b)
	result := re.FindAllStringSubmatch(string(bodyBytes), -1)

	if result != nil {
		var s []string

		for _, match := range result {
			s = append(s, match[1])
		}

		return s
	}

	return nil
}
