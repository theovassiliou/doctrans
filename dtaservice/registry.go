package dtaservice

import (
	"time"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

func (dtas *DocTransServer) InstanceInfo() *eureka.InstanceInfo {
	return dtas.instanceInfo
}
func (dtas *DocTransServer) UnregisterAtRegistry() {
	if dtas.instanceInfo != nil {
		dtas.registrar.UnregisterInstance(dtas.instanceInfo.App, dtas.instanceInfo.HostName) // unregister the instance in your eureka(s)
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Unregister"}).Infof("Unregister service %s with id %s", dtas.instanceInfo.App, dtas.instanceInfo.InstanceID)
	} else {
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Unregister"}).Infof("service %s allready unregistered", dtas.instanceInfo.App)
	}
	dtas.heartBeatJob.Quit <- true
	dtas.instanceInfo = nil
}

// RegisterAtRegistry registers the DocTransServer at the Area Registry
func (dtas *DocTransServer) RegisterAtRegistry(hostname, app, ipAddress, port, dtaType string, ttl uint, isSsl bool) {

	// Build Eureka Configuration
	dtas.registrar = eureka.NewClient([]string{
		dtas.RegistrarURL,
		// add others servers here
	})

	// Create the app instance
	dtas.instanceInfo = eureka.NewInstanceInfo(hostname, app, ipAddress, port, ttl, isSsl) //Create a new instance to register
	// Add some meta data. Currently no meaning
	// TODO: Remove this playground if not further required
	dtas.instanceInfo.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	dtas.instanceInfo.Metadata.Map["DTA-Type"] = dtaType //one of Gateway, Service
	// Register instance and heartbeat for Eureka
	dtas.registrar.RegisterInstance(app, dtas.instanceInfo) // Register new instance in your eureka(s)
	log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Init"}).Infof("Registering service %s", app)

	job := func() {
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Up"}).Trace("sending heartbeat : ", time.Now().UTC())
		dtas.registrar.SendHeartbeat(dtas.instanceInfo.App, dtas.instanceInfo.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
	}

	// Run every 25 seconds but not now.
	// FIXME:0 We have somehow be able to deregister the heartbeat
	dtas.heartBeatJob, _ = scheduler.Every(25).Seconds().NotImmediately().Run(job)
}

// Registrar returns the Eureka intance where the server has registered.
func (dtas *DocTransServer) Registrar() *eureka.Client {
	return dtas.registrar
}
