package gwf

import (
	"net/http"
)

type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request

	// 从uri提取出来的参数
	params map[string]string
}
