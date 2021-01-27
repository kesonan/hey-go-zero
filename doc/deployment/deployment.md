# 服务部署

# 准备工作

详情见[准备工作](./prepare.md)

# etcd&redis&mysql

# goctl
在`jenkins`所在机器上安装`goctl`
```shell
$ goctl -v
```
```text
goctl version 1.1.4 linux/amd64
```

# gitlab

## 添加SSH Key

进入gitlab，点击用户中心，找到`Settings`，在左侧找到`SSH Keys`tab
![ssh key](../../resource/ssh-add-key.png)

* 1、在jenkins所在机器上查看公钥

``` shell
$ cat ~/.ssh/id_rsa.pub
```

* 2、如果没有，则需要生成，如果存在，请跳转到第3步

``` shell
$ ssh-keygen -t rsa -b 2048 -C "email@example.com"
```

> "email@example.com" 可以替换为自己的邮箱
>
完成生成后，重复第一步操作

* 3、将公钥添加到gitlab中

## 上传代码到gitlab仓库

* 1、创建group名称为`zero`
  ![gitlab-group](../../resource/gitlab-group.png)
* 2、创建project
  ![gitlab-project](../../resource/gitlab-project.png)
* 3、上传代码到gitlab
  ![gitlabg-code](../../resource/gitlab-code.png)

# jenkins

## 添加凭据

* 查看jenkins所在机器的私钥，与前面gitlab公钥对应

``` shell
$ cat id_rsa
```

* 进入jenkins，依次点击`Manage Jenkins`-> `Manage Credentials`
  ![credentials](../../resource/jenkins-credentials.png)

* 进入`全局凭据`页面，添加凭据，`Username`是一个标识，后面添加pipeline你知道这个标识是代表gitlab的凭据就行，Private Key`即上面获取的私钥
  ![jenkins-add-credentials](../../resource/jenkins-add-credentials.png)

## 添加全局变量
进入`Manage Jenkins`->`Configure System`，滑动到`全局属性`条目，添加docker私有仓库相关信息，如图为`docker用户名`、`docker用户密码`、`docker私有仓库地址`
> ## 温馨提示
> 这里我使用的是滴滴云的公测阶段私有仓库，如果没有云厂商提供的私有仓库使用，可以自行搭建一个私有仓库，这里就不赘述了，大家自行google。

## 配置git
进入`Manage Jenkins`->`Global Tool Configureation`，找到Git条目，填写jenkins所在机器git可执行文件所在path，如果没有的话，需要在jenkins插件管理中下载Git插件。
![jenkins-git](../../resource/jenkins-git.png)


![jenkins-configure](../../resource/jenkins-configure.png)
## 添加一个Pipeline

> pipeline用于构建项目，从gitlab拉取代码->生成Dockerfile->部署到k8s均在这个步骤去做，这里是演示环境，为了保证部署流程顺利，
> 需要将jenkins安装在和k8s集群的其中过一个节点所在机器上，我这里安装在master上的。

* 获取凭据id 进入凭据页面，找到Username为`gitlab`的凭据id
  ![jenkins-credentials-id](../../resource/jenkins-credentials-id.png)

* 进入jenkins首页，点击`新建Item`，名称为`user`
  ![jenkins-add-item](../../resource/jenkins-new-item.png)

* 查看项目git地址
![gitlab-git-url](../../resource/gitlab-git-url.png)

* 添加服务类型Choice Parameter,在General中勾选`This project is parameterized
  `,点击`添加参数`选择`Choice Parameter`，按照图中添加选择的值常量(api、rpc)及接收值的变量(type)，后续在Pipeline script中会用到。
  ![jenkins-choice-parameter](../../resource/jenkins-choice.png)
  
* 配置`user`
  在`user`配置页面，向下滑动找到`Pipeline script`,填写脚本内容
  ```shell
  pipeline {
    agent any
    parameters {
        gitParameter name: 'branch', 
        type: 'PT_BRANCH',
        branchFilter: 'origin/(.*)',
        defaultValue: 'master',
        selectedValue: 'DEFAULT',
        sortMode: 'ASCENDING_SMART',
        description: '选择需要构建的分支'
    }

    stages {
        stage('服务信息')    {
            steps {
                sh 'echo 分支：$branch'
                sh 'echo 构建服务类型：${JOB_NAME}-$type'
            }
        }


        stage('check out') {
            steps {
                checkout([$class: 'GitSCM', 
                branches: [[name: '$branch']],
                doGenerateSubmoduleConfigurations: false, 
                extensions: [], 
                submoduleCfg: [],
                userRemoteConfigs: [[credentialsId: '${credentialsId}', url: '${gitUrl}']]])
            }   
        }

        stage('获取commit_id') {
            steps {
                echo '获取commit_id'
                git credentialsId: '${credentialsId}', url: '${gitUrl}'
                script {
                    env.commit_id = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
                }
            }
        }


        stage('goctl版本检测') {
            steps{
                sh '/usr/local/bin/goctl -v'
            }
        }

        stage('Dockerfile Build') {
            steps{
                   sh '/usr/local/bin/goctl docker -go service/${JOB_NAME}/${type}/${JOB_NAME}.go'
                   script{
                       env.image = sh(returnStdout: true, script: 'echo ${JOB_NAME}-${type}:${commit_id}').trim()
                   }
                   sh 'echo 镜像名称：${image}'
                   sh 'docker build -t ${image} .'
            }
        }

        stage('上传到镜像仓库') {
            steps{
                sh '/root/dockerlogin.sh'
                sh 'docker tag  ${image} ${dockerServer}/${image}'
                sh 'docker push ${dockerServer}/${image}'
            }
        }

        stage('部署到k8s') {
            steps{
                script{
                    env.deployYaml = sh(returnStdout: true, script: 'echo ${JOB_NAME}-${type}-deploy.yaml').trim()
                    env.port=sh(returnStdout: true, script: '/root/port.sh ${JOB_NAME}-${type}').trim()
                }

                sh 'echo ${port}'

                sh 'rm -f ${deployYaml}'
                sh '/usr/local/bin/goctl kube deploy -secret dockersecret -replicas 2 -nodePort 3${port} -requestCpu 200 -requestMem 50 -limitCpu 300 -limitMem 100 -name ${JOB_NAME}-${type} -namespace hey-go-zero -image ${dockerServer}/${image} -o ${deployYaml} -port ${port}'
                sh '/usr/bin/kubectl apply -f ${deployYaml}'
            }
        }
        
        stage('Clean') {
            steps{
                sh 'docker rmi -f ${image}'
                sh 'docker rmi -f ${dockerServer}/${image}'
                cleanWs notFailBuild: true
            }
        }
    }
  }
  ```
  > ## 温馨提示
  > ${credentialsId}要替换为你的具体凭据值，即【添加凭据】模块中的一串字符串，${gitUrl}需要替换为你代码的git仓库地址，其他的${xxx}形式的变量无需修改，保持原样即可。
  ![user-pipepine-script](../../resource/user-pipeline-script.png)

# 查看pipeline
![build with parameters](../../resource/jenkins-build-with-parameters.png)
![build with parameters](../../resource/pipeline.png)


# 查看k8s服务
![k8s01](../../resource/k8s-01.png)
![k8s02](../../resource/k8s-02.png)
![k8s03](../../resource/k8s-03.png)

# 温馨提示
本文档仅仅是为了演示简单的部署流程，不能实际用于生产环境，本次演示中，请忽略一些应用安装规划不合理的地方，如：jenkins本身占用内存资源比较大，应该尽量单独部署在一台服务器等。