package config

type WebConfig struct {
	Port  string `ini:"port"`
	Https bool   `ini:"https"`
}

type Config struct {
	Web WebConfig `ini:"web"`
}
