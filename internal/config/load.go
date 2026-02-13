package config

import "github.com/go-ini/ini"

func Load(path string) (*Config, error) {
	cfg := &Config{}

	iniFile, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	err = iniFile.MapTo(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
