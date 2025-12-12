import React, { useState } from 'react';
import { Card, Tag, Button, Popconfirm, Typography, Avatar } from '@douyinfe/semi-ui';
import type { AvatarProps } from '@douyinfe/semi-ui/lib/es/avatar';
import { IconDelete, IconEdit } from '@douyinfe/semi-icons';
import type { ModelListItem } from '@/services/models';

const { Text } = Typography;

interface ModelCardProps {
    model: ModelListItem;
    onEdit: (model: ModelListItem) => void;
    onDelete: (id: string) => void; // 使用字符串避免大整数精度丢失
}

// 模型类型配置
const MODEL_TYPE_CONFIG: Record<number, { label: string; color: string }> = {
    0: { label: 'LLM', color: 'blue' },
    1: { label: 'Embedding', color: 'cyan' },
    2: { label: 'Rerank', color: 'violet' },
};

// 能力标签配置
const CAPABILITY_CONFIG = [
    { key: 'functionCall', label: '函数调用', color: 'blue' },
    { key: 'imageUnderstanding', label: '图像理解', color: 'green' },
    { key: 'videoUnderstanding', label: '视频理解', color: 'teal' },
    { key: 'audioUnderstanding', label: '音频理解', color: 'amber' },
    { key: 'cotDisplay', label: '思维链', color: 'purple' },
    { key: 'multiModal', label: '多模态', color: 'orange' },
    { key: 'prefillResp', label: '预填充', color: 'indigo' },
] as const;

// 格式化 token 数量
const formatTokens = (tokens: number) => {
    if (tokens >= 1000000) return `${(tokens / 1000000).toFixed(1)}M`;
    if (tokens >= 1000) return `${(tokens / 1000).toFixed(0)}K`;
    return tokens.toString();
};

// 获取提供商首字母作为头像备选
const getProviderInitial = (name: string) => {
    return name?.charAt(0)?.toUpperCase() || 'M';
};

// 根据提供商名称生成颜色
const AVATAR_COLORS: AvatarProps['color'][] = ['blue', 'cyan', 'green', 'orange', 'red', 'purple', 'violet', 'pink'];
const getAvatarColor = (name: string): AvatarProps['color'] => {
    const index = name?.charCodeAt(0) % AVATAR_COLORS.length || 0;
    return AVATAR_COLORS[index];
};

const ModelCard: React.FC<ModelCardProps> = ({ model, onEdit, onDelete }) => {
    const [imgError, setImgError] = useState(false);

    const modelType = model.type ?? 0;
    const typeConfig = MODEL_TYPE_CONFIG[modelType] || { label: 'Unknown', color: 'grey' };

    // 获取启用的能力列表
    const enabledCapabilities = CAPABILITY_CONFIG.filter(
        cap => model.capabilities[cap.key as keyof typeof model.capabilities]
    );

    return (
        <Card
            shadows="hover"
            style={{
                height: '100%',
                borderRadius: 12,
                overflow: 'hidden',
            }}
            bodyStyle={{ padding: 0 }}
        >
            {/* 头部区域 */}
            <div style={{
                padding: '20px 20px 16px',
                background: 'linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%)',
                borderBottom: '1px solid #e2e8f0',
            }}>
                <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
                    {/* 提供商头像 */}
                    {!imgError && model.providerIcon ? (
                        <img
                            src={model.providerIcon}
                            alt={model.providerName}
                            onError={() => setImgError(true)}
                            style={{
                                width: 40,
                                height: 40,
                                borderRadius: 10,
                                objectFit: 'contain',
                                background: '#fff',
                                padding: 4,
                                boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
                            }}
                        />
                    ) : (
                        <Avatar
                            size="default"
                            color={getAvatarColor(model.providerName)}
                            style={{
                                width: 40,
                                height: 40,
                                borderRadius: 10,
                                fontSize: 18,
                                fontWeight: 600,
                            }}
                        >
                            {getProviderInitial(model.providerName)}
                        </Avatar>
                    )}

                    {/* 模型信息 */}
                    <div style={{ flex: 1, minWidth: 0 }}>
                        <div style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 8,
                            marginBottom: 4,
                        }}>
                            <Text
                                strong
                                style={{
                                    fontSize: 16,
                                    overflow: 'hidden',
                                    textOverflow: 'ellipsis',
                                    whiteSpace: 'nowrap',
                                }}
                                title={model.name}
                            >
                                {model.name}
                            </Text>
                            <Tag size="small" color={typeConfig.color as any}>
                                {typeConfig.label}
                            </Tag>
                        </div>
                        <Text type="tertiary" size="small">
                            {model.providerName}
                        </Text>
                    </div>
                </div>
            </div>

            {/* 内容区域 */}
            <div style={{ padding: '16px 20px' }}>
                {/* 模型 ID */}
                <div style={{ marginBottom: 12 }}>
                    <Text type="tertiary" size="small" style={{ display: 'block', marginBottom: 4 }}>
                        模型 ID
                    </Text>
                    <code style={{
                        display: 'block',
                        background: '#f4f5f7',
                        padding: '8px 12px',
                        borderRadius: 6,
                        fontSize: 13,
                        color: '#1e293b',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                    }}>
                        {model.modelId || '-'}
                    </code>
                </div>

                {/* 上下文长度 */}
                {model.maxTokens ? (
                    <div style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: 6,
                        background: '#eff6ff',
                        padding: '4px 10px',
                        borderRadius: 6,
                        marginBottom: 12,
                    }}>
                        <Text type="tertiary" size="small">上下文</Text>
                        <Text strong size="small" style={{ color: '#3b82f6' }}>
                            {formatTokens(model.maxTokens)}
                        </Text>
                    </div>
                ) : null}

                {/* 能力标签 */}
                {enabledCapabilities.length > 0 && (
                    <div style={{
                        display: 'flex',
                        flexWrap: 'wrap',
                        gap: 6,
                        marginTop: 8,
                    }}>
                        {enabledCapabilities.map(cap => (
                            <Tag
                                key={cap.key}
                                size="small"
                                color={cap.color as any}
                                style={{ borderRadius: 4 }}
                            >
                                {cap.label}
                            </Tag>
                        ))}
                    </div>
                )}
            </div>

            {/* 底部操作区 */}
            <div style={{
                padding: '12px 20px',
                borderTop: '1px solid #f1f5f9',
                display: 'flex',
                justifyContent: 'flex-end',
                gap: 8,
            }}>
                <Button
                    type="tertiary"
                    size="small"
                    icon={<IconEdit />}
                    onClick={() => onEdit(model)}
                >
                    编辑
                </Button>
                <Popconfirm
                    title="确定删除此模型？"
                    content="删除后无法恢复"
                    onConfirm={() => onDelete(model.id)}
                >
                    <Button
                        type="danger"
                        size="small"
                        icon={<IconDelete />}
                    >
                        删除
                    </Button>
                </Popconfirm>
            </div>
        </Card>
    );
};

export default ModelCard;
