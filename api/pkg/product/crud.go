package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// ====================== 商品纯 CRUD ======================
// CreateProduct 创建商品（纯数据操作，返回商品ID/error）
func CreateProduct(ctx context.Context, merchantID int64, name string, ext map[string]string) (int64, error) {
	// 转换 ext 为 JSON 字符串
	extStr := marshalExt(ext)

	// 插入数据库
	result, err := DB.ExecContext(ctx, `
		INSERT INTO product (merchant_id, name, ext) VALUES (?, ?, ?)
	`, merchantID, name, extStr)
	if err != nil {
		return 0, fmt.Errorf("CreateProduct: %w", err)
	}

	// 获取自增 ID
	productID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateProduct: 获取ID失败: %w", err)
	}
	return productID, nil
}

// DeleteProduct 删除商品（校验归属，仅返回 error）
func DeleteProduct(ctx context.Context, merchantID, productID int64) error {
	// 校验商品归属
	var count int64
	err := DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM product WHERE id=? AND merchant_id=?
	`, productID, merchantID).Scan(&count)
	if err != nil {
		return fmt.Errorf("DeleteProduct: 校验归属失败: %w", err)
	}
	if count == 0 {
		return errors.New("DeleteProduct: 商品不存在或无权删除")
	}

	// 删除商品
	_, err = DB.ExecContext(ctx, `DELETE FROM product WHERE id=?`, productID)
	if err != nil {
		return fmt.Errorf("DeleteProduct: 执行删除失败: %w", err)
	}
	return nil
}

// UpdateProduct 更新商品（校验归属，仅返回 error）
func UpdateProduct(ctx context.Context, merchantID, productID int64, name *string, ext *map[string]string) error {
	// 校验归属
	var count int64
	err := DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM product WHERE id=? AND merchant_id=?
	`, productID, merchantID).Scan(&count)
	if err != nil {
		return fmt.Errorf("UpdateProduct: 校验归属失败: %w", err)
	}
	if count == 0 {
		return errors.New("UpdateProduct: 商品不存在或无权更新")
	}

	// 构造动态更新 SQL
	sqlStr := "UPDATE product SET "
	args := []interface{}{}
	if name != nil {
		sqlStr += "name=?, "
		args = append(args, *name)
	}
	if ext != nil {
		sqlStr += "ext=?, "
		args = append(args, marshalExt(*ext))
	}
	if len(args) == 0 {
		return errors.New("UpdateProduct: 无更新字段")
	}
	// 去掉最后一个逗号 + 条件
	sqlStr = sqlStr[:len(sqlStr)-2] + " WHERE id=?"
	args = append(args, productID)

	// 执行更新
	_, err = DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("UpdateProduct: 执行更新失败: %w", err)
	}
	return nil
}

// GetProduct 查询商品（返回商品结构体/error）
func GetProduct(ctx context.Context, merchantID, productID int64) (*Product, error) {
	var (
		p   Product
		ext []byte // 数据库存储的 JSON 字节
	)

	// 查询数据库
	err := DB.QueryRowContext(ctx, `
		SELECT id, merchant_id, name, ext FROM product WHERE id=? AND merchant_id=?
	`, productID, merchantID).Scan(&p.ID, &p.MerchantID, &p.Name, &ext)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("GetProduct: 商品不存在")
		}
		return nil, fmt.Errorf("GetProduct: 查询失败: %w", err)
	}

	// 转换 ext 为 map
	p.Ext = unmarshalExt(ext)
	return &p, nil
}

// ====================== SKU 纯 CRUD ======================
// CreateSku 创建 SKU（返回 SKU ID/error）
func CreateSku(ctx context.Context, merchantID, productID int64, skuCode string, price int64, stock int32) (int64, error) {
	// 先校验商品归属（商品必须属于该商户）
	var count int64
	err := DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM product WHERE id=? AND merchant_id=?
	`, productID, merchantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CreateSku: 校验商品归属失败: %w", err)
	}
	if count == 0 {
		return 0, errors.New("CreateSku: 商品不存在或无权创建 SKU")
	}

	// 插入 SKU
	result, err := DB.ExecContext(ctx, `
		INSERT INTO sku (merchant_id, product_id, sku_code, price, stock)
		VALUES (?, ?, ?, ?, ?)
	`, merchantID, productID, skuCode, price, stock)
	if err != nil {
		return 0, fmt.Errorf("CreateSku: 插入失败: %w", err)
	}

	skuID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateSku: 获取 ID 失败: %w", err)
	}
	return skuID, nil
}

// DeleteSku 删除 SKU（校验归属，仅返回 error）
func DeleteSku(ctx context.Context, merchantID, skuID int64) error {
	// 校验 SKU 归属
	var count int64
	err := DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sku WHERE id=? AND merchant_id=?
	`, skuID, merchantID).Scan(&count)
	if err != nil {
		return fmt.Errorf("DeleteSku: 校验归属失败: %w", err)
	}
	if count == 0 {
		return errors.New("DeleteSku: SKU 不存在或无权删除")
	}

	// 执行删除
	_, err = DB.ExecContext(ctx, `DELETE FROM sku WHERE id=?`, skuID)
	if err != nil {
		return fmt.Errorf("DeleteSku: 执行删除失败: %w", err)
	}
	return nil
}

