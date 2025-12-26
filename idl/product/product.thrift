namespace go product

include "../base/base.thrift" // 仅保留基础响应（Code+Msg）

// ====================== 核心结构体（去掉状态/时间字段） ======================
// 商品主信息
struct Product {
    1: i64 id,                  // 商品ID
    2: i64 merchant_id,         // 商户ID
    3: string name,             // 商品名称
    4: map<string, string> ext, // 扩展字段（可选）
}

// SKU信息（订单核心依赖）
struct Sku {
    1: i64 id,          // SKU ID
    2: i64 product_id,  // 商品ID
    3: string sku_code, // 规格编码
    4: i64 price,       // 单价（分）
    5: i32 stock,       // 库存
}

// ====================== 商户端CRUD请求/响应 ======================
// 1. 商品增删改查
struct CreateProductRequest {
    1: required i64 merchant_id,         // 商户ID（必传）
    2: required string name,             // 商品名称
    3: optional map<string, string> ext, // 扩展字段（可选）
}

struct CreateProductResponse {
    1: required base.BaseResponse BaseResp,
    2: optional i64 product_id,             // 返回商品ID
}

struct DeleteProductRequest {
    1: required i64 merchant_id, // 商户ID（校验归属）
    2: required i64 product_id,  // 商品ID
}

struct DeleteProductResponse {
    1: required base.BaseResponse BaseResp,
}

struct UpdateProductRequest {
    1: required i64 merchant_id,         // 商户ID（校验归属）
    2: required i64 product_id,          // 商品ID
    3: optional string name,             // 可选更新名称
    4: optional map<string, string> ext, // 可选更新扩展字段
}

struct UpdateProductResponse {
    1: required base.BaseResponse BaseResp,
}

struct GetProductRequest {
    1: required i64 merchant_id, // 商户ID（校验归属）
    2: required i64 product_id,  // 商品ID
}

struct GetProductResponse {
    1: required base.BaseResponse BaseResp,
    2: optional Product product,            // 商品信息
}

// 2. SKU增删改查
struct CreateSkuRequest {
    1: required i64 merchant_id, // 商户ID（校验归属）
    2: required i64 product_id,  // 商品ID
    3: required string sku_code, // 规格编码
    4: required i64 price,       // 单价（分）
    5: required i32 stock,       // 库存
}

struct CreateSkuResponse {
    1: required base.BaseResponse BaseResp,
    2: optional i64 sku_id,                 // 返回SKU ID
}

struct DeleteSkuRequest {
    1: required i64 merchant_id, // 商户ID（校验归属）
    2: required i64 sku_id,      // SKU ID
}

struct DeleteSkuResponse {
    1: required base.BaseResponse BaseResp,
}

struct UpdateSkuRequest {
    1: required i64 merchant_id, // 商户ID（校验归属）
    2: required i64 sku_id,      // SKU ID
    3: optional string sku_code, // 可选更新规格编码
    4: optional i64 price,       // 可选更新价格
    5: optional i32 stock,       // 可选更新库存
}

struct UpdateSkuResponse {
    1: required base.BaseResponse BaseResp,
}

// ====================== 订单联动接口（核心） ======================
struct GetSkuRequest {
    1: required i64 sku_id, // SKU ID（订单用，无需商户ID）
}

struct GetSkuResponse {
    1: required base.BaseResponse BaseResp,
    2: optional Sku sku,                    // SKU信息
}

struct DeductSkuStockRequest {
    1: required i64 sku_id, // SKU ID
    2: required i32 count,  // 购买数量（≥1）
}

struct DeductSkuStockResponse {
    1: required base.BaseResponse BaseResp,
    2: optional i32 remain_stock,           // 剩余库存
}

// ====================== 新增：列表接口请求/响应 ======================
// 商户商品列表请求（分页）
struct ListProductRequest {
    1: required i64 merchant_id,    // 商户ID（必传，校验归属）
    2: optional i32 page_num = 1,   // 页码（默认1）
    3: optional i32 page_size = 20, // 页大小（默认20，建议限制最大100）
}

// 商户商品列表响应
struct ListProductResponse {
    1: required base.BaseResponse BaseResp,
    2: optional i64 total,                  // 商品总数（用于分页）
    3: optional list<Product> products,     // 商品列表
}

// 商品SKU列表请求
struct ListSkuRequest {
    1: required i64 merchant_id, // 商户ID（必传，校验归属）
    2: required i64 product_id,  // 商品ID（必传）
}

// 商品SKU列表响应
struct ListSkuResponse {
    1: required base.BaseResponse BaseResp,
    2: optional list<Sku> skus,             // SKU列表
}

// 商户基础结构体（仅保留ID字段，满足需求）
struct Merchant {
    1: i64 id,
    2: string name,
}

// 获取所有商户请求（无入参）
struct ListMerchantRequest {
}

// 获取所有商户响应
struct ListMerchantResponse {
    1: base.BaseResponse baseResp,
    2: list<Merchant> data,
}

// ====================== 商品服务核心接口 ======================
service ProductService {
    // 商户端商品CRUD
    CreateProductResponse CreateProduct(1: CreateProductRequest req) (api.post = "/create_product"),
    DeleteProductResponse DeleteProduct(1: DeleteProductRequest req) (api.post = "/delete_product"),
    UpdateProductResponse UpdateProduct(1: UpdateProductRequest req) (api.post = "/update_product"),
    GetProductResponse GetProduct(1: GetProductRequest req) (api.post = "/get_product"),
    // 商户端SKU CRUD
    CreateSkuResponse CreateSku(1: CreateSkuRequest req) (api.post = "/create_sku"),
    DeleteSkuResponse DeleteSku(1: DeleteSkuRequest req) (api.post = "/delete_sku"),
    UpdateSkuResponse UpdateSku(1: UpdateSkuRequest req) (api.post = "/update_sku"),
    // 订单联动接口
    GetSkuResponse GetSku(1: GetSkuRequest req) (api.post = "/get_sku"),
    DeductSkuStockResponse DeductSkuStock(1: DeductSkuStockRequest req) (api.post = "/deduct_sku"),
    // ====================== 新增：列表接口 ======================
    ListProductResponse ListProduct(1: ListProductRequest req) (api.post = "/list_product"),         // 商户商品列表
    ListSkuResponse ListSku(1: ListSkuRequest req) (api.post = "/list_sku"),                         // 商品SKU列表
    ListMerchantResponse ListMerchant(1: ListMerchantRequest req) (api.get = "/list_merchant"),
}