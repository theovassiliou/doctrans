package main

import (
	"flag"
	"log"
	"os"

	loads "github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	flags "github.com/jessevdk/go-flags"
	"github.com/theovassiliou/doctrans/gen/rest_api"
	"github.com/theovassiliou/doctrans/gen/rest_api/operations"
	"github.com/theovassiliou/doctrans/gen/rest_api/operations/d_t_a_server"
	"github.com/theovassiliou/doctrans/gen/rest_models"
	"jaytaylor.com/html2text"
)

var portFlag = flag.Int("port", 3000, "Port to run this service on")

func main() {

	swaggerSpec, err := loads.Embedded(rest_api.SwaggerJSON, rest_api.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewDtaserviceProtoAPI(swaggerSpec)
	server := rest_api.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "dtaservice.proto"
	parser.LongDescription = swaggerSpec.Spec().Info.Description

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	server.Port = *portFlag

	api.DtaServerTransformDocumentHandler = d_t_a_server.TransformDocumentHandlerFunc(
		func(params d_t_a_server.TransformDocumentParams) middleware.Responder {
			document := params.Body.Document
			fileName := params.Body.FileName
			log.Println(document)
			text, _ := html2text.FromString(string(document), html2text.Options{PrettyTables: true})

			return d_t_a_server.NewTransformDocumentOK().WithPayload(&rest_models.DtaserviceTransformDocumentReply{
				TransDocument: []byte(text),
				TransOutput:   []string{fileName},
			})
		})

	// api.GetGreetingHandler = operations.GetGreetingHandlerFunc(
	// 	func(params operations.GetGreetingParams) middleware.Responder {
	// 		name := swag.StringValue(params.Name)
	// 		if name == "" {
	// 			name = "World"
	// 		}

	// 		greeting := fmt.Sprintf("Hello, %s!", name)
	// 		return operations.NewGetGreetingOK().WithPayload(greeting)
	// 	})
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
