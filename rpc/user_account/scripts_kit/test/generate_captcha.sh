#!/bin/bash

kitexcall \
--idl-path ../../idl/verify_code/verify_code.thrift \
--method VerifyCodeService/GenerateCaptcha \
--endpoint 127.0.0.1:8000 \
-d '{"type": 1, "proj": "user-account-service", "biz_type": "login", "target": "904413"}'

# struct GenerateCaptchaRequest {
#     1: CaptchaType type,
#     2: string target,
#     3: string purpose,                      // discarded
#     4: optional i32 expire_seconds = 300,
#     5: optional i32 max_validate_times = 3,
#     6: string proj,
#     7: string biz_type,
# }

# struct GenerateCaptchaResponse {
#     1: base.BaseResponse baseResp,
# }

# struct ValidateCaptchaRequest {
#     1: string target,
#     2: string purpose,  //deprecated, use biz_type instead
#     3: string captcha,
#     4: string proj,
#     5: string biz_type,
# }