#!/bin/bash
kitexcall \
--idl-path ../../idl/order/order.thrift \
--method OrderService/Create \
--endpoint 127.0.0.1:8002 \
-f scripts_kit/test/create_order.json

# namespace go order

# include "../base/base.thrift"

# // 订单项（内嵌在 Order.items 中，id 用雪花算法 int64，order_id 关联订单 ObjectID）
# struct OrderItem {
#     1: i64 id,                  // 订单项唯一标识（雪花算法生成，全局唯一）
#     2: string order_id,         // 关联订单 id（对齐 Order.id，MongoDB ObjectID 字符串）
#     3: i64 product_id,          // 商品 id，非负
#     4: i64 sku_id,              // 规格 id：蓝/XL/件，仅存储数值，无业务语义，非负
#     5: i64 count,               // 购买数量，≥1（最少1件）
#     6: i64 price,               // 单价（分），≥0（避免浮点数精度问题）
#     7: map<string, string> ext, // 订单项扩展字段
# }

# // 订单主表（id 复用 MongoDB ObjectID 字符串）
# struct Order {
#     1: string id,               // 订单唯一标识（MongoDB ObjectID 字符串）
#     2: i32 type,                // 订单类型（决定 ext 字段语义）：1=普通订单，2=秒杀订单，3=团购订单
#     3: i32 status,              // 订单状态（仅存储数值）：1=待支付，2=已支付，3=已取消，4=已完成
#     4: i64 req_user_id,         // 订单发起人 id（A 向 B 下单，A 的 id），非负
#     5: i64 resp_user_id,        // 订单接收人 id（A 向 B 下单，B 的 id），非负
#     6: list<OrderItem> items,   // 订单明细
#     7: i64 created_at,          // 创建时间戳（秒级）
#     8: i64 updated_at,          // 更新时间戳（秒级）
#     9: map<string, string> ext, // 订单扩展字段（按 type 存储差异化数据）
# }

# // 创建订单时的订单项参数（剥离 id/order_id，由服务端生成）
# struct OrderItemForCreate {
#     1: i64 product_id,          // 商品 id，非负
#     2: i64 sku_id,              // 规格 id，非负
#     3: i64 count,               // 购买数量，≥1
#     4: i64 price,               // 单价（分），≥0
#     5: map<string, string> ext,
# }

# struct CreateRequest {
#     1: i32 type,
#     2: i32 status,
#     3: i64 req_user_id,
#     4: i64 resp_user_id,
#     5: list<OrderItemForCreate> items,
#     6: map<string, string> ext,
# }

# struct CreateResponse {
#     1: base.BaseResponse baseResp,
#     2: string order_id,            // 创建成功的订单 id（MongoDB ObjectID 字符串）
# }

# struct UpdateRequest {
#     1: string id,
#     2: optional i32 status,              // 订单状态，调用方自行确认业务语义，服务端仅存储值
#     3: optional map<string, string> ext, // 扩展字段：覆盖（而非合并）原有值，如需合并请传全量
# }

# struct UpdateResponse {
#     1: base.BaseResponse baseResp,
# }

# struct QueryOrderInfoRequest {
#     1: string id, // 订单 id（MongoDB ObjectID 字符串）
# }

# struct QueryOrderInfoResponse {
#     1: base.BaseResponse baseResp,
#     2: Order order,
# }

# enum QueryOrderIdType {
#     REQ_USER = 1,  // 按 req_user_id（发起人）查询
#     RESP_USER = 2, // 按 resp_user_id（接收人）查询
#     EXT_KEY = 3,   // 按 ext[ext_key] = ext_val 精准查询
# }

# struct QueryOrderIdRequest {
#     1: QueryOrderIdType type,
#     2: optional i64 user_id,        // REQ_USER/RESP_USER 场景必填，EXT_KEY 场景无效，非负
#     3: optional string ext_key,     // EXT_KEY 场景必填：ext 字段的 key（如 seckill_id），不能为空
#     4: optional string ext_val,     // EXT_KEY 场景必填：ext 字段的 value（如 100），不能为空
#     5: optional i32 page = 1,       // 页码，≥1，默认1
#     6: optional i32 page_size = 20, // 页大小，1≤page_size≤100，默认20（避免单次返回过多数据）
# }

# struct QueryOrderIdResponse {
#     1: base.BaseResponse baseResp,
#     2: list<string> order_id,      // 符合条件的订单 id 列表
#     3: i32 total,                  // 符合条件的订单总数（前端计算总页数：total/page_size 向上取整）
#     4: i32 page,                   // 当前页码（和请求一致）
#     5: i32 page_size,              // 当前页大小（和请求一致）
# }

# service OrderService {
#     CreateResponse Create(1: CreateRequest req),
#     UpdateResponse Update(1: UpdateRequest req),
#     QueryOrderInfoResponse QueryOrderInfo(1: QueryOrderInfoRequest req),
#     QueryOrderIdResponse QueryOrderId(1: QueryOrderIdRequest req),
# }