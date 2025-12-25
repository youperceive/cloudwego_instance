-- 1. 商户表（merchant）
CREATE TABLE `merchant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '商户ID',
  `name` varchar(64) NOT NULL COMMENT '商户名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 商品表（product）
CREATE TABLE `product` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `merchant_id` bigint(20) NOT NULL COMMENT '商户ID',
  `name` varchar(128) NOT NULL COMMENT '商品名称',
  `ext` json DEFAULT NULL COMMENT '扩展字段',
  PRIMARY KEY (`id`),
  KEY `idx_merchant_id` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. SKU表（sku）
CREATE TABLE `sku` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `merchant_id` bigint(20) NOT NULL COMMENT '商户ID',
  `product_id` bigint(20) NOT NULL COMMENT '商品ID',
  `sku_code` varchar(64) NOT NULL COMMENT '规格编码',
  `price` bigint(20) NOT NULL COMMENT '单价（分）',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '库存',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_sku_code` (`product_id`,`sku_code`),
  KEY `idx_merchant_id` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入测试商户（方便调试）
INSERT INTO `merchant` (`name`) VALUES ('测试商户');