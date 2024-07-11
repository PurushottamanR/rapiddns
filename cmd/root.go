package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCMD *cobra.Command = &cobra.Command{
	Use:   "RapidDNS",
	Short: "A CLI tool for retrieving, filtering domain information provided by RapidDNS application",
	Long:  "A CLI tool for retrieving, filtering domain information provided by RapidDNS application",
}

func Execute() {
	err := rootCMD.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Setup adds all subcommands to the root command.
func Setup() {
	rootCMD.CompletionOptions.DisableDefaultCmd = true
	rootCMD.AddCommand(subDomainCMD)
	rootCMD.AddCommand(ipCMD)
}
