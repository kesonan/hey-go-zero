# user.rpc
user rpc 用于给其他服务根据其需要提供服务能力，如查询用户信息等。

# 创建rpc目录
在`service/user`下创建rpc目录，得到目录树

```text
user
└── rpc
```

# 新建user.proto
在文件夹`service/user/rpc`上右键->`New`->`New proto file`->`选择Empty file`->`输入user`->`OK`

# 定义proto