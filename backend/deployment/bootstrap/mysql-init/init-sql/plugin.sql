-- Create "agent_tool_draft" table
CREATE TABLE IF NOT EXISTS `airi_go`.`agent_tool_draft` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Primary Key ID",
    `agent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Agent ID",
    `plugin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `tool_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Tool ID",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `sub_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Sub URL Path",
    `method` varchar(64) NOT NULL DEFAULT "" COMMENT "HTTP Request Method",
    `tool_name` varchar(255) NOT NULL DEFAULT "" COMMENT "Tool Name",
    `tool_version` varchar(255) NOT NULL DEFAULT "" COMMENT "Tool Version, e.g. v1.0.0",
    `operation` json NULL COMMENT "Tool Openapi Operation Schema",
    PRIMARY KEY (`id`),
    INDEX `idx_agent_plugin_tool` (`agent_id`, `plugin_id`, `tool_id`),
    INDEX `idx_agent_tool_bind` (`agent_id`, `created_at`),
    UNIQUE INDEX `uniq_idx_agent_tool_id` (`agent_id`, `tool_id`),
    UNIQUE INDEX `uniq_idx_agent_tool_name` (`agent_id`, `tool_name`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Draft Agent Tool";
-- Create "agent_tool_version" table
CREATE TABLE IF NOT EXISTS `airi_go`.`agent_tool_version` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Primary Key ID",
    `agent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Agent ID",
    `plugin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `tool_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Tool ID",
    `agent_version` varchar(255) NOT NULL DEFAULT "" COMMENT "Agent Tool Version",
    `tool_name` varchar(255) NOT NULL DEFAULT "" COMMENT "Tool Name",
    `tool_version` varchar(255) NOT NULL DEFAULT "" COMMENT "Tool Version, e.g. v1.0.0",
    `sub_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Sub URL Path",
    `method` varchar(64) NOT NULL DEFAULT "" COMMENT "HTTP Request Method",
    `operation` json NULL COMMENT "Tool Openapi Operation Schema",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    PRIMARY KEY (`id`),
    INDEX `idx_agent_tool_id_created_at` (`agent_id`, `tool_id`, `created_at`),
    INDEX `idx_agent_tool_name_created_at` (`agent_id`, `tool_name`, `created_at`),
    UNIQUE INDEX `uniq_idx_agent_tool_id_agent_version` (`agent_id`, `tool_id`, `agent_version`),
    UNIQUE INDEX `uniq_idx_agent_tool_name_agent_version` (`agent_id`, `tool_name`, `agent_version`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Agent Tool Version";
-- Create "plugin" table
CREATE TABLE IF NOT EXISTS `airi_go`.`plugin` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `developer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Developer ID",
    `app_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Application ID",
    `icon_uri` varchar(512) NOT NULL DEFAULT "" COMMENT "Icon URI",
    `server_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Server URL",
    `plugin_type` tinyint NOT NULL DEFAULT 0 COMMENT "Plugin Type, 1:http, 6:local",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time in Milliseconds",
    `version` varchar(255) NOT NULL DEFAULT "" COMMENT "Plugin Version, e.g. v1.0.0",
    `version_desc` text NULL COMMENT "Plugin Version Description",
    `manifest` json NULL COMMENT "Plugin Manifest",
    `openapi_doc` json NULL COMMENT "OpenAPI Document, only stores the root",
    PRIMARY KEY (`id`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_updated_at` (`updated_at`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Latest Plugin";
-- Create "plugin_draft" table
CREATE TABLE IF NOT EXISTS `airi_go`.`plugin_draft` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `developer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Developer ID",
    `app_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Application ID",
    `icon_uri` varchar(512) NOT NULL DEFAULT "" COMMENT "Icon URI",
    `server_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Server URL",
    `plugin_type` tinyint NOT NULL DEFAULT 0 COMMENT "Plugin Type, 1:http, 6:local",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time in Milliseconds",
    `deleted_at` datetime NULL COMMENT "Delete Time",
    `manifest` json NULL COMMENT "Plugin Manifest",
    `openapi_doc` json NULL COMMENT "OpenAPI Document, only stores the root",
    PRIMARY KEY (`id`),
    INDEX `idx_app_id` (`app_id`, `id`),
    INDEX `idx_app_created_at` (`app_id`, `created_at`),
    INDEX `idx_app_updated_at` (`app_id`, `updated_at`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Draft Plugin";
-- Create "plugin_oauth_auth" table
CREATE TABLE IF NOT EXISTS `airi_go`.`plugin_oauth_auth` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Primary Key",
    `user_id` varchar(255) NOT NULL DEFAULT "" COMMENT "User ID",
    `plugin_id` bigint NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `is_draft` bool NOT NULL DEFAULT 0 COMMENT "Is Draft Plugin",
    `oauth_config` json NULL COMMENT "Authorization Code OAuth Config",
    `access_token` varchar(1024) NOT NULL DEFAULT "" COMMENT "Access Token",
    `refresh_token` varchar(1024) NOT NULL DEFAULT "" COMMENT "Refresh Token",
    `token_expired_at` bigint NULL COMMENT "Token Expired in Milliseconds",
    `next_token_refresh_at` bigint NULL COMMENT "Next Token Refresh Time in Milliseconds",
    `last_active_at` bigint NULL COMMENT "Last active time in Milliseconds",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time in Milliseconds",
    PRIMARY KEY (`id`),
    INDEX `idx_last_active_at` (`last_active_at`),
    INDEX `idx_last_token_expired_at` (`token_expired_at`),
    INDEX `idx_next_token_refresh_at` (`next_token_refresh_at`),
    UNIQUE INDEX `uniq_idx_user_plugin_is_draft` (`user_id`, `plugin_id`, `is_draft`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Plugin OAuth Authorization Code Info";
-- Create "plugin_version" table
CREATE TABLE IF NOT EXISTS `airi_go`.`plugin_version` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Primary Key ID",
    `developer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Developer ID",
    `plugin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `app_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Application ID",
    `icon_uri` varchar(512) NOT NULL DEFAULT "" COMMENT "Icon URI",
    `server_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Server URL",
    `plugin_type` tinyint NOT NULL DEFAULT 0 COMMENT "Plugin Type, 1:http, 6:local",
    `version` varchar(255) NOT NULL DEFAULT "" COMMENT "Plugin Version, e.g. v1.0.0",
    `version_desc` text NULL COMMENT "Plugin Version Description",
    `manifest` json NULL COMMENT "Plugin Manifest",
    `openapi_doc` json NULL COMMENT "OpenAPI Document, only stores the root",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `deleted_at` datetime NULL COMMENT "Delete Time",
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_idx_plugin_version` (`plugin_id`, `version`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Plugin Version";
-- Create "tool" table
CREATE TABLE IF NOT EXISTS `airi_go`.`tool` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Tool ID",
    `plugin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time in Milliseconds",
    `version` varchar(255) NOT NULL DEFAULT "" COMMENT "Tool Version, e.g. v1.0.0",
    `sub_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Sub URL Path",
    `method` varchar(64) NOT NULL DEFAULT "" COMMENT "HTTP Request Method",
    `operation` json NULL COMMENT "Tool Openapi Operation Schema",
    `activated_status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT "0:activated; 1:deactivated",
    PRIMARY KEY (`id`),
    INDEX `idx_plugin_activated_status` (`plugin_id`, `activated_status`),
    UNIQUE INDEX `uniq_idx_plugin_sub_url_method` (`plugin_id`, `sub_url`, `method`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Latest Tool";
-- Create "tool_draft" table
CREATE TABLE IF NOT EXISTS `airi_go`.`tool_draft` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Tool ID",
    `plugin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Update Time in Milliseconds",
    `sub_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Sub URL Path",
    `method` varchar(64) NOT NULL DEFAULT "" COMMENT "HTTP Request Method",
    `operation` json NULL COMMENT "Tool Openapi Operation Schema",
    `debug_status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT "0:not pass; 1:pass",
    `activated_status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT "0:activated; 1:deactivated",
    PRIMARY KEY (`id`),
    INDEX `idx_plugin_created_at_id` (`plugin_id`, `created_at`, `id`),
    UNIQUE INDEX `uniq_idx_plugin_sub_url_method` (`plugin_id`, `sub_url`, `method`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Draft Tool";
-- Create "tool_version" table
CREATE TABLE IF NOT EXISTS `airi_go`.`tool_version` (
    `id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Primary Key ID",
    `tool_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Tool ID",
    `plugin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Plugin ID",
    `version` varchar(255) NOT NULL DEFAULT "" COMMENT "Tool Version, e.g. v1.0.0",
    `sub_url` varchar(512) NOT NULL DEFAULT "" COMMENT "Sub URL Path",
    `method` varchar(64) NOT NULL DEFAULT "" COMMENT "HTTP Request Method",
    `operation` json NULL COMMENT "Tool Openapi Operation Schema",
    `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT "Create Time in Milliseconds",
    `deleted_at` datetime NULL COMMENT "Delete Time",
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_idx_tool_version` (`tool_id`, `version`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_unicode_ci COMMENT "Tool Version";
-- Create "user" table