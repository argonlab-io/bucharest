package bucharest

import (
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"
)

func NewGinHandlerFunc(ctx Context, handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(g *gin.Context) {
		httpError := handlerFunc(&httpContextWithGin{
			Context:             ctx,
			gin:                 g,
			ginHandlerInfo:      ginHandlerInfo{gin: g},
			ginRequest:          ginRequest{gin: g},
			ginHandlerControl:   ginHandlerControl{gin: g},
			ginSetterAndGetter:  ginSetterAndGetter{gin: g},
			ginParamterAndQuery: ginParamterAndQuery{gin: g},
			ginURLEncodedForm:   ginURLEncodedForm{gin: g},
			ginMultipartForm:    ginMultipartForm{gin: g},
		})
		if httpError != nil {
			g.JSON(httpError.GetStatus(), httpError.GetJSON())
		}
	}
}

func NewGinHandlerFuncWithData(ctx Context, handlerFunc HandlerFuncWithData, data Map) gin.HandlerFunc {
	return func(g *gin.Context) {
		httpError := handlerFunc(&httpContextWithGin{
			Context:             ctx,
			gin:                 g,
			ginHandlerInfo:      ginHandlerInfo{gin: g},
			ginRequest:          ginRequest{gin: g},
			ginHandlerControl:   ginHandlerControl{gin: g},
			ginSetterAndGetter:  ginSetterAndGetter{gin: g},
			ginParamterAndQuery: ginParamterAndQuery{gin: g},
			ginURLEncodedForm:   ginURLEncodedForm{gin: g},
			ginMultipartForm:    ginMultipartForm{gin: g},
		}, data)
		if httpError != nil {
			g.JSON(httpError.GetStatus(), httpError.GetJSON())
		}
	}
}

type httpContextWithGin struct {
	Context
	gin *gin.Context
	ginHandlerInfo
	ginRequest
	ginHandlerControl
	ginSetterAndGetter
	ginParamterAndQuery
	ginURLEncodedForm
	ginMultipartForm
}

type ginHandlerInfo struct {
	gin *gin.Context
}

func (hi *ginHandlerInfo) HandlerName() string {
	return hi.gin.HandlerName()
}

func (hi *ginHandlerInfo) HandlerNames() []string {
	return hi.gin.HandlerNames()
}

type ginRequest struct {
	gin *gin.Context
}

func (r *ginRequest) FullPath() string {
	return r.gin.FullPath()
}

func (r *ginRequest) ClientIP() string {
	return r.gin.ClientIP()
}

func (r *ginRequest) RemoteIP() (net.IP, bool) {
	return r.gin.RemoteIP()
}

func (r *ginRequest) Cookie(name string) (string, error) {
	return r.gin.Cookie(name)
}

func (r *ginRequest) ContentType() string {
	return r.gin.ContentType()
}

func (r *ginRequest) GetHeader(key string) string {
	return r.gin.GetHeader(key)
}

func (r *ginRequest) GetRawData() ([]byte, error) {
	return r.gin.GetRawData()
}

func (r *ginRequest) IsWebsocket() bool {
	return r.gin.IsWebsocket()
}

type ginHandlerControl struct {
	gin *gin.Context
}

func (hc *ginHandlerControl) Next() {
	hc.gin.Next()
}

func (hc *ginHandlerControl) IsAborted() bool {
	return hc.gin.IsAborted()
}

func (hc *ginHandlerControl) Abort() {
	hc.gin.Abort()
}

func (hc *ginHandlerControl) AbortWithStatusJSON(code int, jsonObj interface{}) {
	hc.gin.AbortWithStatusJSON(code, jsonObj)
}

func (hc *ginHandlerControl) AbortWithStatus(code int) {
	hc.gin.AbortWithStatus(code)
}

type ginSetterAndGetter struct {
	gin *gin.Context
}

func (sg *ginSetterAndGetter) Get(key string) (interface{}, bool) {
	return sg.gin.Get(key)
}

func (sg *ginSetterAndGetter) Set(key string, value interface{}) {
	sg.gin.Set(key, value)
}

