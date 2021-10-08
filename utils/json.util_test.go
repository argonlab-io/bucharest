package utils_test

import (
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

func TestJSONMapperMapToMap(t *testing.T) {
	src := map[string]interface{}{
		"foo": "bar",
		"foz": "baz",
	}
	dest := make(map[string]interface{})
	err := JSONMapper(src, &dest)
	assert.NoError(t, err)
	assert.NotEmpty(t, dest)
	assert.NotEmpty(t, dest["foo"])
	assert.Equal(t, dest["foo"], "bar")
	assert.NotEmpty(t, dest["foz"])
	assert.Equal(t, dest["foz"], "baz")
}

type jsonMapperTestStruct struct {
	Foo string `json:"foo"`
	Foz string `json:"foz"`
}

func TestJSONMapperMapToStruct(t *testing.T) {
	src := map[string]interface{}{
		"foo": "bar",
		"foz": "baz",
	}
	dest := &jsonMapperTestStruct{}
	err := JSONMapper(src, dest)
	assert.NoError(t, err)
	assert.NotEmpty(t, dest)
	assert.NotEmpty(t, dest.Foo)
	assert.Equal(t, dest.Foo, "bar")
	assert.NotEmpty(t, dest.Foz)
	assert.Equal(t, dest.Foz, "baz")
}

func TestJSONMapperStructToMap(t *testing.T) {
	src := &jsonMapperTestStruct{Foo: "bar", Foz: "baz"}
	dest := make(map[string]interface{})
	err := JSONMapper(src, &dest)
	assert.NoError(t, err)
	assert.NotEmpty(t, dest)
	assert.NotEmpty(t, dest["foo"])
	assert.Equal(t, dest["foo"], "bar")
	assert.NotEmpty(t, dest["foz"])
	assert.Equal(t, dest["foz"], "baz")
}

func TestJSONMapperStructToStruct(t *testing.T) {
	src := &jsonMapperTestStruct{Foo: "bar", Foz: "baz"}
	dest := &jsonMapperTestStruct{}
	err := JSONMapper(src, dest)
	assert.NoError(t, err)
	assert.NotEmpty(t, dest)
	assert.NotEmpty(t, dest.Foo)
	assert.Equal(t, dest.Foo, "bar")
	assert.NotEmpty(t, dest.Foz)
	assert.Equal(t, dest.Foz, "baz")
}
