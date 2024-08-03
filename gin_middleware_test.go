package bucharest_test

import (
	"net/http"
	"testing"

	. "github.com/argonlab-io/bucharest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GinMiddlewareTestSuite struct {
	GinTestSuite
}

func (ts *GinMiddlewareTestSuite) TestNext() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx
	handler := func(ctx HTTPContext) HTTPError {
		assert.Equal(ts.T(), "my-middleware", ctx.GetString("middleware"))
		ctx.Status(http.StatusNoContent)
		return nil
	}
	middleware := func(ctx HTTPContext) HTTPError {
		ctx.Set("middleware", "my-middleware")
		ctx.Next()
		return nil
	}
	ts.createTestServer(&testServerOption{
		handlers: []*GinTestHandler{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		},
		middlewares: []gin.HandlerFunc{
			NewGinHandlerFunc(&NewHandlerPayload{
				Ctx:  ctx,
				Func: middleware,
			}),
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	resp, err := http.Get(ts.Paths[0])

	// Assert
	ts.assertNoContentResponse(resp, err)
	ts.HTTPStatusCode(ts.Server.Handler.ServeHTTP,
		http.MethodGet,
		ts.Paths[0],
		nil,
		http.StatusNoContent,
	)
}

func (ts *GinMiddlewareTestSuite) TestAbort() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx
	handler := ts.createUnreachableHandler()

	middleware := func(ctx HTTPContext) HTTPError {
		ctx.Abort()
		assert.True(ts.T(), ctx.IsAborted())
		return nil
	}
	ts.createTestServer(&testServerOption{
		handlers: []*GinTestHandler{
			ts.createGetTestHandler(ctx, handler),
		},
		middlewares: []gin.HandlerFunc{
			NewGinHandlerFunc(&NewHandlerPayload{
				Ctx:  ctx,
				Func: middleware,
			}),
		}})

	ts.startTestServer()

	// Act and Assert
	// being set at c.writemem.reset(w), w.status = defaultStatus = http.StatusOK
	ts.HTTPStatusCode(ts.Server.Handler.ServeHTTP,
		http.MethodGet,
		ts.Paths[0],
		nil,
		http.StatusOK,
	)
}

func (ts *GinMiddlewareTestSuite) TestAbortWithStatus() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx
	handler := ts.createUnreachableHandler()

	middleware := func(ctx HTTPContext) HTTPError {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		assert.True(ts.T(), ctx.IsAborted())
		return nil
	}
	ts.createTestServer(&testServerOption{
		handlers: []*GinTestHandler{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		},
		middlewares: []gin.HandlerFunc{
			NewGinHandlerFunc(&NewHandlerPayload{
				Ctx:  ctx,
				Func: middleware,
			}),
		}})

	ts.startTestServer()

	// Act and Assert
	ts.HTTPStatusCode(ts.Server.Handler.ServeHTTP,
		http.MethodGet,
		ts.Paths[0],
		nil,
		http.StatusInternalServerError,
	)
}

func (ts *GinMiddlewareTestSuite) TestAbortWithStatusJSON() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx
	handler := ts.createUnreachableHandler()

	middleware := func(ctx HTTPContext) HTTPError {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": "foobar",
		})
		assert.True(ts.T(), ctx.IsAborted())
		return nil
	}
	ts.createTestServer(&testServerOption{
		handlers: []*GinTestHandler{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		},
		middlewares: []gin.HandlerFunc{
			NewGinHandlerFunc(&NewHandlerPayload{
				Ctx:  ctx,
				Func: middleware,
			}),
		}})

	ts.startTestServer()

	// Act and Assert
	ts.HTTPStatusCode(ts.Server.Handler.ServeHTTP,
		http.MethodGet,
		ts.Paths[0],
		nil,
		http.StatusInternalServerError,
	)
	ts.HTTPBodyContains(ts.Server.Handler.ServeHTTP,
		http.MethodGet,
		ts.Paths[0],
		nil,
		"{\"message\":\"foobar\"}",
	)
}

func TestGinMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(GinMiddlewareTestSuite))
}
