-- 切换到product库
USE product;

-- 存储过程：创建商品（带商户合法性校验）
DROP PROCEDURE IF EXISTS sp_create_product;
DELIMITER // -- 临时修改语句结束符（避免;冲突）
CREATE PROCEDURE sp_create_product(
    IN in_merchant_id BIGINT,  -- 商户ID
    IN in_product_name VARCHAR(128),  -- 商品名称
    IN in_ext JSON,  -- 扩展字段
    OUT out_result INT,  -- 返回结果：0=失败，1=成功
    OUT out_msg VARCHAR(255)  -- 返回提示信息
)
BEGIN
    -- 声明变量：存储商户的用户类型
    DECLARE merchant_user_type INT DEFAULT 0;
    
    -- 1. 校验参数非空
    IF in_merchant_id IS NULL OR in_product_name = '' THEN
        SET out_result = 0;
        SET out_msg = '商户ID和商品名称不能为空';
        LEAVE; -- 退出存储过程
    END IF;
    
    -- 2. 校验商户是否存在，且是商户类型（user_type=2）
    SELECT user_type INTO merchant_user_type 
    FROM `user` 
    WHERE id = in_merchant_id;
    
    -- 商户不存在
    IF merchant_user_type IS NULL THEN
        SET out_result = 0;
        SET out_msg = CONCAT('商户ID ', in_merchant_id, ' 不存在');
        LEAVE;
    END IF;
    
    -- 商户不是商户类型（是普通用户）
    IF merchant_user_type != 2 THEN
        SET out_result = 0;
        SET out_msg = CONCAT('用户ID ', in_merchant_id, ' 不是商户，无法创建商品');
        LEAVE;
    END IF;
    
    -- 3. 校验通过，创建商品
    INSERT INTO product (merchant_id, name, ext) 
    VALUES (in_merchant_id, in_product_name, in_ext);
    
    -- 4. 返回成功结果
    SET out_result = 1;
    SET out_msg = CONCAT('商品创建成功，商品ID：', LAST_INSERT_ID());
END //
DELIMITER ; -- 恢复语句结束符

-- 【作业展示用】调用示例（测试2个场景）
-- 场景1：合法商户（AlfredGit，id=2，user_type=2）创建商品
SET @result = 0;
SET @msg = '';
CALL sp_create_product(2, '华为Mate70 Pro', '{"category":"手机"}', @result, @msg);
SELECT @result AS 结果码, @msg AS 提示信息; -- 应返回：1 + 商品创建成功

-- 场景2：普通用户（youperceive，id=3，user_type=1）创建商品
SET @result = 0;
SET @msg = '';
CALL sp_create_product(3, '苹果16', '{"category":"手机"}', @result, @msg);
SELECT @result AS 结果码, @msg AS 提示信息; -- 应返回：0 + 不是商户，无法创建