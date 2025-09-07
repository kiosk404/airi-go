CREATE TABLE IF NOT EXISTS `airi-go`.`user`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT "Primary Key ID",
    `name` varchar(128) NOT NULL DEFAULT "" COMMENT "User Nickname",
    `unique_name` varchar(128) NOT NULL DEFAULT "" COMMENT "User Unique Name",
    `account` varchar(128) NOT NULL DEFAULT "" COMMENT "Account",
    `password` varchar(128) NOT NULL DEFAULT "" COMMENT "Password (Encrypted)",
    `description` varchar(512) NOT NULL DEFAULT "" COMMENT "User Description",
    `icon_uri` varchar(512) NOT NULL DEFAULT "" COMMENT "Avatar URI",
    `user_verified` bool NOT NULL DEFAULT 0 COMMENT "User Verification Status",
    `locale` varchar(128) NOT NULL DEFAULT "" COMMENT "Locale",
    `session_key` varchar(256) NOT NULL DEFAULT "" COMMENT "Session Key",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Creation Time (Milliseconds)",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time (Milliseconds)",
    `deleted_at` bigint unsigned NULL COMMENT "Deletion Time (Milliseconds)",
     PRIMARY KEY (`id`),
     UNIQUE INDEX `idx_email` (`account`),
     INDEX `idx_session_key` (`session_key`),
     UNIQUE INDEX `idx_unique_name` (`unique_name`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE utf8mb4_general_ci COMMENT "User 用户表";
