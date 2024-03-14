package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antony-with-no-h/go-confluence/config"
	"github.com/antony-with-no-h/go-confluence/markdown"
	"github.com/hexops/valast"
)

type NewlineToBrWriter struct {
	Writer io.Writer
}

func (w *NewlineToBrWriter) Write(data []byte) (n int, err error) {
	re := regexp.MustCompile(`(?mU)<ac:structured-macro.*?</ac:structured-macro>|\\n\\n`)
	str := re.ReplaceAllStringFunc(string(data), func(match string) string {
		if strings.HasPrefix(match, "<ac:structured-macro") {
			return match
		}
		return "<br/>"
	})

	fmt.Printf("data == %s\n\nstr == %s\n", data, str)

	return w.Writer.Write([]byte(str))
}

type Page struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Space     `json:"space"`
	Body      `json:"body"`
	Version   `json:"version,omitempty"`
	Ancestors []Ancestors `json:"ancestors,omitempty"`
}

type Space struct {
	Key string `json:"key"`
}

type Body struct {
	Storage `json:"storage"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type Version struct {
	Number int `json:"number,omitempty"`
}

type Ancestors struct {
	ID int `json:"id,omitempty"`
}

type ResponseBody struct {
	Results []struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Status  string `json:"status"`
		Title   string `json:"title"`
		Version struct {
			By struct {
				Type           string `json:"type"`
				Username       string `json:"username"`
				UserKey        string `json:"userKey"`
				ProfilePicture struct {
					Path      string `json:"path"`
					Width     int    `json:"width"`
					Height    int    `json:"height"`
					IsDefault bool   `json:"isDefault"`
				} `json:"profilePicture"`
				DisplayName string `json:"displayName"`
				Links       struct {
					Self string `json:"self"`
				} `json:"_links"`
				Expandable struct {
					Status string `json:"status"`
				} `json:"_expandable"`
			} `json:"by"`
			When      time.Time `json:"when"`
			Message   string    `json:"message"`
			Number    int       `json:"number"`
			MinorEdit bool      `json:"minorEdit"`
			Hidden    bool      `json:"hidden"`
			Links     struct {
				Self string `json:"self"`
			} `json:"_links"`
			Expandable struct {
				Content string `json:"content"`
			} `json:"_expandable"`
		} `json:"version,omitempty"`
		Links struct {
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

type ResponseError struct {
	StatusCode int `json:"statusCode"`
	Data       struct {
		Authorized            bool  `json:"authorized"`
		Valid                 bool  `json:"valid"`
		AllowedInReadOnlyMode bool  `json:"allowedInReadOnlyMode"`
		Errors                []any `json:"errors"`
		Successful            bool  `json:"successful"`
	} `json:"data"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

func DefaultHeaders() map[string]string {
	cfg, _ := config.LoadConfig()
	headers := make(map[string]string)
	headers["Authorization"] = cfg.AccessToken
	headers["Content-Type"] = "application/json"

	return headers
}

func MakeURL(path string, query map[string]string) string {
	cfg, _ := config.LoadConfig()
	ApiEndpoint, _ := url.Parse(cfg.Target)
	ApiEndpoint = ApiEndpoint.JoinPath(path)

	q := ApiEndpoint.Query()
	for k, v := range query {
		q.Set(k, v)
	}

	ApiEndpoint.RawQuery = q.Encode()

	return fmt.Sprintf("%s", ApiEndpoint)
}

func (p *Page) Exists() bool {
	query := make(map[string]string)
	query["spaceKey"] = p.Space.Key
	query["title"] = p.Title

	url := MakeURL("/content", query)
	header := DefaultHeaders()

	body := Get(url, header)
	if body.Size > 0 {
		return true
	}

	return false
}

func (p *Page) SetVersion() {
	query := make(map[string]string)
	query["spaceKey"] = p.Space.Key
	query["title"] = p.Title
	query["expand"] = "version"

	url := MakeURL("/content", query)
	header := DefaultHeaders()

	body := Get(url, header)
	p.Version.Number = body.Results[0].Version.Number + 1
}

func (p *Page) SetID() {
	query := make(map[string]string)
	query["spaceKey"] = p.Space.Key
	query["title"] = p.Title

	url := MakeURL("/content", query)
	header := DefaultHeaders()

	body := Get(url, header)
	p.ID = body.Results[0].ID
}

func (p *Page) SetParent(title string) {
	query := make(map[string]string)
	query["spaceKey"] = p.Space.Key
	query["title"] = title

	url := MakeURL("/content", query)
	header := DefaultHeaders()

	body := Get(url, header)
	parentID, _ := strconv.Atoi(body.Results[0].ID)

	parent := &Ancestors{
		ID: parentID,
	}

	p.Ancestors = append(p.Ancestors, *parent)
}

func (p *Page) SetStorageValue(filePath string) {
	fd, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal(err)
	}

	p.Body.Storage.Value = markdown.RenderHTML(b)

}

func Get(url string, headers map[string]string) ResponseBody {
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

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var body ResponseBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func Post(url string, headers map[string]string, data io.Reader) {
	req, err := http.NewRequest("POST", url, data)
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

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		var err ResponseError
		json.Unmarshal(b, &err)

		fmt.Println(res.StatusCode)
		fmt.Println(valast.String(err))
		os.Exit(1)
	}

	var body ResponseBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Fatal(err)
	}
}

func Put(url string, headers map[string]string, data io.Reader) {
	req, err := http.NewRequest(http.MethodPut, url, data)
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

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		var err ResponseError
		json.Unmarshal(b, &err)

		fmt.Println(res.StatusCode)
		fmt.Println(valast.String(err))
		os.Exit(1)
	}

	var body ResponseBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Fatal(err)
	}
}
