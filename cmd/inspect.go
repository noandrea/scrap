package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/noandrea/scrap/pkg/scrap"
	"github.com/spf13/cobra"
)

var provider, region, chromeAddress string
var quiet bool

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect ASIN",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run:   inspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	inspectCmd.Flags().StringVarP(&provider, "provider", "p", scrap.AmazonPrime, "set the provider")
	inspectCmd.Flags().StringVarP(&region, "region", "r", "de", "set the region for the query")
	inspectCmd.Flags().StringVarP(&chromeAddress, "chrome-address", "a", "http://127.0.0.1:9222", "set the address for headless chrome")
	inspectCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "only output the json result (to be used with jq)")

}

func inspect(cmd *cobra.Command, args []string) {
	id := args[0]
	start := time.Now()
	// inspect
	scrap.Configure(chromeAddress)
	// run scrap
	var movie scrap.Movie
	err := scrap.Run(provider, id, region, &movie)
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := json.MarshalIndent(movie, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(s))
	if !quiet {
		fmt.Println("took ", time.Since(start))
	}
}
