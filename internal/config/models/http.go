package models

import "time"

type HTTPConfig struct {
	Host               string        `mapstructure:"host"`
	Port               string        `mapstructure:"port"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	SessionExpiration  time.Duration `mapstructure:"sessionExpiration"`
	AllowedOrigins     []string      `mapstructure:"allowedOrigins"`
	AllowedHeaders     []string      `mapstructure:"allowedHeaders"`
	AllowedMethods     []string      `mapstructure:"allowedMethods"`
	ExposedHeaders     []string      `mapstructure:"exposedHeaders"`
	AllowCredentials   bool          `mapstructure:"allowCredentials"`
	ProxyHeader        string        `mapstructure:"proxyHeader"`
}
