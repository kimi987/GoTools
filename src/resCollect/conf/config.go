package conf

import (
	"github.com/BurntSushi/toml"
	"log"
)

type config struct {
	ResPath string
}

var Config *config

func InitConfig() {
	var cg config
	if _, err := toml.DecodeFile("conf/config.toml", &cg); err != nil {
		log.Fatal(err)
	}
	Config = &cg
}
