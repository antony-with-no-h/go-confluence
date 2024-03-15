package add

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/antony-with-no-h/go-confluence/requests"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "page",
		Short: "",
		Long:  `Publish a new page.`,
		Run: func(cmd *cobra.Command, args []string) {
			AddPage(cmd)
		},
	}
)

func AddPage(cmd *cobra.Command) {
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

	if pageExists := page.Exists(); pageExists {
		fmt.Printf("Page already exists, use `edit` to update.\n")
		os.Exit(1)
	}

	if cmd.Parent().Flags().Lookup("parent").Changed {
		page.SetParent(flagParent)
	}

	page.SetStorageValue(flagFile)

	buf := new(bytes.Buffer)
	bufEncoder := json.NewEncoder(buf)
	bufEncoder.SetEscapeHTML(false)
	bufEncoder.Encode(page)

	url := requests.MakeURL("/content", nil)

	requests.Post(url, requests.DefaultHeaders(), buf)
}

func init() {
	AddCmd.AddCommand(addCmd)
}
