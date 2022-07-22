ALTER TABLE `polls_result`
    MODIFY COLUMN `poll` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL AFTER `telegram_id`;