# goctl

## api 
(api服务相关操作)

### -o
(生成api文件)

- 示例：goctl api -o user.api

### new 
(快速创建一个api服务)

- 示例：goctl api new user

### format 
(api格式化，vscode使用)

- -dir
(目标目录)
- -iu
(是否自动更新goctl)
- -stdin
(是否从标准输入读取数据)

### validate
(验证api文件是否有效)

- -api
(指定api文件源)

	- 示例：goctl api validate -api user.api

### doc
(生成doc markdown)

- -dir
(指定目录)

	- 示例：goctl api doc -dir user

### go
(生成golang api服务)

- -dir
(指定代码存放目录)
- -api
(指定api文件源)
- -force
(是否强制覆盖已经存在的文件)
- -style
(指定文件名命名风格，gozero:小写，go_zero:下划线,GoZero:驼峰)

### java
(生成访问api服务代码-java语言)

- -dir
(指定代码存放目录)
- -api
(指定api文件源)

### ts
(生成访问api服务代码-ts语言)

- -dir
(指定代码存放目录)
- -api
(指定api文件源)
- webapi
- caller
- unwrap

### dart
(生成访问api服务代码-dart语言)

- -dir
(指定代码存放目标)
- -api
(指定api文件源)

### kt
(生成访问api服务代码-kotlin语言)

- -dir
(指定代码存放目标)
- -api
(指定api文件源)
- -pkg
(指定包名)

## template
(模板操作)

### init
(缓存api/rpc/model模板)

- 示例：goctl template init

### clean
(清空缓存模板)

- 示例：goctl template clean

### update
(更新模板)

- -category,c 
(指定需要更新的分组名 api|rpc|model)

	- 示例：goctl template update -c api

### revert
(还原指定模板文件)

- -category,c 
(指定需要更新的分组名 api|rpc|model)
- -name,n
(指定模板文件名)

## config
(配置文件生成)

### -path,p
(指定配置文件存放目录)

- 示例：goctl config -p user

## docker 
(生成Dockerfile)

### -go
(指定main函数文件)

### -namespace,n
(指定namespace)

## rpc (rpc服务相关操作)

### new
(快速生成一个rpc服务)

- -idea
(标识命令是否来源于idea插件，用于idea插件开发使用，终端执行请忽略[可选参数])
- -style
(指定文件名命名风格，gozero:小写，go_zero:下划线,GoZero:驼峰)

### template
(创建一个proto模板文件)

- -idea
(标识命令是否来源于idea插件，用于idea插件开发使用，终端执行请忽略[可选参数])
- -out,o
(指定代码存放目录)

### proto
(根据proto生成rpc服务)

- -src,s
(指定proto文件源)
- -proto_path,I
(指定proto import查找目录，protoc原生命令，具体用法可参考protoc -h查看)
- -dir,d
(指定代码存放目录)
- -idea
(标识命令是否来源于idea插件，用于idea插件开发使用，终端执行请忽略[可选参数])
- -style
(指定文件名命名风格，gozero:小写，go_zero:下划线,GoZero:驼峰)

### model
(model层代码操作)

- mysql
(从mysql生成model代码)

	- ddl
(指定数据源为
ddl文件生成model代码)

		- -src,s
(指定包含ddl的sql文件源，支持通配符匹配)
		- -dir,d
(指定代码存放目录)
		- -style
(指定文件名命名风格，gozero:小写，go_zero:下划线,GoZero:驼峰)
		- -cache,c
(生成代码是否带redis缓存逻辑，bool值)
		- -idea
(标识命令是否来源于idea插件，用于idea插件开发使用，终端执行请忽略[可选参数])

	- datasource
(指定数据源从
数据库链接生成model代码)

		- -url
(指定数据库链接)
		- -table,t
(指定表名，支持通配符)
		- -dir,d
(指定代码存放目录)
		- -style
(指定文件名命名风格，gozero:小写，go_zero:下划线,GoZero:驼峰)
		- -cache,c
(生成代码是否带redis缓存逻辑，bool值)
		- -idea
(标识命令是否来源于idea插件，用于idea插件开发使用，终端执行请忽略[可选参数])

