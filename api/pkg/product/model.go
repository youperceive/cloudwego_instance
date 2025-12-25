package product

// Product 商品原生结构体（匹配数据库字段）
type Product struct {
	ID         int64             `db:"id"`
	MerchantID int64             `db:"merchant_id"`
	Name       string            `db:"name"`
	Ext        map[string]string `db:"ext"` // 扩展字段
}

// Sku SKU 原生结构体（匹配数据库字段）
type Sku struct {
	ID         int64  `db:"id"`
	MerchantID int64  `db:"merchant_id"` // 冗余存储，方便归属校验
	ProductID  int64  `db:"product_id"`
	SkuCode    string `db:"sku_code"`
	Price      int64  `db:"price"` // 分
	Stock      int32  `db:"stock"`
}
