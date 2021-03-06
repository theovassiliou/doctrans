// Code generated by go-swagger; DO NOT EDIT.

package d_t_a_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/theovassiliou/doctrans/gen/rest_models"
)

// DTAServerTransformDocumentReader is a Reader for the DTAServerTransformDocument structure.
type DTAServerTransformDocumentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DTAServerTransformDocumentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDTAServerTransformDocumentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDTAServerTransformDocumentDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDTAServerTransformDocumentOK creates a DTAServerTransformDocumentOK with default headers values
func NewDTAServerTransformDocumentOK() *DTAServerTransformDocumentOK {
	return &DTAServerTransformDocumentOK{}
}

/*DTAServerTransformDocumentOK handles this case with default header values.

A successful response.
*/
type DTAServerTransformDocumentOK struct {
	Payload *rest_models.DtaserviceTransformDocumentResponse
}

func (o *DTAServerTransformDocumentOK) Error() string {
	return fmt.Sprintf("[POST /v1/document/transform][%d] dTAServerTransformDocumentOK  %+v", 200, o.Payload)
}

func (o *DTAServerTransformDocumentOK) GetPayload() *rest_models.DtaserviceTransformDocumentResponse {
	return o.Payload
}

func (o *DTAServerTransformDocumentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_models.DtaserviceTransformDocumentResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDTAServerTransformDocumentDefault creates a DTAServerTransformDocumentDefault with default headers values
func NewDTAServerTransformDocumentDefault(code int) *DTAServerTransformDocumentDefault {
	return &DTAServerTransformDocumentDefault{
		_statusCode: code,
	}
}

/*DTAServerTransformDocumentDefault handles this case with default header values.

An unexpected error response
*/
type DTAServerTransformDocumentDefault struct {
	_statusCode int

	Payload *rest_models.RuntimeError
}

// Code gets the status code for the d t a server transform document default response
func (o *DTAServerTransformDocumentDefault) Code() int {
	return o._statusCode
}

func (o *DTAServerTransformDocumentDefault) Error() string {
	return fmt.Sprintf("[POST /v1/document/transform][%d] DTAServer_TransformDocument default  %+v", o._statusCode, o.Payload)
}

func (o *DTAServerTransformDocumentDefault) GetPayload() *rest_models.RuntimeError {
	return o.Payload
}

func (o *DTAServerTransformDocumentDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_models.RuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
