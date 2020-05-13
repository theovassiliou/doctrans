// Code generated by go-swagger; DO NOT EDIT.

package d_t_a_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DTAServerTransformPipeHandlerFunc turns a function with the right signature into a d t a server transform pipe handler
type DTAServerTransformPipeHandlerFunc func(DTAServerTransformPipeParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DTAServerTransformPipeHandlerFunc) Handle(params DTAServerTransformPipeParams) middleware.Responder {
	return fn(params)
}

// DTAServerTransformPipeHandler interface for that can handle valid d t a server transform pipe params
type DTAServerTransformPipeHandler interface {
	Handle(DTAServerTransformPipeParams) middleware.Responder
}

// NewDTAServerTransformPipe creates a new http.Handler for the d t a server transform pipe operation
func NewDTAServerTransformPipe(ctx *middleware.Context, handler DTAServerTransformPipeHandler) *DTAServerTransformPipe {
	return &DTAServerTransformPipe{Context: ctx, Handler: handler}
}

/*DTAServerTransformPipe swagger:route POST /v1/document/transform-pipe DTAServer dTAServerTransformPipe

DTAServerTransformPipe d t a server transform pipe API

*/
type DTAServerTransformPipe struct {
	Context *middleware.Context
	Handler DTAServerTransformPipeHandler
}

func (o *DTAServerTransformPipe) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDTAServerTransformPipeParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}