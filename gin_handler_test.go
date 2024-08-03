package bucharest_test

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	. "github.com/argonlab-io/bucharest"
	"github.com/argonlab-io/bucharest/consts"
	"github.com/argonlab-io/bucharest/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GinHandlerTestSuite struct {
	GinTestSuite
}

func (ts *GinHandlerTestSuite) TestGinHandler() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx
	handler := func(ctx HTTPContext) HTTPError {
		ctx.Status(http.StatusNoContent)
		return nil
	}
	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	resp, err := http.Get(ts.Paths[0])

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestNewGinHandlerFuncWithData() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx
	handler := func(ctx HTTPContext) HTTPError {
		ctx.String(http.StatusOK, "%s", ctx.Value("message"))
		return nil
	}
	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
					Data: map[string]any{"message": "foobar"},
				}),
				method: http.MethodGet,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	resp, err := http.Get(ts.Paths[0])

	// Assert
	ts.assertOkResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestHandlerInfo() {
	// // Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx
	handler := func(ctx HTTPContext) HTTPError {
		assert.Len(ts.T(), ctx.HandlerNames(), 2)
		assert.Equal(ts.T(), ctx.HandlerNames()[0], "UnamedHandler")
		assert.Equal(ts.T(), ctx.HandlerNames()[1], ctx.HandlerName())
		assert.Equal(ts.T(), "MyHandler", ctx.HandlerName())
		ctx.Status(http.StatusNoContent)
		return nil
	}

	handlerName := "MyHandler"
	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
					Name: &handlerName,
				}),
				method: http.MethodGet,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	resp, err := http.Get(ts.Paths[1])

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestFullPath() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx

	var path string
	handler := func(ctx HTTPContext) HTTPError {
		split := strings.Split(path, "/")
		handlerPath := split[len(split)-1]
		assert.Equal(ts.T(), fmt.Sprintf("/%s", handlerPath), ctx.FullPath())
		ctx.Status(http.StatusNoContent)
		return nil
	}

	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		}})
	path = ts.Paths[0]

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	resp, err := http.Get(ts.Paths[0])

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestIp() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx

	handler := func(ctx HTTPContext) HTTPError {
		assert.Equal(ts.T(), ctx.Gin().ClientIP(), ctx.ClientIP())
		remoteIp := ctx.RemoteIP()
		expectedRemoteIp := ctx.Gin().RemoteIP()
		assert.Equal(ts.T(), remoteIp, expectedRemoteIp)
		ctx.Status(http.StatusNoContent)
		return nil
	}

	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	resp, err := http.Get(ts.Paths[0])

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestCookie() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx

	handler := func(ctx HTTPContext) HTTPError {
		foo, err := ctx.Cookie("foo")
		assert.NoError(ts.T(), err)
		assert.Equal(ts.T(), foo, "bar")
		ctx.Status(http.StatusNoContent)
		return nil
	}

	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.Paths[0], nil)
	req.Header.Set("Cookie", "foo=bar;")
	resp, err := client.Do(req)

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestContentType() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx

	handler := func(ctx HTTPContext) HTTPError {
		contentType := ctx.ContentType()
		assert.Equal(ts.T(), contentType, gin.MIMEJSON)
		ctx.Status(http.StatusNoContent)
		return nil
	}

	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.Paths[0], nil)
	req.Header.Set(string(consts.ContentType), gin.MIMEJSON)
	resp, err := client.Do(req)

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestHeaders() {
	// Arrange
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx

	handler := func(ctx HTTPContext) HTTPError {
		myHeader := ctx.GetHeader("X-My-Header")
		assert.Equal(ts.T(), "my-header", myHeader)
		ctx.Status(http.StatusNoContent)
		return nil
	}

	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodGet,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ts.Paths[0], nil)
	req.Header.Set("X-My-Header", "my-header")
	resp, err := client.Do(req)

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func (ts *GinHandlerTestSuite) TestRawData() {
	ctx := NewContextWithOptions(nil)
	ts.Ctx = ctx

	handler := func(ctx HTTPContext) HTTPError {
		rawData, err := ctx.GetRawData()
		assert.NoError(ts.T(), err)

		data := make(map[string]any, 0)
		err = utils.JSONMapper(rawData, &data)

		assert.NoError(ts.T(), err)
		assert.Equal(ts.T(), data["message"], "foobar")

		ctx.Status(http.StatusNoContent)
		return nil
	}

	ts.createTestServer(&testServerOption{
		handlers: []*testHandlers{
			{
				function: NewGinHandlerFunc(&NewHandlerPayload{
					Ctx:  ctx,
					Func: handler,
				}),
				method: http.MethodPost,
			},
		}})

	ts.startTestServer()
	ts.assertServerHealthy()

	// Act
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost,
		ts.Paths[0],
		bytes.NewBufferString("{\"message\":\"foobar\"}"))
	resp, err := client.Do(req)

	// Assert
	ts.assertNoContentResponse(resp, err)
}

func TestGinHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GinHandlerTestSuite))
}

// func TestGetterAndSetter(t *testing.T) {
// 	ctx := NewContextWithOptions(nil)
// 	assert.NotNil(t, ctx)

// 	var path string

// 	handler := func(ctx HTTPContext) HTTPError {
// 		bar, ok := ctx.Get("foo")
// 		assert.Equal(t, bar, "bar")
// 		assert.True(t, ok)

// 		bar, ok = ctx.MustGet("foo").(string)
// 		assert.Equal(t, bar, "bar")
// 		assert.True(t, ok)

// 		bar = ctx.GetString("foo")
// 		assert.Equal(t, bar, "bar")

// 		getInt := ctx.GetInt("int")
// 		assert.Equal(t, getInt, 2046)

// 		getInt64 := ctx.GetInt64("int64")
// 		assert.Equal(t, getInt64, int64(2046))

// 		getUint := ctx.GetUint("uint")
// 		assert.Equal(t, getUint, uint(2046))

// 		getUint64 := ctx.GetUint64("uint64")
// 		assert.Equal(t, getUint64, uint64(2046))

// 		getFloat64 := ctx.GetFloat64("float64")
// 		assert.Equal(t, getFloat64, float64(3.14))

// 		getBool := ctx.GetBool("bool")
// 		assert.True(t, getBool)

// 		getTime := ctx.GetTime("time")
// 		assert.Equal(t, getTime, time.Time{})

// 		getDuration := ctx.GetDuration("duration")
// 		assert.Equal(t, getDuration, time.Second)

// 		getStringSlice := ctx.GetStringSlice("ss")
// 		assert.Len(t, getStringSlice, 2)
// 		assert.Equal(t, getStringSlice[0], "foobar")
// 		assert.Equal(t, getStringSlice[1], "fozbaz")

// 		getStringMap := ctx.GetStringMap("sm")
// 		assert.Len(t, getStringMap, 2)
// 		assert.Equal(t, getStringMap["one"], 1)
// 		assert.Equal(t, getStringMap["true"], true)

// 		getStringMapString := ctx.GetStringMapString("sms")
// 		assert.Len(t, getStringMapString, 2)
// 		assert.Equal(t, getStringMapString["foo"], "bar")
// 		assert.Equal(t, getStringMapString["foz"], "baz")

// 		getStringMapStringSlice := ctx.GetStringMapStringSlice("smss")
// 		assert.Len(t, getStringMapStringSlice, 2)
// 		assert.Equal(t, getStringMapStringSlice["foo"][0], "bar")
// 		assert.Equal(t, getStringMapStringSlice["foz"][0], "baz")

// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}
// 	middleware := func(ctx HTTPContext) HTTPError {
// 		ctx.Set("foo", "bar")
// 		ctx.Set("int", 2046)
// 		ctx.Set("int64", int64(2046))
// 		ctx.Set("uint", uint(2046))
// 		ctx.Set("uint64", uint64(2046))
// 		ctx.Set("float64", float64(3.14))
// 		ctx.Set("bool", true)
// 		ctx.Set("time", time.Time{})
// 		ctx.Set("duration", time.Second)
// 		ctx.Set("ss", []string{"foobar", "fozbaz"})
// 		ctx.Set("sm", map[string]interface{}{"one": 1, "true": true})
// 		ctx.Set("sms", map[string]string{"foo": "bar", "foz": "baz"})
// 		ctx.Set("smss", map[string][]string{"foo": {"bar"}, "foz": {"baz"}})
// 		return nil
// 	}
// 	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
// 	assert.NotNil(t, ginHandlerFunc)

