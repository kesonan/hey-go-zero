# 自定义错误
在很多情况下，为了更好的将错误信息描述和特殊的业务错误码传递给请求方，我们都会增加自定义错误，这样不仅可以定义任何
业务错误码，而且对按照统一格式响应对请求方进行封装处理响应体也比较友好。

# 定义错误
增加一些自定义错误类型
在项目工程下，创建一个`common/errorx`文件夹，创建`codeerror.go`文件，添加如下代码：

``` go
package errorx

import (
	"fmt"
	"net/http"
)

const defaultCode = -1

type Handler struct{}

type ErrorBody struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func (h *Handler) Handle() func(error) (int, interface{}) {
	return func(err error) (int, interface{}) {
		switch v := err.(type) {
		case *CodeError:
			return http.StatusNotAcceptable, ErrorBody{
				Code: v.code,
				Desc: v.desc,
			}
		case *DescriptionError:
			return http.StatusNotAcceptable, ErrorBody{
				Code: defaultCode,
				Desc: v.desc,
			}
		case *InvalidParameterError:
			return http.StatusNotAcceptable, ErrorBody{
				Code: defaultCode,
				Desc: fmt.Sprintf("参数错误: %v", v.parameter),
			}
		default:
			return http.StatusInternalServerError, ErrorBody{
				Code: defaultCode,
				Desc: v.Error(),
			}
		}
	}
}

type CodeError struct {
	code int
	desc string
}

func NewCodeError(code int, desc string) *CodeError {
	return &CodeError{
		code: code,
		desc: desc,
	}
}

func (e *CodeError) Error() string {
	return e.desc
}

type DescriptionError struct {
	desc string
}

func NewDescriptionError(desc string) *DescriptionError {
	return &DescriptionError{
		desc: desc,
	}
}

func (e *DescriptionError) Error() string {
	return e.desc
}

type InvalidParameterError struct {
	parameter string
}

func NewInvalidParameterError(parameter string) *InvalidParameterError {
	return &InvalidParameterError{
		parameter: parameter,
	}
}

func (e *InvalidParameterError) Error() string {
	return e.parameter
}
```

# 设置自定义错误处理函数
在api的main函数中添加两行内容
``` go
errHandler := errorx.Handler{}
httpx.SetErrorHandler(errHandler.Handle())
```

完整代码
``` go
func main() {
    flag.Parse()
    
    var c config.Config
    conf.MustLoad(*configFile, &c)
    
    ctx := svc.NewServiceContext(c)
    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()
    
    errHandler := errorx.Handler{}
    httpx.SetErrorHandler(errHandler.Handle())
    
    handler.RegisterHandlers(server, ctx)
    
    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}
```