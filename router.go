package gwf

import (
	"fmt"
	"net/http"
	"regexp"
)

var defaultRouter *DefaultRouter

func init() {
	defaultRouter = &DefaultRouter{make(map[*regexp.Regexp]ControllerInterface)}
}

// 路由器
type DefaultRouter struct {
	routerMap map[*regexp.Regexp]ControllerInterface
}

// 实现http.Handler接口
func (dr *DefaultRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	matched := false
	for re, value := range dr.routerMap {
		if re.MatchString(path) {

			subexpNames := re.SubexpNames()
			matches := re.FindAllStringSubmatch(path, -1)

			params := make(map[string]string)
			for i, v := range subexpNames {
				if v == "" {
					continue
				}
				params[v] = matches[0][i]
				fmt.Println(i, v, v == "", matches[0][i])
			}

			context := &Context{w, r, params}
			value.Init(context)

			switch method := r.Method; method {
			case "GET":
				value.Get()
			case "POST":
				value.Post()
			case "PUT":
				value.Put()
			default:
				value.Get()
			}
			matched = true
			break
		}
	}

	// 没有对应的路由设置
	if !matched {
		http404(w, r)
	}
}

// 添加一个路由
func (dr *DefaultRouter) add(path string, controller ControllerInterface) {
	re := regexp.MustCompile(path)
	dr.routerMap[re] = controller
}

// 添加一个路由
func AddRouter(path string, controller ControllerInterface) {
	defaultRouter.add(path, controller)
}

// 返回默认的路由
func GetRouter() *DefaultRouter {
	return defaultRouter
}
