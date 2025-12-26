-- 存储过程：扣减SKU库存
DROP PROCEDURE IF EXISTS sp_deduct_sku_stock;
DELIMITER //
CREATE PROCEDURE sp_deduct_sku_stock(
    IN in_sku_id BIGINT,  -- SKU ID
    IN in_deduct_count INT,  -- 扣减数量
    OUT out_result INT,  -- 0=失败，1=成功
    OUT out_msg VARCHAR(255),
    OUT out_remain_stock INT  -- 扣减后剩余库存
)
BEGIN
    -- 声明变量：当前库存
    DECLARE current_stock INT DEFAULT 0;
    
    -- 1. 校验参数
    IF in_sku_id IS NULL OR in_deduct_count <= 0 THEN
        SET out_result = 0;
        SET out_msg = 'SKU ID不能为空，扣减数量必须大于0';
        SET out_remain_stock = 0;
        LEAVE;
    END IF;
    
    -- 2. 查询当前库存
    SELECT stock INTO current_stock 
    FROM sku 
    WHERE id = in_sku_id;
    
    -- SKU不存在
    IF current_stock IS NULL THEN
        SET out_result = 0;
        SET out_msg = CONCAT('SKU ID ', in_sku_id, ' 不存在');
        SET out_remain_stock = 0;
        LEAVE;
    END IF;
    
    -- 3. 校验库存是否充足
    IF current_stock < in_deduct_count THEN
        SET out_result = 0;
        SET out_msg = CONCAT('库存不足，当前库存：', current_stock, '，需扣减：', in_deduct_count);
        SET out_remain_stock = current_stock;
        LEAVE;
    END IF;
    
    -- 4. 扣减库存
    UPDATE sku 
    SET stock = stock - in_deduct_count 
    WHERE id = in_sku_id;
    
    -- 5. 获取扣减后库存
    SELECT stock INTO out_remain_stock 
    FROM sku 
    WHERE id = in_sku_id;
    
    -- 6. 返回成功
    SET out_result = 1;
    SET out_msg = CONCAT('库存扣减成功，扣减数量：', in_deduct_count, '，剩余库存：', out_remain_stock);
END //
DELIMITER ;

-- 【作业展示用】调用示例
SET @result = 0;
SET @msg = '';
SET @remain_stock = 0;
-- 假设SKU ID=1有库存10，扣减3
CALL sp_deduct_sku_stock(1, 3, @result, @msg, @remain_stock);
SELECT @result AS 结果码, @msg AS 提示信息, @remain_stock AS 剩余库存;