package flag

import (
	"log"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Config string `long:"config" short:"c" description:"specify config file" default:"config.toml"`
}

var Opts Options

func init() {
	parser := flags.NewParser(&Opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}
}
