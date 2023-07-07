CREATE TABLE task
(
    `id`      CHAR(26)     NOT NULL COMMENT 'ULID 26bytes',
    `title`   VARCHAR(191) NOT NULL COMMENT 'タイトル',
    `content` TEXT         NOT NULL COMMENT '内容',
    `status`  ENUM('BACKLOG', 'PROGRESS', 'DONE') NOT NULL COMMENT 'ステータス',
    `created` DATETIME     NOT NULL COMMENT '作成時間',
    `updated` DATETIME     NOT NULL COMMENT '更新時間',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;