package bucharest_test

import (
	"bytes"
	"errors"
	"fmt"
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

func getCallingPath(handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) (string, error) {
	path := uuid.New().String()
	g := gin.New()
	gin.SetMode(gin.TestMode)
	if len(middlewares) > 0 {
		for _, middleware := range middlewares {
			g.Use(middleware)
		}
	}
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
		if err != nil {
			return false
		}
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
		if err != nil {
			return false
		}
		req.Header.Set("Cookie", "foo=bar;")
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestContentType(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	applicationJSONContentType := "application/json"
	handler := func(ctx HTTPContext) HTTPError {
		contentType := ctx.ContentType()
		assert.Equal(t, contentType, applicationJSONContentType)
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
		if err != nil {
			return false
		}
		req.Header.Set("Content-Type", applicationJSONContentType)
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGetHeader(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	XAPIKey := "f2dc1bdd32fa41389b0c5670a90081e6"
	handler := func(ctx HTTPContext) HTTPError {
		headerXAPIKey := ctx.GetHeader("X-API-Key")
		assert.Equal(t, XAPIKey, headerXAPIKey)
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
		if err != nil {
			return false
		}
		req.Header.Set("X-API-Key", XAPIKey)
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGetRawData(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		b, err := ctx.GetRawData()
		assert.NoError(t, err)

		m := make(map[string]interface{}, 0)
		err = utils.JSONMapper(b, &m)

		assert.NoError(t, err)
		assert.Equal(t, m["foo"], "bar")

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
		req, err := http.NewRequest(http.MethodGet, path, bytes.NewBufferString("{\"foo\":\"bar\"}"))
		if err != nil {
			return false
		}
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestWebSocketHeader(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		IsWebsocket := ctx.IsWebsocket()
		assert.False(t, IsWebsocket)

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
		if err != nil {
			return false
		}
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestHandlerControlNext(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		ctx.Status(http.StatusNoContent)
		return nil
	}
	middleware := func(ctx HTTPContext) HTTPError {
		ctx.Next()
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc, NewGinHandlerFunc(ctx, middleware))
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return false
		}
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestHandlerControlAbort(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		ctx.Status(http.StatusNoContent)
		return nil
	}
	middleware := func(ctx HTTPContext) HTTPError {
		ctx.Abort()
		assert.True(t, ctx.IsAborted())
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc, NewGinHandlerFunc(ctx, middleware))
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return false
		}
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestHandlerControlAbortWithStatus(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		ctx.Status(http.StatusNoContent)
		return nil
	}
	middleware := func(ctx HTTPContext) HTTPError {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		assert.True(t, ctx.IsAborted())
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc, NewGinHandlerFunc(ctx, middleware))
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return false
		}
		res, err = client.Do(req)
		if res != nil {
			assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
		}
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestHandlerControlAbortWithStatusJSON(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		ctx.Status(http.StatusNoContent)
		return nil
	}
	middleware := func(ctx HTTPContext) HTTPError {
		httpError := NewInternalServerError(errors.New("foobar"))
		ctx.AbortWithStatusJSON(httpError.GetStatus(), httpError.GetJSON())
		assert.True(t, ctx.IsAborted())
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc, NewGinHandlerFunc(ctx, middleware))
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return false
		}
		res, err = client.Do(req)

		if res != nil {
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
			j := make(map[string]interface{}, 0)
			err = utils.JSONMapper(res.Body, &j)
			assert.NoError(t, err)
			assert.Equal(t, j["message"].(string), "foobar")
		}

		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGetterAndSetter(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	var path string

	handler := func(ctx HTTPContext) HTTPError {
		bar, ok := ctx.Get("foo")
		assert.Equal(t, bar, "bar")
		assert.True(t, ok)

		bar, ok = ctx.MustGet("foo").(string)
		assert.Equal(t, bar, "bar")
		assert.True(t, ok)

		bar = ctx.GetString("foo")
		assert.Equal(t, bar, "bar")

		getInt := ctx.GetInt("int")
		assert.Equal(t, getInt, 2046)

		getInt64 := ctx.GetInt64("int64")
		assert.Equal(t, getInt64, int64(2046))

		getUint := ctx.GetUint("uint")
		assert.Equal(t, getUint, uint(2046))

		getUint64 := ctx.GetUint64("uint64")
		assert.Equal(t, getUint64, uint64(2046))

		getFloat64 := ctx.GetFloat64("float64")
		assert.Equal(t, getFloat64, float64(3.14))

		getBool := ctx.GetBool("bool")
		assert.True(t, getBool)

		getTime := ctx.GetTime("time")
		assert.Equal(t, getTime, time.Time{})

		getDuration := ctx.GetDuration("duration")
		assert.Equal(t, getDuration, time.Second)

		getStringSlice := ctx.GetStringSlice("ss")
		assert.Len(t, getStringSlice, 2)
		assert.Equal(t, getStringSlice[0], "foobar")
		assert.Equal(t, getStringSlice[1], "fozbaz")

		getStringMap := ctx.GetStringMap("sm")
		assert.Len(t, getStringMap, 2)
		assert.Equal(t, getStringMap["one"], 1)
		assert.Equal(t, getStringMap["true"], true)

		getStringMapString := ctx.GetStringMapString("sms")
		assert.Len(t, getStringMapString, 2)
		assert.Equal(t, getStringMapString["foo"], "bar")
		assert.Equal(t, getStringMapString["foz"], "baz")

		getStringMapStringSlice := ctx.GetStringMapStringSlice("smss")
		assert.Len(t, getStringMapStringSlice, 2)
		assert.Equal(t, getStringMapStringSlice["foo"][0], "bar")
		assert.Equal(t, getStringMapStringSlice["foz"][0], "baz")

		ctx.Status(http.StatusNoContent)
		return nil
	}
	middleware := func(ctx HTTPContext) HTTPError {
		ctx.Set("foo", "bar")
		ctx.Set("int", 2046)
		ctx.Set("int64", int64(2046))
		ctx.Set("uint", uint(2046))
		ctx.Set("uint64", uint64(2046))
		ctx.Set("float64", float64(3.14))
		ctx.Set("bool", true)
		ctx.Set("time", time.Time{})
		ctx.Set("duration", time.Second)
		ctx.Set("ss", []string{"foobar", "fozbaz"})
		ctx.Set("sm", map[string]interface{}{"one": 1, "true": true})
		ctx.Set("sms", map[string]string{"foo": "bar", "foz": "baz"})
		ctx.Set("smss", map[string][]string{"foo": {"bar"}, "foz": {"baz"}})
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	var err error
	path, err = getCallingPath(ginHandlerFunc, NewGinHandlerFunc(ctx, middleware))
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return false
		}
		res, err = client.Do(req)
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
