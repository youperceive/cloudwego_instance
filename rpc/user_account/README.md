# UserAccountService

## 简介
支持用户通过邮箱、手机号注册账号，标记用户权限类型，并实现登录返还 token。
基于 CloudWeGo Kitex 开发，依赖 MySQL（存储用户数据）、Redis（缓存/验证码）、verify_code_service（验证码校验）。

## 前置准备
### 1. 安装基础环境
确保本地已安装以下工具：
- Git >= 2.30.0
- Docker >= 20.10.0
- Docker Compose >= 2.0.0（或 Docker Desktop 内置版本）
- 可访问的 MySQL 实例（5.7+/8.0+）
- 可访问的 Redis 实例（6.0+）

### 2. 依赖服务部署
该服务依赖 `verify_code_service`，需先部署：
1. 进入仓库的 `rpc/verify_code` 目录，阅读该目录下的 `README.md` 完成部署；
2. 记录 `verify_code_service` 的可访问地址（格式：`IP:端口`，如 `172.17.0.2:8000`）。

### 3. 数据库准备
在 MySQL 中创建用户服务数据库，并执行以下表结构脚本（核心表）：
```sql
CREATE DATABASE IF NOT EXISTS user_account DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE user_account;

CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(64) DEFAULT '' COMMENT '用户名',
  `email` varchar(64) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) DEFAULT '' COMMENT '手机号',
  `password` varchar(128) NOT NULL COMMENT '哈希后的密码',
  `ext` json DEFAULT '{}' COMMENT '扩展字段',
  `user_type` tinyint(4) DEFAULT 1 COMMENT '1-普通用户 2-管理员 3-第三方用户',
  `status` tinyint(4) DEFAULT 1 COMMENT '1-正常 2-禁用 3-注销',
  `created_at` bigint(20) DEFAULT UNIX_TIMESTAMP() COMMENT '创建时间（秒级）',
  `updated_at` bigint(20) DEFAULT UNIX_TIMESTAMP() COMMENT '更新时间（秒级）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`),
  UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

## 部署步骤
### 1. 克隆仓库
```shell
git clone https://github.com/youperceive/cloudwego_instance.git
cd cloudwego_instance
```

### 2. 构建 Docker 镜像
进入用户服务目录，执行构建脚本：
```shell
# 进入 user_account 服务目录
cd rpc/user_account

# 添加脚本执行权限
chmod +x scripts_kit/docker_build.sh

# 执行镜像构建脚本
./scripts_kit/docker_build.sh
```

### 3. 验证镜像构建结果
执行以下命令，确认镜像存在：
```shell
docker images | grep user-account-service
# 预期输出：user-account-service   latest    [镜像ID]   [构建时间]
```

### 4. 编写 Docker Compose 配置
在 `rpc/user_account` 目录下创建 `docker-compose.yml` 文件，内容如下（替换占位符为实际值）：
```yml
version: '3.8'

# 网络配置：确保与 verify_code_service 在同一网络（若verify_code_service用Docker部署）
networks:
  app-network:
    driver: bridge

services:
  user-account-service:
    image: user-account-service:latest
    container_name: user-account-service
    restart: always
    networks:
      - app-network
    environment:
      # MySQL连接串：格式为 用户名:密码@tcp(MySQL地址:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
      - MYSQL_DSN="root:123456@tcp(192.168.1.100:3306)/user_account?charset=utf8mb4&parseTime=True&loc=Local"
      # Redis地址：格式为 IP:端口
      - REDIS_ADDR="192.168.1.100:6379"
      # 验证码服务地址：需与前置部署的 verify_code_service 地址一致
      - VERIFY_CODE_SERVICE_ADDR="verify-code-service:8000"
      # JWT密钥：建议至少32位随机字符串（必填，用于生成token）
      - JWT_SECRET="your_32_bit_random_jwt_secret_12345678"
    ports:
      - "8001:8000"  # 宿主机端口:容器端口（可根据需要修改）
    # 依赖：确保Redis/MySQL/验证码服务先启动（若这些服务也在Docker Compose中）
    depends_on:
      - verify-code-service  # 若verify_code_service在该Compose中，需添加其配置

  # （可选）若verify_code_service未单独部署，可在此添加其配置
  verify-code-service:
    image: verify-code-service:latest
    container_name: verify-code-service
    restart: always
    networks:
      - app-network
    environment:
      - REDIS_ADDR="192.168.1.100:6379"
    ports:
      - "8000:8000"
```

