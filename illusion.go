package main

import (
	"log"
	"net/http"

	"github.com/k0kubun/pp"
	"github.com/mix3/illusion/config"
	"github.com/mix3/illusion/proxy"
)

func main() {
	conf := config.Conf
	pp.Println(conf)
	http.Handle("/", proxy.NewProxy(conf))
	log.Fatal(http.ListenAndServe(conf.ListenAddr, nil))
}
