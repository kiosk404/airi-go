CREATE TABLE IF NOT EXISTS `airi_go`.`api_key` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT "Primary Key ID",
    `api_key` varchar(255) NOT NULL DEFAULT "" COMMENT "API Key hash",
    `ak_type` tinyint NOT NULL DEFAULT 0 COMMENT "AK Type",
    `name` varchar(255) NOT NULL DEFAULT "" COMMENT "API Key Name",
    `status` tinyint NOT NULL DEFAULT 0 COMMENT "0 normal, 1 deleted",
    `user_id` bigint NOT NULL DEFAULT 0 COMMENT "API Key Owner",
    `expired_at` bigint NOT NULL DEFAULT 0 COMMENT "API Key Expired Time",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time in Milliseconds",
    `last_used_at` bigint NOT NULL DEFAULT 0 COMMENT "Used Time in Milliseconds",
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "api key table";
