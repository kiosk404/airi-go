import { Button, Typography, Space, Select, Tag, Collapse } from "tdesign-react"
import { AddIcon, HelpCircleIcon, CloudUploadIcon } from "tdesign-icons-react"
import React, { useState, useMemo } from "react"
import type { BotInfo, DraftBotDisplayInfoData, BotOptionData, ModelDetail } from '@/services/draftbot';

const { Title, Paragraph } = Typography;
const { Option } = Select;

// 分组标题组件
const SectionTitle = ({ title, extra }: { title: string; extra?: React.ReactNode }) => (
    <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '16px 0 8px 0',
        color: '#666',
        fontSize: '14px',
        fontWeight: 500,
    }}>
        <span>{title}</span>
        {extra}
    </div>
);

// 折叠项头部组件
const PanelHeader = ({ title, extra }: { title: string; extra?: React.ReactNode }) => (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%'}}>
        <span style={{ lineHeight: '22px' }}>{title}</span>
        {extra || <Button shape="square" variant="text" icon={<AddIcon />} onClick={(e) => e.stopPropagation()} />}
    </div>
);

interface ModelListItem {
    label: string;
    value: string;
    description: string;
    avator: string;
}

interface SkillPanelProps {
    botInfo?: BotInfo | null;
    displayInfo?: DraftBotDisplayInfoData | null;
    botOptionData?: BotOptionData | null;
}