func (sg *ginSetterAndGetter) MustGet(key string) interface{} {
	return sg.gin.MustGet(key)
}

func (sg *ginSetterAndGetter) GetString(key string) (s string) {
	return sg.gin.GetString(key)
}

func (sg *ginSetterAndGetter) GetBool(key string) (b bool) {
	return sg.gin.GetBool(key)
}

func (sg *ginSetterAndGetter) GetInt(key string) (i int) {
	return sg.gin.GetInt(key)
}

func (sg *ginSetterAndGetter) GetInt64(key string) (i64 int64) {
	return sg.gin.GetInt64(key)
}

func (sg *ginSetterAndGetter) GetUint(key string) (ui uint) {
	return sg.gin.GetUint(key)
}

func (sg *ginSetterAndGetter) GetUint64(key string) (ui64 uint64) {
	return sg.gin.GetUint64(key)
}

func (sg *ginSetterAndGetter) GetFloat64(key string) (f64 float64) {
	return sg.gin.GetFloat64(key)
}

func (sg *ginSetterAndGetter) GetTime(key string) (t time.Time) {
	return sg.gin.GetTime(key)
}

func (sg *ginSetterAndGetter) GetDuration(key string) (d time.Duration) {
	return sg.gin.GetDuration(key)
}

func (sg *ginSetterAndGetter) GetStringSlice(key string) (ss []string) {
	return sg.gin.GetStringSlice(key)
}

func (sg *ginSetterAndGetter) GetStringMap(key string) (sm map[string]interface{}) {
	return sg.gin.GetStringMap(key)
}

func (sg *ginSetterAndGetter) GetStringMapString(key string) (sms map[string]string) {
	return sg.gin.GetStringMapString(key)
}

func (sg *ginSetterAndGetter) GetStringMapStringSlice(key string) (smss map[string][]string) {
	return sg.gin.GetStringMapStringSlice(key)
}

type ginParamterAndQuery struct {
	gin *gin.Context
}

func (pq *ginParamterAndQuery) Param(key string) string {
	return pq.gin.Param(key)
}

func (pq *ginParamterAndQuery) Query(key string) string {
	return pq.gin.Query(key)
}

func (pq *ginParamterAndQuery) DefaultQuery(key, defaultValue string) string {
	return pq.gin.DefaultQuery(key, defaultValue)
}

func (pq *ginParamterAndQuery) GetQuery(key string) (string, bool) {
	return pq.gin.GetQuery(key)
}

func (pq *ginParamterAndQuery) QueryArray(key string) []string {
	return pq.gin.QueryArray(key)
}

func (pq *ginParamterAndQuery) GetQueryArray(key string) ([]string, bool) {
	return pq.gin.GetQueryArray(key)
}

func (pq *ginParamterAndQuery) QueryMap(key string) map[string]string {
	return pq.gin.QueryMap(key)
}

func (pq *ginParamterAndQuery) GetQueryMap(key string) (map[string]string, bool) {
	return pq.gin.GetQueryMap(key)
}

type ginURLEncodedForm struct {
	gin *gin.Context
}

func (uef *ginURLEncodedForm) PostForm(key string) string {
	return uef.gin.PostForm(key)
}

func (uef *ginURLEncodedForm) DefaultPostForm(key, defaultValue string) string {
	return uef.gin.DefaultPostForm(key, defaultValue)
}

func (uef *ginURLEncodedForm) GetPostForm(key string) (string, bool) {
	return uef.gin.GetPostForm(key)
}

func (uef *ginURLEncodedForm) PostFormArray(key string) []string {
	return uef.gin.PostFormArray(key)
}

func (uef *ginURLEncodedForm) GetPostFormArray(key string) ([]string, bool) {
	return uef.gin.GetPostFormArray(key)
}

func (uef *ginURLEncodedForm) PostFormMap(key string) map[string]string {
	return uef.gin.PostFormMap(key)
}

func (uef *ginURLEncodedForm) GetPostFormMap(key string) (map[string]string, bool) {
	return uef.gin.GetPostFormMap(key)
}

