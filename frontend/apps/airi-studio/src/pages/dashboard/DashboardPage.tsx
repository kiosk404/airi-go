import React from 'react';
import { Card, Row, Col, Statistic } from 'antd';
import { UserOutlined, ProjectOutlined, ApiOutlined } from '@ant-design/icons';

const DashboardPage: React.FC = () => {
    return (
        <div className="space-y-6">
        <h1 className="text-2xl font-bold text-gray-900">仪表板</h1>

            <Row gutter={[16, 16]}>
    <Col xs={24} sm={12} lg={8}>
    <Card>
        <Statistic
            title="总用户数"
    value={1128}
    prefix={<UserOutlined />}
    valueStyle={{ color: '#3f8600' }}
    />
    </Card>
    </Col>
    <Col xs={24} sm={12} lg={8}>
    <Card>
        <Statistic
            title="项目数量"
    value={93}
    prefix={<ProjectOutlined />}
    valueStyle={{ color: '#1890ff' }}
    />
    </Card>
    </Col>
    <Col xs={24} sm={12} lg={8}>
    <Card>
        <Statistic
            title="API 调用"
    value={112893}
    prefix={<ApiOutlined />}
    valueStyle={{ color: '#cf1322' }}
    />
    </Card>
    </Col>
    </Row>

    <Row gutter={[16, 16]}>
    <Col xs={24} lg={12}>
    <Card title="最近活动" className="h-64">
    <div className="space-y-4">
    <div className="flex items-center space-x-3">
    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
        <span className="text-sm text-gray-600">用户登录</span>
        <span className="text-xs text-gray-400 ml-auto">2分钟前</span>
    </div>
    <div className="flex items-center space-x-3">
    <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
        <span className="text-sm text-gray-600">创建新项目</span>
        <span className="text-xs text-gray-400 ml-auto">5分钟前</span>
    </div>
    </div>
    </Card>
    </Col>
    <Col xs={24} lg={12}>
    <Card title="系统状态" className="h-64">
    <div className="space-y-4">
    <div className="flex justify-between items-center">
    <span className="text-sm text-gray-600">CPU 使用率</span>
    <span className="text-sm font-medium">45%</span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2">
    <div className="bg-green-500 h-2 rounded-full" style={{ width: '45%' }}></div>
    </div>
    <div className="flex justify-between items-center">
    <span className="text-sm text-gray-600">内存使用率</span>
        <span className="text-sm font-medium">67%</span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2">
    <div className="bg-yellow-500 h-2 rounded-full" style={{ width: '67%' }}></div>
    </div>
    </div>
    </Card>
    </Col>
    </Row>
    </div>
);
};

export default DashboardPage;
