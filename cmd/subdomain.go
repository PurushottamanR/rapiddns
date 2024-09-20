package cmd

import (
	"fmt"
	
	"github.com/PurushottamanR/rapiddns/subdomains"
	"github.com/spf13/cobra"
)

var subDomainoptions *subdomains.Options = &subdomains.Options{}

func subDomainSetup(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&subDomainoptions.All, "all", "a", false, "Fetch all records or just first 100 records")
	cmd.Flags().IntVarP(&subDomainoptions.Page, "page", "p", 1, "Fetch records until page")
	cmd.Flags().BoolVarP(&subDomainoptions.Verbose, "verbose", "v", false, "Dump records as obtained")
}

var subDomainCMD *cobra.Command = &cobra.Command{
	Use:   "subdomain <domain>",
	Short: "list subdomains for a given domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subDomainoptions.Domain = args[0]
		recs, err := subdomains.SubDomains(subDomainoptions)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(recs)
		
	},
}