type ginMultipartForm struct {
	gin *gin.Context
}

func (mf *ginMultipartForm) FormFile(name string) (*multipart.FileHeader, error) {
	return mf.gin.FormFile(name)
}

func (mf *ginMultipartForm) MultipartForm() (*multipart.Form, error) {
	return mf.gin.MultipartForm()
}

func (mf *ginMultipartForm) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	return mf.gin.SaveUploadedFile(file, dst)
}

func (h *httpContextWithGin) Bind(obj interface{}) error {
	return h.gin.Bind(obj)
}

func (h *httpContextWithGin) BindJSON(obj interface{}) error {
	return h.gin.BindJSON(obj)
}

func (h *httpContextWithGin) BindXML(obj interface{}) error {
	return h.gin.BindXML(obj)
}

func (h *httpContextWithGin) BindQuery(obj interface{}) error {
	return h.gin.BindQuery(obj)
}

func (h *httpContextWithGin) BindYAML(obj interface{}) error {
	return h.gin.BindYAML(obj)
}

func (h *httpContextWithGin) BindHeader(obj interface{}) error {
	return h.gin.BindHeader(obj)
}

func (h *httpContextWithGin) BindUri(obj interface{}) error {
	return h.gin.BindUri(obj)
}

func (h *httpContextWithGin) MustBindWith(obj interface{}, b binding.Binding) error {
	return h.gin.MustBindWith(obj, b)
}

func (h *httpContextWithGin) ShouldBind(obj interface{}) error {
	return h.gin.ShouldBind(obj)
}

func (h *httpContextWithGin) ShouldBindJSON(obj interface{}) error {
	return h.gin.ShouldBindJSON(obj)
}

func (h *httpContextWithGin) ShouldBindXML(obj interface{}) error {
	return h.gin.ShouldBindXML(obj)
}

func (h *httpContextWithGin) ShouldBindQuery(obj interface{}) error {
	return h.gin.ShouldBindQuery(obj)
}

func (h *httpContextWithGin) ShouldBindYAML(obj interface{}) error {
	return h.gin.ShouldBindYAML(obj)
}

func (h *httpContextWithGin) ShouldBindHeader(obj interface{}) error {
	return h.gin.ShouldBindHeader(obj)
}

func (h *httpContextWithGin) ShouldBindUri(obj interface{}) error {
	return h.gin.ShouldBindUri(obj)
}

func (h *httpContextWithGin) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return h.gin.ShouldBindWith(obj, b)
}

func (h *httpContextWithGin) ShouldBindBodyWith(obj interface{}, bb binding.BindingBody) (err error) {
	return h.gin.ShouldBindBodyWith(obj, bb)
}

func (h *httpContextWithGin) Status(code int) {
	h.gin.Status(code)
}

func (h *httpContextWithGin) Header(key, value string) {
	h.gin.Header(key, value)
}

func (h *httpContextWithGin) SetSameSite(samesite http.SameSite) {
	h.gin.SetSameSite(samesite)
}

func (h *httpContextWithGin) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	h.gin.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (h *httpContextWithGin) Render(code int, r render.Render) {
	h.gin.Render(code, r)
}

func (h *httpContextWithGin) HTML(code int, name string, obj interface{}) {
	h.gin.HTML(code, name, obj)
}

func (h *httpContextWithGin) IndentedJSON(code int, obj interface{}) {
	h.gin.IndentedJSON(code, obj)
}

func (h *httpContextWithGin) SecureJSON(code int, obj interface{}) {
	h.gin.SecureJSON(code, obj)
}

func (h *httpContextWithGin) JSONP(code int, obj interface{}) {
	h.gin.SecureJSON(code, obj)
}

func (h *httpContextWithGin) JSON(code int, obj interface{}) {
	h.gin.JSON(code, obj)
}

func (h *httpContextWithGin) AsciiJSON(code int, obj interface{}) {
	h.gin.JSON(code, obj)
}

func (h *httpContextWithGin) PureJSON(code int, obj interface{}) {
	h.gin.PureJSON(code, obj)
}

