# Scenarios on localhost

Here we provide different start scripts to start different scenarios easily on localhost.

![Our universe](./../scenarios.png)

## What it is good for

While we tried to keep the scripts as generic as possible, their goals for us was to facilitate the creation of a *universe*.

Our universe (as seen in the above sketch) consists of services that can live in a global scope and service that can liven in a local scope.  
The services in each scope can be addressed via the scopes registry. In our case an EUREKA server (E).

Services (SRV1 .. SRVn) in each scope are prefixed with the scope name, e.g. `DE.TU-BERLIN.QDS`. The global scope's prefix is `<empty>`.

In the following we assume that our *universe*

- lives on a single machine,
- that this single machine serves both
  - the global scope as well as
  - all local scopes, via docker-containers.

### Global scope

To create a useful global scope we need a registry E as well as some services.
We offer support to start all QDS-services, i.e. all services that are implemented under [doctrans/services](../../services/) as own service on `localhost`

[x] eureka.sh
  Starts eureka registry server on localhost listening the port 8761. While the server itself runs in a docker container, the server behaves as it would be running directly on the host

[x] qds_<serviceName>.sh / all_services.sh
  starts the respective service with grpc and REST support individually on localhost directly. Port is selected dynamically and the service is registering itself at the local eureka service.

[x] killall_services.sh
  send a -9 signal to all services started by `qds_<serviceName>.sh`  or `all_services.sh`

#### Usage example

To create a usable global scope environment you could try

```shell
$> cd $(PRJHOME)/
$> scenarios/global/eureka.sh && scenarios/global/qds_all.sh &
$> go run client/client.go
```

### Local scope

As mentioned each scope requires it's "local" scope registry. To access the services that live in a local scope the services have to be used via the scope's gateway.

Consequently we support the setup of a following scenario with the help of a `docker-compose` scenario.
- scope prefix: `DE.TU-BERLIN.QDS`
- Eureka server
- Gateway registered at global scope registry. Exposed port is statically defined
- All QDS-services, i.e. all services that are implemented under [doctrans/services](../../services/) registered at local scope registry

#### Usage example

To create a usable global scope environment you could try

```shell
$> cd $(PRJHOME)/
$> scenarios/localhost/eureka.sh && scenarios/localhost/qds_all.sh &
$> go run client/client.go
```
