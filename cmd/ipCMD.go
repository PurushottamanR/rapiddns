package cmd

import (
	"fmt"	
	
	"github.com/PurushottamanR/rapiddns/iptools"
	"github.com/spf13/cobra"
)

var ipOpts *iptools.Options = iptools.NewOptions()

func ipToolsSetup(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ipOpts.IPaddr, "address", "a", "", "The IPv4 address to query for")
	cmd.MarkFlagRequired("address")
}


var ipCMD *cobra.Command = &cobra.Command{
	Use:   "ip",
	Short: "list ip details",
	Run: func(cmd *cobra.Command, args []string) {
		details := iptools.GetIPDetails(ipOpts)
		fmt.Println(details)
	},
}
