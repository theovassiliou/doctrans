# Microservices and Name Services / Registries

To use a micro services the adress has to be know apriory. Such an adress can either be configured

- statically (config files, environment variables, etc)
- dynamically (configuration servers)
- resolved (using name services)

The first two options (statically / dynamically) are typically applied in applications with a variable in general fixed set-up.

Registries, or name services as they are called some times, serve a different purpose. They resolve an abstract name to a technical adresses (like protocol, IP adresses, ports, urls).

In the context of micro services, and in particular to DTA, this means that a service name, as retrieved by `ListServices` has to be resolved to a IP-adresss and port, and for multi-protocol services like DTA (GRPC and HTTP) also the used protocol.

The DTA Framework uses the [EUREKA](https://github.com/Netflix/eureka) registry. Different docker-images are available for easy deployment. While we are using a pretty old [EUREKA Container](https://hub.docker.com/r/netflixoss/eureka) it serves well our purpose.

- It resolves service names (Application Name for Eureka) to IP-Adresses
- It offers a RESTful API for registering services and resolving the addresses
- It offers meta-data storage beyond adresses in order to store information like protocol

As for the DTA framework we define the following instance definition for a EUREKA registered service

- hostName: \[protocol://\]hostName _if protocol is not given `http` is assumed_
- app: The fully qualified application name. See also [Universe/Galaxy/Service](https://github.com/theovassiliou/doctrans/blob/master/scenarios/Universe.de.md)
- metadata
  - .dtaType: Service | Gateway _if none given `Service` assumed_
  - .dtaProto: http | grpc | (other) _if none given `http` is assumed_

Here is an excerpt of an instance defintiion as returned by 
http://eureka:8761/eureka/apps/ for a service named COUNT in the Galaxy DE.TU-BERLIN.QDS being accessible via HTTP at port 60001 and GRPC at port 60000.

```XML
<applications>
...
 <application>
    <name>DE.TU-BERLIN.QDS.COUNT</name>
    <instance>
        <hostName>http://localPC1.fritz.box</hostName>
        <app>DE.TU-BERLIN.QDS.COUNT</app>
        <ipAddr>192.168.178.60</ipAddr>
        <status>UP</status>
        <port enabled="true">60001</port>
        <securePort enabled="false">0</securePort>
        <metadata>
            <dtaType>Service</dtaType>
            <dtaProto>http</dtaProto>
        </metadata>
        ...
    </instance>
    <instance>
        <hostName>grpc://localPC1.fritz.box</hostName>
        <app>DE.TU-BERLIN.QDS.COUNT</app>
        <ipAddr>192.168.178.60</ipAddr>
        <status>UP</status>
        <port enabled="true">60000</port>
        <securePort enabled="false">0</securePort>
        <metadata>
            <dtaType>Service</dtaType>
            <dtaProto>grpc</dtaProto>
        </metadata>
        ...
    </instance>
</application>
</applications>
```