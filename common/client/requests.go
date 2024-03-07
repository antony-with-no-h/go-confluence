package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

func Get(url string, headers map[string]string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func Post(url string, body io.Reader, headers map[string]string) *http.Response {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", req)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func PageExists(space string, title string) bool {
	query := make(map[string]string)
	query["spaceKey"] = space
	query["title"] = title

	endpoint := MakeURL("/content", query)

	headers := make(map[string]string)
	headers["Authorization"] = viper.GetString("AccessToken")
	headers["Content-Type"] = "application/json"

	res := Get(endpoint, headers)

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var body QueryBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Fatal(err)
	}

	if body.Size > 0 {
		return true
	}

	return false
}

func DefaultHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Authorization"] = viper.GetString("AccessToken")
	headers["Content-Type"] = "application/json"

	return headers
}

func MakeURL(path string, query map[string]string) string {
	ApiEndpoint, _ := url.Parse(viper.GetString("Target"))
	ApiEndpoint = ApiEndpoint.JoinPath(path)

	q := ApiEndpoint.Query()
	for k, v := range query {
		q.Set(k, v)
	}

	ApiEndpoint.RawQuery = q.Encode()

	return fmt.Sprintf("%s", ApiEndpoint)
}

type QueryBody struct {
	Results []struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Title  string `json:"title"`
		Links  struct {
			Webui  string `json:"webui"`
			Edit   string `json:"edit"`
			Tinyui string `json:"tinyui"`
			Self   string `json:"self"`
		} `json:"_links"`
		Expandable struct {
			Container    string `json:"container"`
			Metadata     string `json:"metadata"`
			Operations   string `json:"operations"`
			Children     string `json:"children"`
			Restrictions string `json:"restrictions"`
			History      string `json:"history"`
			Ancestors    string `json:"ancestors"`
			Body         string `json:"body"`
			Version      string `json:"version"`
			Descendants  string `json:"descendants"`
			Space        string `json:"space"`
		} `json:"_expandable"`
	} `json:"results"`
	Start int `json:"start"`
	Limit int `json:"limit"`
	Size  int `json:"size"`
	Links struct {
		Self    string `json:"self"`
		Base    string `json:"base"`
		Context string `json:"context"`
	} `json:"_links"`
}