// 	var err error
// 	path, err = createTestServer(http.MethodGet, ginHandlerFunc, NewGinHandlerFunc(ctx, middleware))
// 	assert.NoError(t, err)

// 	var res *http.Response
// 	fn := func() bool {
// 		client := &http.Client{}
// 		req, err := http.NewRequest(http.MethodGet, path, nil)
// 		if err != nil {
// 			return false
// 		}
// 		res, err = client.Do(req)
// 		return err == nil
// 	}

// 	utils.RunUntil(fn, time.Second*4)
// 	assert.NotNil(t, res)
// }

// func TestParam(t *testing.T) {
// 	ctx := NewContextWithOptions(nil)
// 	assert.NotNil(t, ctx)

// 	var path string
// 	handler := func(ctx HTTPContext) HTTPError {
// 		param := ctx.Param("param")
// 		assert.Equal(t, param, "param")
// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}
// 	middleware := func(ctx HTTPContext) HTTPError {
// 		return nil
// 	}
// 	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
// 	assert.NotNil(t, ginHandlerFunc)

// 	var err error
// 	path, err = getCallingPathWithParamterAndQuery(ginHandlerFunc, "param", map[string]string{}, NewGinHandlerFunc(ctx, middleware))
// 	assert.NoError(t, err)

// 	var res *http.Response
// 	fn := func() bool {
// 		client := &http.Client{}
// 		req, err := http.NewRequest(http.MethodGet, path, nil)
// 		if err != nil {
// 			return false
// 		}
// 		res, err = client.Do(req)

// 		if res != nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, res.StatusCode, http.StatusNoContent)
// 		}

// 		return err == nil
// 	}

// 	utils.RunUntil(fn, time.Second*4)
// 	assert.NotNil(t, res)
// }

// func TestQuery(t *testing.T) {
// 	ctx := NewContextWithOptions(nil)
// 	assert.NotNil(t, ctx)

// 	var path string
// 	handler := func(ctx HTTPContext) HTTPError {
// 		value := ctx.Query("key")
// 		assert.Equal(t, value, "value")

// 		bar := ctx.Query("foo")
// 		assert.Equal(t, bar, "bar")

// 		blank, ok := ctx.GetQuery("blank")
// 		assert.False(t, ok)
// 		assert.Empty(t, blank)

// 		defaultQuery := ctx.DefaultQuery("blank", "default")
// 		assert.Equal(t, defaultQuery, "default")

// 		queryArray := ctx.QueryArray("arr")
// 		assert.Len(t, queryArray, 2)
// 		_, containFoobar := utils.ArrayStringContians(queryArray, "foobar")
// 		_, containFozbaz := utils.ArrayStringContians(queryArray, "fozbaz")
// 		assert.True(t, containFoobar)
// 		assert.True(t, containFozbaz)

// 		queryArray, ok = ctx.GetQueryArray("arr")
// 		assert.Equal(t, queryArray, ctx.QueryArray("arr"))
// 		assert.True(t, ok)

// 		queryMap := ctx.QueryMap("map")
// 		assert.Len(t, queryMap, 2)
// 		assert.Equal(t, queryMap["foo"], "bar")
// 		assert.Equal(t, queryMap["foz"], "baz")

// 		queryMap, ok = ctx.GetQueryMap("map")
// 		assert.Equal(t, queryMap, ctx.QueryMap("map"))
// 		assert.True(t, ok)

// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}
// 	middleware := func(ctx HTTPContext) HTTPError {
// 		return nil
// 	}
// 	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
// 	assert.NotNil(t, ginHandlerFunc)

// 	path, err := getCallingPathWithParamterAndQuery(ginHandlerFunc, "", map[string]string{
// 		"key":      "value",
// 		"foo":      "bar",
// 		"arr{0}":   "foobar",
// 		"arr{1}":   "fozbaz",
// 		"map[foo]": "bar",
// 		"map[foz]": "baz",
// 	}, NewGinHandlerFunc(ctx, middleware))
// 	assert.NoError(t, err)

// 	var res *http.Response
// 	fn := func() bool {
// 		client := &http.Client{}
// 		req, err := http.NewRequest(http.MethodGet, path, nil)
// 		if err != nil {
// 			return false
// 		}
// 		res, err = client.Do(req)

// 		if res != nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, res.StatusCode, http.StatusNoContent)
// 		}

// 		return err == nil
// 	}

