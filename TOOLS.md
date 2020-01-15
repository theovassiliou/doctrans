# Tools used

For the creation and maintenance of the DTA infrastructure we are using the following tool

- go / golang
    Just a cool language, that we are happy to able to use.
- [grpc.io](http://grpc.ip)
- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
    The grpc-gateway is a plugin of the Google protocol buffers compiler protoc. It reads protobuf service definitions and generates a reverse-proxy server which 'translates a RESTful HTTP API into gRPC. This server is generated according to the google.api.http annotations in your service definitions.

    We are using this to create bi-protocol servers, that support both GRPC and REST (optional) communication. And to create the swagger specification out of the proto specification.
- [github.com/go-swagger](https://github.com/go-swagger/go-swagger)
    go-swagger brings to the go community a complete suite of fully-featured, high-performance, API components to work with a Swagger API: server, client and data model.

    We are using this to create REST clients, as well to pilot REST only servers


swagger generate client -h