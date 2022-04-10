package bucharest_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/argonlab-io/bucharest"
	"github.com/argonlab-io/bucharest/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var DEFAULT_TEST_PORT int = 9000

func getCallingPath(handler gin.HandlerFunc) (string, error) {
	path := uuid.New().String()
	g := gin.New()
	gin.SetMode(gin.TestMode)
	g.GET(fmt.Sprintf("/%s", path), handler)
	var err error
	avaiable_port := fmt.Sprint(DEFAULT_TEST_PORT)
	DEFAULT_TEST_PORT++
	go func() { err = g.Run(fmt.Sprintf(":%s", avaiable_port)) }()
	return fmt.Sprintf("http://0.0.0.0:%s/%s", avaiable_port, path), err
}

func TestNewGinHandlerFunc(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	handler := func(ctx HTTPContext) HTTPError {
		return NewBadRequestError(errors.New("foobar"))
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err := getCallingPath(ginHandlerFunc)
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		res, err = http.Get(path)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestNewGinHandlerFuncWithData(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	handler := func(ctx HTTPContext, data Map) HTTPError {
		return NewBadRequestError(errors.New(data["message"].(string)))
	}
	ginHandlerFunc := NewGinHandlerFuncWithData(ctx, handler, Map{"message": "foobar"})
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err := getCallingPath(ginHandlerFunc)
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		res, err = http.Get(path)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	j := make(map[string]interface{}, 0)
	err = utils.JSONMapper(res.Body, &j)
	assert.NoError(t, err)
	log.Println(j)
	assert.Equal(t, j["message"].(string), "foobar")
}

func TestHandleInfo(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	handler := func(ctx HTTPContext) HTTPError {
		assert.Equal(t, ctx.HandlerNames()[0], ctx.HandlerName())
		ctx.Status(http.StatusNoContent)
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err := getCallingPath(ginHandlerFunc)
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		res, err = http.Get(path)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestHandleFullPath(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		split := strings.Split(path, "/")
		handler_path := split[len(split)-1]
		assert.Equal(t, fmt.Sprintf("/%s", handler_path), ctx.FullPath())
		ctx.Status(http.StatusNoContent)
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc)
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		res, err = http.Get(path)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestIP(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		assert.Equal(t, ctx.Gin().ClientIP(), ctx.ClientIP())
		remote_ip, trust := ctx.RemoteIP()
		expected_remote_ip, expected_trust := ctx.Gin().RemoteIP()
		assert.Equal(t, remote_ip, expected_remote_ip)
		assert.Equal(t, trust, expected_trust)
		ctx.Status(http.StatusNoContent)
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc)
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCookie(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		foo, err := ctx.Cookie("foo")
		assert.NoError(t, err)
		assert.Equal(t, foo, "bar")
		ctx.Status(http.StatusNoContent)
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc)
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Cookie", "foo=bar;")
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
