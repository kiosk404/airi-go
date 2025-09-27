-- Create "conversation" table
CREATE TABLE IF NOT EXISTS `airi_go`.`conversation` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "主键ID",
    `agent_id` bigint NOT NULL DEFAULT 0 COMMENT "agent_id",
    `scene` tinyint NOT NULL DEFAULT 0 COMMENT "会话场景",
    `section_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "最新section_id",
    `creator_id` bigint unsigned NULL DEFAULT 0 COMMENT "创建者id",
    `ext` text NULL COMMENT "扩展字段",
    `name` varchar(255) NOT NULL DEFAULT "" COMMENT "conversation name",
    `status` tinyint NOT NULL DEFAULT 1 COMMENT "status: 1-normal 2-deleted",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "创建时间",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "更新时间",
    PRIMARY KEY (`id`),
    INDEX `idx_bot_status` (`agent_id`, `creator_id`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "会话信息表";
-- Create "message" table
CREATE TABLE IF NOT EXISTS `airi_go`.`message` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "主键ID",
    `run_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "对应的run_id",
    `conversation_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "conversation id",
    `user_id` varchar(60) NOT NULL DEFAULT "" COMMENT "user id",
    `agent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "agent_id",
    `role` varchar(100) NOT NULL DEFAULT "" COMMENT "角色: user、assistant、system",
    `content_type` varchar(100) NOT NULL DEFAULT "" COMMENT "内容类型 1 text",
    `content` mediumtext NULL COMMENT "内容",
    `message_type` varchar(100) NOT NULL DEFAULT "" COMMENT "消息类型：",
    `display_content` text NULL COMMENT "展示内容",
    `ext` text NULL COMMENT "message 扩展字段" COLLATE utf8mb4_general_ci,
    `section_id` bigint unsigned NULL COMMENT "段落id",
    `broken_position` int NULL DEFAULT -1 COMMENT "打断位置",
    `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT "消息状态 1 Available 2 Deleted 3 Replaced 4 Broken 5 Failed 6 Streaming 7 Pending",
    `model_content` mediumtext NULL COMMENT "模型输入内容",
    `meta_info` text NULL COMMENT "引用、高亮等文本标记信息",
    `reasoning_content` text NULL COMMENT "思考内容" COLLATE utf8mb4_general_ci,
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "创建时间",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "更新时间",
    PRIMARY KEY (`id`),
    INDEX `idx_conversation_id` (`conversation_id`),
    INDEX `idx_run_id` (`run_id`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "消息表";
-- Create "run_record" table
CREATE TABLE IF NOT EXISTS `airi_go`.`run_record` (
    `id` bigint unsigned NOT NULL COMMENT "主键ID",
    `conversation_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "会话 ID",
    `section_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "section ID",
    `agent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "agent_id",
    `user_id` varchar(255) NOT NULL DEFAULT "" COMMENT "user id",
    `source` tinyint unsigned NOT NULL DEFAULT 0 COMMENT "执行来源 0 API,",
    `token_count` int NOT NULL DEFAULT 0 COMMENT "token 消耗",
    `usage` json NULL COMMENT "usage",
    `output_tokens` int NOT NULL DEFAULT 0 COMMENT "消耗的 output token 数",
    `input_tokens` int NOT NULL DEFAULT 0 COMMENT "消耗的 input token 数",
    `status` varchar(255) NOT NULL DEFAULT "" COMMENT "状态,0 Unknown, 1-Created,2-InProgress,3-Completed,4-Failed,5-Expired,6-Cancelled,7-RequiresAction", `creator_id` bigint NOT NULL DEFAULT 0 COMMENT "创建者标识",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "创建时间",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "更新时间",
    `failed_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "失败时间",
    `last_error` text NULL COMMENT "error message" COLLATE utf8mb4_general_ci,
    `completed_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "结束时间",
    `chat_request` text NULL COMMENT "保存原始请求的部分字段" COLLATE utf8mb4_general_ci,
    `ext` text NULL COMMENT "扩展字段" COLLATE utf8mb4_general_ci,
    PRIMARY KEY (`id`),
    INDEX `idx_c_s` (`conversation_id`, `section_id`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "执行记录表";
