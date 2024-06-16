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

func (suite *GinTestSuite) SetupSuite() {
	suite.Port = 9000
}

func (suite *GinTestSuite) SetupTest() {
	suite.Port++
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

func (suite *GinTestSuite) createTestServer(option *testServerOption) {
	g := gin.New()
	gin.SetMode(gin.TestMode)
	g.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	if len(option.middlewares) > 0 {
		for _, middleware := range option.middlewares {
			g.Use(middleware)
		}
	}

	paths := make([]string, 0)
	port := fmt.Sprint(suite.Port)

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

	suite.Server = server
	suite.Paths = paths
}

func (suite *GinTestSuite) startTestServer() {
	go func() {
		if err := suite.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			assert.NoError(suite.T(), err)
		}
	}()
}

func (suite *GinTestSuite) shutdownTestServer() {
	err := suite.Server.Shutdown(suite.Ctx)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.Ctx.Err())
}

func (suite *GinTestSuite) assertServerHealthy() {
	assert.NoError(suite.T(), utils.Await(&utils.AwaitOptions{
		Condition: func() bool {
			res, err := http.Get(fmt.Sprintf("http://%s", suite.Server.Addr))
			return err == nil && res.StatusCode == http.StatusNoContent
		},
		Timeout: 5 * time.Second,
	}))
}

func (suite *GinTestSuite) assertNoContentResponse(resp *http.Response, err error) {
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)
}

func (suite *GinTestSuite) assertOkResponse(resp *http.Response, err error) {
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}
