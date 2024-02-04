package config

import "github.com/spf13/viper"

func Initialize(path string) error {
	v = viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
