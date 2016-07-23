package utils

import (
	"fmt"
	"log"

	"gopkg.in/gcfg.v1"
)

type GConfig struct {
	Database  DatabaseConfig
	Scheduler SchedulerConfig
	Token     TokenConfig
}

type TokenConfig struct {
	VerifyToken string
	Verify      bool
	PageToken   string
	AccessToken string
	PageID      string
}

type DatabaseConfig struct {
	Order   string
	User    string
	Payment string
}

type SchedulerConfig struct {
	Run bool
}

var ConfigG *GConfig

func NewConfig(filePath string) (*GConfig, error) {
	if ConfigG == nil {
		var c GConfig

		err := gcfg.ReadFileInto(&c, filePath)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Could not find configuration file")
		}

		ConfigG = &c
	}

	return ConfigG, nil
}
