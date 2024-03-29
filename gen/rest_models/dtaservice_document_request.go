// Code generated by go-swagger; DO NOT EDIT.

package rest_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DtaserviceDocumentRequest The request message containing the document to be transformed
//
// swagger:model dtaserviceDocumentRequest
type DtaserviceDocumentRequest struct {

	// document
	// Format: byte
	Document strfmt.Base64 `json:"document,omitempty"`

	// file name
	FileName string `json:"file_name,omitempty"`

	// options
	Options interface{} `json:"options,omitempty"`

	// service name
	ServiceName string `json:"service_name,omitempty"`
}

// Validate validates this dtaservice document request
func (m *DtaserviceDocumentRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this dtaservice document request based on context it is used
func (m *DtaserviceDocumentRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DtaserviceDocumentRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DtaserviceDocumentRequest) UnmarshalBinary(b []byte) error {
	var res DtaserviceDocumentRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
