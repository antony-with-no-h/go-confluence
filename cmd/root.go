package cmd

import (
	"os"

	"github.com/antony-with-no-h/go-confluence/cmd/add"
	"github.com/antony-with-no-h/go-confluence/cmd/edit"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "go-confluence",
		Short: "",
		Long:  ``,
	}
)

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("space", "s", "", "")
	RootCmd.PersistentFlags().StringP("title", "t", "", "")
	RootCmd.PersistentFlags().StringP("file", "f", "", "")
	RootCmd.PersistentFlags().StringP("parent", "p", "", "")

	RootCmd.AddCommand(edit.EditCmd)
	RootCmd.AddCommand(add.AddCmd)
}
