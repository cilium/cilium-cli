// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/extensions/filters/http/basic_auth/v3/basic_auth.proto

package basic_authv3

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

// Validate checks the field values on BasicAuth with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *BasicAuth) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on BasicAuth with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in BasicAuthMultiError, or nil
// if none found.
func (m *BasicAuth) ValidateAll() error {
	return m.validate(true)
}

func (m *BasicAuth) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetUsers()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BasicAuthValidationError{
					field:  "Users",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BasicAuthValidationError{
					field:  "Users",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUsers()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BasicAuthValidationError{
				field:  "Users",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if !_BasicAuth_ForwardUsernameHeader_Pattern.MatchString(m.GetForwardUsernameHeader()) {
		err := BasicAuthValidationError{
			field:  "ForwardUsernameHeader",
			reason: "value does not match regex pattern \"^[^\\x00\\n\\r]*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return BasicAuthMultiError(errors)
	}

	return nil
}

// BasicAuthMultiError is an error wrapping multiple validation errors returned
// by BasicAuth.ValidateAll() if the designated constraints aren't met.
type BasicAuthMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BasicAuthMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BasicAuthMultiError) AllErrors() []error { return m }

// BasicAuthValidationError is the validation error returned by
// BasicAuth.Validate if the designated constraints aren't met.
type BasicAuthValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BasicAuthValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BasicAuthValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BasicAuthValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BasicAuthValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BasicAuthValidationError) ErrorName() string { return "BasicAuthValidationError" }

// Error satisfies the builtin error interface
func (e BasicAuthValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBasicAuth.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BasicAuthValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BasicAuthValidationError{}

var _BasicAuth_ForwardUsernameHeader_Pattern = regexp.MustCompile("^[^\x00\n\r]*$")

// Validate checks the field values on BasicAuthPerRoute with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *BasicAuthPerRoute) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on BasicAuthPerRoute with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// BasicAuthPerRouteMultiError, or nil if none found.
func (m *BasicAuthPerRoute) ValidateAll() error {
	return m.validate(true)
}

func (m *BasicAuthPerRoute) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUsers() == nil {
		err := BasicAuthPerRouteValidationError{
			field:  "Users",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetUsers()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BasicAuthPerRouteValidationError{
					field:  "Users",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BasicAuthPerRouteValidationError{
					field:  "Users",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUsers()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BasicAuthPerRouteValidationError{
				field:  "Users",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return BasicAuthPerRouteMultiError(errors)
	}

	return nil
}

// BasicAuthPerRouteMultiError is an error wrapping multiple validation errors
// returned by BasicAuthPerRoute.ValidateAll() if the designated constraints
// aren't met.
type BasicAuthPerRouteMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BasicAuthPerRouteMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BasicAuthPerRouteMultiError) AllErrors() []error { return m }

// BasicAuthPerRouteValidationError is the validation error returned by
// BasicAuthPerRoute.Validate if the designated constraints aren't met.
type BasicAuthPerRouteValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BasicAuthPerRouteValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BasicAuthPerRouteValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BasicAuthPerRouteValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BasicAuthPerRouteValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BasicAuthPerRouteValidationError) ErrorName() string {
	return "BasicAuthPerRouteValidationError"
}

// Error satisfies the builtin error interface
func (e BasicAuthPerRouteValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBasicAuthPerRoute.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BasicAuthPerRouteValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BasicAuthPerRouteValidationError{}
