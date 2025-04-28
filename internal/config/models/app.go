package models

type Application struct {
	Site           string    `mapstructure:"site"`
	Https          bool      `mapstructure:"https"`
	Mode           string    `mapstructure:"mode"`
	Secret         string    `mapstructure:"sercret"`
	DevMode        bool      `mapstructure:"devMode"`
	Managment      Managment `mapstructure:"managment"`
	MigrationsPath string    `mapstructure:"migrationsPath"`
}
