package bucharest_test

import (
	"encoding/json"
	"errors"
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
	test_err := errors.New("fooooooooooo")
	httpError := NewBadRequestError(test_err)
	assert.NotEmpty(t, httpError)
	assert.True(t, errors.Is(httpError.OriginalError(), test_err))
}

func TestNewBadRequestErrorWithJSON(t *testing.T) {
	test_err := &jsonError{
		"foo": "bar",
	}
	httpError := NewBadRequestError(test_err)
	assert.NotEmpty(t, httpError)
	assert.ErrorIs(t, httpError.OriginalError(), test_err)
	serializable := httpError.GetJSON()
	mapper := make(map[string]interface{})
	err := utils.JSONMapper(serializable, &mapper)
	assert.NoError(t, err)
	jerr := &jsonError{}
	err = utils.JSONMapper(mapper["message"], jerr)
	assert.NoError(t, err)
	assert.Equal(t, jerr, test_err)
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
