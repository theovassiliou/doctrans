/*
Copyright (c) 2019 Theofanis Vassiliou-Gioles

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

//go:generate protoc -I ../dtaservice --go_out=plugins=grpc:../dtaservice ../dtaservice/dtaservice.proto

// Package main implements a client for DtaService.
package main

import (
	"github.com/jpillora/opts"

	log "github.com/sirupsen/logrus"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"
var conf = config{}

type config struct {
	FileName       string    `type:"arg" name:"file" help:" the file to be uploaded"`
	EurekaURL      string    `help:"if set the indicated eureka server will be used to find DTA-GW"`
	ServiceName    string    `help:"The service to be used"`
	GatewayAddress string    `help:"Address and port of the server to connect"`
	ListServices   bool      `help:"List all the services accessible"`
	LogLevel       log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

type message struct {
	Document []byte `json:document`
	FileName string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const dtaGwID = "DE.TU-BERLIN.QDS.DTA.DTA-GW"

func main() {

	conf = config{
		GatewayAddress: "127.0.0.1:50051",
		ServiceName:    "DE.TU-BERLIN.QDS.DTA.COUNT",
		EurekaURL:      "http://eureka:8761/eureka",
	}

	//parse config
	opts.New(&conf).
		Repo("github.com/theovassiliou/dta").
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)
}
