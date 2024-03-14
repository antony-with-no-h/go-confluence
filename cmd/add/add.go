package add

import (
	"github.com/spf13/cobra"
)

var (
	flagSpace  string
	flagTitle  string
	flagFile   string
	flagParent string

	AddCmd = &cobra.Command{
		Use:   "add",
		Short: "",
		Long:  ``,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flagSpace = cmd.Parent().Flags().Lookup("space").Value.String()
			flagTitle = cmd.Parent().Flags().Lookup("title").Value.String()
			flagFile = cmd.Parent().Flags().Lookup("file").Value.String()
			flagParent = cmd.Parent().Flags().Lookup("parent").Value.String()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {}
