package logger

type Config struct {
	Channels []ChannelConfig
}

type ChannelConfig struct {
	Name       string
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
	Level      string
}
