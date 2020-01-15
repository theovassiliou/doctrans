build: dtaservice/dtaservice.pb.go swagger/index.html dtaservice/dtaservice.validator.pb.go gen/rest_client/dtaservice_proto_client.go gen/rest_api/configure_dtaservice_proto.go
	go build ./...


.PHONY: gateway 

gateway: 
	go run gateway/gateway.go

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