// 	utils.RunUntil(fn, time.Second*4)
// 	assert.NotNil(t, res)
// }

// func TestURLEncodedForm(t *testing.T) {
// 	ctx := NewContextWithOptions(nil)
// 	assert.NotNil(t, ctx)

// 	var path string
// 	handler := func(ctx HTTPContext) HTTPError {
// 		contentType := ctx.ContentType()
// 		assert.Equal(t, contentType, gin.MIMEPOSTForm)

// 		bar := ctx.PostForm("foo")
// 		assert.Equal(t, bar, "bar")

// 		baz := ctx.DefaultPostForm("foz", "baz")
// 		assert.Equal(t, baz, "baz")

// 		bar, fooExist := ctx.GetPostForm("foo")
// 		assert.True(t, fooExist)
// 		assert.Equal(t, bar, "bar")

// 		baz, fozExist := ctx.GetPostForm("foz")
// 		assert.False(t, fozExist)
// 		assert.Empty(t, baz)

// 		arr := ctx.PostFormArray("arr")
// 		assert.Len(t, arr, 2)
// 		_, contianFoobar := utils.ArrayStringContians(arr, "foobar")
// 		assert.True(t, contianFoobar)
// 		_, contianFozbaz := utils.ArrayStringContians(arr, "fozbaz")
// 		assert.True(t, contianFozbaz)

// 		getArr, ok := ctx.GetPostFormArray("arr")
// 		assert.True(t, ok)
// 		assert.Equal(t, arr, getArr)

// 		formMap := ctx.PostFormMap("map")
// 		assert.Len(t, formMap, 2)
// 		assert.Equal(t, formMap["foo"], "bar")
// 		assert.Equal(t, formMap["foz"], "baz")

// 		getFormMap, ok := ctx.GetPostFormMap("map")
// 		assert.True(t, ok)
// 		assert.Equal(t, formMap, getFormMap)

// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}
// 	middleware := func(ctx HTTPContext) HTTPError {
// 		return nil
// 	}
// 	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
// 	ginMiddleware := NewGinHandlerFunc(ctx, middleware)
// 	assert.NotNil(t, ginHandlerFunc)
// 	path, err := createTestServer(http.MethodPost, ginHandlerFunc, ginMiddleware)
// 	assert.NoError(t, err)

// 	var res *http.Response
// 	fn := func() bool {
// 		client := &http.Client{}
// 		data := &url.Values{}
// 		data.Add("foo", "bar")
// 		data.Add("arr{0}", "foobar")
// 		data.Add("arr{1}", "fozbaz")
// 		data.Add("map[foo]", "bar")
// 		data.Add("map[foz]", "baz")
// 		re, _ := regexp.Compile(`\{(.*?)\}`)
// 		decodedValue, _ := url.QueryUnescape(data.Encode())
// 		formData := re.ReplaceAllString(decodedValue, "")
// 		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(formData))
// 		if err != nil {
// 			return false
// 		}
// 		req.Header.Add(string(consts.ContentType), gin.MIMEPOSTForm)

// 		res, err = client.Do(req)
// 		if res != nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, res.StatusCode, http.StatusNoContent)
// 		}

// 		return err == nil
// 	}

// 	utils.RunUntil(fn, time.Second*4)
// 	assert.NotNil(t, res)
// }

// func TestMultipartForm(t *testing.T) {
// 	ctx := NewContextWithOptions(nil)
// 	assert.NotNil(t, ctx)

// 	fileContent := utils.NewEncoder(nil).Random(8).Bytes()
// 	file := utils.NewFile("file", fileContent)
// 	var path string
// 	handler := func(ctx HTTPContext) HTTPError {
// 		contentType := ctx.ContentType()
// 		assert.Equal(t, contentType, gin.MIMEMultipartPOSTForm)

// 		multipartFormFileHeader, err := ctx.FormFile(file.Name())
// 		assert.NotEmpty(t, multipartFormFileHeader)
// 		assert.NoError(t, err)

// 		multipartFormFile, err := multipartFormFileHeader.Open()
// 		assert.NotEmpty(t, multipartFormFile)
// 		assert.NoError(t, err)

// 		buffer := bytes.NewBuffer(make([]byte, 0))
// 		_, err = io.Copy(buffer, multipartFormFile)
// 		assert.NoError(t, err)

