namespace go user

include "./base.thrift"

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

enum RegisterType {
    EMAIL = 1,
    PHONE = 2,
}

struct RegisterRequest {
    1: optional string username,
    2: string target,              // phone or email, determined by register_type
    3: RegisterType register_type,
    4: string password,
    5: string captcha,             // before register, need to get a captcha. can be arranged in frontend, not here
}

struct RegisterResponse {
    1: base.BaseResponse baseResp,
    2: optional i64 user_id,
}

service UserAccountService {
    RegisterResponse Register(1: RegisterRequest req),
}