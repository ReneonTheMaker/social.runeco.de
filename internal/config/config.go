package config

type WebConfig struct {
	Port string `ini:"port"`
}

type Config struct {
	Web WebConfig 			`ini:"web"`
}
