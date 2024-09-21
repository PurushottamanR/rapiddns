package cmd

import (
	"fmt"
	
	"github.com/PurushottamanR/rapiddns/subdomain"
	"github.com/spf13/cobra"
)

var options *subdomain.Options = &subdomain.Options{}

func subDomainSetup(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&options.All, "all", "a", false, "Fetch all records or just first 100 records")
	cmd.Flags().IntVarP(&options.Page, "page", "p", 1, "Fetch records until page")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Dump records as obtained")
}

var subDomainCMD *cobra.Command = &cobra.Command{
	Use:   "subdomain <domain>",
	Short: "list subdomains for a given domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		
		options.Domain = args[0]
		
		result := subdomain.SubDomains(options)
		fmt.Println(result)
		
	},
}
