package cmd

import (
	"fmt"
	
	"github.com/PurushottamanR/rapiddns/subdomains"
	"github.com/spf13/cobra"
)

var subDomainCMD *cobra.Command = &cobra.Command{
	Use:   "subdomain <domain>",
	Short: "list subdomains for a given domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		records, err := subdomains.NewDomain(domain).SubDomains(all, page)
		for _, record := range records {
			fmt.Println(record)
		}
		if err != nil {
			fmt.Println(err)
		}
	},
}
