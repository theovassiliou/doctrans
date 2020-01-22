servers := "./services/qds_count" "./services/qds_echo" "./services/qds_html2text" "./services/rest_html2text" 
clients := "./clients/grpcClient" "./clients/restClient"


gofiles := $(subst ./services/,,$(servers))
dockerexes := $(subst ./services/,./docker/,$(servers))
dockerfiles := $(subst ./services/,./docker/Dockerfile.,$(servers))

build: dtaservice/dtaservice.pb.go swagger/index.html dtaservice/dtaservice.validator.pb.go gen/rest_client/dtaservice_proto_client.go gen/rest_api/configure_dtaservice_proto.go
	go build ./...


.PHONY: gateway 

all: server clients

server: $(servers)
clients: $(clients)

docker: $(dockerexes)

$(servers):
	go build -o bin/$(subst ./services/,,$@) $@

$(clients):
	go build -o bin/$(subst ,,$@) $@


$(dockerexes):
	rm -f docker/Dockerfile.$(subst ./docker/,,$@); \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $@ $(subst ./docker/,./services/,$@) ; \
	echo "FROM scratch" >> docker/Dockerfile.$(subst ./docker/,,$@) ; \
	echo "EXPOSE 50051" >> docker/Dockerfile.$(subst ./docker/,,$@) ; \
	echo "ADD \"$@\" "/ >> docker/Dockerfile.$(subst ./docker/,,$@) ; \
	echo "CMD [\"/$(subst ./docker/,,$@)\"]" >> docker/Dockerfile.$(subst ./docker/,,$@) ; \
	docker build -t $(subst ./docker/,,$@) -f docker/Dockerfile.$(subst ./docker/,,$@) . 

./docker/gateway: 
	rm -rf docker/Dockerfile.gateway; \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./docker/gateway ./gateway ; \
	echo "FROM scratch" >> docker/Dockerfile.gateway ; \
	echo "EXPOSE 50051" >> docker/Dockerfile.gateway ; \
	echo "ADD \"./docker/gateway\" "/ >> docker/Dockerfile.gateway ; \
	echo "CMD [\"/gateway\"]" >> docker/Dockerfile.gateway ; 


clean: 
ifneq ($(gofiles),)
	rm -f $(gofiles)
endif
ifneq ($(dockerexes),)
	rm -f $(dockerexes)
endif
ifneq ($(dockerfiles),)
	rm -f $(dockerfiles)
endif
	rm -rf bin/
	rm -rf docker/gateway

print-%  : ; @echo $* = $($*)

gateway: 
	go build -o bin/gateway ./gateway/

swagger/index.html: swagger/dtaservice.swagger.json
	swagger-codegen generate -o swagger -i swagger/dtaservice.swagger.json -l html

dtaservice/dtaservice.pb.go dtaservice/dtaservice.validator.pb.go: dtaservice/dtaservice.proto
	protoc -Idtaservice/ \
		-I/usr/local/include -I. \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:dtaservice \
  		--grpc-gateway_out=logtostderr=true:dtaservice \
  		--swagger_out=logtostderr=true:swagger \
		--govalidators_out=dtaservice \
	dtaservice.proto 


gen/rest_client/dtaservice_proto_client.go: swagger/dtaservice.swagger.json
	swagger generate client -c gen/rest_client -s gen/rest_api -m gen/rest_models -f swagger/dtaservice.swagger.json

gen/rest_api/configure_dtaservice_proto.go: swagger/dtaservice.swagger.json
	swagger generate server -c gen/rest_client -s gen/rest_api -m gen/rest_models -f swagger/dtaservice.swagger.json

