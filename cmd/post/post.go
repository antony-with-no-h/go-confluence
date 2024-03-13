/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package post

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var (
	PostCmd = &cobra.Command{
		Use:     "post",
		Aliases: []string{"add"},
		Short:   "",
		Long:    ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	PostCmd.PersistentFlags().StringP("space", "s", "",
		"Confluence space name (QA, HR, Engineering etc)")
	PostCmd.PersistentFlags().StringP("title", "t", "",
		"Page title")
	PostCmd.PersistentFlags().StringP("file", "f", "",
		"Path to file containing page contents")

	PostCmd.MarkPersistentFlagRequired("space")
}
