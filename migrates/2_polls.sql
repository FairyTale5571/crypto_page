CREATE TABLE `polls_result`  (
     `id` bigint(20) UNSIGNED NOT NULL,
     `telegram_id` bigint(20) UNSIGNED NOT NULL,
     `poll` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
     `result` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
     `insert_time` datetime NULL DEFAULT NULL,
     PRIMARY KEY (`id`, `telegram_id`) USING BTREE,
     INDEX `telegram_id`(`telegram_id`) USING BTREE,
     CONSTRAINT `tg_ifbk2` FOREIGN KEY (`telegram_id`) REFERENCES `users` (`telegram_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;