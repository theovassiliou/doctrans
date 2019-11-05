dtaservice/dtaservice.pb.go: dtaservice/dtaservice.proto
	protoc -I dtaservice/ dtaservice.proto --go_out=plugins=grpc:dtaservice