// UpdateSku 更新 SKU（校验归属，仅返回 error）
func UpdateSku(ctx context.Context, merchantID, skuID int64, skuCode *string, price *int64, stock *int32) error {
	// 校验归属
	var count int64
	err := DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sku WHERE id=? AND merchant_id=?
	`, skuID, merchantID).Scan(&count)
	if err != nil {
		return fmt.Errorf("UpdateSku: 校验归属失败: %w", err)
	}
	if count == 0 {
		return errors.New("UpdateSku: SKU 不存在或无权更新")
	}

	// 构造动态更新 SQL
	sqlStr := "UPDATE sku SET "
	args := []interface{}{}
	if skuCode != nil {
		sqlStr += "sku_code=?, "
		args = append(args, *skuCode)
	}
	if price != nil {
		sqlStr += "price=?, "
		args = append(args, *price)
	}
	if stock != nil {
		sqlStr += "stock=?, "
		args = append(args, *stock)
	}
	if len(args) == 0 {
		return errors.New("UpdateSku: 无更新字段")
	}
	sqlStr = sqlStr[:len(sqlStr)-2] + " WHERE id=?"
	args = append(args, skuID)

	// 执行更新
	_, err = DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("UpdateSku: 执行更新失败: %w", err)
	}
	return nil
}

// GetSku 查询 SKU（订单用，无需商户ID，返回 SKU 结构体/error）
func GetSku(ctx context.Context, skuID int64) (*Sku, error) {
	var s Sku

	err := DB.QueryRowContext(ctx, `
		SELECT id, merchant_id, product_id, sku_code, price, stock FROM sku WHERE id=?
	`, skuID).Scan(&s.ID, &s.MerchantID, &s.ProductID, &s.SkuCode, &s.Price, &s.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("GetSku: SKU 不存在")
		}
		return nil, fmt.Errorf("GetSku: 查询失败: %w", err)
	}
	return &s, nil
}

// DeductSkuStock 扣减 SKU 库存（事务防超卖，返回剩余库存/error）
func DeductSkuStock(ctx context.Context, skuID int64, count int32) (int32, error) {
	if count <= 0 {
		return 0, errors.New("DeductSkuStock: 扣减数量必须 > 0")
	}

	// 开启事务
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("DeductSkuStock: 开启事务失败: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// 查询库存并加行锁（防超卖）
	var stock int32
	err = tx.QueryRowContext(ctx, `
		SELECT stock FROM sku WHERE id=? FOR UPDATE
	`, skuID).Scan(&stock)
	if err != nil {
		return 0, fmt.Errorf("DeductSkuStock: 查询库存失败: %w", err)
	}

	// 校验库存
	if stock < count {
		return 0, errors.New("DeductSkuStock: 库存不足")
	}

	// 扣减库存
	_, err = tx.ExecContext(ctx, `
		UPDATE sku SET stock = stock - ? WHERE id=?
	`, count, skuID)
	if err != nil {
		return 0, fmt.Errorf("DeductSkuStock: 扣减失败: %w", err)
	}

	// 返回剩余库存
	return stock - count, nil
}

// ====================== 内部工具函数 ======================
// marshalExt map 转 JSON 字符串（存储到数据库）
func marshalExt(ext map[string]string) string {
	if ext == nil {
		return ""
	}
	data, _ := json.Marshal(ext)
	return string(data)
}

// unmarshalExt JSON 字符串转 map（从数据库读取）
func unmarshalExt(data []byte) map[string]string {
	if len(data) == 0 {
		return nil
	}
	var ext map[string]string
	_ = json.Unmarshal(data, &ext)
	return ext
}

// ListProductByMerchant 按商户ID查询商品列表（分页）
func ListProductByMerchant(ctx context.Context, merchantID int64, pageNum, pageSize int) ([]*Product, int64, error) {
	// 前置检查：DB是否初始化
	if DB == nil {
		return nil, 0, errors.New("ListProductByMerchant: 数据库连接未初始化（DB为nil）")
	}

	// 1. 查询总数（用于分页）
	var total int64
	err := DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM product WHERE merchant_id=?
	`, merchantID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("ListProductByMerchant: 查询总数失败: %w", err)
	}

	// 2. 计算分页偏移量（pageNum从1开始）
	offset := (pageNum - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 3. 查询商品列表
	rows, err := DB.QueryContext(ctx, `
		SELECT id, merchant_id, name, ext FROM product 
		WHERE merchant_id=? 
		ORDER BY id DESC 
		LIMIT ? OFFSET ?
	`, merchantID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ListProductByMerchant: 查询列表失败: %w", err)
	}
	defer rows.Close()

	// 4. 解析结果
	products := make([]*Product, 0, pageSize)
	for rows.Next() {
		var (
			p   Product
			ext []byte
		)
		err := rows.Scan(&p.ID, &p.MerchantID, &p.Name, &ext)
		if err != nil {
			return nil, 0, fmt.Errorf("ListProductByMerchant: 解析商品失败: %w", err)
		}
		p.Ext = unmarshalExt(ext)
		products = append(products, &p)
	}

	// 5. 检查rows遍历是否出错
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ListProductByMerchant: 遍历商品失败: %w", err)
	}

	return products, total, nil
}