func (h *httpContextWithGin) XML(code int, obj interface{}) {
	h.gin.XML(code, obj)
}

func (h *httpContextWithGin) YAML(code int, obj interface{}) {
	h.gin.YAML(code, obj)
}

func (h *httpContextWithGin) ProtoBuf(code int, obj interface{}) {
	h.gin.ProtoBuf(code, obj)
}

func (h *httpContextWithGin) String(code int, format string, values ...interface{}) {
	h.gin.String(code, format, values...)
}

func (h *httpContextWithGin) Redirect(code int, location string) {
	h.gin.Redirect(code, location)
}

func (h *httpContextWithGin) Data(code int, contentType string, data []byte) {
	h.gin.Data(code, contentType, data)
}

func (h *httpContextWithGin) DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	h.gin.DataFromReader(code, contentLength, contentType, reader, extraHeaders)
}

func (h *httpContextWithGin) File(filepath string) {
	h.gin.File(filepath)
}

func (h *httpContextWithGin) FileFromFS(filepath string, fs http.FileSystem) {
	h.gin.FileFromFS(filepath, fs)
}

func (h *httpContextWithGin) FileAttachment(filepath, filename string) {
	h.gin.FileAttachment(filepath, filename)
}

func (h *httpContextWithGin) SSEvent(name string, message interface{}) {
	h.gin.SSEvent(name, message)
}

func (h *httpContextWithGin) Stream(step func(w io.Writer) bool) bool {
	return h.gin.Stream(step)
}

func (h *httpContextWithGin) SetAccepted(formats ...string) {
	h.gin.SetAccepted(formats...)
}

func (h *httpContextWithGin) Deadline() (deadline time.Time, ok bool) {
	return h.Context.Deadline()
}

func (h *httpContextWithGin) Done() <-chan struct{} {
	return h.Context.Done()
}

func (h *httpContextWithGin) Err() error {
	return h.Context.Err()
}

func (h *httpContextWithGin) Value(key interface{}) interface{} {
	fromGin := h.gin.Value(key)
	if fromGin != nil {
		return fromGin
	}
	return h.Context.Value(key)
}

func (h *httpContextWithGin) originalContext() interface{} {
	return h.gin
}

func (h *httpContextWithGin) GetGin() (*gin.Context, bool) {
	ctx, ok := h.originalContext().(*gin.Context)
	if !ok {
		return nil, false
	}
	return ctx, true
}

func (h *httpContextWithGin) Gin() *gin.Context {
	ctx, ok := h.originalContext().(*gin.Context)
	if !ok {
		panic("Gin engine does not exist in this context")
	}
	return ctx
}

func GinLoggerWithConfig(ctx HTTPContext, data Map) HTTPError {
	var conf gin.LoggerConfig
	if data["conf"] == nil {
		conf = gin.LoggerConfig{}
	} else {
		conf = data["conf"].(gin.LoggerConfig)
	}

	formatter := conf.Formatter

	out := conf.Output
	if out == nil {
		out = os.Stdout
	}

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	// Start timer
	start := time.Now()
	path := ctx.Gin().Request.URL.Path
	raw := ctx.Gin().Request.URL.RawQuery

	// Process request
	ctx.Next()

	// Log only when path is not being skipped
	if _, ok := skip[path]; !ok {
		param := gin.LogFormatterParams{
			Request: ctx.Gin().Request,
			Keys:    ctx.Gin().Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = ctx.ClientIP()
		param.Method = ctx.Gin().Request.Method
		param.StatusCode = ctx.Gin().Writer.Status()
		param.ErrorMessage = ctx.Gin().Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = ctx.Gin().Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		ctx.Log().Out = out
		if formatter != nil {
			ctx.Log().Info(formatter(param))
		}

		var statusColor, methodColor, resetColor string
		if ctx.Log().Formatter.(*logrus.TextFormatter).ForceColors {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency - param.Latency%time.Second
		}

		logLevel := logrus.InfoLevel
		if param.StatusCode >= 500 {
			logLevel = logrus.ErrorLevel
		}
		ctx.Log().Logf(logLevel, "[GIN]| %s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}

	return nil
}
