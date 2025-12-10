import React from 'react';
import { Card, Tag, Button, Popconfirm, Space, Typography } from '@douyinfe/semi-ui';
import type { ModelListItem } from '@/services/models';

const { Text } = Typography;

interface ModelCardProps {
    model: ModelListItem;
    onEdit: (model: ModelListItem) => void;
    onDelete: (id: number) => void;
}

// 获取模型类型标签
const getModelTypeTag = (type: number) => {
    const config: Record<number, { label: string; color: 'blue' | 'green' | 'orange' }> = {
        0: { label: 'LLM', color: 'blue' },
        1: { label: 'Embedding', color: 'green' },
        2: { label: 'Rerank', color: 'orange' },
    };
    const { label, color } = config[type] || { label: 'Unknown', color: 'blue' as const };
    return <Tag color={color}>{label}</Tag>;
};

const ModelCard: React.FC<ModelCardProps> = ({ model, onEdit, onDelete }) => {
    return (
        <Card
            shadows="hover"
            style={{ height: '100%' }}
            headerStyle={{ padding: '16px 20px' }}
            bodyStyle={{ padding: '12px 20px' }}
            title={
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', width: '100%' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                        <img
                            src={model.providerIcon}
                            alt={model.providerName}
                            style={{ width: 28, height: 28, borderRadius: 4, objectFit: 'contain' }}
                        />
                        <div>
                            <div style={{ fontWeight: 600, fontSize: 16 }}>{model.name}</div>
                            <Text type="tertiary" size="small">{model.providerName}</Text>
                        </div>
                    </div>
                    {getModelTypeTag(model.type)}
                </div>
            }
            footerLine
            footerStyle={{ padding: '12px 20px' }}
            footer={
                <Space>
                    <Button
                        type="tertiary"
                        onClick={() => onEdit(model)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确定删除"
                        content="确定要删除这个模型吗？"
                        onConfirm={() => onDelete(model.id)}
                    >
                        <Button type="danger">
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            }
        >
            <div style={{ fontSize: 13, color: '#666' }}>
                <div style={{ marginBottom: 8 }}>
                    <Text type="tertiary">模型 ID：</Text>
                    <code style={{
                        background: '#f5f5f5',
                        padding: '2px 6px',
                        borderRadius: 4,
                        fontSize: 12
                    }}>
                        {model.modelId || '-'}
                    </code>
                </div>

                <div style={{ marginBottom: 8 }}>
                    <Text type="tertiary">Base URL：</Text>
                    <Text style={{ wordBreak: 'break-all' }}>
                        {model.baseUrl || '-'}
                    </Text>
                </div>

                {model.capabilities && (
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4, marginTop: 12 }}>
                        {model.capabilities.functionCall && (
                            <Tag size="small" color="blue">✓ 函数调用</Tag>
                        )}
                        {model.capabilities.imageUnderstanding && (
                            <Tag size="small" color="green">✓ 图像理解</Tag>
                        )}
                        {model.capabilities.cotDisplay && (
                            <Tag size="small" color="purple">✓ 思维链</Tag>
                        )}
                        {model.capabilities.multiModal && (
                            <Tag size="small" color="orange">✓ 多模态</Tag>
                        )}
                    </div>
                )}
            </div>
        </Card>
    );
};

export default ModelCard;
