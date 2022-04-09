package bucharest

import (
	"net/http"

	"github.com/argonlab-io/bucharest/utils"
	"github.com/go-playground/validator/v10"
)

type HttpError struct {
	status        int
	Message       interface{} `json:"message"`
	originalError error
}

func (e *HttpError) GetStatus() int {
	return e.status
}

func (e *HttpError) GetJSON() interface{} {
	return e
}

func (e *HttpError) OriginalError() error {
	return e.originalError
}

type ValidationErrors struct {
	Error string `json:"error"`
	Param string `json:"param,omitempty"`
}

func getErrorMapper(err error) map[string]interface{} {
	mapper := make(map[string]interface{})
	jerr := utils.JSONMapper(err, &mapper)
	if jerr != nil {
		return make(map[string]interface{})
	}
	return mapper
}

func newHttpErrorFromMapper(status int, mapper map[string]interface{}, err error) *HttpError {
	return &HttpError{
		status:        status,
		Message:       mapper,
		originalError: err,
	}
}

func NewBadRequestError(err error) HTTPError {
	mapper := getErrorMapper(err)
	if len(mapper) != 0 {
		return newHttpErrorFromMapper(http.StatusBadRequest, mapper, err)
	}

	validatorErrors, ok := err.(validator.ValidationErrors)
	if ok {
		mapper = make(map[string]interface{})
		for _, validatorError := range validatorErrors {
			mapper[validatorError.Field()] = &ValidationErrors{
				Error: validatorError.Error(),
				Param: validatorError.Param(),
			}
		}

		return &HttpError{
			status:        http.StatusBadRequest,
			Message:       mapper,
			originalError: err,
		}
	}

	return &HttpError{
		status:        http.StatusBadRequest,
		Message:       err.Error(),
		originalError: err,
	}
}

func NewInternalServerError(err error) HTTPError {
	mapper := getErrorMapper(err)
	if len(mapper) != 0 {
		return newHttpErrorFromMapper(http.StatusInternalServerError, mapper, err)
	}

	return &HttpError{
		status:        http.StatusInternalServerError,
		Message:       err.Error(),
		originalError: err,
	}
}
