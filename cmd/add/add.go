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
	AddCmd = &cobra.Command{
		Use:              "add",
		Short:            "",
		Long:             ``,
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	pageCmd = &cobra.Command{
		Use:   "page",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			addPage(cmd)
		},
	}
)

func init() {
	AddCmd.AddCommand(pageCmd)
}

func addPage(cmd *cobra.Command) {
	flagTitle := cmd.Parent().Flags().Lookup("title").Value.String()
	flagSpace := cmd.Parent().Flags().Lookup("space").Value.String()
	flagFile := cmd.Parent().Flags().Lookup("file").Value.String()

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
		fmt.Printf("Page '%s/%s' already exists, use `edit` to update.\n", flagSpace, flagTitle)
		os.Exit(1)
	}

	if cmd.Parent().Flags().Lookup("parent").Changed {
		page.SetParent(cmd.Parent().Flags().Lookup("parent").Value.String())
	}

	page.SetStorageValue(flagFile)
	page.Storage.Value = wrapWithLayout(page.Storage.Value)

	buf := new(bytes.Buffer)
	bufEncoder := json.NewEncoder(buf)
	bufEncoder.SetEscapeHTML(false)
	bufEncoder.Encode(page)

	URL := requests.MakeURL("/content", nil)
	requests.Post(URL, requests.DefaultHeaders(), buf)
}

func wrapWithLayout(wrap string) string {
	return fmt.Sprintf(`<ac:layout>
	<ac:layout-section ac:type="two_equal">
		<ac:layout-cell>
		%s
		</ac:layout-cell>
		<ac:layout-cell>
		<p>
			<br/>
		</p>
		</ac:layout-cell>
	</ac:layout-section>
</ac:layout>`, wrap)
}
