/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/antony-with-no-h/go-confluence/client"
	"github.com/spf13/cobra"
)

// contentCmd represents the content command
var (
	body       client.QueryBody
	contentCmd = &cobra.Command{
		Use:     "content",
		Args:    cobra.ArbitraryArgs,
		Short:   "GET request to /content and print result",
		Long:    ``,
		Example: "  get content spaceKey=Engineering title=\"How-to guides\"",
		Run: func(cmd *cobra.Command, args []string) {
			browse(args)
		},
	}
)

func browse(args []string) {
	URL := client.MakeURL("/content", parseArgs(args))
	res := client.Get(URL, client.DefaultHeaders())

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("io.ReadAll: %s\n", err)
	}
	json.Unmarshal(b, &body)

	stdout, _ := json.MarshalIndent(&body, "", "    ")
	fmt.Println(string(stdout))
}

func parseArgs(args []string) map[string]string {
	queryMap := make(map[string]string)

	for arg := range args {
		keyValue := strings.Split(args[arg], "=")
		queryMap[keyValue[0]] = keyValue[1]
	}

	return queryMap
}

func init() {
	GetCmd.AddCommand(contentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
