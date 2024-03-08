/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/antony-with-no-h/go-confluence/client"
	"github.com/antony-with-no-h/go-confluence/config"
	convert_markdown "github.com/antony-with-no-h/go-confluence/convert-markdown"
	"github.com/hexops/valast"
	"github.com/spf13/cobra"
)

// pageCmd represents the page command
var (
	tmplPath  string
	title     string
	mdConvert bool
	pageCmd   = &cobra.Command{
		Use:   "page",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			postPage(cmd, args)
		},
	}
)

func init() {
	PostCmd.AddCommand(pageCmd)
	pageCmd.Flags().StringVarP(&tmplPath, "filename", "f", "", "")
	pageCmd.Flags().StringVarP(&title, "title", "t", "", "")
	pageCmd.Flags().BoolVar(&mdConvert, "md", false, "Treat source file as Markdown")
}

func postPage(cmd *cobra.Command, args []string) {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	requestBody := &Data{
		Type:  "page",
		Title: title,
		Space: Space{
			Key: PostCmd.PersistentFlags().Lookup("space").Value.String(),
		},
		Body: Body{
			Storage{
				Representation: "storage",
			},
		},
	}

	file := filepath.Join(tmplPath)
	if mdConvert {
		requestBody.StorageFromMarkdown(file, cfg)
	} else {
		requestBody.SetStorage(file)
	}

	if pageExists := client.PageExists(requestBody.Space.Key, requestBody.Title); pageExists == true {
		log.Printf("Page already exists")
		os.Exit(1)
	}

	bodyBuf := new(bytes.Buffer)
	bodyBufEncoder := json.NewEncoder(&NewLineToBrWriter{bodyBuf})
	bodyBufEncoder.SetEscapeHTML(false)
	bodyBufEncoder.Encode(requestBody)

	URL := client.MakeURL("/content", nil)
	res := client.Post(URL,
		bodyBuf,
		client.DefaultHeaders(),
	)

	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		var jsonErr ResponseError
		json.Unmarshal(resBody, &jsonErr)

		errJson, _ := json.MarshalIndent(jsonErr, "", "  ")
		log.Fatalf("%s\n", errJson)
	}

	var jsonRes Response
	json.Unmarshal(resBody, &jsonRes)

	fmt.Println(valast.String(jsonRes))

}

type NewLineToBrWriter struct {
	Writer io.Writer
}

func (w *NewLineToBrWriter) Write(p []byte) (n int, err error) {
	re := regexp.MustCompile(`(?s)\[CDATA\[(.*?)]]|\\n`)
	str := re.ReplaceAllStringFunc(string(p), func(match string) string {
		if strings.HasPrefix(match, "[CDATA[") {
			return match
		}
		return "<br/>"
	})
	//str := strings.Replace(string(p), `\n`, "<br/>", -1)
	return w.Writer.Write([]byte(str))
}

type Data struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Space `json:"space"`
	Body  `json:"body"`
}

type Body struct {
	Storage `json:"storage"`
}

type Space struct {
	Key string `json:"key"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type Response struct {
	Links struct {
		Base       string `json:"base"`
		Collection string `json:"collection"`
		Self       string `json:"self"`
		Tinyui     string `json:"tinyui"`
		Webui      string `json:"webui"`
	} `json:"_links"`
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
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

func (d *Data) SetStorage(file string) {
	d.Body.Storage.Value = string(open(file))
}

func (d *Data) StorageFromMarkdown(file string, cfg config.Config) {
	d.Body.Storage.Value = convert_markdown.RenderHTML(open(file), cfg)
}

func open(file string) []byte {
	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
