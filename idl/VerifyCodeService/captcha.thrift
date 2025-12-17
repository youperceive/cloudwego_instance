namespace go captcha

include "../base/base.thrift"

/*
This mean a request to generate a captcha to a target for a purpose.
A possible example:
    redis-cli.Insert("{purpose}:{target}", "{captcha_code}")
*/
struct GenerateCaptchaRequest {
    1: base.TargetType type,
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

service CaptchaService {
    GenerateCaptchaResponse GenerateCaptcha(1: GenerateCaptchaRequest req),
    ValidateCaptchaResponse ValidateCaptcha(1: ValidateCaptchaRequest req),
}