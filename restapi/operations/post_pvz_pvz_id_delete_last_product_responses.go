// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/totorialman/go-task-avito/models"
)

// PostPvzPvzIDDeleteLastProductOKCode is the HTTP code returned for type PostPvzPvzIDDeleteLastProductOK
const PostPvzPvzIDDeleteLastProductOKCode int = 200

/*
PostPvzPvzIDDeleteLastProductOK Товар удален

swagger:response postPvzPvzIdDeleteLastProductOK
*/
type PostPvzPvzIDDeleteLastProductOK struct {
}

// NewPostPvzPvzIDDeleteLastProductOK creates PostPvzPvzIDDeleteLastProductOK with default headers values
func NewPostPvzPvzIDDeleteLastProductOK() *PostPvzPvzIDDeleteLastProductOK {

	return &PostPvzPvzIDDeleteLastProductOK{}
}

// WriteResponse to the client
func (o *PostPvzPvzIDDeleteLastProductOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// PostPvzPvzIDDeleteLastProductBadRequestCode is the HTTP code returned for type PostPvzPvzIDDeleteLastProductBadRequest
const PostPvzPvzIDDeleteLastProductBadRequestCode int = 400

/*
PostPvzPvzIDDeleteLastProductBadRequest Неверный запрос, нет активной приемки или нет товаров для удаления

swagger:response postPvzPvzIdDeleteLastProductBadRequest
*/
type PostPvzPvzIDDeleteLastProductBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostPvzPvzIDDeleteLastProductBadRequest creates PostPvzPvzIDDeleteLastProductBadRequest with default headers values
func NewPostPvzPvzIDDeleteLastProductBadRequest() *PostPvzPvzIDDeleteLastProductBadRequest {

	return &PostPvzPvzIDDeleteLastProductBadRequest{}
}

// WithPayload adds the payload to the post pvz pvz Id delete last product bad request response
func (o *PostPvzPvzIDDeleteLastProductBadRequest) WithPayload(payload *models.Error) *PostPvzPvzIDDeleteLastProductBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post pvz pvz Id delete last product bad request response
func (o *PostPvzPvzIDDeleteLastProductBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostPvzPvzIDDeleteLastProductBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostPvzPvzIDDeleteLastProductForbiddenCode is the HTTP code returned for type PostPvzPvzIDDeleteLastProductForbidden
const PostPvzPvzIDDeleteLastProductForbiddenCode int = 403

/*
PostPvzPvzIDDeleteLastProductForbidden Доступ запрещен

swagger:response postPvzPvzIdDeleteLastProductForbidden
*/
type PostPvzPvzIDDeleteLastProductForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostPvzPvzIDDeleteLastProductForbidden creates PostPvzPvzIDDeleteLastProductForbidden with default headers values
func NewPostPvzPvzIDDeleteLastProductForbidden() *PostPvzPvzIDDeleteLastProductForbidden {

	return &PostPvzPvzIDDeleteLastProductForbidden{}
}

// WithPayload adds the payload to the post pvz pvz Id delete last product forbidden response
func (o *PostPvzPvzIDDeleteLastProductForbidden) WithPayload(payload *models.Error) *PostPvzPvzIDDeleteLastProductForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post pvz pvz Id delete last product forbidden response
func (o *PostPvzPvzIDDeleteLastProductForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostPvzPvzIDDeleteLastProductForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
