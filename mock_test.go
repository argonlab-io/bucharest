package bucharest_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/argonlab-io/bucharest"
	"github.com/stretchr/testify/assert"
)

func TestNewMockContext(t *testing.T) {
	_, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	assert.NotEmpty(t, mock)

	ctxMock := NewMockContext(nil, mock)
	assert.NotEmpty(t, ctxMock)
	assert.Equal(t, mock, ctxMock.SQLMock())
}
