namespace go account

include "./base.thrift"

// simplest account struct, a mvp model
struct Account {
    1: i64 id 
    2: string username
    3: string password 
    4: string email
}

struct CreateAccountRequest {

}