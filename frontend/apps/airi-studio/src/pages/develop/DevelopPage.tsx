import React from 'react';
import { Card, Button, Space, Row, Col } from 'antd';
import { PlayCircleOutlined, PauseCircleOutlined, SettingOutlined } from '@ant-design/icons';

const DevelopPage: React.FC = () => {
    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-2xl font-bold text-gray-900">开发环境</h1>
                <Space>
                    <Button
                        type="primary"
                        icon={<PlayCircleOutlined />}
                        style={{
                            backgroundColor: '#10b981',
                            borderColor: '#10b981',
                            color: 'white'
                        }}
                        onMouseEnter={(e) => {
                            e.currentTarget.style.backgroundColor = '#059669';
                            e.currentTarget.style.borderColor = '#059669';
                        }}
                        onMouseLeave={(e) => {
                            e.currentTarget.style.backgroundColor = '#10b981';
                            e.currentTarget.style.borderColor = '#10b981';
                        }}
                    >
                        运行
                    </Button>
                    <Button
                        icon={<PauseCircleOutlined />}
                        style={{
                            color: '#ea580c',
                            borderColor: '#ea580c'
                        }}
                        onMouseEnter={(e) => {
                            e.currentTarget.style.backgroundColor = '#fed7aa';
                        }}
                        onMouseLeave={(e) => {
                            e.currentTarget.style.backgroundColor = 'transparent';
                        }}
                    >
                        暂停
                    </Button>
                    <Button
                        icon={<SettingOutlined />}
                        style={{
                            color: '#6b7280',
                            borderColor: '#6b7280'
                        }}
                        onMouseEnter={(e) => {
                            e.currentTarget.style.backgroundColor = '#f3f4f6';
                        }}
                        onMouseLeave={(e) => {
                            e.currentTarget.style.backgroundColor = 'transparent';
                        }}
                    >
                        设置
                    </Button>
                </Space>
            </div>

            <Row gutter={[16, 16]}>
                <Col xs={24} lg={16}>
                    <Card title="代码编辑器" className="h-96">
                        <div className="bg-gray-900 text-green-400 p-4 rounded font-mono text-sm">
                            <div>// 在这里编写你的 AI Agent 代码</div>
                            <div>function createAgent() {'{'}</div>
                            <div>  return {'{'}</div>
                            <div>    name: "My AI Agent",</div>
                            <div>    description: "一个智能助手"</div>
                            <div>  {'}'}</div>
                            <div>{'}'}</div>
                        </div>
                    </Card>
                </Col>
                <Col xs={24} lg={8}>
                    <Card title="控制台" className="h-96">
                        <div className="space-y-2">
                            <div className="text-sm text-green-600">✓ 系统启动成功</div>
                            <div className="text-sm text-blue-600">ℹ 正在加载模型...</div>
                            <div className="text-sm text-yellow-600">⚠ 检测到配置变更</div>
                        </div>
                    </Card>
                </Col>
            </Row>

            <Card title="项目信息">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                        <div className="text-sm text-gray-500">项目名称</div>
                        <div className="font-medium">My AI Agent</div>
                    </div>
                    <div>
                        <div className="text-sm text-gray-500">状态</div>
                        <div className="font-medium text-green-600">运行中</div>
                    </div>
                    <div>
                        <div className="text-sm text-gray-500">最后更新</div>
                        <div className="font-medium">2分钟前</div>
                    </div>
                </div>
            </Card>
        </div>
    );
};

export default DevelopPage;
