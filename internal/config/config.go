package config

import (
	"time"

	"github.com/AidlyTeam/Aidly-Backend/internal/config/models"

	"github.com/spf13/viper"
)

const (
	defaultConfigDir              = "./config"
	defaultHTTPPort               = "8081"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultSessionExpiration      = 24 * time.Hour
	defaultManagmentPath          = "/management"
	defaultAppMode                = "self"
)

var Version string

type Config struct {
	HTTP           models.HTTPConfig     `mapstructure:"http"`
	DatabaseConfig models.DatabaseConfig `mapstructure:"database"`
	Application    models.Application    `mapstructure:"app"`
	Defaults       models.Defaults
}

func Init(configsDir ...string) (cfg *Config, err error) {
	cfg = new(Config)
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("http.session_expiration", defaultSessionExpiration)
	viper.SetDefault("managment.managmentPath", defaultManagmentPath)
	viper.SetDefault("mode", defaultAppMode)

	dir := ""
	if len(configsDir) > 0 {
		dir = configsDir[0]
	} else {
		dir = defaultConfigDir
	}

	// Viper Getting configs.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.MergeInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}

	cfg.Defaults = models.NewDefaults()

	return
}
