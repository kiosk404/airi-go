import React from 'react';
import { Card, Tabs, List, Tag, Button, Space } from 'antd';
import { StarOutlined, DownloadOutlined, EyeOutlined } from '@ant-design/icons';

const { TabPane } = Tabs;

const ExplorePage: React.FC = () => {
    const plugins = [
        {
            id: '1',
            name: '天气查询插件',
            description: '可以查询实时天气信息的插件',
            downloads: 1234,
            rating: 4.8,
            tags: ['天气', 'API', '实用'],
        },
        {
            id: '2',
            name: '翻译助手',
            description: '支持多语言翻译的智能助手',
            downloads: 2567,
            rating: 4.9,
            tags: ['翻译', 'AI', '多语言'],
        },
        {
            id: '3',
            name: '代码生成器',
            description: '根据描述自动生成代码的插件',
            downloads: 1890,
            rating: 4.7,
            tags: ['代码', '生成', '开发'],
        },
    ];

    const templates = [
        {
            id: '1',
            name: '客服机器人模板',
            description: '一个完整的客服机器人解决方案',
            downloads: 3456,
            rating: 4.9,
            tags: ['客服', '机器人', '模板'],
        },
        {
            id: '2',
            name: '教育助手模板',
            description: '专为教育场景设计的 AI 助手',
            downloads: 2134,
            rating: 4.8,
            tags: ['教育', '助手', '学习'],
        },
    ];

    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold text-gray-900">探索</h1>

            <Tabs defaultActiveKey="plugins">
                <TabPane tab="插件" key="plugins">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {plugins.map((plugin) => (
                            <Card
                                key={plugin.id}
                                hoverable
                                actions={[
                                    <Button key="view" type="link" icon={<EyeOutlined />}>
                                        查看
                                    </Button>,
                                    <Button key="download" type="link" icon={<DownloadOutlined />}>
                                        安装
                                    </Button>,
                                ]}
                            >
                                <Card.Meta
                                    title={plugin.name}
                                    description={plugin.description}
                                />
                                <div className="mt-4 space-y-2">
                                    <div className="flex items-center justify-between text-sm text-gray-500">
                                        <span>下载量: {plugin.downloads}</span>
                                        <span className="flex items-center">
                      <StarOutlined className="text-yellow-400 mr-1" />
                                            {plugin.rating}
                    </span>
                                    </div>
                                    <div className="flex flex-wrap gap-1">
                                        {plugin.tags.map((tag) => (
                                            <Tag key={tag}>
                                                {tag}
                                            </Tag>
                                        ))}
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </div>
                </TabPane>
                <TabPane tab="模板" key="templates">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {templates.map((template) => (
                            <Card
                                key={template.id}
                                hoverable
                                actions={[
                                    <Button key="view" type="link" icon={<EyeOutlined />}>
                                        查看
                                    </Button>,
                                    <Button key="download" type="link" icon={<DownloadOutlined />}>
                                        使用
                                    </Button>,
                                ]}
                            >
                                <Card.Meta
                                    title={template.name}
                                    description={template.description}
                                />
                                <div className="mt-4 space-y-2">
                                    <div className="flex items-center justify-between text-sm text-gray-500">
                                        <span>下载量: {template.downloads}</span>
                                        <span className="flex items-center">
                      <StarOutlined className="text-yellow-400 mr-1" />
                                            {template.rating}
                    </span>
                                    </div>
                                    <div className="flex flex-wrap gap-1">
                                        {template.tags.map((tag) => (
                                            <Tag key={tag}>
                                                {tag}
                                            </Tag>
                                        ))}
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </div>
                </TabPane>
            </Tabs>
        </div>
    );
};

export default ExplorePage;
