// Code generated by go-swagger; DO NOT EDIT.

package rest_models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DtaserviceOptionsResponse dtaservice options response
//
// swagger:model dtaserviceOptionsResponse
type DtaserviceOptionsResponse struct {

	// services
	Services string `json:"services,omitempty"`
}

// Validate validates this dtaservice options response
func (m *DtaserviceOptionsResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DtaserviceOptionsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DtaserviceOptionsResponse) UnmarshalBinary(b []byte) error {
	var res DtaserviceOptionsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}