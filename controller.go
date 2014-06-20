package gwf

import (
	"net/http"
)

// 控制器接口，统一控制器的操作
type ControllerInterface interface {
	Init(*Context)
	Get()
	Post()
	Put()
}

// 控制器
type Controller struct {
	context        *Context
	responseWriter http.ResponseWriter
	request        *http.Request
}

func (c *Controller) Init(context *Context) {
	c.context = context
}
func (c *Controller) Get() {
	http404(c.context.responseWriter, c.context.request)
}
func (c *Controller) Post() {
	http404(c.context.responseWriter, c.context.request)
}
func (c *Controller) Put() {
	http404(c.context.responseWriter, c.context.request)
}

// 获取路由器路径匹配到的参数
func (c *Controller) GetParam(name string) string {
	ret := ""
	if _, ok := c.context.params[name]; ok {
		ret = c.context.params[name]
	}
	return ret
}

func (c *Controller) GetResponseWriter() http.ResponseWriter {
	return c.context.responseWriter
}

func (c *Controller) GetRequest() *http.Request {
	return c.context.request
}

func http404(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 Not Found"))
}

func HTTP404(w http.ResponseWriter, r *http.Request) {
	http404(w, r)
}
