package utils

import (
	"fmt"
	"gopkg.in/gcfg.v1"
)

type GConfig struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Finmate string
}

var ConfigG *GConfig

func NewConfig(filePath string) (*GConfig, error) {
	if ConfigG == nil {
		var c GConfig

		err := gcfg.ReadFileInto(&c, filePath)
		if err != nil {
			return nil, fmt.Errorf("Could not find configuration file")
		}

		ConfigG = &c
	}

	return ConfigG, nil
}
