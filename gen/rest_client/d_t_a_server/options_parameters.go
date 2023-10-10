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

// NewOptionsParams creates a new OptionsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewOptionsParams() *OptionsParams {
	return &OptionsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewOptionsParamsWithTimeout creates a new OptionsParams object
// with the ability to set a timeout on a request.
func NewOptionsParamsWithTimeout(timeout time.Duration) *OptionsParams {
	return &OptionsParams{
		timeout: timeout,
	}
}

// NewOptionsParamsWithContext creates a new OptionsParams object
// with the ability to set a context for a request.
func NewOptionsParamsWithContext(ctx context.Context) *OptionsParams {
	return &OptionsParams{
		Context: ctx,
	}
}

// NewOptionsParamsWithHTTPClient creates a new OptionsParams object
// with the ability to set a custom HTTPClient for a request.
func NewOptionsParamsWithHTTPClient(client *http.Client) *OptionsParams {
	return &OptionsParams{
		HTTPClient: client,
	}
}

/*
OptionsParams contains all the parameters to send to the API endpoint

	for the options operation.

	Typically these are written to a http.Request.
*/
type OptionsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the options params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *OptionsParams) WithDefaults() *OptionsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the options params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *OptionsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the options params
func (o *OptionsParams) WithTimeout(timeout time.Duration) *OptionsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the options params
func (o *OptionsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the options params
func (o *OptionsParams) WithContext(ctx context.Context) *OptionsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the options params
func (o *OptionsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the options params
func (o *OptionsParams) WithHTTPClient(client *http.Client) *OptionsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the options params
func (o *OptionsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *OptionsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}