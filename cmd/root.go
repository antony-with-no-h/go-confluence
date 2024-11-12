package cmd

import (
	"os"

	"github.com/antony-with-no-h/go-confluence/cmd/add"
	"github.com/antony-with-no-h/go-confluence/cmd/edit"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:              "go-confluence",
		Short:            "",
		Long:             ``,
		TraverseChildren: true,
	}
)

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("space", "s", "", "The Confluence space where the page should be published (e.g. Engineering, QA)")
	RootCmd.PersistentFlags().StringP("title", "t", "", "Page title")
	RootCmd.PersistentFlags().StringP("file", "f", "", "Path to file containing Page Markdown")
	RootCmd.PersistentFlags().StringP("parent", "p", "", "Title of the page that will act as the parent (e.g. Support, Backup and Restore)")
	RootCmd.PersistentFlags().Bool("xml", false, "<file> is Confluence XML, do not convert")

	RootCmd.AddCommand(edit.EditCmd)
	RootCmd.AddCommand(add.AddCmd)
	RootCmd.AddCommand(dryRunCmd)
}
