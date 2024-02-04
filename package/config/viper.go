package config

import (
	"github.com/spf13/viper"
)

func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return v.UnmarshalKey(key, rawVal, opts...)
}

func GetString(key string) string {
	return v.GetString(key)
}
