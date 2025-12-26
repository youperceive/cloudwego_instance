-- 触发器：更新user表自动同步updated_at
DROP TRIGGER IF EXISTS trg_user_before_update;
DELIMITER //
CREATE TRIGGER trg_user_before_update
BEFORE UPDATE ON `user`
FOR EACH ROW
BEGIN
    -- 把updated_at设为当前时间戳（UNIX秒数，和你的表结构一致）
    SET NEW.updated_at = UNIX_TIMESTAMP(NOW());
END //
DELIMITER ;

-- 【作业展示用】测试触发器
-- 更新前查看updated_at（比如更新id=2的用户）
SELECT id, username, updated_at FROM `user` WHERE id=2;

-- 更新用户信息
UPDATE `user` SET phone = '18777411453' WHERE id=2;

-- 更新后查看，updated_at已变为当前时间戳
SELECT id, username, phone, updated_at FROM `user` WHERE id=2;