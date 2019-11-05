package dtaservice

import (
	"time"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

type config struct {
	RegistrarURL string
}

// DocTransServer is a generic server
type DocTransServer struct {
	UnimplementedDTAServerServer
	conf config
}

// RegisterAtRegistry registers the DocTransServer at the Area Registry
func (dtas *DocTransServer) RegisterAtRegistry(hostname, app, ipAddress, port, dtaType string, ttl uint, isSsl bool) {

	// Build Eureka Configuration
	client := eureka.NewClient([]string{
		dtas.conf.RegistrarURL,
		// add others servers here
	})

	// Create the app instance
	instance := eureka.NewInstanceInfo(hostname, app, ipAddress, port, ttl, isSsl) //Create a new instance to register
	// Add some meta data. Currently no meaning
	// TODO: Remove this playground if not further required
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	instance.Metadata.Map["DTA-Type"] = dtaType //one of Gateway, Service

	// Register instance and heartbeat for Eureka
	client.RegisterInstance(app, instance) // Register new instance in your eureka(s)
	log.WithFields(log.Fields{"Service": "Eureka", "Status": "Init"}).Infof("Registering service %s", app)

	job := func() {
		log.WithFields(log.Fields{"Service": "Eureka", "Status": "Up"}).Debug("sending heartbeat : ", time.Now().UTC())
		client.SendHeartbeat(instance.App, instance.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
	}

	// Run every 25 seconds but not now.
	// FIXME: We have somehow be able to deregister the heartbeat
	_, _ = scheduler.Every(25).Seconds().NotImmediately().Run(job)

}

func (dtas *DocTransServer) FindWorker() {

}
