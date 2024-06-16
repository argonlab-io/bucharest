package bucharest_test

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	. "github.com/argonlab-io/bucharest"
	"github.com/argonlab-io/bucharest/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type binderTest struct {
	Foo string `form:"foo" json:"foo" xml:"foo"  binding:"required"`
}

type GinTestSuite struct {
	suite.Suite
	Port   int
	Server *http.Server
	Paths  []string
	Ctx    Context
}

func (ts *GinTestSuite) SetupSuite() {
	ts.Port = 9000
}

func (ts *GinTestSuite) SetupTest() {
	ts.Port++
}

func (suite *GinTestSuite) TearDownTest() {

	suite.shutdownTestServer()
}

type testHandlers struct {
	function gin.HandlerFunc
	method   string
}

type testServerOption struct {
	handlers    []*testHandlers
	middlewares []gin.HandlerFunc
}

func (ts *GinTestSuite) createTestServer(option *testServerOption) {
	g := gin.New()
	g.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
	gin.SetMode(gin.TestMode)
	if len(option.middlewares) > 0 {
		for _, middleware := range option.middlewares {
			g.Use(middleware)
		}
	}

	paths := make([]string, 0)
	port := fmt.Sprint(ts.Port)

	for _, handler := range option.handlers {
		randomPath := uuid.New().String()
		paths = append(paths, fmt.Sprintf("http://0.0.0.0:%s/%s", port, randomPath))
		switch handler.method {
		case http.MethodGet:
			{
				g.GET(fmt.Sprintf("/%s", randomPath), handler.function)
			}
		case http.MethodPost:
			{
				g.POST(fmt.Sprintf("/%s", randomPath), handler.function)
			}
		default:
			{
				panic("not implemented")
			}
		}
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: g.Handler(),
	}

	ts.Server = server
	ts.Paths = paths
}

func (ts *GinTestSuite) startTestServer() {
	go func() {
		if err := ts.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			assert.NoError(ts.T(), err)
		}
	}()
}

func (ts *GinTestSuite) shutdownTestServer() {
	err := ts.Server.Shutdown(ts.Ctx)
	assert.NoError(ts.T(), err)
	assert.NoError(ts.T(), ts.Ctx.Err())
}

func (ts *GinTestSuite) assertServerHealthy() {
	assert.NoError(ts.T(), utils.Await(&utils.AwaitOptions{
		Condition: func() bool {
			res, err := http.Get(fmt.Sprintf("http://%s", ts.Server.Addr))
			return err == nil && res.StatusCode == http.StatusNoContent
		},
		Timeout: 5 * time.Second,
	}))
}

func (ts *GinTestSuite) assertNoContentResponse(resp *http.Response, err error) {
	assert.NoError(ts.T(), err)
	assert.NotNil(ts.T(), resp)
	assert.Equal(ts.T(), http.StatusNoContent, resp.StatusCode)
}

func (ts *GinTestSuite) assertOkResponse(resp *http.Response, err error) {
	assert.NoError(ts.T(), err)
	assert.NotNil(ts.T(), resp)
	assert.Equal(ts.T(), http.StatusOK, resp.StatusCode)
}
