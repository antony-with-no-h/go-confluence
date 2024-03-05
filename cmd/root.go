/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/antony-with-no-h/go-confluence/cmd/get"
	"github.com/antony-with-no-h/go-confluence/cmd/post"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "go-confluence",
		Short: "",
		Long:  ``,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//rootCmd.PersistentFlags().String("space", "", "Confluence space name (QA, HR, Engineering etc)")
	//rootCmd.MarkPersistentFlagRequired("space")
	//viper.BindPFlag("space", rootCmd.PersistentFlags().Lookup("space"))

	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(post.PostCmd)
}

func initConfig() {

	var configPath string

	dir, dirErr := os.UserConfigDir()
	if dirErr == nil {
		configPath = filepath.Join(dir, "go-confluence")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %s\n", err)
	}
}
