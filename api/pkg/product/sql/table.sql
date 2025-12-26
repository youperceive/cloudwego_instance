-- user 表
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `password` varchar(100) NOT NULL,
  `register_type` int NOT NULL DEFAULT '1' COMMENT '1=2=',
  `user_type` int NOT NULL DEFAULT '1' COMMENT '1=2=',
  `ext` json DEFAULT NULL,
  `status` int NOT NULL DEFAULT '1' COMMENT '1=0=',
  `created_at` int NOT NULL,
  `updated_at` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

-- product 表
CREATE TABLE `product` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `merchant_id` bigint NOT NULL COMMENT 'ID',
  `name` varchar(128) NOT NULL,
  `ext` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_merchant_id` (`merchant_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

-- sku 表
CREATE TABLE `sku` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `merchant_id` bigint NOT NULL COMMENT 'ID',
  `product_id` bigint NOT NULL COMMENT 'ID',
  `sku_code` varchar(64) NOT NULL,
  `price` bigint NOT NULL,
  `stock` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_sku_code` (`product_id`,`sku_code`),
  KEY `idx_merchant_id` (`merchant_id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci