# Api语法
在上一章节[《服务目录》](./service-structure.md)我们提到过.api文件，该文件是对api服务的一个定义文件，像proto定义rpc服务一样，
我们可以通过.api文件声明服务名称，请求路由，请求方法，请求入参，请求出参，中间件声明，jwt声明等。利用该文件，我们可以通过goctl工具快速
生成api服务。

# 为什么需要Api文件
这里我们来假设，在没有.api文件的情况下，你会通过什么方式去省时省力的创建一个服务，我想没有几个人无不会用到工具其替代，如果说你的手速能达到这个效率，
那我们另当别论，请你跳过这一章节。我们需要.api文件可以带来哪些好处？
* 便于快速生成api服务
* 0学习成本
* 阅读性强
* 能生成各种语言的调用代码
* 替代api文档
* 扩展性强

# Api概览
api语法及其简单易懂，0学习成本，我们来看一下下列api定义

```text
import "common.api" [1]

type LoginReq {
    Username string `json:"username"`
    password string `json:"password"`
}

type LoginReply {
    Token string `json:"token"`
    Expire int64 `json:"expire"`
    RefreshToken string `json:"refreshToken"`
}

type UserInfoReply {
    Name string `json:"name"`
    Age int `json:"age"`
    Gender string `json:"gender"`
    Birthday string `json:"birthday"`
}

type UserPasswordEditReq {
    OldPassword string `json:"oldPassword"`
    NewPassword string `json:"newPassword"`
}

@server(
    jwt: Auth   [2]
    middleware: MetricMiddleware,LogMiddleware  [3]
    group: user [4]
)
service user-api {      [5]
    @handler healthCheck    [6]
    post /user/health/check     [7]

    @handler login
    post /user/login (LoginReq) returns (LoginReply)

    @handler userInfo
    get /user/info returns (UserInfoReply)

    @handler passwordEdit
    post /user/password/edit (UserPasswordEditReq)
}
```

在没了解api语法之前，我们可以从中简单的看出该文件定了4个路由，分别为

* `/user/health/check` 健康检查
* `/user/login` 登录
* `/user/info` 获取用户信息
* `/user/password/edit` 编辑用户密码

在service block head定义了user-api服务需要`jwt`授权，需要`MetricMiddleware`和`LogMiddleware`中间件和`user`分组。

# Api结构

```tetx
api
├── import block
├── info block
├── type block
└── service block
```

### import block
import block包含了对api文件引入的定义

语法定义

```antlrv4
api: importStatement
importStatement:importSpec+;
importSpec: 'import' importValue;
importValue:VALUE;
```

语法格式

```text
import $path1
import $path2
```

* `import`标志着一个api文件引入的开始
* `$path`则为具体的api文件路径（绝对/相对路径）

语法示例

```text
import "common.api"
import "empty.api"
```

> 说明：目前暂不支持Golang语法形式的group import。

### info block

info block是对api文件的一个描述块，便于其他开发人员能够快速了解该api信息，该block不参与api服务生成。

语法定义

```antlrv4
api: infoStatement;
infoStatement: 'info' '(' pair ')';
pair:(key ':' VALUE?)*;
```

语法格式

```text
info(
    $key: $value
)
```

语法示例

```text
info(
    title: "用户模块api"
    desc: "本api包含用户注册、登录、操作用户信息等"
    author: "songmeizi"
    email: "songmeizi@xx.com"
    version: "1.0"
)
```

* `info()`标志一个info描述块
* `$key`为描述信息声明
* `$value`为描述信息内容，可选

### type block

语法定义

```antlrv4
api: typeStatement;
typeStatement: (typeSingleSpec|typeGroupSpec);
typeGroupSpec:'type' '(' typeGroupBody ')';
typeGroupBody:(typeGroupAlias|structType)*;
typeGroupAlias:structNameId normalFieldType;
// eg: type xx struct {...}
typeSingleSpec: typeAlias|typeStruct;
typeStruct:'type' structType;
// eg: type Integer int
typeAlias:'type' structNameId '='? normalFieldType;
typeFiled:anonymousField|normalField |structType;
normalField:fieldName fieldType  tag?;
fieldType:normalFieldType|starFieldType|mapFieldType|arrayOrSliceType;
anonymousField: '*'? referenceId;
normalFieldType: GOTYPE|referenceId|('interface' '{' '}');
starFieldType: '*' normalFieldType;
mapFieldType: 'map' '[' GOTYPE ']' objType;
arrayOrSliceType: ('[' ']')+ objType;
structType: structNameId 'struct'? '{' (typeFiled)* '}';
objType: normalFieldType|starFieldType;
structNameId:IDENT;
fieldName:IDENT;
referenceId:IDENT;
tag: RAW_STRING;
```

语法格式

