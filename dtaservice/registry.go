package dtaservice

import (
	"strconv"
	"time"

	"github.com/carlescere/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/go-eureka-client/eureka"
)

func (dtas *GenDocTransServer) InstanceInfo() *eureka.InstanceInfo {
	return dtas.instanceInfo
}

func (dtas *GenDocTransServer) UnregisterAtRegistry() {
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
func (dtas *GenDocTransServer) RegisterAtRegistry(registerURL string) {
	// Build Eureka Configuration
	dtas.registrar = eureka.NewClient([]string{
		registerURL,
	})
	dtas.registrar.CheckRetry = eureka.ExpBackOffCheckRetry
	// Create the app instance

	// Register instance and heartbeat for Eureka
	dtas.registrar.RegisterInstance(dtas.instanceInfo.App, dtas.instanceInfo) // Register new instance in your eureka(s)
	log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Init"}).Infof("Registering service %s\n", dtas.instanceInfo.App)
	log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Init"}).Tracef("InstanceInfo %v\n", dtas.instanceInfo)

	job := func() {
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Up"}).Tracef("sending heartbeat : %v\n", time.Now().UTC())
		log.WithFields(log.Fields{"Service": "->Registrar", "Status": "Upt"}).Tracef("InstanceInfo %v\n", dtas.instanceInfo)
		dtas.registrar.SendHeartbeat(dtas.instanceInfo.App, dtas.instanceInfo.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
	}

	// Run every 25 seconds but not now.
	// FIXME:0 We have somehow be able to deregister the heartbeat
	dtas.heartBeatJob, _ = scheduler.Every(25).Seconds().NotImmediately().Run(job)
}

// Registrar returns the Eureka intance where the server has registered.
func (dtas *GenDocTransServer) Registrar() *eureka.Client {
	return dtas.registrar
}

// NewInstanceInfo
func (dtas *GenDocTransServer) NewInstanceInfo(hostName, app, ip string, port int, ttl uint, isSsl bool, dtaType, proto string, homePageURL, statusURL, healthURL string) *eureka.InstanceInfo {
	dataCenterInfo := &eureka.DataCenterInfo{
		Name:     "MyOwn",
		Class:    "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		Metadata: nil,
	}
	leaseInfo := &eureka.LeaseInfo{
		EvictionDurationInSecs: ttl,
	}
	instanceInfo := &eureka.InstanceInfo{
		HostName:       hostName,
		App:            app,
		IpAddr:         ip,
		Status:         eureka.UP,
		DataCenterInfo: dataCenterInfo,
		LeaseInfo:      leaseInfo,
		Metadata:       nil,
	}
	instanceInfo.Port = &eureka.Port{
		Port:    strconv.Itoa(port),
		Enabled: "true",
	}

	instanceInfo.SecurePort = &eureka.Port{
		Port:    "0",
		Enabled: "false",
	}

	instanceInfo.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	instanceInfo.Metadata.Map["dtaType"] = dtaType //one of Gateway, Service
	instanceInfo.Metadata.Map["dtaProto"] = proto
	instanceInfo.HomePageUrl = homePageURL
	instanceInfo.StatusPageUrl = statusURL
	instanceInfo.HealthCheckUrl = healthURL
	dtas.instanceInfo = instanceInfo
	return instanceInfo
}
