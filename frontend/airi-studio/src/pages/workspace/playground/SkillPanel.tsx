import { Button, Typography, Space, Select, Tag, Collapse } from "tdesign-react"
import { AddIcon } from "tdesign-icons-react"
import React, { useState, useMemo } from "react"
import type { BotInfo, DraftBotDisplayInfoData, BotOptionData, ModelDetail } from '@/services/draftbot';

const { Title, Paragraph } = Typography;
const { Option } = Select;

const PanelHeader = ({ title }: { title: string }) => (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%'}}>
        <span style={{ lineHeight: '22px' }}>{title}</span>
        <Button shape="square" variant="text" icon={<AddIcon />}/>
    </div>
)

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
            value: model.model_id || key,
            description: model.model_name || '',
            avator: model.model_icon_url || 'https://ollama.com/public/ollama.png',
        }));
    }, [botOptionData?.model_detail_map]);

    // 根据 botInfo 中的 ModelInfo 选择默认模型
    const defaultModelValue = botInfo?.ModelInfo?.model_id || modelList[0]?.value || '';
    const [selectedValue, setSelectedValue] = useState(defaultModelValue);
    const selectedModel = modelList.find((model) => model.value === selectedValue) || modelList[0];

    // 从 botInfo 中获取插件列表
    const plugins = botInfo?.PluginInfoList || [];
    // 从 botInfo 中获取变量列表
    const variables = botInfo?.VariableList || [];
    // 从 botInfo 中获取工作流列表
    const workflows = botInfo?.WorkflowInfoList || [];
    // 从 botInfo 中获取数据库列表
    const databases = botInfo?.DatabaseList || [];

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%', background: '#fff' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0 16px', height: '50px', borderBottom: '1px solid #e0e0e0' }}>
                <Title level="h5" style={{ margin: 0, lineHeight: '50px'}}>技能</Title>
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
                                            <div style={{ color: 'var(--td-gray-color-9', fontSize: '13px' }}>{model.description}</div>
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

            <div style={{ flex: 1, overflowY: 'auto', padding: '16px', borderLeft: '1px solid #e0e0e0'}}>
                <Collapse defaultValue={['plugin']}>
                    <Collapse.Panel header={<PanelHeader title={`插件`} />} value="plugin">
                        {plugins.length > 0 ? (
                            plugins.map((plugin, idx) => (
                                <div key={idx} style={{ display: 'flex', alignItems: 'center', marginBottom: idx < plugins.length - 1 ? 12 : 0 }}>
                                    <img src={plugin.icon || "https://g1.itc.cn/msfe-pch-prod/300000000000/assets/images/a5df49ba69.png"} alt={plugin.name} style={{ width: 32, height: 32, marginRight: 12 }} />
                                    <div>
                                        <Title level="h6" style={{margin: 0}}>{plugin.name || '未命名插件'}</Title>
                                        <Paragraph style={{ margin: 0, color: '#888'}}>{plugin.description || '暂无描述'}</Paragraph>
                                    </div>
                                </div>
                            ))
                        ) : (
                            <Paragraph style={{ margin: 0, color: '#888'}}>暂无插件，点击 + 添加插件</Paragraph>
                        )}
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title={`工作流`} />} value="workflow">
                        {workflows.length > 0 ? (
                            workflows.map((wf, idx) => (
                                <div key={idx} style={{ marginBottom: idx < workflows.length - 1 ? 8 : 0 }}>
                                    <Tag>{wf.name || `工作流 ${wf.workflow_id}`}</Tag>
                                </div>
                            ))
                        ) : (
                            <Paragraph style={{ margin: 0, color: '#888'}}>工作流支持通过可视化的方式，对插件、大语言模型、代码块等功能进行组合，从而实现复杂、稳定的业务流程编排，例如旅行规划、报告分析等。</Paragraph>
                        )}
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="知识" />} value="knowledge">
                        <Paragraph style={{ margin: 0, color: '#888'}}>将文档、URL、三方数据源上传为文本知识库后，用户发送消息时，智能体能够引用文本知识中的内容回答用户问题。</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title={`表格`} />} value="table">
                        {databases.length > 0 ? (
                            databases.map((db, idx) => (
                                <div key={idx} style={{ marginBottom: idx < databases.length - 1 ? 8 : 0 }}>
                                    <Tag>{db.name || `数据库 ${db.database_id}`}</Tag>
                                </div>
                            ))
                        ) : (
                            <Paragraph style={{ margin: 0, color: '#888'}}>用户上传表格后，支持按照表格的某列来匹配合适的行给智能体引用，同时也支持基于自然语言对数据库进行查询和计算。</Paragraph>
                        )}
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="照片" />} value="photot">
                        <Paragraph style={{ margin: 0, color: '#888'}}>照片上传到知识库后自动/手动添加语义描述，智能体可以基于照片的描述匹配到最合适的照片</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="记忆" />} value="memory">
                        <Paragraph style={{ margin: 0, color: '#888'}}>表格存储支持智能体开发者定义表结构，并将用户数据存储在表中。实现收藏夹、todo list、书籍管理、财务管理等功能。了解更多</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title={`变量`} />} value="variables">
                        <Paragraph style={{ margin: 0, color: '#888'}}>
                            <Space>
                                {variables.length > 0 ? (
                                    variables.map((v, idx) => (
                                        <Tag key={idx}>{v.name || '未命名变量'}</Tag>
                                    ))
                                ) : (
                                    <>
                                        <Tag>sys_uuid</Tag>
                                        <Tag>master_name</Tag>
                                    </>
                                )}
                            </Space>
                        </Paragraph>
                    </Collapse.Panel>

                </Collapse>
            </div>

        </div>
    )
}

export default SkillPanel;