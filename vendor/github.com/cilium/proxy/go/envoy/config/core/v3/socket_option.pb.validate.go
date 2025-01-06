// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/core/v3/socket_option.proto

package corev3

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on SocketOption with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SocketOption) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SocketOption with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SocketOptionMultiError, or
// nil if none found.
func (m *SocketOption) ValidateAll() error {
	return m.validate(true)
}

func (m *SocketOption) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Description

	// no validation rules for Level

	// no validation rules for Name

	if _, ok := SocketOption_SocketState_name[int32(m.GetState())]; !ok {
		err := SocketOptionValidationError{
			field:  "State",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetType()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SocketOptionValidationError{
					field:  "Type",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SocketOptionValidationError{
					field:  "Type",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetType()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SocketOptionValidationError{
				field:  "Type",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	oneofValuePresent := false
	switch v := m.Value.(type) {
	case *SocketOption_IntValue:
		if v == nil {
			err := SocketOptionValidationError{
				field:  "Value",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
		oneofValuePresent = true
		// no validation rules for IntValue
	case *SocketOption_BufValue:
		if v == nil {
			err := SocketOptionValidationError{
				field:  "Value",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
		oneofValuePresent = true
		// no validation rules for BufValue
	default:
		_ = v // ensures v is used
	}
	if !oneofValuePresent {
		err := SocketOptionValidationError{
			field:  "Value",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return SocketOptionMultiError(errors)
	}

	return nil
}

// SocketOptionMultiError is an error wrapping multiple validation errors
// returned by SocketOption.ValidateAll() if the designated constraints aren't met.
type SocketOptionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SocketOptionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SocketOptionMultiError) AllErrors() []error { return m }

// SocketOptionValidationError is the validation error returned by
// SocketOption.Validate if the designated constraints aren't met.
type SocketOptionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SocketOptionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SocketOptionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SocketOptionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SocketOptionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SocketOptionValidationError) ErrorName() string { return "SocketOptionValidationError" }

// Error satisfies the builtin error interface
func (e SocketOptionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSocketOption.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SocketOptionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SocketOptionValidationError{}

// Validate checks the field values on SocketOptionsOverride with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SocketOptionsOverride) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SocketOptionsOverride with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SocketOptionsOverrideMultiError, or nil if none found.
func (m *SocketOptionsOverride) ValidateAll() error {
	return m.validate(true)
}

func (m *SocketOptionsOverride) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetSocketOptions() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, SocketOptionsOverrideValidationError{
						field:  fmt.Sprintf("SocketOptions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, SocketOptionsOverrideValidationError{
						field:  fmt.Sprintf("SocketOptions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SocketOptionsOverrideValidationError{
					field:  fmt.Sprintf("SocketOptions[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return SocketOptionsOverrideMultiError(errors)
	}

	return nil
}

// SocketOptionsOverrideMultiError is an error wrapping multiple validation
// errors returned by SocketOptionsOverride.ValidateAll() if the designated
// constraints aren't met.
type SocketOptionsOverrideMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SocketOptionsOverrideMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SocketOptionsOverrideMultiError) AllErrors() []error { return m }

// SocketOptionsOverrideValidationError is the validation error returned by
// SocketOptionsOverride.Validate if the designated constraints aren't met.
type SocketOptionsOverrideValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SocketOptionsOverrideValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SocketOptionsOverrideValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SocketOptionsOverrideValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SocketOptionsOverrideValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SocketOptionsOverrideValidationError) ErrorName() string {
	return "SocketOptionsOverrideValidationError"
}

// Error satisfies the builtin error interface
func (e SocketOptionsOverrideValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSocketOptionsOverride.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SocketOptionsOverrideValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SocketOptionsOverrideValidationError{}

// Validate checks the field values on SocketOption_SocketType with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SocketOption_SocketType) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SocketOption_SocketType with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SocketOption_SocketTypeMultiError, or nil if none found.
func (m *SocketOption_SocketType) ValidateAll() error {
	return m.validate(true)
}

func (m *SocketOption_SocketType) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetStream()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SocketOption_SocketTypeValidationError{
					field:  "Stream",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SocketOption_SocketTypeValidationError{
					field:  "Stream",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetStream()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SocketOption_SocketTypeValidationError{
				field:  "Stream",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetDatagram()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SocketOption_SocketTypeValidationError{
					field:  "Datagram",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SocketOption_SocketTypeValidationError{
					field:  "Datagram",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDatagram()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SocketOption_SocketTypeValidationError{
				field:  "Datagram",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SocketOption_SocketTypeMultiError(errors)
	}

	return nil
}

// SocketOption_SocketTypeMultiError is an error wrapping multiple validation
// errors returned by SocketOption_SocketType.ValidateAll() if the designated
// constraints aren't met.
type SocketOption_SocketTypeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SocketOption_SocketTypeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SocketOption_SocketTypeMultiError) AllErrors() []error { return m }

// SocketOption_SocketTypeValidationError is the validation error returned by
// SocketOption_SocketType.Validate if the designated constraints aren't met.
type SocketOption_SocketTypeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SocketOption_SocketTypeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SocketOption_SocketTypeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SocketOption_SocketTypeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SocketOption_SocketTypeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SocketOption_SocketTypeValidationError) ErrorName() string {
	return "SocketOption_SocketTypeValidationError"
}

// Error satisfies the builtin error interface
func (e SocketOption_SocketTypeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSocketOption_SocketType.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SocketOption_SocketTypeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SocketOption_SocketTypeValidationError{}

// Validate checks the field values on SocketOption_SocketType_Stream with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SocketOption_SocketType_Stream) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SocketOption_SocketType_Stream with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// SocketOption_SocketType_StreamMultiError, or nil if none found.
func (m *SocketOption_SocketType_Stream) ValidateAll() error {
	return m.validate(true)
}

func (m *SocketOption_SocketType_Stream) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return SocketOption_SocketType_StreamMultiError(errors)
	}

	return nil
}

// SocketOption_SocketType_StreamMultiError is an error wrapping multiple
// validation errors returned by SocketOption_SocketType_Stream.ValidateAll()
// if the designated constraints aren't met.
type SocketOption_SocketType_StreamMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SocketOption_SocketType_StreamMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SocketOption_SocketType_StreamMultiError) AllErrors() []error { return m }

// SocketOption_SocketType_StreamValidationError is the validation error
// returned by SocketOption_SocketType_Stream.Validate if the designated
// constraints aren't met.
type SocketOption_SocketType_StreamValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SocketOption_SocketType_StreamValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SocketOption_SocketType_StreamValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SocketOption_SocketType_StreamValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SocketOption_SocketType_StreamValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SocketOption_SocketType_StreamValidationError) ErrorName() string {
	return "SocketOption_SocketType_StreamValidationError"
}

// Error satisfies the builtin error interface
func (e SocketOption_SocketType_StreamValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSocketOption_SocketType_Stream.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SocketOption_SocketType_StreamValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SocketOption_SocketType_StreamValidationError{}

// Validate checks the field values on SocketOption_SocketType_Datagram with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *SocketOption_SocketType_Datagram) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SocketOption_SocketType_Datagram with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// SocketOption_SocketType_DatagramMultiError, or nil if none found.
func (m *SocketOption_SocketType_Datagram) ValidateAll() error {
	return m.validate(true)
}

func (m *SocketOption_SocketType_Datagram) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return SocketOption_SocketType_DatagramMultiError(errors)
	}

	return nil
}

// SocketOption_SocketType_DatagramMultiError is an error wrapping multiple
// validation errors returned by
// SocketOption_SocketType_Datagram.ValidateAll() if the designated
// constraints aren't met.
type SocketOption_SocketType_DatagramMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SocketOption_SocketType_DatagramMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SocketOption_SocketType_DatagramMultiError) AllErrors() []error { return m }

// SocketOption_SocketType_DatagramValidationError is the validation error
// returned by SocketOption_SocketType_Datagram.Validate if the designated
// constraints aren't met.
type SocketOption_SocketType_DatagramValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SocketOption_SocketType_DatagramValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SocketOption_SocketType_DatagramValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SocketOption_SocketType_DatagramValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SocketOption_SocketType_DatagramValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SocketOption_SocketType_DatagramValidationError) ErrorName() string {
	return "SocketOption_SocketType_DatagramValidationError"
}

// Error satisfies the builtin error interface
func (e SocketOption_SocketType_DatagramValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSocketOption_SocketType_Datagram.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SocketOption_SocketType_DatagramValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SocketOption_SocketType_DatagramValidationError{}
