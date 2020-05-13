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

// NewDTAServerTransformPipeParams creates a new DTAServerTransformPipeParams object
// with the default values initialized.
func NewDTAServerTransformPipeParams() *DTAServerTransformPipeParams {
	var ()
	return &DTAServerTransformPipeParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDTAServerTransformPipeParamsWithTimeout creates a new DTAServerTransformPipeParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDTAServerTransformPipeParamsWithTimeout(timeout time.Duration) *DTAServerTransformPipeParams {
	var ()
	return &DTAServerTransformPipeParams{

		timeout: timeout,
	}
}

// NewDTAServerTransformPipeParamsWithContext creates a new DTAServerTransformPipeParams object
// with the default values initialized, and the ability to set a context for a request
func NewDTAServerTransformPipeParamsWithContext(ctx context.Context) *DTAServerTransformPipeParams {
	var ()
	return &DTAServerTransformPipeParams{

		Context: ctx,
	}
}

// NewDTAServerTransformPipeParamsWithHTTPClient creates a new DTAServerTransformPipeParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDTAServerTransformPipeParamsWithHTTPClient(client *http.Client) *DTAServerTransformPipeParams {
	var ()
	return &DTAServerTransformPipeParams{
		HTTPClient: client,
	}
}

/*DTAServerTransformPipeParams contains all the parameters to send to the API endpoint
for the d t a server transform pipe operation typically these are written to a http.Request
*/
type DTAServerTransformPipeParams struct {

	/*Body*/
	Body *rest_models.DtaserviceTransformPipeRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) WithTimeout(timeout time.Duration) *DTAServerTransformPipeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) WithContext(ctx context.Context) *DTAServerTransformPipeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) WithHTTPClient(client *http.Client) *DTAServerTransformPipeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) WithBody(body *rest_models.DtaserviceTransformPipeRequest) *DTAServerTransformPipeParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the d t a server transform pipe params
func (o *DTAServerTransformPipeParams) SetBody(body *rest_models.DtaserviceTransformPipeRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *DTAServerTransformPipeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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