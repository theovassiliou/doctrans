// Code generated by go-swagger; DO NOT EDIT.

package d_t_a_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/theovassiliou/doctrans/gen/rest_models"
)

// TransformDocumentOKCode is the HTTP code returned for type TransformDocumentOK
const TransformDocumentOKCode int = 200

/*TransformDocumentOK A successful response.

swagger:response transformDocumentOK
*/
type TransformDocumentOK struct {

	/*
	  In: Body
	*/
	Payload *rest_models.DtaserviceTransformDocumentResponse `json:"body,omitempty"`
}

// NewTransformDocumentOK creates TransformDocumentOK with default headers values
func NewTransformDocumentOK() *TransformDocumentOK {

	return &TransformDocumentOK{}
}

// WithPayload adds the payload to the transform document o k response
func (o *TransformDocumentOK) WithPayload(payload *rest_models.DtaserviceTransformDocumentResponse) *TransformDocumentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the transform document o k response
func (o *TransformDocumentOK) SetPayload(payload *rest_models.DtaserviceTransformDocumentResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *TransformDocumentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
