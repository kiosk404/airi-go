
import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Card, Form, Button, Toast, Tabs, TabPane, Typography } from '@douyinfe/semi-ui';
import { IconUser, IconLock } from '@douyinfe/semi-icons';
import { useAuth } from '@/contexts/AuthContext';
import { ApiError } from '@/services/core';

const { Title, Text } = Typography;

export default function LoginPage() {
    const navigate = useNavigate();
    const location = useLocation();
    const { login, register } = useAuth();

    const [activeTab, setActiveTab] = useState<string>('login');
    const [loading, setLoading] = useState(false);

    // 获取登录后要跳转的页面
    const from = (location.state as { from?: string })?.from || '/apps';

    const handleLogin = async (values: { account: string; password: string }) => {
        setLoading(true);
        try {
            await login(values.account, values.password);
            Toast.success('登录成功');
            navigate(from, { replace: true });
        } catch (error) {
            if (error instanceof ApiError) {
                Toast.error(error.message || '登录失败');
            } else {
                Toast.error('登录失败，请稍后重试');
            }
        } finally {
            setLoading(false);
        }
    };

    const handleRegister = async (values: { account: string; password: string; confirmPassword: string }) => {
        if (values.password !== values.confirmPassword) {
            Toast.error('两次输入的密码不一致');
            return;
        }

        setLoading(true);
        try {
            await register(values.account, values.password);
            Toast.success('注册成功');
            navigate(from, { replace: true });
        } catch (error) {
            if (error instanceof ApiError) {
                Toast.error(error.message || '注册失败');
            } else {
                Toast.error('注册失败，请稍后重试');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
            <Card
                className="w-full max-w-md"
                shadows="always"
                style={{ padding: '32px' }}
            >
                <div className="text-center mb-8">
                    <div className="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center">
                        <svg className="w-10 h-10 text-white" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
                        </svg>
                    </div>
                    <Title heading={3} style={{ marginBottom: 4 }}>Airi Studio</Title>
                    <Text type="tertiary">AI 应用开发平台</Text>
                </div>

                <Tabs activeKey={activeTab} onChange={(key) => setActiveTab(key)}>
                    <TabPane tab="登录" itemKey="login">
                        <Form
                            onSubmit={handleLogin}
                            style={{ marginTop: 16 }}
                        >
                            <Form.Input
                                field="account"
                                label="账号"
                                placeholder="请输入账号或邮箱"
                                prefix={<IconUser />}
                                size="large"
                                showClear
                                rules={[
                                    { required: true, message: '请输入账号' },
                                    { min: 2, message: '账号长度至少2位' },
                                ]}
                            />
                            <Form.Input
                                field="password"
                                label="密码"
                                mode="password"
                                placeholder="请输入密码"
                                prefix={<IconLock />}
                                size="large"
                                rules={[
                                    { required: true, message: '请输入密码' },
                                    { min: 6, message: '密码长度至少6位' },
                                ]}
                            />
                            <Button
                                htmlType="submit"
                                theme="solid"
                                block
                                size="large"
                                loading={loading}
                                style={{ marginTop: 8 }}
                            >
                                登录
                            </Button>
                        </Form>
                    </TabPane>

                    <TabPane tab="注册" itemKey="register">
                        <Form
                            onSubmit={handleRegister}
                            style={{ marginTop: 16 }}
                        >
                            <Form.Input
                                field="account"
                                label="账号"
                                placeholder="请输入账号或邮箱"
                                prefix={<IconUser />}
                                size="large"
                                showClear
                                rules={[
                                    { required: true, message: '请输入账号' },
                                    { min: 2, max: 20, message: '账号长度为2-20位' },
                                    {
                                        pattern: /^[a-zA-Z0-9_]+$|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
                                        message: '账号只能包含字母、数字、下划线，或使用邮箱格式'
                                    },
                                ]}
                            />
                            <Form.Input
                                field="password"
                                label="密码"
                                mode="password"
                                placeholder="请输入密码"
                                prefix={<IconLock />}
                                size="large"
                                rules={[
                                    { required: true, message: '请输入密码' },
                                    { min: 6, message: '密码长度至少6位' },
                                ]}
                            />
                            <Form.Input
                                field="confirm_password"
                                label="再次确认密码"
                                mode="password"
                                placeholder="请再次输入密码"
                                prefix={<IconLock />}
                                size="large"
                                rules={[
                                    { required: true, message: '请确认密码' },
                                ]}
                            />
                            <Button
                                htmlType="submit"
                                theme="solid"
                                block
                                size="large"
                                loading={loading}
                                style={{ marginTop: 8 }}
                            >
                                注册
                            </Button>
                        </Form>
                    </TabPane>
                </Tabs>

                <div className="text-center text-gray-400 text-xs mt-8">
                    © 2024 Airi Studio. All rights reserved.
                </div>
            </Card>
        </div>
    );
}