# DocTransService ECHO

## FQSN: DE.TU-BERLIN.QDS.ECHO
 
- SERVICE-NAME: ECHO
- DEFAULT-GALAXY: DE.TU-BERLIN.QDS

## PURPOSE:

To send a document back that has been received.

## OPTIONS:

NONE

## OPERATIONS IMPLEMENTED:

- TransformDocument
- ListServices
- TransformPipe
- Options

## SUPPORTED PROTCOLS:

- GRPC
- HTTP

## PORTS:

Starting from 9230 the next available ports

## SUPPORTED SERVICE-DIRECTORY

Eureka

## OPERATIONS

### TransformDocument

Receives a document and send the same document back. 

