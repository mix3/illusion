package main

import (
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/mix3/illusion/config"
	"github.com/mix3/illusion/proxy"
)

func main() {
	conf := config.Conf
	spew.Dump(conf)
	http.Handle("/", proxy.NewProxy(conf))
	log.Fatal(http.ListenAndServe(conf.ListenAddr, nil))
}
