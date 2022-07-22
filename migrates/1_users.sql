CREATE TABLE `users`  (
        `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        `telegram_id` bigint(20) NULL,
        `user_name` varchar(255) NULL,
        `user_first_name` varchar(255) NULL,
        `user_last_name` varchar(255) NULL,
        `instagram` varchar(255) NULL,
        `twitter` int NULL,
        `status` varchar(255) NULL,
        `referred_by` varchar(255) NULL,
        `registered_at` datetime NOT NULL,
        PRIMARY KEY (`id`),
        UNIQUE INDEX `telegram_id`(`telegram_id`) USING BTREE
);