### 5. 启动服务
```shell
# 进入docker-compose.yml所在目录
cd rpc/user_account

# 启动服务（后台运行）
docker-compose up -d

# 查看服务日志，确认无报错
docker-compose logs -f user-account-service
```

## 验证服务可用性
### 1. 检查容器状态
```shell
docker-compose ps
# 预期 user-account-service 状态为 Up (healthy)
```

### 2. 接口调用示例（可选）
通过 Kitex 客户端调用注册接口，示例代码：
```go
package main

import (
    "context"
    "fmt"

    "github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/user_account"
    "github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/user_account/useraccountservice"
    "github.com/youperceive/cloudwego_instance/rpc/base/kitex_gen/base"
    "github.com/cloudwego/kitex/client"
)

func main() {
    // 创建客户端（替换为实际服务地址）
    cli, err := useraccountservice.NewClient(
        "user-account-service",
        client.WithHostPorts("127.0.0.1:8001"),
    )
    if err != nil {
        panic(err)
    }

    // 构造注册请求
    req := &user_account.RegisterRequest{
        Username:   "test_user",
        Target:     "13800138000",
        TargetType: base.TargetType_PHONE,
        Password:   "e10adc3949ba59abbe56e057f20f883e", // 123456的MD5哈希（前端需哈希后传输）
        Captcha:    "123456", // 需先从verify_code_service获取有效验证码
    }

    // 调用注册接口
    resp, err := cli.Register(context.Background(), req)
    if err != nil {
        panic(err)
    }
    fmt.Printf("注册结果：%+v\n", resp)
}
```

## 常见问题排查
1. 镜像构建失败：
   - 检查 `scripts_kit/docker_build.sh` 脚本是否有编译步骤，确保本地Docker可访问Go镜像源；
   - 确认仓库代码完整，无缺失文件。
2. 服务启动后无法连接MySQL/Redis：
   - 检查环境变量中的地址是否正确，容器是否能访问外部MySQL/Redis（可进入容器 `docker exec -it user-account-service bash` 测试ping）；
   - 确认MySQL/Redis的防火墙/安全组已开放端口。
3. 验证码校验失败：
   - 确认 `VERIFY_CODE_SERVICE_ADDR` 地址正确，verify_code_service 已正常运行；
   - 检查验证码是否过期、是否与手机号/邮箱匹配。
4. JWT token生成失败：
   - 确认 `JWT_SECRET` 长度足够（建议32位以上），格式为字符串。

## 客户端使用
在Go项目中导入客户端，调用服务：
```go
import (
    "github.com/youperceive/cloudwego_instance/rpc/user_account/kitex_gen/user_account/useraccountservice"
    "github.com/cloudwego/kitex/client"
)

// 创建客户端
cli, err := useraccountservice.NewClient(
    "user-account-service",
    client.WithHostPorts("服务地址:端口"), // 如 "127.0.0.1:8001"
)
if err != nil {
    // 处理错误
}

// 调用Login/Register接口（示例见「验证服务可用性」章节）
```

## IDL
```thrift
namespace go user_account

include "../base/base.thrift"

enum UserType {
    USER = 1,
    ADMIN = 2,
    THIRD_PARTY = 3,
}

struct User {
    1: i64 id,
    2: string username,
    3: string email,
    4: string phone,
    5: optional map<string, string> ext = { },
    6: optional UserType user_type = UserType.USER,
    7: optional i64 created_at,
    8: optional i64 updated_at,
    9: optional i32 status = 1,                     // 1-正常，2-禁用，3-注销
}

struct RegisterRequest {
    1: optional string username,
    2: string target,               // phone or email, determined by register_type
    3: base.TargetType target_type,
    4: string password,             // frontend need to transmit password after hash it
    5: string captcha,              // before register, need to get a captcha. can be arranged in frontend, not here
}

struct RegisterResponse {
    1: base.BaseResponse baseResp,
    2: optional i64 user_id,
}

struct LoginRequest {
    1: string target,
    2: base.TargetType target_type,
    3: string password,             // frontend need to transmit password after hash it
}

struct LoginResponse {
    1: base.BaseResponse baseResp,
    2: string token,
}

service UserAccountService {
    RegisterResponse Register(1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
}
```