// Code generated by go-swagger; DO NOT EDIT.

package d_t_a_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DTAServerOptionsHandlerFunc turns a function with the right signature into a d t a server options handler
type DTAServerOptionsHandlerFunc func(DTAServerOptionsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DTAServerOptionsHandlerFunc) Handle(params DTAServerOptionsParams) middleware.Responder {
	return fn(params)
}

// DTAServerOptionsHandler interface for that can handle valid d t a server options params
type DTAServerOptionsHandler interface {
	Handle(DTAServerOptionsParams) middleware.Responder
}

// NewDTAServerOptions creates a new http.Handler for the d t a server options operation
func NewDTAServerOptions(ctx *middleware.Context, handler DTAServerOptionsHandler) *DTAServerOptions {
	return &DTAServerOptions{Context: ctx, Handler: handler}
}

/*DTAServerOptions swagger:route GET /v1/service/options DTAServer dTAServerOptions

DTAServerOptions d t a server options API

*/
type DTAServerOptions struct {
	Context *middleware.Context
	Handler DTAServerOptionsHandler
}

func (o *DTAServerOptions) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDTAServerOptionsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
