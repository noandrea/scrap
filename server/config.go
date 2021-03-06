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
	viper.SetDefault("chrome_address", "http://127.0.0.1:9222")
	viper.SetDefault("scrap_region", "de")
	viper.SetDefault("cache_enabled", true)   // enable cache
	viper.SetDefault("cache_max_size", 10)    // 10 megabytes
	viper.SetDefault("cache_lifetime", 86400) // one day
}

// Validate a configuration
func Validate(schema *ConfigSchema) (err []error) {
	return
}

// Settings general settings
var Settings ConfigSchema
