package cmd

import (
	"github.com/spf13/cobra"
)


var ipCMD *cobra.Command = &cobra.Command{
	Use:   "ip <address>",
	Short: "list domains for a given IP",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//ip := args[0]
		
	},
}
