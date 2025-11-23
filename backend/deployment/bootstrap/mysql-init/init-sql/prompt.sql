-- Create "prompt_resource" table
CREATE TABLE IF NOT EXISTS `airi_go`.`prompt_resource` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "主键ID",
    `name` varchar(255) NOT NULL COMMENT "名称",
    `description` varchar(255) NOT NULL COMMENT "描述",
    `prompt_text` mediumtext NULL COMMENT "prompt正文",
    `status` int NOT NULL COMMENT "状态,0无效,1有效",
    `creator_id` bigint NOT NULL COMMENT "创建者ID",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "创建时间",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "更新时间",
    PRIMARY KEY (`id`),
    INDEX `idx_creator_id` (`creator_id`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_general_ci COMMENT "prompt_resource";
