# The Document Transformation Application

The Document Transformation Application is a QDS microservice based web service application, with the goal to provide one single interface to document transformation applications. 

The  protocol between a [DTA user](#user-content-dtauser) and the [DTA Server](#user-content-dtaserver) is defined using [gRPC](https://grpc.io/) and we call it [DTA Server Protocol](#user-content-thedtaserverprotocol). A [DTA worker](#user-content-dtaworker) might also  

# The DTA Server protocol
The server protocol is defined using [gRPC](https://grpc.io).

We have build the protocol using the [protobuf v3.9.1](https://github.com/protocolbuffers/protobuf/releases/tag/v3.9.1) tool.
 

# Glossary

 * DTA  - Document Transformation Application
     -  DTA client - Synonym for DTA user.
     -  DTA server - The DTA server provides an API for document transformation. The DTA server might use [DTA worker](#user-content-dtaworker) to perform the task, or other means.
     -  DTA server protocol - The protocol between DTA server and DTA user.
     -  DTA user - Is a entity that uses the DTA Server API to transform a document. Also called DTA client.
     -  DTA worker - A microservice providing *one* transformation application, potentially parametrised