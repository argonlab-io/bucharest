package bucharest_test

import (
	"net/http"

	. "github.com/argonlab-io/bucharest"
)

type GinMiddlewareTestSuite struct {
	GinTestSuite
}

func (suite *GinMiddlewareTestSuite) TestNext() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	suite.Ctx = ctx
	handler := func(ctx HTTPContext) HTTPError {
		ctx.Status(http.StatusNoContent)
		return nil
	}
	suite.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		}})

	suite.startTestServer()
	suite.assertServerHealthy()

	// Act
	resp, err := http.Get(suite.Paths[0])

	// Assert
	suite.assertNoContentResponse(resp, err)
}
