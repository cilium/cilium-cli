// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/extensions/access_loggers/fluentd/v3/fluentd.proto

package fluentdv3

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

// Validate checks the field values on FluentdAccessLogConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *FluentdAccessLogConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FluentdAccessLogConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FluentdAccessLogConfigMultiError, or nil if none found.
func (m *FluentdAccessLogConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *FluentdAccessLogConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetCluster()) < 1 {
		err := FluentdAccessLogConfigValidationError{
			field:  "Cluster",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetTag()) < 1 {
		err := FluentdAccessLogConfigValidationError{
			field:  "Tag",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetStatPrefix()) < 1 {
		err := FluentdAccessLogConfigValidationError{
			field:  "StatPrefix",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if d := m.GetBufferFlushInterval(); d != nil {
		dur, err := d.AsDuration(), d.CheckValid()
		if err != nil {
			err = FluentdAccessLogConfigValidationError{
				field:  "BufferFlushInterval",
				reason: "value is not a valid duration",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {

			gt := time.Duration(0*time.Second + 0*time.Nanosecond)

			if dur <= gt {
				err := FluentdAccessLogConfigValidationError{
					field:  "BufferFlushInterval",
					reason: "value must be greater than 0s",
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}
	}

	if all {
		switch v := interface{}(m.GetBufferSizeBytes()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FluentdAccessLogConfigValidationError{
					field:  "BufferSizeBytes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FluentdAccessLogConfigValidationError{
					field:  "BufferSizeBytes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetBufferSizeBytes()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FluentdAccessLogConfigValidationError{
				field:  "BufferSizeBytes",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.GetRecord() == nil {
		err := FluentdAccessLogConfigValidationError{
			field:  "Record",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetRecord()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FluentdAccessLogConfigValidationError{
					field:  "Record",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FluentdAccessLogConfigValidationError{
					field:  "Record",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRecord()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FluentdAccessLogConfigValidationError{
				field:  "Record",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetRetryOptions()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FluentdAccessLogConfigValidationError{
					field:  "RetryOptions",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FluentdAccessLogConfigValidationError{
					field:  "RetryOptions",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRetryOptions()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FluentdAccessLogConfigValidationError{
				field:  "RetryOptions",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetFormatters() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, FluentdAccessLogConfigValidationError{
						field:  fmt.Sprintf("Formatters[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, FluentdAccessLogConfigValidationError{
						field:  fmt.Sprintf("Formatters[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FluentdAccessLogConfigValidationError{
					field:  fmt.Sprintf("Formatters[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return FluentdAccessLogConfigMultiError(errors)
	}

	return nil
}

// FluentdAccessLogConfigMultiError is an error wrapping multiple validation
// errors returned by FluentdAccessLogConfig.ValidateAll() if the designated
// constraints aren't met.
type FluentdAccessLogConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FluentdAccessLogConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FluentdAccessLogConfigMultiError) AllErrors() []error { return m }

// FluentdAccessLogConfigValidationError is the validation error returned by
// FluentdAccessLogConfig.Validate if the designated constraints aren't met.
type FluentdAccessLogConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FluentdAccessLogConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FluentdAccessLogConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FluentdAccessLogConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FluentdAccessLogConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FluentdAccessLogConfigValidationError) ErrorName() string {
	return "FluentdAccessLogConfigValidationError"
}

// Error satisfies the builtin error interface
func (e FluentdAccessLogConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFluentdAccessLogConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FluentdAccessLogConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FluentdAccessLogConfigValidationError{}

// Validate checks the field values on FluentdAccessLogConfig_RetryOptions with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *FluentdAccessLogConfig_RetryOptions) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FluentdAccessLogConfig_RetryOptions
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// FluentdAccessLogConfig_RetryOptionsMultiError, or nil if none found.
func (m *FluentdAccessLogConfig_RetryOptions) ValidateAll() error {
	return m.validate(true)
}

func (m *FluentdAccessLogConfig_RetryOptions) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetMaxConnectAttempts()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FluentdAccessLogConfig_RetryOptionsValidationError{
					field:  "MaxConnectAttempts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FluentdAccessLogConfig_RetryOptionsValidationError{
					field:  "MaxConnectAttempts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMaxConnectAttempts()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FluentdAccessLogConfig_RetryOptionsValidationError{
				field:  "MaxConnectAttempts",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetBackoffOptions()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FluentdAccessLogConfig_RetryOptionsValidationError{
					field:  "BackoffOptions",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FluentdAccessLogConfig_RetryOptionsValidationError{
					field:  "BackoffOptions",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetBackoffOptions()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FluentdAccessLogConfig_RetryOptionsValidationError{
				field:  "BackoffOptions",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return FluentdAccessLogConfig_RetryOptionsMultiError(errors)
	}

	return nil
}

// FluentdAccessLogConfig_RetryOptionsMultiError is an error wrapping multiple
// validation errors returned by
// FluentdAccessLogConfig_RetryOptions.ValidateAll() if the designated
// constraints aren't met.
type FluentdAccessLogConfig_RetryOptionsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FluentdAccessLogConfig_RetryOptionsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FluentdAccessLogConfig_RetryOptionsMultiError) AllErrors() []error { return m }

// FluentdAccessLogConfig_RetryOptionsValidationError is the validation error
// returned by FluentdAccessLogConfig_RetryOptions.Validate if the designated
// constraints aren't met.
type FluentdAccessLogConfig_RetryOptionsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FluentdAccessLogConfig_RetryOptionsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FluentdAccessLogConfig_RetryOptionsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FluentdAccessLogConfig_RetryOptionsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FluentdAccessLogConfig_RetryOptionsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FluentdAccessLogConfig_RetryOptionsValidationError) ErrorName() string {
	return "FluentdAccessLogConfig_RetryOptionsValidationError"
}

// Error satisfies the builtin error interface
func (e FluentdAccessLogConfig_RetryOptionsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFluentdAccessLogConfig_RetryOptions.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FluentdAccessLogConfig_RetryOptionsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FluentdAccessLogConfig_RetryOptionsValidationError{}
