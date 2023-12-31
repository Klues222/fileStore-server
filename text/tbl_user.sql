CREATE TABLE `tbl_user`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `user_pwd` varchar(256) not null DEFAULT '' comment '用户encode密码',
    `email` varchar(64) DEFAULT '' COMMENT '邮箱',
    `phone` varchar(128) DEFAULT '' comment '手机号',
    `email_validated` tinyint(1) DEFAULT 0 comment '邮箱是否已验证',
    'phone_validated' tinyint(1) DEFAULT 0 comment '邮箱是否已验证',
    `signup_at` datetime default CURRENT_TIMESTAMP COMMENT '注册日期',
    `last_active` datetime default current_timestamp on update CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
    `profile` text comment '用户属性',
    `status` int(11) NOT NULL  DEFAULT '0'COMMENT '账户状态（启用 禁用 锁定 标记 标记删除等）',
    primary key (`id`),
    UNIQUE KEY `idx_phone` (`phone`),
    KEY `idx_status` (`status`)
)ENGINE = InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET =utf8mb4;