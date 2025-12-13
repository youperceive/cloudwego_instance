# VerifyCodeService

## 启动方法

```bash
make
```

## 数据结构

``` go
enum CaptchaType {
    Phone = 1,
    Email = 2,
}

struct GenerateCaptchaRequest {
    1: CaptchaType type,
    2: string target,
    3: string purpose,
    4: optional i32 expire_seconds = 300,
    5: optional i32 max_validate_times = 3,
}

struct GenerateCaptchaResponse {
    1: base.BaseResponse baseResp,
    2: string captcha,
}

struct ValidateCaptchaRequest {
    1: string target,
    2: string purpose,
    3: string captcha,
}

struct ValidateCaptchaResponse {
    1: base.BaseResponse baseResp,
    2: bool valid,
}
```

## 方法

```go
service CaptchaService {
    GenerateCaptchaResponse GenerateCaptcha(1: GenerateCaptchaRequest req),
    ValidateCaptchaResponse ValidateCaptcha(1: ValidateCaptchaRequest req),
}
```