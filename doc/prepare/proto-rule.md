# Proto使用说明
你心里是否有一个为什么，明明已经有了proto的官方文档，为什么还要介绍单独介绍一下proto文件的使用，
没错，我们这里不会介绍proto的语法、怎么使用protoc，我们这里只介绍proto配合goctl使用的一些规范，
如果你看完本章节觉得proto这样写是一种约束，那么建议你在选择goctl生成rpc服务时需要慎重考虑，但同时
也请你相信我们这样设计有其必要的原因和好处。

# proto编写规则
除proto本身的用法外，为了配合goctl生成，我们做了一下变动：
* 一个rpc服务的proto中有且只有一个service
* service中的rpc定义的入参和出参不能import外部的元素

这样做的原因是方便goctl生成rpc client层代码，见[服务目录#Rpc Structure](./service-structure.md)

> 生成rpc client层的目的是为了屏蔽rpc内部细节，这样一来，即使rpc内部怎么变动，rpc client均可以达到以不变应万变的效果。

# End

上一篇 [《Api语法介绍》](./api-grammar.md)

下一篇 [《创建工程》](./project-create.md)

# 猜你想

* [《目录说明》](../index.md)