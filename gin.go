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
		httpError := handlerFunc(&httpContextWithGin{Context: ctx, gin: g})
		if httpError != nil {
			g.JSON(httpError.GetStatus(), httpError.GetJSON())
		}
	}
}

func NewGinHandlerFuncWithData(ctx Context, handlerFunc HandlerFuncWithData, data Map) gin.HandlerFunc {
	return func(g *gin.Context) {
		httpError := handlerFunc(&httpContextWithGin{Context: ctx, gin: g}, data)
		if httpError != nil {
			g.JSON(httpError.GetStatus(), httpError.GetJSON())
		}
	}
}

type httpContextWithGin struct {
	Context
	gin *gin.Context
}

func (h *httpContextWithGin) HandlerName() string {
	return h.gin.HandlerName()
}

func (h *httpContextWithGin) HandlerNames() []string {
	return h.gin.HandlerNames()
}

func (h *httpContextWithGin) FullPath() string {
	return h.gin.FullPath()
}

func (h *httpContextWithGin) Next() {
	h.gin.Next()
}

func (h *httpContextWithGin) IsAborted() bool {
	return h.gin.IsAborted()
}

func (h *httpContextWithGin) Abort() {
	h.gin.Abort()
}

func (h *httpContextWithGin) AbortWithStatusJSON(code int, jsonObj interface{}) {
	h.gin.AbortWithStatusJSON(code, jsonObj)
}

func (h *httpContextWithGin) AbortWithStatus(code int) {
	h.gin.AbortWithStatus(code)
}

func (h *httpContextWithGin) Get(key string) (interface{}, bool) {
	return h.gin.Get(key)
}

func (h *httpContextWithGin) Set(key string, value interface{}) {
	h.gin.Set(key, value)
}

func (h *httpContextWithGin) MustGet(key string) interface{} {
	return h.gin.MustGet(key)
}

func (h *httpContextWithGin) GetString(key string) (s string) {
	return h.gin.GetString(key)
}

func (h *httpContextWithGin) GetBool(key string) (b bool) {
	return h.gin.GetBool(key)
}

func (h *httpContextWithGin) GetInt(key string) (i int) {
	return h.gin.GetInt(key)
}

func (h *httpContextWithGin) GetInt64(key string) (i64 int64) {
	return h.gin.GetInt64(key)
}

func (h *httpContextWithGin) GetUint(key string) (ui uint) {
	return h.gin.GetUint(key)
}

func (h *httpContextWithGin) GetUint64(key string) (ui64 uint64) {
	return h.gin.GetUint64(key)
}

func (h *httpContextWithGin) GetFloat64(key string) (f64 float64) {
	return h.gin.GetFloat64(key)
}

func (h *httpContextWithGin) GetTime(key string) (t time.Time) {
	return h.gin.GetTime(key)
}

func (h *httpContextWithGin) GetDuration(key string) (d time.Duration) {
	return h.gin.GetDuration(key)
}

func (h *httpContextWithGin) GetStringSlice(key string) (ss []string) {
	return h.gin.GetStringSlice(key)
}

func (h *httpContextWithGin) GetStringMap(key string) (sm map[string]interface{}) {
	return h.gin.GetStringMap(key)
}

func (h *httpContextWithGin) GetStringMapString(key string) (sms map[string]string) {
	return h.gin.GetStringMapString(key)
}

func (h *httpContextWithGin) GetStringMapStringSlice(key string) (smss map[string][]string) {
	return h.gin.GetStringMapStringSlice(key)
}

func (h *httpContextWithGin) Param(key string) string {
	return h.gin.Param(key)
}

func (h *httpContextWithGin) Query(key string) string {
	return h.gin.Query(key)
}

func (h *httpContextWithGin) DefaultQuery(key, defaultValue string) string {
	return h.gin.DefaultQuery(key, defaultValue)
}

func (h *httpContextWithGin) GetQuery(key string) (string, bool) {
	return h.gin.GetQuery(key)
}

func (h *httpContextWithGin) QueryArray(key string) []string {
	return h.gin.QueryArray(key)
}

func (h *httpContextWithGin) GetQueryArray(key string) ([]string, bool) {
	return h.gin.GetQueryArray(key)
}

func (h *httpContextWithGin) QueryMap(key string) map[string]string {
	return h.gin.QueryMap(key)
}

func (h *httpContextWithGin) GetQueryMap(key string) (map[string]string, bool) {
	return h.gin.GetQueryMap(key)
}

func (h *httpContextWithGin) PostForm(key string) string {
	return h.gin.PostForm(key)
}

func (h *httpContextWithGin) DefaultPostForm(key, defaultValue string) string {
	return h.gin.DefaultPostForm(key, defaultValue)
}

func (h *httpContextWithGin) GetPostForm(key string) (string, bool) {
	return h.gin.GetPostForm(key)
}

func (h *httpContextWithGin) PostFormArray(key string) []string {
	return h.gin.PostFormArray(key)
}

func (h *httpContextWithGin) GetPostFormArray(key string) ([]string, bool) {
	return h.gin.GetPostFormArray(key)
}

func (h *httpContextWithGin) PostFormMap(key string) map[string]string {
	return h.gin.PostFormMap(key)
}

func (h *httpContextWithGin) GetPostFormMap(key string) (map[string]string, bool) {
	return h.gin.GetPostFormMap(key)
}

func (h *httpContextWithGin) FormFile(name string) (*multipart.FileHeader, error) {
	return h.gin.FormFile(name)
}

func (h *httpContextWithGin) MultipartForm() (*multipart.Form, error) {
	return h.gin.MultipartForm()
}

func (h *httpContextWithGin) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	return h.gin.SaveUploadedFile(file, dst)
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

func (h *httpContextWithGin) ClientIP() string {
	return h.gin.ClientIP()
}

func (h *httpContextWithGin) RemoteIP() (net.IP, bool) {
	return h.gin.RemoteIP()
}

func (h *httpContextWithGin) ContentType() string {
	return h.gin.ContentType()
}

func (h *httpContextWithGin) IsWebsocket() bool {
	return h.gin.IsWebsocket()
}

func (h *httpContextWithGin) Status(code int) {
	h.gin.Status(code)
}

func (h *httpContextWithGin) Header(key, value string) {
	h.gin.Header(key, value)
}

func (h *httpContextWithGin) GetHeader(key string) string {
	return h.gin.GetHeader(key)
}

func (h *httpContextWithGin) GetRawData() ([]byte, error) {
	return h.gin.GetRawData()
}

func (h *httpContextWithGin) SetSameSite(samesite http.SameSite) {
	h.gin.SetSameSite(samesite)
}

func (h *httpContextWithGin) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	h.gin.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (h *httpContextWithGin) Cookie(name string) (string, error) {
	return h.gin.Cookie(name)
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
