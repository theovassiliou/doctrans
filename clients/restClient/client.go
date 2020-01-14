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

// Package main implements a client for DtaService.
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/doctrans/models"
	apiclient "github.com/theovassiliou/doctrans/restClient"
	"github.com/theovassiliou/doctrans/restClient/d_t_a_server"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = ".1"

// VERSION is the current version number.
var VERSION = "0.0" + version + "-src"
var conf = config{
	LogLevel: log.TraceLevel,
}

type config struct {
	FileName string    `help:"file name of the file to be translated"`
	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	conf = config{}

	//parse config
	opts.New(&conf).
		Repo("github.com/theovassiliou/doctrans").
		Version(VERSION).
		Parse()

	log.SetLevel(conf.LogLevel)

	if conf.FileName == "" {
		log.Fatalln("No file provided.")
	}
	// make the request to get all items
	params := d_t_a_server.NewTransformDocumentParams()
	_, fileContent := readFile(conf.FileName)

	params.SetBody(&models.DtaserviceDocumentRequest{
		Document: fileContent,
		FileName: "conf.FileName",
	})
	resp, err := apiclient.Default.DtaServer.TransformDocument(params)

	if err != nil {
		log.Fatal(err)
	}
	doc := string(resp.GetPayload().TransDocument)
	fmt.Printf("%s\n", doc)
}

func readFile(path string) (int, []byte) {
	dat, readErr := ioutil.ReadFile(path)

	if readErr != nil {
		log.Fatal(readErr)
	}

	file, openErr := os.Open(path)
	if openErr != nil {
		log.Fatal(openErr)
	}
	defer file.Close()

	var noOfLines int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		noOfLines++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return noOfLines, dat
}
