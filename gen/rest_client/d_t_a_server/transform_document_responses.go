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

// TransformDocumentReader is a Reader for the TransformDocument structure.
type TransformDocumentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *TransformDocumentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewTransformDocumentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[POST /v1/document/transform] TransformDocument", response, response.Code())
	}
}

// NewTransformDocumentOK creates a TransformDocumentOK with default headers values
func NewTransformDocumentOK() *TransformDocumentOK {
	return &TransformDocumentOK{}
}

/*
TransformDocumentOK describes a response with status code 200, with default header values.

A successful response.
*/
type TransformDocumentOK struct {
	Payload *rest_models.DtaserviceTransformDocumentResponse
}

// IsSuccess returns true when this transform document o k response has a 2xx status code
func (o *TransformDocumentOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this transform document o k response has a 3xx status code
func (o *TransformDocumentOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this transform document o k response has a 4xx status code
func (o *TransformDocumentOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this transform document o k response has a 5xx status code
func (o *TransformDocumentOK) IsServerError() bool {
	return false
}

// IsCode returns true when this transform document o k response a status code equal to that given
func (o *TransformDocumentOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the transform document o k response
func (o *TransformDocumentOK) Code() int {
	return 200
}

func (o *TransformDocumentOK) Error() string {
	return fmt.Sprintf("[POST /v1/document/transform][%d] transformDocumentOK  %+v", 200, o.Payload)
}

func (o *TransformDocumentOK) String() string {
	return fmt.Sprintf("[POST /v1/document/transform][%d] transformDocumentOK  %+v", 200, o.Payload)
}

func (o *TransformDocumentOK) GetPayload() *rest_models.DtaserviceTransformDocumentResponse {
	return o.Payload
}

func (o *TransformDocumentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_models.DtaserviceTransformDocumentResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
