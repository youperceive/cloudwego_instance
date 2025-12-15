USE user_account_db

CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID（主键，对应IDL的i64）',
  `username` varchar(64) DEFAULT '' COMMENT '用户名（可选，对应IDL的username）',
  `email` varchar(128) DEFAULT NULL COMMENT '邮箱（可选，注册类型为EMAIL时必填）',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号（可选，注册类型为PHONE时必填）',
  `password` varchar(255) NOT NULL COMMENT '密码（BCrypt加密存储，对应IDL的password）',
  `register_type` tinyint NOT NULL COMMENT '注册类型：1-EMAIL(邮箱)，2-PHONE(手机号)（对应IDL的RegisterType）',
  `user_type` tinyint NOT NULL DEFAULT 1 COMMENT '用户类型：1-普通用户，2-管理员，3-第三方用户（对应IDL的UserType）',
  `ext` json COMMENT '扩展字段（键值对，对应IDL的map<string,string>）',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '用户状态：1-正常，2-禁用，3-注销（对应IDL的status）',
  `created_at` bigint NOT NULL COMMENT '创建时间戳（秒级，对应IDL的created_at）',
  `updated_at` bigint NOT NULL COMMENT '更新时间戳（秒级，对应IDL的updated_at）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_email` (`email`) COMMENT '邮箱唯一（避免重复注册）',
  UNIQUE KEY `uk_phone` (`phone`) COMMENT '手机号唯一（避免重复注册）',
  KEY `idx_username` (`username`) COMMENT '用户名索引（便于按用户名查询）',
  KEY `idx_status` (`status`) COMMENT '状态索引（便于筛选正常/禁用用户）'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户账户表';