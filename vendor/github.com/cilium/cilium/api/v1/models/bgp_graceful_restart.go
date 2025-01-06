// Code generated by go-swagger; DO NOT EDIT.

// Copyright Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// BgpGracefulRestart BGP graceful restart parameters negotiated with the peer.
//
// swagger:model BgpGracefulRestart
type BgpGracefulRestart struct {

	// When set, graceful restart capability is negotiated for all AFI/SAFIs of
	// this peer.
	Enabled bool `json:"enabled,omitempty"`

	// This is the time advertised to peer for the BGP session to be re-established
	// after a restart. After this period, peer will remove stale routes.
	// (RFC 4724 section 4.2)
	RestartTimeSeconds int64 `json:"restart-time-seconds,omitempty"`
}

// Validate validates this bgp graceful restart
func (m *BgpGracefulRestart) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this bgp graceful restart based on context it is used
func (m *BgpGracefulRestart) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *BgpGracefulRestart) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BgpGracefulRestart) UnmarshalBinary(b []byte) error {
	var res BgpGracefulRestart
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
