namespace go base

enum Code {
    SUCCESS = 0,
    INVALID_PARAM = 1,
    DB_ERR = 2,
    SERVICE_ERR = 3,
    NOT_FOUND = 4,
}

struct BaseResponse {
    1: Code code,
    2: string msg,
}

enum TargetType {
    Email = 1,
    Phone = 2,
}