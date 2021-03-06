// Code generated by go-swagger; DO NOT EDIT.

package d_t_a_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewDTAServerListServicesParams creates a new DTAServerListServicesParams object
// with the default values initialized.
func NewDTAServerListServicesParams() *DTAServerListServicesParams {

	return &DTAServerListServicesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDTAServerListServicesParamsWithTimeout creates a new DTAServerListServicesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDTAServerListServicesParamsWithTimeout(timeout time.Duration) *DTAServerListServicesParams {

	return &DTAServerListServicesParams{

		timeout: timeout,
	}
}

// NewDTAServerListServicesParamsWithContext creates a new DTAServerListServicesParams object
// with the default values initialized, and the ability to set a context for a request
func NewDTAServerListServicesParamsWithContext(ctx context.Context) *DTAServerListServicesParams {

	return &DTAServerListServicesParams{

		Context: ctx,
	}
}

// NewDTAServerListServicesParamsWithHTTPClient creates a new DTAServerListServicesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDTAServerListServicesParamsWithHTTPClient(client *http.Client) *DTAServerListServicesParams {

	return &DTAServerListServicesParams{
		HTTPClient: client,
	}
}

/*DTAServerListServicesParams contains all the parameters to send to the API endpoint
for the d t a server list services operation typically these are written to a http.Request
*/
type DTAServerListServicesParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the d t a server list services params
func (o *DTAServerListServicesParams) WithTimeout(timeout time.Duration) *DTAServerListServicesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the d t a server list services params
func (o *DTAServerListServicesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the d t a server list services params
func (o *DTAServerListServicesParams) WithContext(ctx context.Context) *DTAServerListServicesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the d t a server list services params
func (o *DTAServerListServicesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the d t a server list services params
func (o *DTAServerListServicesParams) WithHTTPClient(client *http.Client) *DTAServerListServicesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the d t a server list services params
func (o *DTAServerListServicesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *DTAServerListServicesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
