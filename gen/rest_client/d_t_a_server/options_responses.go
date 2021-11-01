// Code generated by go-swagger; DO NOT EDIT.

package d_t_a_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/theovassiliou/doctrans/gen/rest_models"
)

// OptionsReader is a Reader for the Options structure.
type OptionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *OptionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewOptionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewOptionsOK creates a OptionsOK with default headers values
func NewOptionsOK() *OptionsOK {
	return &OptionsOK{}
}

/*OptionsOK handles this case with default header values.

A successful response.
*/
type OptionsOK struct {
	Payload *rest_models.DtaserviceOptionsResponse
}

func (o *OptionsOK) Error() string {
	return fmt.Sprintf("[GET /v1/service/options][%d] optionsOK  %+v", 200, o.Payload)
}

func (o *OptionsOK) GetPayload() *rest_models.DtaserviceOptionsResponse {
	return o.Payload
}

func (o *OptionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_models.DtaserviceOptionsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
