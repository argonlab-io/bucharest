package bucharest

import (
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
)

type HTTPContext interface {
	Context
	// Handler info
	HandlerName() string
	HandlerNames() []string

	// request
	FullPath() string
	ContentType() string
	Cookie(name string) (string, error)
	ClientIP() string
	RemoteIP() (net.IP, bool)
	GetHeader(key string) string
	GetRawData() ([]byte, error)
	IsWebsocket() bool

	// handler control
	Next()
	IsAborted() bool
	Abort()
	AbortWithStatusJSON(code int, jsonObj interface{})
	AbortWithStatus(code int)

	// Setter and Getter
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	MustGet(key string) interface{}
	GetString(key string) (s string)
	GetBool(key string) (b bool)
	GetInt(key string) (i int)
	GetInt64(key string) (i64 int64)
	GetUint(key string) (ui uint)
	GetUint64(key string) (ui64 uint64)
	GetFloat64(key string) (f64 float64)
	GetTime(key string) (t time.Time)
	GetDuration(key string) (d time.Duration)
	GetStringSlice(key string) (ss []string)
	GetStringMap(key string) (sm map[string]interface{})
	GetStringMapString(key string) (sms map[string]string)
	GetStringMapStringSlice(key string) (smss map[string][]string)

	// parameter and query
	Param(key string) string
	Query(key string) string
	DefaultQuery(key, defaultValue string) string
	GetQuery(key string) (string, bool)
	QueryArray(key string) []string
	GetQueryArray(key string) ([]string, bool)
	QueryMap(key string) map[string]string
	GetQueryMap(key string) (map[string]string, bool)

	// urlencoded
	PostForm(key string) string
	DefaultPostForm(key, defaultValue string) string
	GetPostForm(key string) (string, bool)
	PostFormArray(key string) []string
	GetPostFormArray(key string) ([]string, bool)
	PostFormMap(key string) map[string]string
	GetPostFormMap(key string) (map[string]string, bool)

	// multipart
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)
	SaveUploadedFile(file *multipart.FileHeader, dst string) error

	// binder
	Bind(obj interface{}) error
	BindJSON(obj interface{}) error
	BindXML(obj interface{}) error
	BindQuery(obj interface{}) error
	BindYAML(obj interface{}) error
	BindHeader(obj interface{}) error
	BindUri(obj interface{}) error
	MustBindWith(obj interface{}, b binding.Binding) error
	ShouldBind(obj interface{}) error
	ShouldBindJSON(obj interface{}) error
	ShouldBindXML(obj interface{}) error
	ShouldBindQuery(obj interface{}) error
	ShouldBindYAML(obj interface{}) error
	ShouldBindHeader(obj interface{}) error
	ShouldBindUri(obj interface{}) error
	ShouldBindWith(obj interface{}, b binding.Binding) error
	ShouldBindBodyWith(obj interface{}, bb binding.BindingBody) (err error)

	// response header
	Status(code int)
	Header(key, value string)
	SetSameSite(samesite http.SameSite)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	SetAccepted(formats ...string)

	// response body
	Render(code int, r render.Render)
	HTML(code int, name string, obj interface{})
	IndentedJSON(code int, obj interface{})
	SecureJSON(code int, obj interface{})
	JSONP(code int, obj interface{})
	JSON(code int, obj interface{})
	AsciiJSON(code int, obj interface{})
	PureJSON(code int, obj interface{})
	XML(code int, obj interface{})
	YAML(code int, obj interface{})
	ProtoBuf(code int, obj interface{})
	String(code int, format string, values ...interface{})
	Redirect(code int, location string)
	Data(code int, contentType string, data []byte)
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
	File(filepath string)
	FileFromFS(filepath string, fs http.FileSystem)
	FileAttachment(filepath, filename string)
	SSEvent(name string, message interface{})
	Stream(step func(w io.Writer) bool) bool

	originalContext() interface{}
	GetGin() (*gin.Context, bool)
	Gin() *gin.Context
}

type HTTPError interface {
	OriginalError() error
	GetStatus() int
	GetJSON() interface{}
}

type HandlerFunc func(HTTPContext) HTTPError
type HandlerFuncWithData func(HTTPContext, map[string]any) HTTPError
type HandlersChain []HandlerFunc
