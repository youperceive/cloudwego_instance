# VerifyCodeService

## 命令

### 重新生成 rpc 代码

```bash
# kitex 命令生成
kitex -module CloudWeGoInstance/VerifyCodeService -service VerifyCodeService ./idl/captcha.thrift

# kitex 脚本
script_kit/kitex_gen.sh
```

### 调试源码时启动方法

会重新构建可执行文件，并启动。

```bash
make
```

### docker构建镜像

```bash
# docker 命令构建
docker build -f Dockerfile -t verify-code-service .

# docker 脚本构建
script_kit/docker_build.sh
```

### docker-compose 启动服务

```bash
docker compose up
```

## api

### 数据结构

```thrift
/*
There can be add new types in the future, such as wechat validation.
*/
enum CaptchaType {
    Phone = 1,
    Email = 2,
}

/*
This mean a request to generate a captcha to a target for a purpose.
A possible example:
    redis-cli.Insert("{purpose}:{target}", "{captcha_code}")
*/
struct GenerateCaptchaRequest {
    1: CaptchaType type,
    2: string target,
    3: string purpose,                      // discarded
    4: optional i32 expire_seconds = 300,
    5: optional i32 max_validate_times = 3,
    6: string proj,
    7: string biz_type,
}

struct GenerateCaptchaResponse {
    1: base.BaseResponse baseResp,
}

struct ValidateCaptchaRequest {
    1: string target,
    2: string purpose,  //deprecated, use biz_type instead
    3: string captcha,
    4: string proj,
    5: string biz_type,
}

struct ValidateCaptchaResponse {
    1: base.BaseResponse baseResp,
    2: bool valid,
}
```

## 接口

```thrift
service CaptchaService {
    GenerateCaptchaResponse GenerateCaptcha(1: GenerateCaptchaRequest req),
    ValidateCaptchaResponse ValidateCaptcha(1: ValidateCaptchaRequest req),
}
```

## 环境变量

当你拉取了这个服务的镜像，请这样启动一个实例。

```yml
  verify-code-service:
    image: verify-code-service:latest
    environment:
      - PRINT_CAPTCHA=true # 是否在容器日志中打印验证码
      - REDIS_ADDR=redis:6379 # redis 地址
      - ETCD_ADDR=etcd:2379 # etcd 地址
    depends_on:
      - redis
      - etcd
    restart: always
    ports:
      - "8000:8000"
```