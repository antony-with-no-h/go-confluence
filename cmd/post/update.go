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
	"path/filepath"

	"github.com/antony-with-no-h/go-confluence/client"
	"github.com/antony-with-no-h/go-confluence/config"
	"github.com/spf13/cobra"
)

type Version struct {
	Number int
}

// updateCmd represents the update command
var (
	flagSpace, flagTitle, flagFile string

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			flagSpace = cmd.Parent().PersistentFlags().Lookup("space").Value.String()
			flagTitle = cmd.Parent().PersistentFlags().Lookup("title").Value.String()
			flagFile = cmd.Parent().PersistentFlags().Lookup("file").Value.String()
		},
		Run: func(cmd *cobra.Command, args []string) {
			update()
		},
	}
)

func init() {
	PostCmd.AddCommand(updateCmd)
}

func update() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	pageID, version := details()
	requestBody := &Data{
		ID:    pageID,
		Type:  "page",
		Title: flagTitle,
		Space: Space{
			Key: flagSpace,
		},
		Body: Body{
			Storage{
				Representation: "storage",
			},
		},
		Version: Version{
			Number: version,
		},
	}

	file := filepath.Join(flagFile)
	requestBody.StorageFromMarkdown(file, cfg)

	bodyBuf := new(bytes.Buffer)
	bodyBufEncoder := json.NewEncoder(&NewLineToBrWriter{bodyBuf})
	bodyBufEncoder.SetEscapeHTML(false)
	bodyBufEncoder.Encode(requestBody)

	fmt.Printf("\n\nbody = %#v\n", bodyBuf.String())

}

func details() (string, int) {

	query := map[string]string{
		"spaceKey": flagSpace,
		"title":    flagTitle,
		"expand":   "version",
	}

	URL := client.MakeURL("/content", query)
	res := client.Get(URL, client.DefaultHeaders())

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)

	var body client.QueryBody
	json.Unmarshal(b, &body)

	return body.Results[0].ID, body.Results[0].Version.Number

}
