CREATE TABLE `chats`  (
      `id` bigint NOT NULL,
      `name` varchar(255) NULL,
      `username` varchar(255) NULL,
      PRIMARY KEY (`id`),
      UNIQUE INDEX `id`(`id`) USING BTREE
);