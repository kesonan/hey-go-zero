# 部署准备
在正式开始部署服务前我们还是需要做一些准备工作，在本次演示工程中，我们将演示怎么将服务部署到k8s集群中，因此，
在正式进入部署前，我们需要一些服务器（虚拟机中）来搭建k8s集群。

> ## ⚠️ 注意
> 本项目属于演示项目，仅共学习，不建议直接将此后续步骤直接应用于生产环境。

# 准备工作
* 虚拟机安装
  * 虚拟机安装这里就不做详细描述了，推荐安装VMware Fusion
* [centOS安装](../share/centos_install.md)
* [简单的k8s集群安装](../share/k8s_install.md)
* [gitlab安装](../share/gitlab.md)
* [jenkins安装](../share/jenkins-install.md)
* [redis&mysql&nginx&Etcd安装](../share/data.md)

> ## 温馨提示
> 由于后续我们利用jenkins来部署服务到k8s，请将jenkins安装在k8s的某一个节点上。

# End

下一篇 [《服务部署》](./deployment.md)

# 猜你想

返回 [《目录说明》](../index.md)