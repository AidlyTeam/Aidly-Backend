package models

type DatabaseConfig struct {
	Managment Managment `mapstructure:"managment"`
	SSLMode   string    `mapstructure:"sslmode"`
	Driver    string    `mapstructure:"driver"`
	DBName    string    `mapstructure:"dbName"`
	Host      string    `mapstructure:"host"`
	Port      string    `mapstructure:"port"`
	Timezone  string    `mapstructure:"UTC"`
}
