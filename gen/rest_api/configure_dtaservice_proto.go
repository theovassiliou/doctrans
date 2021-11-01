// This file is safe to edit. Once it exists it will not be overwritten

package rest_api

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/theovassiliou/doctrans/gen/rest_api/operations"
	"github.com/theovassiliou/doctrans/gen/rest_api/operations/d_t_a_server"
)

//go:generate swagger generate server --target ../../../doctrans --name DtaserviceProto --spec ../../swagger/dtaservice.swagger.json --model-package gen/rest_models --server-package gen/rest_api

func configureFlags(api *operations.DtaserviceProtoAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.DtaserviceProtoAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.DtaServerListServicesHandler == nil {
		api.DtaServerListServicesHandler = d_t_a_server.ListServicesHandlerFunc(func(params d_t_a_server.ListServicesParams) middleware.Responder {
			return middleware.NotImplemented("operation d_t_a_server.ListServices has not yet been implemented")
		})
	}
	if api.DtaServerOptionsHandler == nil {
		api.DtaServerOptionsHandler = d_t_a_server.OptionsHandlerFunc(func(params d_t_a_server.OptionsParams) middleware.Responder {
			return middleware.NotImplemented("operation d_t_a_server.Options has not yet been implemented")
		})
	}
	if api.DtaServerTransformDocumentHandler == nil {
		api.DtaServerTransformDocumentHandler = d_t_a_server.TransformDocumentHandlerFunc(func(params d_t_a_server.TransformDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation d_t_a_server.TransformDocument has not yet been implemented")
		})
	}
	if api.DtaServerTransformPipeHandler == nil {
		api.DtaServerTransformPipeHandler = d_t_a_server.TransformPipeHandlerFunc(func(params d_t_a_server.TransformPipeParams) middleware.Responder {
			return middleware.NotImplemented("operation d_t_a_server.TransformPipe has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
