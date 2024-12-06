// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/extensions/tracers/opentelemetry/samplers/v3/always_on_sampler.proto

package samplersv3

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

// Validate checks the field values on AlwaysOnSamplerConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *AlwaysOnSamplerConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AlwaysOnSamplerConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AlwaysOnSamplerConfigMultiError, or nil if none found.
func (m *AlwaysOnSamplerConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *AlwaysOnSamplerConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return AlwaysOnSamplerConfigMultiError(errors)
	}

	return nil
}

// AlwaysOnSamplerConfigMultiError is an error wrapping multiple validation
// errors returned by AlwaysOnSamplerConfig.ValidateAll() if the designated
// constraints aren't met.
type AlwaysOnSamplerConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AlwaysOnSamplerConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AlwaysOnSamplerConfigMultiError) AllErrors() []error { return m }

// AlwaysOnSamplerConfigValidationError is the validation error returned by
// AlwaysOnSamplerConfig.Validate if the designated constraints aren't met.
type AlwaysOnSamplerConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AlwaysOnSamplerConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AlwaysOnSamplerConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AlwaysOnSamplerConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AlwaysOnSamplerConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AlwaysOnSamplerConfigValidationError) ErrorName() string {
	return "AlwaysOnSamplerConfigValidationError"
}

// Error satisfies the builtin error interface
func (e AlwaysOnSamplerConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAlwaysOnSamplerConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AlwaysOnSamplerConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AlwaysOnSamplerConfigValidationError{}
