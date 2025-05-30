// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Reception reception
//
// swagger:model Reception
type Reception struct {

	// date time
	// Required: true
	// Format: date-time
	DateTime *strfmt.DateTime `json:"dateTime"`

	// id
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// pvz Id
	// Required: true
	// Format: uuid
	PvzID *strfmt.UUID `json:"pvzId"`

	// status
	// Required: true
	// Enum: ["in_progress","close"]
	Status *string `json:"status"`
}

// Validate validates this reception
func (m *Reception) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDateTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePvzID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Reception) validateDateTime(formats strfmt.Registry) error {

	if err := validate.Required("dateTime", "body", m.DateTime); err != nil {
		return err
	}

	if err := validate.FormatOf("dateTime", "body", "date-time", m.DateTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Reception) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Reception) validatePvzID(formats strfmt.Registry) error {

	if err := validate.Required("pvzId", "body", m.PvzID); err != nil {
		return err
	}

	if err := validate.FormatOf("pvzId", "body", "uuid", m.PvzID.String(), formats); err != nil {
		return err
	}

	return nil
}

var receptionTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["in_progress","close"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		receptionTypeStatusPropEnum = append(receptionTypeStatusPropEnum, v)
	}
}

const (

	// ReceptionStatusInProgress captures enum value "in_progress"
	ReceptionStatusInProgress string = "in_progress"

	// ReceptionStatusClose captures enum value "close"
	ReceptionStatusClose string = "close"
)

// prop value enum
func (m *Reception) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, receptionTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Reception) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", *m.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this reception based on context it is used
func (m *Reception) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Reception) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Reception) UnmarshalBinary(b []byte) error {
	var res Reception
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
