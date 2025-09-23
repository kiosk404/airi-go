import { Button, Typography, Space, Select, Tag, Collapse } from "tdesign-react"
import { AddIcon } from "tdesign-icons-react"
import React, { useState } from "react"

const { Title, Paragraph } = Typography;
const { Option } = Select;

const PanelHeader = ({ title }: { title: string }) => (
  <div style={{ display: 'flex', justifyContent: 'space-between', width: '100%'}}>
    <span>{title}</span>
    <Button shape="square" variant="text" icon={<AddIcon />}/>
  </div>
)

const modelList = [
    { label: "ollama", value: "1", description: "deepseek-r1:1.5b", avator: "https://ollama.com/public/ollama.png"},
    { label: "qwen", value: "2", description: "qwen3-coder", avator: "https://assets.alicdn.com/g/qwenweb/qwen-webui-fe/0.0.191/static/favicon.png"},
]

const SkillPanel: React.FC = () => {
    const [selectedValue, setSelectedValue] = useState(modelList[0].value); 
    const selectedModel = modelList.find((model) => model.value === selectedValue) || modelList[0];

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%', background: '#fff' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '8px 16px', borderBottom: '1px solid #e0e0e0' }}>
                <Title level="h5" style={{ margin: 0 }}>技能</Title>
                <Space>
                    <Select
                        value={selectedModel.description}
                        onChange={(value) => setSelectedValue(value as string)}
                        borderless
                        showArrow
                        size="medium"
                        valueType="value"
                        filterable
                        prefixIcon={<img src={selectedModel.avator} style={{ marginRight: '8px', width: '21px', borderRadius: '25%' }} />}
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
                </Space>
            </div>

            <div style={{ flex: 1, overflowY: 'auto', padding: '16px', borderLeft: '1px solid #e0e0e0'}}>
                <Collapse defaultValue={['plugin']}>
                    <Collapse.Panel header={<PanelHeader title="插件" />} value="plugin">
                        <div style={{ display: 'flex', alignItems: 'center' }}>
                            <img src="https://g1.itc.cn/msfe-pch-prod/300000000000/assets/images/a5df49ba69.png" alt="Sohu News" style={{ width: 32, height: 32, marginRight: 12 }} />
                            <div>
                                <Title level="h6" style={{margin: 0}}>搜狐新闻 / top_news</Title>
                                <Paragraph style={{ margin: 0, color: '#888'}}>帮助用户获取搜狐网上的每日热闻</Paragraph>
                            </div>
                        </div>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="工作流" />} value="workflow">
                         <Paragraph style={{ margin: 0, color: '#888'}}>工作流支持通过可视化的方式，对插件、大语言模型、代码块等功能进行组合，从而实现复杂、稳定的业务流程编排，例如旅行规划、报告分析等。</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="知识" />} value="knowledge">
                         <Paragraph style={{ margin: 0, color: '#888'}}>将文档、URL、三方数据源上传为文本知识库后，用户发送消息时，智能体能够引用文本知识中的内容回答用户问题。</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="表格" />} value="table">
                         <Paragraph style={{ margin: 0, color: '#888'}}>用户上传表格后，支持按照表格的某列来匹配合适的行给智能体引用，同时也支持基于自然语言对数据库进行查询和计算。</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="照片" />} value="photot">
                         <Paragraph style={{ margin: 0, color: '#888'}}>照片上传到知识库后自动/手动添加语义描述，智能体可以基于照片的描述匹配到最合适的照片</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="记忆" />} value="memory">
                         <Paragraph style={{ margin: 0, color: '#888'}}>表格存储支持智能体开发者定义表结构，并将用户数据存储在表中。实现收藏夹、todo list、书籍管理、财务管理等功能。了解更多</Paragraph>
                    </Collapse.Panel>
                    <Collapse.Panel header={<PanelHeader title="变量" />} value="variables">
                         <Paragraph style={{ margin: 0, color: '#888'}}>
                            <Space>
                                <Tag>sys_uuid</Tag>
                                <Tag>master_name</Tag>
                            </Space>
                         </Paragraph>
                    </Collapse.Panel>
                    
                </Collapse>
            </div>
            
        </div>
    )
}

export default SkillPanel;