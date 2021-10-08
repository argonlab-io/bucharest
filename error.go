package bucharest

import (
	"net/http"

	"github.com/argonlab-io/bucharest/utils"
	"github.com/go-playground/validator/v10"
)

type httpError struct {
	status  int
	Message interface{} `json:"message"`
}

func (e *httpError) GetStatus() int {
	return e.status
}

func (e *httpError) GetJSON() interface{} {
	return e
}

type ValidationErrors struct {
	Error string `json:"error"`
	Param string `json:"param,omitempty"`
}

func NewBadRequestError(err error) HTTPError {
	mapper := make(map[string]interface{})
	jerr := utils.JSONMapper(err, mapper)
	if jerr != nil && len(mapper) != 0 {
		return &httpError{
			status:  http.StatusBadRequest,
			Message: mapper,
		}
	}

	validatorErrors, ok := err.(validator.ValidationErrors)
	if ok {
		for _, validatorError := range validatorErrors {
			mapper[validatorError.Field()] = &ValidationErrors{
				Error: validatorError.Error(),
				Param: validatorError.Param(),
			}
		}
		return &httpError{
			status:  http.StatusBadRequest,
			Message: mapper,
		}
	}

	return &httpError{
		status:  http.StatusBadRequest,
		Message: err.Error(),
	}
}

func NewInternalServerError(err error) HTTPError {
	mapper := make(map[string]interface{})
	jerr := utils.JSONMapper(err, mapper)
	if jerr != nil && len(mapper) != 0 {
		return &httpError{
			status:  http.StatusInternalServerError,
			Message: mapper,
		}
	}

	return &httpError{
		status:  http.StatusInternalServerError,
		Message: err.Error(),
	}
}
