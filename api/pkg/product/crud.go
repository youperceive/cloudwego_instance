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
