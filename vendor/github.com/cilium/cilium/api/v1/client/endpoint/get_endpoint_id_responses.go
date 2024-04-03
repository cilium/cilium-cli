// Code generated by go-swagger; DO NOT EDIT.

// Copyright Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package endpoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/cilium/cilium/api/v1/models"
)

// GetEndpointIDReader is a Reader for the GetEndpointID structure.
type GetEndpointIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetEndpointIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetEndpointIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetEndpointIDInvalid()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetEndpointIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 429:
		result := NewGetEndpointIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /endpoint/{id}] GetEndpointID", response, response.Code())
	}
}

// NewGetEndpointIDOK creates a GetEndpointIDOK with default headers values
func NewGetEndpointIDOK() *GetEndpointIDOK {
	return &GetEndpointIDOK{}
}

/*
GetEndpointIDOK describes a response with status code 200, with default header values.

Success
*/
type GetEndpointIDOK struct {
	Payload *models.Endpoint
}

// IsSuccess returns true when this get endpoint Id o k response has a 2xx status code
func (o *GetEndpointIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get endpoint Id o k response has a 3xx status code
func (o *GetEndpointIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get endpoint Id o k response has a 4xx status code
func (o *GetEndpointIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get endpoint Id o k response has a 5xx status code
func (o *GetEndpointIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get endpoint Id o k response a status code equal to that given
func (o *GetEndpointIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get endpoint Id o k response
func (o *GetEndpointIDOK) Code() int {
	return 200
}

func (o *GetEndpointIDOK) Error() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdOK  %+v", 200, o.Payload)
}

func (o *GetEndpointIDOK) String() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdOK  %+v", 200, o.Payload)
}

func (o *GetEndpointIDOK) GetPayload() *models.Endpoint {
	return o.Payload
}

func (o *GetEndpointIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Endpoint)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetEndpointIDInvalid creates a GetEndpointIDInvalid with default headers values
func NewGetEndpointIDInvalid() *GetEndpointIDInvalid {
	return &GetEndpointIDInvalid{}
}

/*
GetEndpointIDInvalid describes a response with status code 400, with default header values.

Invalid endpoint ID format for specified type
*/
type GetEndpointIDInvalid struct {
	Payload models.Error
}

// IsSuccess returns true when this get endpoint Id invalid response has a 2xx status code
func (o *GetEndpointIDInvalid) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get endpoint Id invalid response has a 3xx status code
func (o *GetEndpointIDInvalid) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get endpoint Id invalid response has a 4xx status code
func (o *GetEndpointIDInvalid) IsClientError() bool {
	return true
}

// IsServerError returns true when this get endpoint Id invalid response has a 5xx status code
func (o *GetEndpointIDInvalid) IsServerError() bool {
	return false
}

// IsCode returns true when this get endpoint Id invalid response a status code equal to that given
func (o *GetEndpointIDInvalid) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the get endpoint Id invalid response
func (o *GetEndpointIDInvalid) Code() int {
	return 400
}

func (o *GetEndpointIDInvalid) Error() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdInvalid  %+v", 400, o.Payload)
}

func (o *GetEndpointIDInvalid) String() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdInvalid  %+v", 400, o.Payload)
}

func (o *GetEndpointIDInvalid) GetPayload() models.Error {
	return o.Payload
}

func (o *GetEndpointIDInvalid) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetEndpointIDNotFound creates a GetEndpointIDNotFound with default headers values
func NewGetEndpointIDNotFound() *GetEndpointIDNotFound {
	return &GetEndpointIDNotFound{}
}

/*
GetEndpointIDNotFound describes a response with status code 404, with default header values.

Endpoint not found
*/
type GetEndpointIDNotFound struct {
}

// IsSuccess returns true when this get endpoint Id not found response has a 2xx status code
func (o *GetEndpointIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get endpoint Id not found response has a 3xx status code
func (o *GetEndpointIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get endpoint Id not found response has a 4xx status code
func (o *GetEndpointIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get endpoint Id not found response has a 5xx status code
func (o *GetEndpointIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get endpoint Id not found response a status code equal to that given
func (o *GetEndpointIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get endpoint Id not found response
func (o *GetEndpointIDNotFound) Code() int {
	return 404
}

func (o *GetEndpointIDNotFound) Error() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdNotFound ", 404)
}

func (o *GetEndpointIDNotFound) String() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdNotFound ", 404)
}

func (o *GetEndpointIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetEndpointIDTooManyRequests creates a GetEndpointIDTooManyRequests with default headers values
func NewGetEndpointIDTooManyRequests() *GetEndpointIDTooManyRequests {
	return &GetEndpointIDTooManyRequests{}
}

/*
GetEndpointIDTooManyRequests describes a response with status code 429, with default header values.

Rate-limiting too many requests in the given time frame
*/
type GetEndpointIDTooManyRequests struct {
}

// IsSuccess returns true when this get endpoint Id too many requests response has a 2xx status code
func (o *GetEndpointIDTooManyRequests) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get endpoint Id too many requests response has a 3xx status code
func (o *GetEndpointIDTooManyRequests) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get endpoint Id too many requests response has a 4xx status code
func (o *GetEndpointIDTooManyRequests) IsClientError() bool {
	return true
}

// IsServerError returns true when this get endpoint Id too many requests response has a 5xx status code
func (o *GetEndpointIDTooManyRequests) IsServerError() bool {
	return false
}

// IsCode returns true when this get endpoint Id too many requests response a status code equal to that given
func (o *GetEndpointIDTooManyRequests) IsCode(code int) bool {
	return code == 429
}

// Code gets the status code for the get endpoint Id too many requests response
func (o *GetEndpointIDTooManyRequests) Code() int {
	return 429
}

func (o *GetEndpointIDTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdTooManyRequests ", 429)
}

func (o *GetEndpointIDTooManyRequests) String() string {
	return fmt.Sprintf("[GET /endpoint/{id}][%d] getEndpointIdTooManyRequests ", 429)
}

func (o *GetEndpointIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