const SkillPanel: React.FC<SkillPanelProps> = ({ botInfo, botOptionData }) => {
    // 从 bot_option_data 中获取模型列表
    const modelList = useMemo<ModelListItem[]>(() => {
        if (!botOptionData?.model_detail_map) {
            return [];
        }
        return Object.entries(botOptionData.model_detail_map).map(([key, model]: [string, ModelDetail]) => ({
            label: model.name || model.model_name || '未知模型',
            value: String(model.model_id || key),
            description: model.model_name || '',
            avator: model.model_icon_url || 'https://huggingface.co/front/assets/huggingface_logo-noborder.svg',
        }));
    }, [botOptionData?.model_detail_map]);

    // 根据 botInfo 中的 ModelInfo 选择默认模型
    const defaultModelValue = botInfo?.ModelInfo?.ModelId ? String(modelList[0]?.value) : (modelList[0]?.value || '');
    const [selectedValue, setSelectedValue] = useState(defaultModelValue);
    const selectedModel = modelList.find((model) => model.value === selectedValue) || modelList[0];

    // 知识调用模式
    const [knowledgeMode, setKnowledgeMode] = useState('on_demand');
    // 用户问题建议开关
    const [suggestionEnabled, setSuggestionEnabled] = useState('enabled');

    // 从 botInfo 中获取数据
    const plugins = botInfo?.PluginInfoList || [];
    const variables = botInfo?.VariableList || [];
    const workflows = botInfo?.WorkflowInfoList || [];
    const databases = botInfo?.DatabaseList || [];

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%', background: '#fff' }}>
            {/* Header */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0 16px', height: '56px', borderBottom: '1px solid #e0e0e0' }}>
                <Title level="h5" style={{ margin: 0, lineHeight: '56px' }}>技能</Title>
                <Space>
                    {modelList.length > 0 ? (
                        <Select
                            value={selectedValue}
                            onChange={(value) => setSelectedValue(value as string)}
                            borderless
                            showArrow
                            size="medium"
                            valueType="value"
                            filterable
                            prefixIcon={selectedModel?.avator ? <img src={selectedModel.avator} style={{ marginRight: '8px', width: '21px', borderRadius: '25%' }} /> : undefined}
                        >
                            {modelList.map((model, idx) =>(
                                <Option style={{ height: '60px'}} key={idx} value={model.value} label={model.label}>
                                    <div style={{ display: 'flex', alignItems: 'center' }}>
                                        <img src={model.avator} style={{ maxWidth: '25px', borderRadius: '25%'}} />
                                        <div style={{marginLeft: '16px'}}>
                                            <div>{model.label}</div>
                                            <div style={{ color: 'var(--td-gray-color-9)', fontSize: '13px' }}>{model.description}</div>
                                        </div>
                                    </div>
                                </Option>
                            ))}
                        </Select>
                    ) : (
                        <span style={{ color: '#888', fontSize: '14px' }}>暂无可用模型</span>
                    )}
                </Space>
            </div>

            {/* Content */}
            <div style={{ flex: 1, overflowY: 'auto', padding: '0 16px', borderLeft: '1px solid #e0e0e0'}}>

                {/* 技能分组 */}
                <SectionTitle title="技能" />
                <Collapse expandIconPlacement="left" borderless>
                    <Collapse.Panel header={<PanelHeader title="插件" />} value="plugin">
                        {plugins.length > 0 ? (
                            plugins.map((plugin, idx) => (
                                <div key={idx} style={{ display: 'flex', alignItems: 'center', marginBottom: idx < plugins.length - 1 ? 12 : 0 }}>
                                    <img src={"https://g1.itc.cn/msfe-pch-prod/300000000000/assets/images/a5df49ba69.png"} alt={plugin.ApiName} style={{ width: 32, height: 32, marginRight: 12, borderRadius: 4 }} />
                                    <div>
                                        <Title level="h6" style={{margin: 0}}>{plugin.ApiName || '未命名插件'}</Title>
                                        <Paragraph style={{ margin: 0, color: '#888', fontSize: 12 }}>{'暂无描述'}</Paragraph>
                                    </div>
                                </div>
                            ))
                        ) : (
                            <Paragraph style={{ margin: 0, color: '#888'}}>暂无插件，点击 + 添加插件</Paragraph>
                        )}
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="工作流" />} value="workflow">
                        {workflows.length > 0 ? (
                            workflows.map((wf, idx) => (
                                <div key={idx} style={{ marginBottom: idx < workflows.length - 1 ? 8 : 0 }}>
                                    <Tag>{wf.WorkflowName || `工作流 ${wf.WorkflowId}`}</Tag>
                                </div>
                            ))
                        ) : (
                            <Paragraph style={{ margin: 0, color: '#888'}}>通过可视化方式编排插件、模型、代码块等，实现复杂业务流程。</Paragraph>
                        )}
                    </Collapse.Panel>
                </Collapse>

                {/* 知识分组 */}
                <SectionTitle
                    title="知识"
                    extra={
                        <Select
                            value={knowledgeMode}
                            onChange={(v) => setKnowledgeMode(v as string)}
                            borderless
                            size="small"
                            style={{ width: 110 }}
                            prefixIcon={<CloudUploadIcon />}
                        >
                            <Option value="on_demand" label="按需调用" />
                            <Option value="always" label="每次调用" />
                        </Select>
                    }
                />
                <Collapse expandIconPlacement="left" borderless>
                    <Collapse.Panel header={<PanelHeader title="文本" />} value="text">
                        <Paragraph style={{ margin: 0, color: '#888'}}>上传文档、URL 等文本知识，智能体可引用回答问题。</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="表格" />} value="table">
                        {databases.length > 0 ? (
                            databases.map((db, idx) => (
                                <div key={idx} style={{ marginBottom: idx < databases.length - 1 ? 8 : 0 }}>
                                    <Tag>{db.TableName || `数据库 ${db.TableId}`}</Tag>
                                </div>
                            ))
                        ) : (
                            <Paragraph style={{ margin: 0, color: '#888'}}>上传表格后，支持按列匹配或自然语言查询。</Paragraph>
                        )}
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="照片" />} value="photo">
                        <Paragraph style={{ margin: 0, color: '#888'}}>上传照片并添加语义描述，智能体可匹配合适的照片。</Paragraph>
                    </Collapse.Panel>
                </Collapse>

                {/* 记忆分组 */}
                <SectionTitle title="记忆" />
                <Collapse expandIconPlacement="left" borderless>
                    <Collapse.Panel header={<PanelHeader title="变量" />} value="variables">
                        <Space>
                            {variables.length > 0 ? (
                                variables.map((v, idx) => (
                                    <Tag key={idx}>{v.keyword || '未命名变量'}</Tag>
                                ))
                            ) : (
                                <>
                                    <Tag>sys_uuid</Tag>
                                    <Tag>master_name</Tag>
                                </>
                            )}
                        </Space>
                    </Collapse.Panel>
                    <Collapse.Panel
                        header={
                            <PanelHeader
                                title="数据库"
                                extra={
                                    <Space>
                                        <HelpCircleIcon style={{ color: '#888', cursor: 'pointer' }} onClick={(e) => e.stopPropagation()} />
                                        <Button shape="square" variant="text" icon={<AddIcon />} onClick={(e) => e.stopPropagation()} />
                                    </Space>
                                }
                            />
                        }
                        value="database"
                    >
                        <Paragraph style={{ margin: 0, color: '#888'}}>定义表结构存储用户数据，实现收藏夹、待办等功能。</Paragraph>
                    </Collapse.Panel>
                </Collapse>

                {/* 对话体验分组 */}
                <SectionTitle title="对话体验" />
                <Collapse expandIconPlacement="left" borderless>
                    <Collapse.Panel header={<PanelHeader title="开场白" />} value="opening">
                        <Paragraph style={{ margin: 0, color: '#888'}}>设置智能体的开场白内容。</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel
                        header={
                            <PanelHeader
                                title="用户问题建议"
                                extra={
                                    <Select
                                        value={suggestionEnabled}
                                        onChange={(v) => setSuggestionEnabled(v as string)}
                                        borderless
                                        size="small"
                                        style={{ width: 80 }}
                                    >
                                        <Option value="enabled" label="开启" />
                                        <Option value="disabled" label="关闭" />
                                    </Select>
                                }
                            />
                        }
                        value="suggestion"
                    >
                        <Paragraph style={{ margin: 0, color: '#888'}}>开启后，智能体会根据对话上下文推荐问题。</Paragraph>
                    </Collapse.Panel>
                </Collapse>

            </div>
        </div>
    )
}

export default SkillPanel;
