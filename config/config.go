package config

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/mix3/illusion/flag"
)

type Config struct {
	Domain         string `toml:"domain"`
	ListenAddr     string `toml:"listen_addr"`
	DockerEndpoint string `toml:"docker_endpoint"`
}

var Conf Config

func init() {
	opts := flag.Opts
	if _, err := toml.DecodeFile(opts.Config, &Conf); err != nil {
		log.Fatal(err)
	}
}
