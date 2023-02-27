# KOP 公用框架工具

## 模块

`server` http server 包含默认启动项

`client` http请求 包含链路信息、重试、状态判断、中间件等

`perrors` 错误处理

`plog` 日志

`orm` 数据库连接、初始化等

`redis` redis连接、初始化等

`meta` 封装trace链路信息

`util` 应用工具类

`middleware` gin框架通用middleware 进行域名校验 登录身份校验等

`wrapper` http client通用wrapper

`celery` 消息中心

## 使用方法
go 1.17

需要配置使用公司仓库

`go env -w GOPRIVATE=git.im30.lan`

公司gitlab不支持https需要配置使用http

`go env -w GOINSECURE=git.im30.lan`

私有仓库使用git协议访问(需在gitlab上配置ssh-key)
```yaml
执行命令：
git config --global url."git@git.im30.lan:".insteadOf "http://git.im30.lan/"

或编辑 ~/.gitconfig 添加以下内容
[url "git@git.im30.lan:"]
    insteadOf = http://git.im30.lan/
```

添加依赖

`go get -v git.im30.lan/kop/common


## 建议反馈

飞书 haer
