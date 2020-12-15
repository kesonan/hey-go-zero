# 创建工程
在前面我们已经对[软件环境](../index.md)进行了说明，本次演示项目环境是以
* mac OS
* Goland
* Go Module
* Goctl
* Idea Goctl插件

形式展开进行的，除此之外我们操作会在Goland+Terminal混合进行，如果准备工作还没做好的话，请先阅读一下[概要说明](../index.md)和[准备工作](../prepare/prepare.md) 章节。

# New
* 打开Goland->`File`->`New`->`Project...`
* 选择`Go modules`
* 选择项目存放目录为`~/goland/go`，并命名工程名为`hey-go-zero`，即Location位置值应该为`/Users/xxx/goland/go/hey-go-zero`
* `GOROOT`选择自己安装的go版本，我们这里选择的是`Go 1.15.1`
* `Vgo Executable`默认
* `Environment`由于在前面的环境准备中已经对go module设置了proxy,因此这里忽略为空就行。
* 点击`Create`即可。

# 目录结构创建
接下来我们按照之前提到过的[服务目录](./service-structure.md)去创建一个`service`目录，专门用于存放服务代码，然后在`service`子目录创建`user`,`selection`,`course`目录，创建后其目录树如下:

```text
hey-go-zero
├── go.mod
└── service
    ├── course
    ├── selection
    └── user
```

> 由于这是演示项目，这里就没有项目组一说，因此就不创建项目组目录了，后续进入每个服务模块时直接创建模块目录。

# End

上一篇 [《Proto使用说明》](./proto-rule.md)

下一篇 [《数据库准备工作》](./db-create.md)