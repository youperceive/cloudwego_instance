namespace go captcha

include "./base.thrift"

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
    3: string purpose,
    4: optional i32 expire_seconds = 300,
    5: optional i32 max_validate_times = 3,
}

struct GenerateCaptchaResponse {
    1: base.BaseResponse baseResp,
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

service CaptchaService {
    GenerateCaptchaResponse GenerateCaptcha(1: GenerateCaptchaRequest req),
    ValidateCaptchaResponse ValidateCaptcha(1: ValidateCaptchaRequest req),
}