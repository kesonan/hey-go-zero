# rest api 服务基本配置说明
在启动一个rest api服务前，我们需要对api服务做一些常规性配置，接下来我们来看一下`go-zero`中的api基本
配置有哪些。

# 配置定义
[RestConf定义源码](https://github.com/tal-tech/go-zero/blob/master/rest/config.go)

```yaml
Name: "user-api" # 服务名称
Log:  # 日志相关配置
  ServiceName: "user-api" # 日志服务名称
  Mode: console # 日志输出模式 console|file|volume
  Path: logs  # 日志输出路径
  Level: info # 日志打印级别
  Compress: true  # 日志滚动-是否开启gzip压缩
  KeepDays: 7 # 日志滚动-保留天数
  StackCooldownMillis: 100  # stack日志write频率
Mode: dev # 服务环境 dev|test|pre|pro
MetricsUrl: url # 指标上报地址，post json
Prometheus: # prometheus相关
  Host: 127.0.0.1 #  监听主机
  Port: 9101  # 监听端口
  Path: /metrics  # 路由地址
Host: 127.0.0.1 # 服务监听主机
Port: 8080  # 服务监听宽口
CertFile: # https证书
KeyFile:  # https key
Verbose:  # http请求日志是否详细输出
MaxConns: 10000 # 最大连接数
MaxBytes: 1048576 # 最大传输的content-length,范围 0-8388608 bytes
Timeout: 3000 # 超时时长，单位：毫秒
CpuThreshold: 900 # 最大允许使用的cpu，范围 0-1000，超出90%就会drop掉多余的请求
Signature:  # 签名相关配置
  Strict: true  # 开启则需要传PrivateKeys配置
  Expiry: 1h  # 有效期，默认1个小时
  PrivateKeys:  # key相关
    - Fingerprint: test # 指纹
      KeyFile: ./key  # key文件
```