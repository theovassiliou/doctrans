package main

import (
	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	ds "github.com/theovassiliou/doctrans/dtaservice"
)

var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"

type config struct {
	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

var conf = config{}

func ATestFunction() {
	var dts ds.DocTransServer

	log.Println(dts)
}

func main() {
	conf = config{}
	opts.New(&conf).
		Repo("github.com/theovassiliou/doctrans").
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)

	ATestFunction()
}
