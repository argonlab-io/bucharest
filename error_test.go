package bucharest_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	. "github.com/argonlab-io/bucharest"
	"github.com/argonlab-io/bucharest/utils"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type jsonError map[string]interface{}

func (je *jsonError) Error() string {
	b, _ := json.Marshal(je)
	return string(b[:])
}

func TestNewBadRequestError(t *testing.T) {
	testErr := errors.New("fooooooooooo")
	httpError := NewBadRequestError(testErr)
	assert.NotEmpty(t, httpError)
	assert.True(t, errors.Is(httpError.OriginalError(), testErr))
}

func TestNewBadRequestErrorWithJSON(t *testing.T) {
	testErr := &jsonError{
		"foo": "bar",
	}
	httpError := NewBadRequestError(testErr)
	assert.NotEmpty(t, httpError)
	assert.ErrorIs(t, httpError.OriginalError(), testErr)
	serializable := httpError.GetJSON()
	mapper := make(map[string]interface{})
	err := utils.JSONMapper(serializable, &mapper)
	assert.NoError(t, err)
	jerr := &jsonError{}
	err = utils.JSONMapper(mapper["message"], jerr)
	assert.NoError(t, err)
	assert.Equal(t, jerr, testErr)
}

func TestNewBadRequestErrorWithFromValidateError(t *testing.T) {
	validate := validator.New()

	type myStruct struct {
		Foo string `validate:"required"`
		Foz string
	}

	errNooFoo := &myStruct{
		Foz: "baz",
	}
	valErr := validate.Struct(errNooFoo).(validator.ValidationErrors)
	httpError := NewBadRequestError(valErr)
	assert.Equal(t, httpError.OriginalError(), valErr)
	assert.Equal(t, httpError.GetStatus(), http.StatusBadRequest)
	serializable := httpError.GetJSON()
	mapper := make(map[string]interface{})
	err := utils.JSONMapper(serializable, &mapper)
	assert.NoError(t, err)
	verr, ok := mapper["message"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, verr)
	assert.NotEmpty(t, verr["Foo"])
	assert.Equal(t, verr["Foo"], map[string]interface{}{"error": "Key: 'myStruct.Foo' Error:Field validation for 'Foo' failed on the 'required' tag"})
}

func TestInternalServerError(t *testing.T) {
	valErr := errors.New("error: foo")
	httpError := NewInternalServerError(valErr)
	assert.Equal(t, httpError.OriginalError(), valErr)
	assert.Equal(t, httpError.GetStatus(), http.StatusInternalServerError)
	serializable := httpError.GetJSON()
	mapper := make(map[string]interface{})
	err := utils.JSONMapper(serializable, &mapper)
	assert.NoError(t, err)
	message, ok := mapper["message"]
	assert.True(t, ok)
	assert.Equal(t, message, valErr.Error())
}

func TestInternalServerErrorWithSerializableError(t *testing.T) {
	valErr := &jsonError{"foo": "bar"}
	httpError := NewInternalServerError(valErr)
	assert.Equal(t, httpError.OriginalError(), valErr)
	assert.Equal(t, httpError.GetStatus(), http.StatusInternalServerError)
	serializable := httpError.GetJSON()
	mapper := make(map[string]interface{})
	err := utils.JSONMapper(serializable, &mapper)
	assert.NoError(t, err)
	message, ok := mapper["message"]
	assert.True(t, ok)
	_, ok = message.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, message.(map[string]interface{})["foo"], "bar")
	assert.Equal(t, valErr.Error(), fmt.Sprint(valErr))
}
