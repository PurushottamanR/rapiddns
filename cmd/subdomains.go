package cmd

import (
	"github.com/spf13/cobra"
)

var subDomainCMD *cobra.Command = &cobra.Command{
	Use:   "subdomain <domain>",
	Short: "list subdomains for a given domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//domain := args[0]
		
	},
}
