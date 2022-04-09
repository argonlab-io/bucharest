package bucharest_test

import (
	"errors"
	"log"
	"net/http"
	"testing"
	"time"

	. "github.com/argonlab-io/bucharest"
	"github.com/argonlab-io/bucharest/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewGinHandlerFunc(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	assert.NotNil(t, ctx)

	handler := func(ctx HTTPContext) HTTPError {
		return NewBadRequestError(errors.New("foobar"))
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	g := gin.Default()
	g.GET("/foobar", ginHandlerFunc)
	var err error
	go func() { err = g.Run(":9000") }()
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		res, err = http.Get("http://0.0.0.0:9000/foobar")
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

	g := gin.Default()
	g.GET("/foobar", ginHandlerFunc)
	var err error
	go func() { err = g.Run(":9001") }()
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		res, err = http.Get("http://0.0.0.0:9001/foobar")
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
		// handler[0] is logger
		// handler[1] is recover
		assert.Equal(t, ctx.HandlerNames()[2], ctx.HandlerName())
		ctx.Status(http.StatusNoContent)
		return nil
	}
	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
	assert.NotNil(t, ginHandlerFunc)

	g := gin.Default()
	g.GET("/foobar", ginHandlerFunc)
	var err error
	go func() { err = g.Run(":9002") }()
	assert.NoError(t, err)

	var res *http.Response
	fn := func() bool {
		res, err = http.Get("http://0.0.0.0:9002/foobar")
		return err == nil
	}

	utils.RunUntil(fn, time.Second*4)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