// 		err = multipartFormFile.Close()
// 		assert.NoError(t, err)

// 		assert.Equal(t, buffer.Bytes(), file.Value())

// 		multipartForm, err := ctx.MultipartForm()
// 		assert.NoError(t, err)
// 		fileHeader := multipartForm.File[file.Name()][0]
// 		assert.Equal(t, multipartFormFileHeader, fileHeader)

// 		savedPath := fmt.Sprintf("/tmp/%s", uuid.New().String())
// 		err = ctx.SaveUploadedFile(fileHeader, savedPath)
// 		assert.NoError(t, err)

// 		openFile, err := os.Open(savedPath)
// 		assert.NoError(t, err)

// 		savedFile := bytes.NewBuffer(make([]byte, 0))
// 		_, err = io.Copy(savedFile, openFile)
// 		openFile.Close()
// 		assert.NoError(t, err)
// 		assert.Equal(t, savedFile.Bytes(), file.Value())

// 		err = os.Remove(savedPath)
// 		assert.NoError(t, err)

// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}
// 	middleware := func(ctx HTTPContext) HTTPError {
// 		return nil
// 	}
// 	ginHandlerFunc := NewGinHandlerFunc(ctx, handler)
// 	ginMiddleware := NewGinHandlerFunc(ctx, middleware)
// 	assert.NotNil(t, ginHandlerFunc)
// 	path, err := createTestServer(http.MethodPost, ginHandlerFunc, ginMiddleware)
// 	assert.NoError(t, err)

// 	var res *http.Response
// 	fn := func() bool {
// 		client := &http.Client{}
// 		buffer := bytes.NewBuffer([]byte{})
// 		form := multipart.NewWriter(buffer)

// 		part, err := form.CreateFormFile(file.Name(), "file")
// 		assert.NoError(t, err)

// 		_, err = io.Copy(part, file.Reader())
// 		assert.NoError(t, err)

// 		err = form.Close()
// 		assert.NoError(t, err)

// 		req, err := http.NewRequest(http.MethodPost, path, buffer)
// 		if err != nil {
// 			return false
// 		}
// 		req.Header.Add(string(consts.ContentType), form.FormDataContentType())

// 		res, err = client.Do(req)
// 		if res != nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, res.StatusCode, http.StatusNoContent)
// 		}

// 		return err == nil
// 	}

// 	utils.RunUntil(fn, time.Second*4)
// 	assert.NotNil(t, res)
// }

// func TestBinder(t *testing.T) {
// 	ctx := NewContextWithOptions(nil)
// 	assert.NotNil(t, ctx)

// 	middleware := func(ctx HTTPContext) HTTPError {
// 		return nil
// 	}
// 	ginMiddleware := NewGinHandlerFunc(ctx, middleware)

// 	// TestBind
// 	handlerBind := func(ctx HTTPContext) HTTPError {
// 		m := make(map[string]interface{}, 0)
// 		err := ctx.Bind(&m)
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, m)
// 		assert.Equal(t, m["foo"], "bar")
// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}
// 	ginHandlerFunc := NewGinHandlerFunc(ctx, handlerBind)
// 	assert.NotNil(t, ginHandlerFunc)

// 	path, err := createTestServer(http.MethodPost, ginHandlerFunc, ginMiddleware)
// 	assert.NoError(t, err)

// 	var res *http.Response
// 	jsonRequest := func() bool {
// 		client := &http.Client{}
// 		buffer := bytes.NewBuffer([]byte(`{"foo":"bar"}`))

// 		req, err := http.NewRequest(http.MethodPost, path, buffer)
// 		if err != nil {
// 			return false
// 		}
// 		req.Header.Set(string(consts.ContentType), gin.MIMEJSON)
// 		res, err = client.Do(req)
// 		if res != nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, res.StatusCode, http.StatusNoContent)
// 		}

// 		return err == nil
// 	}

// 	utils.RunUntil(jsonRequest, time.Second*4)
// 	assert.NotNil(t, res)

// 	// TestBindJSON
// 	handlerBindJSON := func(ctx HTTPContext) HTTPError {
// 		m := make(map[string]interface{}, 0)
// 		err := ctx.BindJSON(&m)
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, m)
// 		assert.Equal(t, m["foo"], "bar")
// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}

// 	ginHandlerFunc = NewGinHandlerFunc(ctx, handlerBindJSON)
// 	assert.NotNil(t, ginHandlerFunc)

