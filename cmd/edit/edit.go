package edit

import (
	"github.com/spf13/cobra"
)

var (
	flagSpace string
	flagTitle string
	flagFile  string
	flagXML   bool
	EditCmd   = &cobra.Command{
		Use:   "edit",
		Short: "",
		Long:  ``,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flagSpace = cmd.Parent().Flags().Lookup("space").Value.String()
			flagTitle = cmd.Parent().Flags().Lookup("title").Value.String()
			flagFile = cmd.Parent().Flags().Lookup("file").Value.String()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {}
