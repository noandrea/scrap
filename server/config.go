package server

import (
	"github.com/spf13/viper"
)

// ConfigSchema main configuration for the news room
type ConfigSchema struct {
	ListenAddress  string `mapstructure:"listen_address"`
	CacheEnabled   bool   `mapstructure:"cache_enabled"`
	CacheMaxSize   int    `mapstructure:"cache_max_size"`
	CacheLifetime  int    `mapstructure:"cache_lifetime"`
	ChromeAddress  string `mapstructure:"chrome_address"`
	ScrapRegion    string `mapstructure:"scrap_region"`
	RuntimeVersion string `mapstructure:"-"`
}

// Defaults configure defaults for the configuration
func Defaults() {
	// web
	viper.SetDefault("listen_address", ":8080")
	viper.SetDefault("cache_enabled", true)
	viper.SetDefault("cache_max_size", 5000)
	viper.SetDefault("cache_lifetime", 86400)
	viper.SetDefault("chrome_address", ":9222")
	viper.SetDefault("scrap_region", "de")
}

// Validate a configuration
func Validate(schema *ConfigSchema) (err []error) {
	return
}

// Settings general settings
var Settings ConfigSchema
