CREATE TABLE `im_user` (
   `user_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
   `head_url` varchar(255) DEFAULT NULL COMMENT '用户头像',
   `mobile` varchar(11) NOT NULL COMMENT '用户手机',
   `password` varchar(150) NOT NULL COMMENT '密码串',
   `token` varchar(255) DEFAULT NULL COMMENT '用户密钥',
   `status` tinyint(4) DEFAULT NULL COMMENT '用户状态',
   `nick_name` varchar(30) NOT NULL COMMENT '用户昵称',
   `bind_birth` tinyint(1) DEFAULT '0' COMMENT '是否已绑定生日',
   `gender` varchar(10) DEFAULT 'male' COMMENT 'female女,male男,unknown未知',
   `login_ip` varchar(60) DEFAULT NULL COMMENT '登录IP地址',
   `login_time` datetime DEFAULT NULL COMMENT '登录时间',
   `user_growth` tinyint(4) DEFAULT '0' COMMENT '用户成长值',
   `user_picture` varchar(255) DEFAULT NULL COMMENT '用户头像',
   `created_at` datetime DEFAULT NULL COMMENT '创建时间',
   `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`user_id`),
   UNIQUE KEY `mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE `im_contact` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `ownerid` bigint(20) DEFAULT '0' COMMENT '用户ID',
  `dstobj` bigint(20) DEFAULT '0' COMMENT '对端ID',
  `cate` tinyint(4) DEFAULT '0' COMMENT '对端类型,0好友,1群聊',
  `memo` varchar(30) DEFAULT NULL COMMENT '备注信息',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `im_user` (
   `user_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
   `head_url` varchar(255) DEFAULT NULL COMMENT '用户头像',
   `mobile` varchar(11) NOT NULL COMMENT '用户手机',
   `password` varchar(150) NOT NULL COMMENT '密码串',
   `token` varchar(255) DEFAULT NULL COMMENT '用户密钥',
   `status` tinyint(4) DEFAULT NULL COMMENT '用户状态',
   `nick_name` varchar(30) NOT NULL COMMENT '用户昵称',
   `bind_birth` tinyint(1) DEFAULT '0' COMMENT '是否已绑定生日',
   `gender` varchar(10) DEFAULT 'male' COMMENT 'female女,male男,unknown未知',
   `login_ip` varchar(60) DEFAULT NULL COMMENT '登录IP地址',
   `login_time` datetime DEFAULT NULL COMMENT '登录时间',
   `user_growth` tinyint(4) DEFAULT '0' COMMENT '用户成长值',
   `user_picture` varchar(255) DEFAULT NULL COMMENT '用户头像',
   `created_at` datetime DEFAULT NULL COMMENT '创建时间',
   `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`user_id`),
   UNIQUE KEY `mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

INSERT INTO `im_user` (`user_id`, `head_url`, `mobile`, `password`, `token`, `status`, `nick_name`, `bind_birth`, `gender`, `login_ip`, `login_time`, `user_growth`, `user_picture`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, '/asset/upload/16837048221298498081.png', '13678910090', '123456', '0ca175b9c0f726a831d895e269332461', 0, '136****0090', 0, 'male', '', NULL, 0, NULL, NULL, '2023-05-10 15:48:19', NULL);
INSERT INTO `im_user` (`user_id`, `head_url`, `mobile`, `password`, `token`, `status`, `nick_name`, `bind_birth`, `gender`, `login_ip`, `login_time`, `user_growth`, `user_picture`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, '', '13678910091', '123456', '2bf2ac69acb234d24fa9a353b54d0573', 0, '136****0091', 0, 'female', '', NULL, 0, NULL, '2023-05-07 15:55:20', '2023-05-07 15:55:20', NULL);
