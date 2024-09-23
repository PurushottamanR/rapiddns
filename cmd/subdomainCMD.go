package cmd

import (
	"fmt"
	
	"github.com/PurushottamanR/rapiddns/subdomain"
	"github.com/spf13/cobra"
)

var options *subdomain.Options = subdomain.NewOptions()

func subDomainSetup(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&options.Domain, "domain", "d", "", "The domain to query for")
	cmd.MarkFlagRequired("domain")
	cmd.Flags().BoolVarP(&options.All, "all", "a", false, "Fetch all records or just first 100 records")
	cmd.Flags().IntVarP(&options.Pages, "pages", "p", 1, "Fetch records until page")
	cmd.Flags().BoolVarP(&options.Total, "total", "t", false, "Just get total no of records not the list of domains")
	cmd.Flags().IntVarP(&options.Threads, "threads", "T", 15, "No of threads to use")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Dump verbose info like the errors and pages missed")
}

var subDomainCMD *cobra.Command = &cobra.Command{
	Use:   "subdomain",
	Short: "list subdomains for a given domain",
	Run: func(cmd *cobra.Command, args []string) {
		result := subdomain.SubDomains(options)
		fmt.Println(result)
	},
}
