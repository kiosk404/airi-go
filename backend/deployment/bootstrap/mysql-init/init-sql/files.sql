-- Create "files" table
CREATE TABLE IF NOT EXISTS `airi_go`.`files` (
    `id` bigint unsigned NOT NULL COMMENT "id",
    `name` varchar(255) NOT NULL DEFAULT "" COMMENT "file name",
    `file_size` bigint unsigned NOT NULL DEFAULT 0 COMMENT "file size",
    `tos_uri` varchar(1024) NOT NULL DEFAULT "" COMMENT "TOS URI",
    `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT "status，0invalid，1valid",
    `comment` varchar(1024) NOT NULL DEFAULT "" COMMENT "file comment",
    `source` tinyint unsigned NOT NULL DEFAULT 0 COMMENT "source：1 from API,",
    `creator_id` varchar(512) NOT NULL DEFAULT "" COMMENT "creator id",
    `content_type` varchar(255) NOT NULL DEFAULT "" COMMENT "content type",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time in Milliseconds",
    `deleted_at` datetime(3) NULL COMMENT "Delete Time", PRIMARY KEY (`id`),
    INDEX `idx_creator_id` (`creator_id`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_general_ci COMMENT "file resource table";
