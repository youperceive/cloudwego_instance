namespace go user_account

include "../base/base.thrift"

struct User {
    1: i64 id,
    2: string username,
    3: string email,
    4: string phone,
    5: optional map<string, string> ext = { },
    6: optional i8 user_type = 1,              // no means
    7: optional i64 created_at,
    8: optional i64 updated_at,
    9: optional i32 status = 1,                // 1-正常，2-禁用，3-注销
}

struct RegisterRequest {
    1: optional string username,
    2: string target,               // phone or email, determined by register_type
    3: base.TargetType target_type,
    4: string password,             // frontend need to transmit password after hash it
    5: string captcha,              // before register, need to get a captcha. can be arranged in frontend, not here
    6: optional i8 user_type = 1,
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