package bucharest

import (
	"github.com/DATA-DOG/go-sqlmock"
)

type MockContext interface {
	Context
	SQLMock() sqlmock.Sqlmock
}

type mockContext struct {
	Context
	sqlMock sqlmock.Sqlmock
}

func NewMockContext(parentContext Context, mock sqlmock.Sqlmock) MockContext {
	return &mockContext{
		Context: parentContext,
		sqlMock: mock,
	}
}

func (ctx mockContext) SQLMock() sqlmock.Sqlmock {
	return ctx.sqlMock
}
