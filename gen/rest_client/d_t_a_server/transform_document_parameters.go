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

	"github.com/theovassiliou/doctrans/gen/rest_models"
)

// NewTransformDocumentParams creates a new TransformDocumentParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewTransformDocumentParams() *TransformDocumentParams {
	return &TransformDocumentParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewTransformDocumentParamsWithTimeout creates a new TransformDocumentParams object
// with the ability to set a timeout on a request.
func NewTransformDocumentParamsWithTimeout(timeout time.Duration) *TransformDocumentParams {
	return &TransformDocumentParams{
		timeout: timeout,
	}
}

// NewTransformDocumentParamsWithContext creates a new TransformDocumentParams object
// with the ability to set a context for a request.
func NewTransformDocumentParamsWithContext(ctx context.Context) *TransformDocumentParams {
	return &TransformDocumentParams{
		Context: ctx,
	}
}

// NewTransformDocumentParamsWithHTTPClient creates a new TransformDocumentParams object
// with the ability to set a custom HTTPClient for a request.
func NewTransformDocumentParamsWithHTTPClient(client *http.Client) *TransformDocumentParams {
	return &TransformDocumentParams{
		HTTPClient: client,
	}
}

/*
TransformDocumentParams contains all the parameters to send to the API endpoint

	for the transform document operation.

	Typically these are written to a http.Request.
*/
type TransformDocumentParams struct {

	// Body.
	Body *rest_models.DtaserviceDocumentRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the transform document params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *TransformDocumentParams) WithDefaults() *TransformDocumentParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the transform document params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *TransformDocumentParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the transform document params
func (o *TransformDocumentParams) WithTimeout(timeout time.Duration) *TransformDocumentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the transform document params
func (o *TransformDocumentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the transform document params
func (o *TransformDocumentParams) WithContext(ctx context.Context) *TransformDocumentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the transform document params
func (o *TransformDocumentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the transform document params
func (o *TransformDocumentParams) WithHTTPClient(client *http.Client) *TransformDocumentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the transform document params
func (o *TransformDocumentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the transform document params
func (o *TransformDocumentParams) WithBody(body *rest_models.DtaserviceDocumentRequest) *TransformDocumentParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the transform document params
func (o *TransformDocumentParams) SetBody(body *rest_models.DtaserviceDocumentRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *TransformDocumentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
