// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netlify/open-api/v2/go/models"
)

// DeleteEnvVarValueReader is a Reader for the DeleteEnvVarValue structure.
type DeleteEnvVarValueReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteEnvVarValueReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteEnvVarValueNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDeleteEnvVarValueDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteEnvVarValueNoContent creates a DeleteEnvVarValueNoContent with default headers values
func NewDeleteEnvVarValueNoContent() *DeleteEnvVarValueNoContent {
	return &DeleteEnvVarValueNoContent{}
}

/*DeleteEnvVarValueNoContent handles this case with default header values.

No Content (success)
*/
type DeleteEnvVarValueNoContent struct {
}

func (o *DeleteEnvVarValueNoContent) Error() string {
	return fmt.Sprintf("[DELETE /accounts/{account_id}/env/{key}/value/{id}][%d] deleteEnvVarValueNoContent ", 204)
}

func (o *DeleteEnvVarValueNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteEnvVarValueDefault creates a DeleteEnvVarValueDefault with default headers values
func NewDeleteEnvVarValueDefault(code int) *DeleteEnvVarValueDefault {
	return &DeleteEnvVarValueDefault{
		_statusCode: code,
	}
}

/*DeleteEnvVarValueDefault handles this case with default header values.

error
*/
type DeleteEnvVarValueDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the delete env var value default response
func (o *DeleteEnvVarValueDefault) Code() int {
	return o._statusCode
}

func (o *DeleteEnvVarValueDefault) Error() string {
	return fmt.Sprintf("[DELETE /accounts/{account_id}/env/{key}/value/{id}][%d] deleteEnvVarValue default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteEnvVarValueDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteEnvVarValueDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
