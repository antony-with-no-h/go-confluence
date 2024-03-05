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

	"github.com/antony-with-no-h/go-confluence/common/client"
	"github.com/hexops/valast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pageCmd represents the page command
var (
	tmplPath string
	title    string
	pageCmd  = &cobra.Command{
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
}

func postPage(cmd *cobra.Command, args []string) {

	requestBody := &Data{
		Type:  "page",
		Title: title,
		Space: Space{
			Key: viper.GetString("space"),
		},
		Body: Body{
			Storage{
				Representation: "storage",
			},
		},
	}

	requestBody.SetStorage(filepath.Join(tmplPath))

	if pageExists := client.PageExists(requestBody.Space.Key, requestBody.Title); pageExists == true {
		log.Printf("Page already exists")
		os.Exit(1)
	}

	bodyBuf := new(bytes.Buffer)
	json.NewEncoder(bodyBuf).Encode(requestBody)

	URL := client.MakeURL("/content", nil)
	res := client.Post(URL,
		bodyBuf,
		client.DefaultHeaders(),
	)

	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)

	var jsonRes Response
	json.Unmarshal(resBody, &jsonRes)

	fmt.Println(valast.String(jsonRes))

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

func (d *Data) SetStorage(file string) {
	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal(err)
	}

	d.Body.Storage.Value = string(b)
}
