// Code generated by go-swagger; DO NOT EDIT.

package d_t_a_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/theovassiliou/doctrans/gen/rest_models"
)

// TransformPipeOKCode is the HTTP code returned for type TransformPipeOK
const TransformPipeOKCode int = 200

/*
TransformPipeOK A successful response.

swagger:response transformPipeOK
*/
type TransformPipeOK struct {

	/*
	  In: Body
	*/
	Payload *rest_models.DtaserviceTransformPipeResponse `json:"body,omitempty"`
}

// NewTransformPipeOK creates TransformPipeOK with default headers values
func NewTransformPipeOK() *TransformPipeOK {

	return &TransformPipeOK{}
}

// WithPayload adds the payload to the transform pipe o k response
func (o *TransformPipeOK) WithPayload(payload *rest_models.DtaserviceTransformPipeResponse) *TransformPipeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the transform pipe o k response
func (o *TransformPipeOK) SetPayload(payload *rest_models.DtaserviceTransformPipeResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *TransformPipeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
