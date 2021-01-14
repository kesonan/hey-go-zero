# FAQ集合
在前面已经把准备工作介绍得差不多了，如果你耐心的看到了这里，恭喜你，对后面我们演示项目的进入
会进展非常顺利，不至于在开发过程中发现这个工具没安装，那么环境该怎么设置，接下来我们通过QA的形式
将一些没有讲到的一些问题进行介绍

> 本章节会随着使用者发现的问题不断扩充，如果你后续在开发过程中遇到什么问题，请优先尝试来这里找找答案吧!

# 怎么反馈文档错误或者提建议？
> 通过在github提Issue进行问题反馈，在本项目中有三个特殊label，你可以选择对应标签给我们反馈错误或者建议。
> * `feature`新功能推荐
> * `doc`文档错误
> * `good idea`好的建议

# 怎么贡献`hey-go-zero`？
你可以fork [hey-go-zero](https://github.com/songmeizi/hey-go-zero) 然后修复或者添加新内容后进行pull request，
我们会认真阅读你的pull request并在内部团队达成一致后进行merge，这个过程也许会比较漫长，请你保持耐心，不管任何结果我们都会给到你反馈。

> 说明: 提交pull request建议对本次提交进行一个简要概括，这样可以方便我们快速了解你的想法。

# 怎么控制生成代码文件的命名风格？
在[《Goctl介绍》](./goctl-intro.md) 章节中我们介绍来`Goctl`工具的命令大全，在其中我们的api、rpc、model服务生成都有一个`--style`参数，
其代表一个命名风格格式化符，就像日期格式化符`yyy-MM-dd`类似，我们通过`gozero`单词来控制文件命名风格，如：
* `gozero`代表golang命名风格，即全小写
* `go_zero`代表snake命名风格

常见格式化符生成示例
源字符：welcome_to_go_zero

| 格式化符   | 格式化结果            | 说明                      |
|------------|-----------------------|---------------------------|
| gozero     | welcometogozero       | 小写                      |
| goZero     | welcomeToGoZero       | 驼峰                      |
| go_zero    | welcome_to_go_zero    | snake                     |
| Go#zero    | Welcome#to#go#zero    | #号分割Title类型          |
| GOZERO     | WELCOMETOGOZERO       | 大写                      |
| \_go#zero_ | \_welcome#to#go#zero_ | 下划线做前后缀，并且#分割 |

错误格式化符示例
* go
* gOZero
* zero
* goZEro
* goZERo
* goZeRo
* tal

### 使用方法

目前可通过在生成api、rpc、model时通过`--style`参数指定format格式，如：
``` shell script
goctl api go test.api -dir . -style gozero
```
``` shell script
 goctl rpc proto -src test.proto -dir . -style go_zero
```
``` shell script
goctl model mysql datasource -url="" -table="*" -dir ./snake -style GoZero
```

### 默认值
当不指定`--style`时默认值为`gozero`

# Rpc服务运行遇到错误

* 错误一:

  ``` golang
  pb/xx.pb.go:220:7: undefined: grpc.ClientConnInterface
  pb/xx.pb.go:224:11: undefined: grpc.SupportPackageIsVersion6
  pb/xx.pb.go:234:5: undefined: grpc.ClientConnInterface
  pb/xx.pb.go:237:24: undefined: grpc.ClientConnInterface
  ```

  解决方法：请将`protoc-gen-go`版本降至v1.3.2及一下，然后重新生成。

* 错误二:

  ``` golang

  # go.etcd.io/etcd/clientv3/balancer/picker
  ../../../go/pkg/mod/go.etcd.io/etcd@v0.0.0-20200402134248-51bdeb39e698/clientv3/balancer/picker/err.go:25:9: cannot use &errPicker literal (type *errPicker) as type Picker in return argument:*errPicker does not implement Picker (wrong type for Pick method)
    have Pick(context.Context, balancer.PickInfo) (balancer.SubConn, func(balancer.DoneInfo), error)
    want Pick(balancer.PickInfo) (balancer.PickResult, error)
    ../../../go/pkg/mod/go.etcd.io/etcd@v0.0.0-20200402134248-51bdeb39e698/clientv3/balancer/picker/roundrobin_balanced.go:33:9: cannot use &rrBalanced literal (type *rrBalanced) as type Picker in return argument:
    *rrBalanced does not implement Picker (wrong type for Pick method)
		have Pick(context.Context, balancer.PickInfo) (balancer.SubConn, func(balancer.DoneInfo), error)
    want Pick(balancer.PickInfo) (balancer.PickResult, error)
    #github.com/tal-tech/go-zero/zrpc/internal/balancer/p2c
    ../../../go/pkg/mod/github.com/tal-tech/go-zero@v1.0.12/zrpc/internal/balancer/p2c/p2c.go:41:32: not enough arguments in call to base.NewBalancerBuilder
	have (string, *p2cPickerBuilder)
  want (string, base.PickerBuilder, base.Config)
  ../../../go/pkg/mod/github.com/tal-tech/go-zero@v1.0.12/zrpc/internal/balancer/p2c/p2c.go:58:9: cannot use &p2cPicker literal (type *p2cPicker) as type balancer.Picker in return argument:
	*p2cPicker does not implement balancer.Picker (wrong type for Pick method)
		have Pick(context.Context, balancer.PickInfo) (balancer.SubConn, func(balancer.DoneInfo), error)
		want Pick(balancer.PickInfo) (balancer.PickResult, error)
  ```

  解决方法：
  
    ``` golang
    replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
    ```
  
# End

上一篇 [《数据库准备工作》](./db-create.md)

下一篇 [《需求说明》](../requirement/summary.md)