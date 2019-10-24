# Flogo Services

Flogo services module consists of two Services

## State Service
The Flogo State Service is a service for managing the state of process flows executed on Flogo Engine. This service's primary job is to store process (incremental and full) state for flows that are executed on an engine. This service will also facilitate front-end introspection and debugging of process flows.

## Flow Service
The Flogo Process Service is a service for managing process definitions that are to be executed by the Flogo Engine. This service's primary job is to store process definitions designed in the front-end and provide those definitions to a Flogo Engine on demand.


## Building the services

The State service and Flow service are both write by Golang, so please make sure you have installed Golang on your machine

* Clone services code from github
```bash
cd project-flogo/services
#Please create directory if it is not exist
git clone https://github.com/project-flogo/services.git

```
* Build those 2 services
```bash
cd services/flow-store
go test ./...
go build
cd services/flow-state
go test ./...
go build
```

## Running the Services
* Under each service where generated 2 executable file, flow-store and flow-state
* Start flow and state server
```bash
cd GOPATH/bin
./flow-store -p 9090
./flow-state -p 9190

#-p means the port of flow or state service
```

## License
services is licensed under a BSD-type license. See [LICENSE](LICENSE) for license text.
