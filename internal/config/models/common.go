package models

type Managment struct {
	WalletAddress     string `mapstructure:"walletAddress"`
	ManagmentUsername string `mapstructure:"username"`
	ManagmentPassword string `mapstructure:"password"`
}
