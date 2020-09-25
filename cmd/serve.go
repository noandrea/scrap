package cmd

import (
	"fmt"

	"github.com/noandrea/scrap/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the REST API endpoints",
	Long:  ``,
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "load configuration from file")
}

func serve(cmd *cobra.Command, args []string) {

	fmt.Printf(`>>>>>>>>>>>>>>>
╔═╗┌─┐┬─┐┌─┐┌─┐
╚═╗│  ├┬┘├─┤├─┘
╚═╝└─┘┴└─┴ ┴┴ 	v%s
>>>>>  started!`, rootCmd.Version)
	fmt.Println()

	var settings server.ConfigSchema
	server.Defaults()
	viper.AutomaticEnv() // read in environment variables that match
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("cannot read from %s: %v", cfgFile, err)
			return
		}
	}
	viper.Unmarshal(&settings)
	// make the version available via settings
	settings.RuntimeVersion = rootCmd.Version
	log.Debugf("settings %#v", settings)
	err := server.Start(settings)
	if err != nil {
		fmt.Printf("cannot start web server: %v", err)
	}
}
