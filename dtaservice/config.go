package dtaservice

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	aux "github.com/theovassiliou/doctrans/ipaux"
)

func newDefaultDTS() *DocTransServer {
	return &DocTransServer{
		Register:     false,
		REST:         true,
		HTTPPort:     defaultOrNot("80", os.Getenv("DTS_HTTPPort")),
		HostName:     defaultOrNot(aux.GetHostname(), os.Getenv("DTS_HostName")),
		AppName:      defaultOrNot("", os.Getenv("DTS_AppName")),
		PortToListen: defaultOrNot("50051", os.Getenv("DTS_PortToListen")),
		RegistrarURL: defaultOrNot("http://127.0.0.1:8761/eureka", os.Getenv("DTS_RegistrarURL")),
		DtaType:      "Service",
		LogLevel:     log.WarnLevel,
	}
}

// NewConfigFile creates a new example config file and terminates.
func (dtas *DocTransServer) NewConfigFile() error {

	dir := path.Dir(dtas.CfgFile)

	_, err := os.Open(dir)
	if err != nil {
		os.MkdirAll(dir, os.ModePerm)
		_, err := os.Open(dir)
		if err != nil {
			return err
		}
	}

	configJSON, _ := json.MarshalIndent(dtas, "", "  ")
	err = ioutil.WriteFile(dtas.CfgFile, configJSON, 0644)
	log.Infof("Wrote example configuration file to %s. Exiting.", dtas.CfgFile)
	return nil
}

func SetupConfiguration(config *DocTransServer, workingHomeDir, version string) {
	opts.New(config).
		Repo("github.com/theovassiliou/doctrans").
		Version(version).
		Parse()

	if config.LogLevel != 0 {
		log.SetLevel(config.LogLevel)
	}

	if config.AppName != "" && config.CfgFile != "" {
		config.CfgFile = workingHomeDir + "/.dta/" + config.AppName + "/config.json"
	}

	if config.Init {
		config.CfgFile = config.CfgFile + ".example"
		err := config.NewConfigFile()
		if err != nil {
			log.Fatalln(err)
		}
		log.Exit(0)
	}

	// Parse config file
	config, err := NewDocTransFromFile(config.CfgFile)
	if err != nil {
		log.Infoln("No config file found. Consider creating one using --init option.")
	}

	// Parse command line parameters again to insist on config parameters
	opts.New(config).Parse()
	if config.LogLevel != 0 {
		log.SetLevel(config.LogLevel)
	}

}

// NewDocTransFromReader creates a Client configured from a given reader.
// The configuration is expected to use the JSON format.
func NewDocTransFromReader(reader io.Reader) (*DocTransServer, error) {
	d := newDefaultDTS()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// NewDocTransFromFile creates a DocTransServer from a given file path.
// The given file is expected to use the JSON format.
func NewDocTransFromFile(fpath string) (*DocTransServer, error) {
	fi, err := os.Open(fpath)
	if err != nil {
		return newDefaultDTS(), err
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	return NewDocTransFromReader(fi)
}

func defaultOrNot(d, v string) string {
	if v == "" {
		return d
	}
	return v
}
