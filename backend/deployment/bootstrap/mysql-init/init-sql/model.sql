-- Create "model_entity" table
CREATE TABLE IF NOT EXISTS `airi_go`.`model_entity`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `meta_id`           bigint unsigned NOT NULL            COMMENT '模型元信息 id',
    `name`              varchar(128)    NOT NULL            COMMENT '名称',
    `is_selected`       bool            NOT NULL DEFAULT 0  COMMENT '是否选中',
    `description`       text            NULL                COMMENT '描述',
    `default_params`    json            NOT NULL            COMMENT '默认参数',
    `scenario`          bigint unsigned NOT NULL            COMMENT '模型应用场景',
    `status`            int             NOT NULL DEFAULT 1  COMMENT '模型状态',
    `created_at`        bigint unsigned NOT NULL DEFAULT 0  COMMENT 'Create Time in Milliseconds',
    `updated_at`        bigint unsigned NOT NULL DEFAULT 0  COMMENT 'Update Time in Milliseconds',
    `deleted_at`        bigint unsigned NULL                COMMENT 'Delete Time in Milliseconds',
    PRIMARY KEY (`id`),
    INDEX `idx_scenario` (`scenario`),
    INDEX `idx_status` (`status`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_general_ci COMMENT "模型信息";
-- Create "model_meta" table
CREATE TABLE IF NOT EXISTS `airi_go`.`model_meta`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `model_name`        varchar(128)    NOT NULL            COMMENT '模型名称',
    `protocol`          varchar(128)    NOT NULL            COMMENT '模型协议',
    `icon_uri`          varchar(255)    NOT NULL DEFAULT '' COMMENT 'Icon URI',
    `icon_url`          varchar(255)    NOT NULL DEFAULT '' COMMENT 'Icon URL',
    `capability`        json            NULL                COMMENT '模型能力',
    `conn_config`       json            NULL                COMMENT '模型连接配置',
    `status`            int             NOT NULL DEFAULT 1  COMMENT '模型状态',
    `description`       varchar(2048)   NOT NULL DEFAULT '' COMMENT '模型描述',
    `created_at`        bigint unsigned NOT NULL DEFAULT 0  COMMENT 'Create Time in Milliseconds',
    `updated_at`        bigint unsigned NOT NULL DEFAULT 0  COMMENT 'Update Time in Milliseconds',
    `deleted_at`        bigint unsigned NULL                COMMENT 'Delete Time in Milliseconds',
    PRIMARY KEY (`id`),
    INDEX `idx_status` (`status`)
) ENGINE = InnoDB
DEFAULT CHARSET utf8mb4
COLLATE utf8mb4_general_ci COMMENT "模型元信息";
-- Create "model_request_record" table
CREATE TABLE IF NOT EXISTS `airi_go`.`model_request_record`
(
    `id`                    bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键ID',
    `user_id`               varchar(256)    NOT NULL DEFAULT '' COMMENT 'user id',
    `usage_scene`           varchar(128)    NOT NULL DEFAULT '' COMMENT '场景',
    `usage_scene_entity_id` varchar(256)    NOT NULL DEFAULT '' COMMENT '场景实体id',
    `protocol`              varchar(128)    NOT NULL DEFAULT '' COMMENT '使用的协议，如ark/deepseek等',
    `model_identification`  varchar(1024)   NOT NULL DEFAULT '' COMMENT '模型唯一标识',
    `model_ak`              varchar(1024)   NOT NULL DEFAULT '' COMMENT '模型的AK',
    `model_id`              varchar(256)    NOT NULL DEFAULT '' COMMENT 'model id',
    `model_name`            varchar(1024)   NOT NULL DEFAULT '' COMMENT '模型展示名称',
    `input_token`           bigint unsigned NOT NULL DEFAULT '0' COMMENT '输入token数量',
    `output_token`          bigint unsigned NOT NULL DEFAULT '0' COMMENT '输出token数量',
    `logid`                 varchar(128)    NOT NULL DEFAULT '' COMMENT 'logid',
    `error_code`            varchar(128)    NOT NULL DEFAULT '' COMMENT 'error_code',
    `error_msg`             text COLLATE utf8mb4_general_ci COMMENT 'error_msg',
    `created_at`            datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`            datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_create_time` (`created_at`) USING BTREE COMMENT 'create_time'
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_general_ci COMMENT ='模型流量记录表';
-- Create "model_instance" table
CREATE TABLE IF NOT EXISTS `airi_go`.`model_instance` (
    `id`                    bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `type`                  tinyint         NOT NULL                COMMENT 'Model Type 0-LLM 1-TextEmbedding 2-Rerank',
    `provider`              json            NOT NULL                COMMENT 'Provider Information',
    `display_info`          json            NOT NULL                COMMENT 'Display Information',
    `is_selected`           bool            NOT NULL DEFAULT 0      COMMENT 'Selected',
    `connection`            json            NOT NULL                COMMENT 'Connection Information',
    `capability`            json            NOT NULL                COMMENT 'Model Capability',
    `parameters`            json            NOT NULL                COMMENT 'Model Parameters',
    `extra`                 json            NULL                    COMMENT 'Extra Information',
    `created_at`            bigint unsigned NOT NULL DEFAULT 0      COMMENT 'Create Time in Milliseconds',
    `updated_at`            bigint unsigned NOT NULL DEFAULT 0      COMMENT 'Update Time in Milliseconds',
    `deleted_at`            datetime(3)     NULL                    COMMENT 'Delete Time',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_general_ci COMMENT = "Model Instance Management Table";
