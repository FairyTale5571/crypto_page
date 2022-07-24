CREATE TABLE `polls_result`  (
    `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `telegram_id` bigint(20) UNSIGNED NOT NULL,
    `poll` varchar(255) NULL,
    `result` json NULL,
    `insert_time` datetime NULL,
    PRIMARY KEY (`id`, `telegram_id`),
    INDEX `telegram_id`(`telegram_id`) USING BTREE,
    CONSTRAINT `tg_ifbk2` FOREIGN KEY (`telegram_id`) REFERENCES `users` (`telegram_id`) ON DELETE CASCADE ON UPDATE RESTRICT
);