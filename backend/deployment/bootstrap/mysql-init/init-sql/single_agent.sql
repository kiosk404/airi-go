-- Create 'single_agent_draft' table
CREATE TABLE IF NOT EXISTS `airi_go`.`single_agent_draft` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary Key ID',
    `agent_id` bigint NOT NULL DEFAULT 0 COMMENT 'Agent ID',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT 'Agent Name',
    `description` text NULL COMMENT 'Agent Description',
    `icon_uri` varchar(255) NOT NULL DEFAULT '' COMMENT 'Icon URI',
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Create Time in Milliseconds',
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Update Time in Milliseconds',
    `deleted_at` datetime(3) NULL COMMENT 'delete time in millisecond',
    `variable` json NULL COMMENT 'variable',
    `model_info` json NULL COMMENT 'Model Configuration Information',
    `onboarding_info` json NULL COMMENT 'Onboarding Information',
    `prompt` json NULL COMMENT 'Agent Prompt Configuration',
    `plugin` json NULL COMMENT 'Agent Plugin Base Configuration',
    `knowledge` json NULL COMMENT 'Agent Knowledge Base Configuration',
    `workflow` json NULL COMMENT 'Agent Workflow Configuration',
    `suggest_reply` json NULL COMMENT 'Suggested Replies',
    `jump_config` json NULL COMMENT 'Jump Configuration',
    `background_image_info_list` json NULL COMMENT 'Background image',
    `database_config` json NULL COMMENT 'Agent Database Base Configuration',
    `bot_mode` tinyint NOT NULL DEFAULT 0 COMMENT 'bot mode,0:single mode 2:chatflow mode',
    `layout_info` text NULL COMMENT 'chatflow layout info',
    `shortcut_command` json NULL COMMENT 'shortcut command',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_agent_id` (`agent_id`)
) ENGINE=InnoDB CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT 'Single Agent Draft Copy Table';

-- Create 'single_agent_publish' table
CREATE TABLE IF NOT EXISTS `airi_go`.`single_agent_publish` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `agent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'agent_id',
    `publish_id` varchar(50) NOT NULL DEFAULT '' COMMENT 'publish id' COLLATE utf8mb4_general_ci,
    `version` varchar(255) NOT NULL DEFAULT '' COMMENT 'Agent Version',
    `publish_info` text NULL COMMENT 'publish info' COLLATE utf8mb4_general_ci,
    `publish_time` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'publish time',
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Create Time in Milliseconds',
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Update Time in Milliseconds',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT 'Status 0: In use 1: Delete 3: Disabled',
    `extra` json NULL COMMENT 'extra',
    PRIMARY KEY (`id`),
    INDEX `idx_agent_id_version` (`agent_id`, `version`),
    INDEX `idx_publish_id` (`publish_id`)
) ENGINE=InnoDB CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT 'Bot release version info';

-- Create 'single_agent_version' table
CREATE TABLE IF NOT EXISTS `airi_go`.`single_agent_version` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary Key ID',
    `agent_id` bigint NOT NULL DEFAULT 0 COMMENT 'Agent ID',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT 'Agent Name',
    `description` text NULL COMMENT 'Agent Description',
    `icon_uri` varchar(255) NOT NULL DEFAULT '' COMMENT 'Icon URI',
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Create Time in Milliseconds',
    `bot_mode` tinyint NOT NULL DEFAULT 0 COMMENT 'bot mode,0:single mode 2:chatflow mode',
    `layout_info` text NULL COMMENT 'chatflow layout info',
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Update Time in Milliseconds',
    `deleted_at` datetime(3) NULL COMMENT 'delete time in millisecond',
    `variable` json NULL COMMENT 'variable',
    `model_info` json NULL COMMENT 'Model Configuration Information',
    `onboarding_info` json NULL COMMENT 'Onboarding Information',
    `prompt` json NULL COMMENT 'Agent Prompt Configuration',
    `plugin` json NULL COMMENT 'Agent Plugin Base Configuration',
    `knowledge` json NULL COMMENT 'Agent Knowledge Base Configuration',
    `workflow` json NULL COMMENT 'Agent Workflow Configuration',
    `suggest_reply` json NULL COMMENT 'Suggested Replies',
    `jump_config` json NULL COMMENT 'Jump Configuration',
    `version` varchar(255) NOT NULL DEFAULT '' COMMENT 'Agent Version',
    `background_image_info_list` json NULL COMMENT 'Background image',
    `database_config` json NULL COMMENT 'Agent Database Base Configuration',
    `shortcut_command` json NULL COMMENT 'shortcut command',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_agent_id_and_version_id` (`agent_id`, `version`)
) ENGINE=InnoDB CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT 'Single Agent Version Copy Table';