```text
type $name $defineType
type (
    $name {
        $filedName $defineType $rawString
    }
)
type $name {
     $filedName $defineType $rawString
 }

```

* `type`为固定字段，标志着一个type定义的开始
* `$name`为一个结构体的名称，必须满足【ID命名规则】，见下文
* `$defineType`为结构体声明类型，包含已经定义过的类型和go系统类型
* `$filedName`为字段名称，必须满足【ID命名规则】，见下文
* `$rawString`为raw string，即通过"`"引起来的字符串

    | key  | 描述                                             | 支持                    | 生效范围             | 示例        |
    |------|--------------------------------------------------|-------------------------|----------------------|-------------|
    | json | json序列化tag                                    | Golang	|request、response | json:"name,optional" |             |
    | path | 路由path值，配合路由使用                         | go-zero                 | request              | path:"id"   |
    | form | form标签，form 消息体和quertString的值均会被绑定 | go-zero                 | request              | form:"name" |


> 说明：这里你可以理解为golang中type定义的变体，对于struct省去了`struct`关键字。

语法示例

```go
type alias int

type (
    User {
        Name string `json:"name"`
    }

    Person {
        Name string `json:"name"`
    }
)

type Student {
    Name string `json:"name"`
}

type Teacher {
    Name string `json:"name"`
    User
} 
```

### service block

语法定义

```antlrv4
api: serviceStatement;
serviceStatement: (serviceServerSpec? serviceSpec);
serviceServerSpec: '@server' '(' identPair ')';
serviceSpec: 'service' serviceName '{' serviceBody+ '}';
serviceName:IDENT;
serviceBody:serviceDoc? (serviceHandler|serviceHandlerNew) serviceRoute;
serviceDoc: '@doc' '(' pair ')';
serviceHandler: '@server' '(' handlerPair ')';
serviceHandlerNew: '@handler' handlerValue;
serviceRoute:httpRoute ('(' referenceId? ')')? ('returns' '(' referenceId? ')'? ';'?;
httpRoute:HTTPMETHOD PATH;
identPair:(key ':' identValue)*;
handlerPair:(key ':' handlerValue)+;
identValue:(IDENT ','?)+;
handlerValue:IDENT;
pair:(key ':' VALUE?)*;
key:IDENT;
```

语法格式

```text
@server(
    $key: $value
)
service $serviceName {
    @handler $handlerName
    $httpMethod $path ($request) returns ($response)
}
```

* `@server`为固定字段，定义一个生成服务需要的一些附加属性的开始标志
* `service`为固定字段，定义一个服务的开始标志
* `@handler`为固定字段，定义一个handler文件的名称，生成代码后为handler的go文件名
* `returns`为固定字段，定义一个请求有响应体的标志
* `$key`生成服务的附加属性定义，为了扩展，从语法角度来讲支持任何满足【ID命名规则】的值，但目前在生成时仅用到了一下值

|    key     |                             描述                             |
|:----------:|:------------------------------------------------------------:|
|    jwt     |               声明该组协议需要生成jwt鉴权代码                |
| middleware |                声明该组协议需要生成中间件代码                |
|   group    | 标志该组协议生成代码需要按照group定义的value值进行文件夹分组 |

* `$value`生成服务的附加属性值
* `$serviceName`定义服务名称
* `$handlerName`定义生成代码时handler的go文件名称
* `$httpMethod`定义请求方法，请求方法仅支持小写，支持`get`、`head`、`post`、`put`、`patch`、`delete`、`connect`、`options`、`trace`
* `$path`定义请求path，请求路由需要满足正则`(('/'|'/:'|'-') [a-zA-Z0-9_-])+`
* `$request`定义请求体，前提必须要在上下文中已定义
* `$response`定义响应体，前提必须要在上下文中已定义

# 关键字
`info`,`map`,`struct`,`interface`,`type`,`server`,`doc`,`handler`,`service`,`returns`,`import`,`get`、`head`、`post`、`put`、`patch`、`delete`、`connect`、`options`、`trace`

除此之外还包括go语言关键字

> 说明：除关键字本身用途外，其他地方强烈建议不要使用关键字进行定义，如field名称，handler值，中间件值等

# ID命名规则
* 必须以下划线(_)、字母开头
* 以下划线(_)、字母、数字组合

# 补充说明
上述api语法中"语法定义"模块的相关代码为antlr4语法，更多信息可参考[官方文档](https://www.antlr.org/)

为了快速编写api文件，我们开发了idea、vsCode插件，该插件提供了高亮、跳转、模板、格式化等功能，[点击这里](https://github.com/tal-tech/goctl-plugins) 查看更多详情

# End

上一篇 [《服务目录》](./service-structure.md)

下一篇 [《Proto使用说明》](./proto-rule.md)