-- 触发器：插入SKU前校验商户和商品合法性
DROP TRIGGER IF EXISTS trg_sku_before_insert;
DELIMITER //
CREATE TRIGGER trg_sku_before_insert
BEFORE INSERT ON sku
FOR EACH ROW  -- 行级触发器（每插入一行触发）
BEGIN
    DECLARE product_merchant_id BIGINT DEFAULT 0;
    DECLARE merchant_user_type INT DEFAULT 0;
    
    -- 1. 校验SKU的merchant_id是商户（user_type=2）
    SELECT user_type INTO merchant_user_type 
    FROM `user` 
    WHERE id = NEW.merchant_id;
    
    IF merchant_user_type != 2 THEN
        SIGNAL SQLSTATE '45000' -- 自定义异常码
        SET MESSAGE_TEXT = CONCAT('SKU的商户ID ', NEW.merchant_id, ' 不是商户类型，禁止插入');
    END IF;
    
    -- 2. 校验SKU的product_id属于该商户
    SELECT merchant_id INTO product_merchant_id 
    FROM product 
    WHERE id = NEW.product_id;
    
    IF product_merchant_id != NEW.merchant_id THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = CONCAT('商品ID ', NEW.product_id, ' 不属于商户ID ', NEW.merchant_id, '，禁止插入SKU');
    END IF;
END //
DELIMITER ;

-- 【作业展示用】测试触发器
-- 场景1：合法插入（商户2的商品，插入SKU）
INSERT INTO sku (merchant_id, product_id, sku_code, price, stock) 
VALUES (2, 1, 'HUAWEI-MATE70-256G', 599900, 100); -- 成功

-- 场景2：非法插入（普通用户3的商品，触发异常）
INSERT INTO sku (merchant_id, product_id, sku_code, price, stock) 
VALUES (3, 1, 'IPHONE16-256G', 699900, 50); -- 触发触发器，提示“不是商户类型”