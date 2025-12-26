# CloudWeGoInstance

这是一个 “小型电商系统”，用于学习 CloudWeGo 生态。

## 技术栈

Hertz: 字节跳动开源的 Go 微服务 HTTP 框架，具有高易用性、高性能、高扩展性等特点。

Kitex: 字节跳动开源的 Go 微服务 RPC 框架，具有高性能、强可扩展的特点，在字节内部已广泛使用。

Docker: 一个用于开发，交付和运行应用程序的开放平台。

## 工作流

### 构建 rpc 服务

1. 在 idl 文件夹下定义服务。
2. 在 rpc 文件夹下，使用 Kitex 生成代码框架。
3. 实现 rpc 服务，使用 docker-compose.yml 构建依赖的环境。
4. 使用 Kitexcall 以及 gotest 测试 rpc 服务功能。
5. 打包为镜像。

### 构建 http 服务

Hertz 生成代码框架有两种方式：

1. 根据已有的 IDL 文件生成，这种情况下 Hertz 一般起到透传的功能，即 rpc 转 http。
2. 定义新的 IDL 文件，并生成，这种情况赋予了 Hertz 开发人员更高的灵活性，但在小型项目中不必要。

然后实现功能，并打包为镜像。

### 运行项目

当然是使用 docker-compose.yml 进行编排，

在此过程中，你只需在文件中声明所需的环境变量即可。

环境变量一般包含其所依赖的其他微服务的地址、数据库的地址、以及一些参数等。

也可以在 rpc 服务中添加 etcd 的逻辑，这样就只依赖 etcd 的地址，其他依赖的微服务地址通过 etcd 访问得到。

## 优势

微服务架构的主要优势先不谈（如负载均衡、服务治理等）。

在本项目中， verify_code_service, user_account_service 以及 order_service 是可以复用的。

如果你在其他项目需要一个用户账户服务，只需运行一个 user_account_service 的容器，而不必重新开发。

如果你担心这样不够灵活--环境变量就是做这种事的，

user_account_service 的许多都不是硬编码的，如 mysql 和 verify_code_service 的地址，生成 token 的密钥等。

在主逻辑不变的情况下，基于环境变量，尽量给出灵活性。

Hertz 也很灵活，来自两个方面：

1. 如果 Hertz 是对 rpc 服务的透传，这意味着 Hertz 项目也可以复用，前端可以基于 Hertz 暴露的 api 任意编排流程。
2. 如果 Hertz 是对 rpc 服务的编排，那么它的复用性不高，但这意味着前端可以简单调用 Hertz 暴露的 api，而不必自己拼装业务流。

实际上，编排业务流程的任务，前端和 Hertz 都有，也无所谓谁多谁少，因为这取决于 Hertz 想暴露什么，透传还是编排好的流程。

## 部署本项目

这实际上是一个 Hertz 项目和 3 个 rpc 服务，我只能建议你阅读各个模块的 README.md。

但可以这么总结：构建 3 个 rpc 服务的镜像，然后运行 Hertz 项目的 Makefile 即可。