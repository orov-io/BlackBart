package response

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	production = "prod"
	envKey     = "ENV"
)

// Logger is the needed interface to log entries by the response package
type Logger interface {
	WithError(err error) *logrus.Entry
}

var logger Logger

func init() {
	logger = logrus.New()
}

// SetLogger change the logger use by the response package
func SetLogger(theLogger Logger) {
	logger = theLogger
}

// Response models standard response
type Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
	data    interface{}
	ctx     *gin.Context
}

// NewResponse returns a response struct with a context attached
func newResponse(c *gin.Context) *Response {
	r := new(Response)
	r.ctx = c
	return r
}

// AddData adds objects to the data response field
func (r *Response) addData(data interface{}) {
	r.data = data
}

// AddError adds objects to the errors response field
func (r *Response) addError(errors ...error) {
	if r.Errors == nil {
		r.Errors = make([]string, 0)
	}

	for _, err := range errors {
		if err == nil {
			err = fmt.Errorf("Unknown error")
		}

		if os.Getenv(envKey) == production {
			traceID := uuid.New().String()
			logger.WithError(err).WithFields(
				logrus.Fields{
					"trace_id": traceID,
				},
			).Warningf("SERVER ERROR")
			err = NewHiddenError(traceID)
		}
		fmt.Printf("r.Errors: %v\n", r.Errors)
		fmt.Printf("err: %v\n", err)
		fmt.Printf("err.Error(): %v\n", err.Error())
		r.Errors = append(r.Errors, err.Error())
	}
}

// Unauthorized sends a 401 code to the client and ask for re-loggin
func (r *Response) unauthorized() {
	r.Message = "You are no logged-in. Please, loggin"
	r.ctx.JSON(http.StatusUnauthorized, r)
	r.ctx.Abort()
}

func (r *Response) forbidden() {
	r.Message = "User has no enough permissions"
	r.ctx.JSON(http.StatusForbidden, r)
	r.ctx.Abort()
}

func (r *Response) badRequest() {
	r.ctx.JSON(http.StatusBadRequest, r)
	r.ctx.Abort()
}

func (r *Response) internalError() {
	r.ctx.JSON(http.StatusInternalServerError, r)
	r.ctx.Abort()
}

func (r *Response) ok() {
	r.ctx.JSON(http.StatusOK, r.data)
	r.ctx.Abort()
}

// here, we expose the response catalog of the server

// SendUnauthorizedAccess returns a 401 with error info if any
func SendUnauthorizedAccess(c *gin.Context, errors ...error) {
	r := newResponse(c)
	r.addError(errors...)
	r.unauthorized()
}

// SendBadRequest returns a 400 code with the errors info to the client
func SendBadRequest(c *gin.Context, errors ...error) {
	r := newResponse(c)
	r.addError(errors...)
	r.badRequest()
}

// SendInternalError send a standard 500 response
func SendInternalError(c *gin.Context, errors ...error) {
	r := newResponse(c)
	r.addError(errors...)
	r.internalError()
}

// SendCreated returns a 201 http code with the location header attached
func SendCreated(c *gin.Context, location string, data ...interface{}) {
	c.Header("Location", location)
	if data == nil {
		c.AbortWithStatus(http.StatusCreated)
		return
	}

	var response interface{}
	switch len(data) {
	case 0:
		c.AbortWithStatus(http.StatusCreated)
		return

	case 1:
		if data[0] == nil {
			response = gin.H{}
		} else {
			response = data[0]
		}

	default:
		response = data
	}

	c.JSON(http.StatusCreated, response)
}

// SendNoContent returns a 204 http code with no body.
func SendNoContent(c *gin.Context) {
	c.AbortWithStatus(http.StatusNoContent)
}

// SendOK sends the provided data with a 200 http code status.
func SendOK(c *gin.Context, data interface{}) {
	r := newResponse(c)
	r.addData(data)
	r.ok()
}

// SendForbidden sends an http 403 code to the client.
func SendForbidden(c *gin.Context, errors ...error) {
	r := newResponse(c)
	r.addError(errors...)
	r.forbidden()
}

// Parse parses the body of a http.response to a Response struct
func Parse(response *http.Response) (*Response, error) {
	if !isValidResponse(response) {
		r, err := parseError(response)
		if err != nil {
			return r, NewError(response.StatusCode, "Unknown response format")
		}
		return r, NewError(response.StatusCode, r.Message)
	}

	return parseValid(response)
}

// Error models a http error with a message and code
type Error struct {
	message string
	code    int
}

// NewError returns a new Error
func NewError(code int, message string) *Error {
	return &Error{
		message: message,
		code:    code,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("With code: %v -> %s", e.code, e.message)
}

// Code returns the response error inner code
func (e *Error) Code() int {
	return e.code
}

// IsError checks if the error is a ResponseError error
func IsError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

func parseError(response *http.Response) (r *Response, err error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	r = new(Response)
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func parseValid(response *http.Response) (r *Response, err error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var data interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	r.data = data
	return

}

// ParseTo parses the body response and unmarshal content of the data
// field into the receiver
func ParseTo(response *http.Response, receiver interface{}) error {
	if !isAPointer(receiver) {
		return NewNotAPointerError()
	}

	ParsedResponse, err := Parse(response)
	if err != nil {
		return err
	}

	ResponseBytes, err := json.Marshal(ParsedResponse.data)
	if err != nil {
		return fmt.Errorf("Error: %v\nCan't marshal response data: %v", err, ParsedResponse)
	}
	err = json.Unmarshal(ResponseBytes, receiver)
	if err != nil {
		return fmt.Errorf("Error: %v\nCan't unmarshal response data: %v", err, ParsedResponse)
	}

	return nil
}

func isAPointer(i interface{}) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr
}

func isValidResponse(response *http.Response) bool {
	return response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusBadRequest
}

// NotAPointerError is used to send a hidden error
type NotAPointerError struct{}

func (e *NotAPointerError) Error() string {
	return fmt.Sprintf("Receiver is not a pointer")
}

// NewNotAPointerError returns a new NotAPointerErrorError error
func NewNotAPointerError() error {
	return &NotAPointerError{}
}

// IsNotAPointerError checks if the error is a NotAPointerError error
func IsNotAPointerError(err error) bool {
	_, ok := err.(*NotAPointerError)
	return ok
}

// HiddenError is used to send a hidden error
type HiddenError struct {
	TraceUUID string
}

func (e *HiddenError) Error() string {
	return fmt.Sprintf("An error occurs on the server. Please, see trace: %v", e.TraceUUID)
}

// NewHiddenError returns a new HiddenErrorError error
func NewHiddenError(traceUUID string) error {
	return &HiddenError{traceUUID}
}

// IsHiddenError checks if the error is a HiddenError error
func IsHiddenError(err error) bool {
	_, ok := err.(*HiddenError)
	return ok
}