// 	path, err = createTestServer(http.MethodPost, ginHandlerFunc, ginMiddleware)
// 	assert.NoError(t, err)

// 	utils.RunUntil(jsonRequest, time.Second*4)
// 	assert.NotNil(t, res)

// 	// TestBindXML
// 	handlerBindXML := func(ctx HTTPContext) HTTPError {
// 		m := &binderTest{}
// 		err := ctx.BindXML(m)
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, m)
// 		assert.Equal(t, m.Foo, "bar")
// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}

// 	xmlRequest := func() bool {
// 		client := &http.Client{}
// 		buffer := bytes.NewBuffer([]byte(`<?xml version="1.0" encoding="UTF-8"?><root><foo>bar</foo></root>`))

// 		req, err := http.NewRequest(http.MethodPost, path, buffer)
// 		if err != nil {
// 			return false
// 		}

// 		req.Header.Set(string(consts.ContentType), gin.MIMEXML)
// 		res, err = client.Do(req)
// 		if res != nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, res.StatusCode, http.StatusNoContent)
// 		}

// 		return err == nil
// 	}

// 	ginHandlerFunc = NewGinHandlerFunc(ctx, handlerBindXML)
// 	assert.NotNil(t, ginHandlerFunc)

// 	path, err = createTestServer(http.MethodPost, ginHandlerFunc, ginMiddleware)
// 	assert.NoError(t, err)

// 	utils.RunUntil(xmlRequest, time.Second*4)
// 	assert.NotNil(t, res)

// 	// TestBindQuery
// 	handlerBindQuery := func(ctx HTTPContext) HTTPError {
// 		m := &binderTest{}
// 		err := ctx.BindQuery(m)
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, m)
// 		assert.Equal(t, m.Foo, "bar")
// 		ctx.Status(http.StatusNoContent)
// 		return nil
// 	}

// 	ginHandlerFunc = NewGinHandlerFunc(ctx, handlerBindQuery)
// 	assert.NotNil(t, ginHandlerFunc)

// 	path, err = getCallingPathWithParamterAndQuery(ginHandlerFunc, "", map[string]string{
// 		"foo": "bar",
// 	}, ginMiddleware)
// 	assert.NoError(t, err)

// 	queryRequest := func() bool {
// 		client := &http.Client{}

// 		req, err := http.NewRequest(http.MethodGet, path, nil)
// 		if err != nil {
// 			return false
// 		}

// 		res, err = client.Do(req)
// 		if res != nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, res.StatusCode, http.StatusNoContent)
// 		}

// 		return err == nil
// 	}

// 	utils.RunUntil(queryRequest, time.Second*4)
// 	assert.NotNil(t, res)
// }

// func getCallingPathWithParamterAndQuery(handler gin.HandlerFunc,
// 	param string,
// 	query map[string]string,
// 	middlewares ...gin.HandlerFunc) (string, error) {
// 	path := uuid.New().String()
// 	g := gin.New()
// 	gin.SetMode(gin.TestMode)
// 	if len(middlewares) > 0 {
// 		for _, middleware := range middlewares {
// 			g.Use(middleware)
// 		}
// 	}
// 	endpoint := fmt.Sprintf("/%s", path)
// 	if param != "" {
// 		endpoint = fmt.Sprintf("%s/:param", endpoint)
// 	}
// 	g.GET(endpoint, handler)
// 	var err error
// 	avaiablePort := fmt.Sprint(DEFAULT_TEST_PORT)
// 	DEFAULT_TEST_PORT++
// 	go func() { err = g.Run(fmt.Sprintf(":%s", avaiablePort)) }()
// 	url := fmt.Sprintf("http://0.0.0.0:%s/%s/%s", avaiablePort, path, param)
// 	if param == "" {
// 		url = strings.TrimSuffix(url, "/")
// 	}
// 	if len(query) > 0 {
// 		params := "?"
// 		for key, value := range query {
// 			params = fmt.Sprintf("%s%s=%s&", params, key, value)
// 		}
// 		params = strings.TrimSuffix(params, "&")
// 		re, _ := regexp.Compile(`\{(.*?)\}`)
// 		params = re.ReplaceAllString(params, "")
// 		url = fmt.Sprintf("%s%s", url, params)
// 	}
// 	return url, err
// }
