package edit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/antony-with-no-h/go-confluence/requests"
	"github.com/spf13/cobra"
)

var (
	editCmd = &cobra.Command{
		Use:   "page",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			EditPage()
		},
	}
)

func EditPage() {
	page := &requests.Page{
		Type:  "page",
		Title: flagTitle,
		Space: requests.Space{
			Key: flagSpace,
		},
		Body: requests.Body{
			Storage: requests.Storage{
				Representation: "storage",
			},
		},
	}

	if pageExists := page.Exists(); !pageExists {
		fmt.Printf("Cannot access '%s': No such page", page.Title)
		os.Exit(1)
	} else {
		page.SetID()
		page.SetVersion()
	}

	page.SetStorageValue(flagFile)

	buf := new(bytes.Buffer)
	bufEncoder := json.NewEncoder(buf)
	bufEncoder.SetEscapeHTML(false)
	bufEncoder.Encode(page)

	fmt.Println(buf.String())

	url := requests.MakeURL(fmt.Sprintf("/content/%s", page.ID), nil)

	requests.Put(url, requests.DefaultHeaders(), buf)
}

func init() {
	EditCmd.AddCommand(editCmd)
}