// ListSkuByProduct 按商品ID查询SKU列表（校验商户归属）
func ListSkuByProduct(ctx context.Context, merchantID, productID int64) ([]*Sku, error) {
	// 前置检查：DB是否初始化
	if DB == nil {
		return nil, errors.New("ListSkuByProduct: 数据库连接未初始化（DB为nil）")
	}

	// 1. 校验商品归属（确保商户只能查自己的商品SKU）
	var count int64
	err := DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM product WHERE id=? AND merchant_id=?
	`, productID, merchantID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ListSkuByProduct: 校验商品归属失败: %w", err)
	}
	if count == 0 {
		return nil, errors.New("ListSkuByProduct: 商品不存在或无权查询")
	}

	// 2. 查询该商品下的所有SKU
	rows, err := DB.QueryContext(ctx, `
		SELECT id, merchant_id, product_id, sku_code, price, stock FROM sku 
		WHERE product_id=? 
		ORDER BY id ASC
	`, productID)
	if err != nil {
		return nil, fmt.Errorf("ListSkuByProduct: 查询SKU列表失败: %w", err)
	}
	defer rows.Close()

	// 3. 解析结果
	skus := make([]*Sku, 0)
	for rows.Next() {
		var s Sku
		err := rows.Scan(&s.ID, &s.MerchantID, &s.ProductID, &s.SkuCode, &s.Price, &s.Stock)
		if err != nil {
			return nil, fmt.Errorf("ListSkuByProduct: 解析SKU失败: %w", err)
		}
		skus = append(skus, &s)
	}

	// 4. 检查rows遍历是否出错
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ListSkuByProduct: 遍历SKU失败: %w", err)
	}

	return skus, nil
}

// Merchant 商户结构体（精准匹配user表：ID对应id，Name对应username）
type Merchant struct {
	ID   int64  `db:"id"`   // 商户ID（user表的id，int类型转int64）
	Name string `db:"name"` // 商户名称（对应user表的username字段）
}

// ListMerchant 查询所有商户列表（适配真实user表结构）
// 逻辑：从user表筛选user_type=2的记录，提取id和username作为商户信息
func ListMerchant(ctx context.Context) ([]*Merchant, error) {
	// 前置检查：数据库连接是否初始化
	if DB == nil {
		return nil, errors.New("ListMerchant: 数据库连接未初始化（DB为nil）")
	}

	// 核心修正1：SQL查询字段改为 id + username（匹配user表的真实列名）
	// 筛选条件：user_type=2（商户）、status=1（正常状态，可选，根据业务补充）
	rows, err := DB.QueryContext(ctx, `
		SELECT id, username FROM user_account_db.user 
		WHERE user_type=2 AND status=1  -- 仅查正常状态的商户
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("ListMerchant: 查询商户列表失败: %w", err)
	}
	defer rows.Close() // 释放结果集，避免连接泄漏

	// 解析查询结果
	merchants := make([]*Merchant, 0)
	for rows.Next() {
		// 核心修正2：临时变量适配user表的字段类型（id是int，username是string）
		var (
			id       int    // 匹配user表的id（int类型）
			username string // 匹配user表的username（商户名称）
		)
		// Scan顺序严格匹配SQL的 id → username
		err := rows.Scan(&id, &username)
		if err != nil {
			return nil, fmt.Errorf("ListMerchant: 解析商户数据失败（id=%d）: %w", id, err)
		}

		// 转换为Merchant结构体（id转int64，username赋值给Name）
		merchants = append(merchants, &Merchant{
			ID:   int64(id), // int → int64，兼容结构体定义
			Name: username,  // user表的username = 商户名称
		})
	}

	// 检查遍历过程中的错误
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ListMerchant: 遍历商户列表失败: %w", err)
	}

	// 兜底：无商户数据时返回空切片（前端无需判nil）
	if len(merchants) == 0 {
		return []*Merchant{}, nil
	}

	return merchants, nil
}
