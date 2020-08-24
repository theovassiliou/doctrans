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

// NewDTAServerTransformDocumentParams creates a new DTAServerTransformDocumentParams object
// with the default values initialized.
func NewDTAServerTransformDocumentParams() *DTAServerTransformDocumentParams {
	var ()
	return &DTAServerTransformDocumentParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDTAServerTransformDocumentParamsWithTimeout creates a new DTAServerTransformDocumentParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDTAServerTransformDocumentParamsWithTimeout(timeout time.Duration) *DTAServerTransformDocumentParams {
	var ()
	return &DTAServerTransformDocumentParams{

		timeout: timeout,
	}
}

// NewDTAServerTransformDocumentParamsWithContext creates a new DTAServerTransformDocumentParams object
// with the default values initialized, and the ability to set a context for a request
func NewDTAServerTransformDocumentParamsWithContext(ctx context.Context) *DTAServerTransformDocumentParams {
	var ()
	return &DTAServerTransformDocumentParams{

		Context: ctx,
	}
}

// NewDTAServerTransformDocumentParamsWithHTTPClient creates a new DTAServerTransformDocumentParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDTAServerTransformDocumentParamsWithHTTPClient(client *http.Client) *DTAServerTransformDocumentParams {
	var ()
	return &DTAServerTransformDocumentParams{
		HTTPClient: client,
	}
}

/*DTAServerTransformDocumentParams contains all the parameters to send to the API endpoint
for the d t a server transform document operation typically these are written to a http.Request
*/
type DTAServerTransformDocumentParams struct {

	/*Body*/
	Body *rest_models.DtaserviceDocumentRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) WithTimeout(timeout time.Duration) *DTAServerTransformDocumentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) WithContext(ctx context.Context) *DTAServerTransformDocumentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) WithHTTPClient(client *http.Client) *DTAServerTransformDocumentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) WithBody(body *rest_models.DtaserviceDocumentRequest) *DTAServerTransformDocumentParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the d t a server transform document params
func (o *DTAServerTransformDocumentParams) SetBody(body *rest_models.DtaserviceDocumentRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *DTAServerTransformDocumentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
