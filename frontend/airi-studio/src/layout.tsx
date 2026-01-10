import { Layout, Nav, Avatar, Dropdown, Toast, Button } from '@douyinfe/semi-ui';
import { IconBell, IconHelpCircle, IconExit, IconUser, IconHome, IconFile, IconSetting, IconWrench } from '@douyinfe/semi-icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';

const { Header, Content } = Layout;

export default () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { user, logout } = useAuth();

    const handleLogout = async () => {
        try {
            await logout();
            Toast.success('已退出登录');
            navigate('/login', { replace: true });
        } catch {
            Toast.error('退出失败');
        }
    };

    // 根据当前路径获取选中的菜单项
    const getSelectedKey = () => {
        const path = location.pathname;
        if (path.startsWith('/workspace')) return 'workspace';
        if (path.startsWith('/knowledge')) return 'knowledge';
        if (path.startsWith('/models')) return 'models';
        if (path.startsWith('/tool')) return 'tool';
        return '';
    };

    return (
        <Layout style={{ height: "100vh", overflow: "hidden" }}>
            <Header style={{ backgroundColor: 'var(--semi-color-bg-1)' }}>
                <Nav
                    mode="horizontal"
                    selectedKeys={[getSelectedKey()]}
                    onSelect={(data) => {
                        const key = data.itemKey as string;
                        navigate(`/${key}`);
                    }}
                    header={{
                        logo: (
                            <img src="/logo.png"
                            alt="Logo"
                            className="w-8 h-8 rounded-lg object-cover"
                            style={{ cursor: 'pointer' }}
                            onClick={() => navigate('/')}
                            />
                        ),
                        text: <span style={{ cursor: 'pointer'}} onClick={() => navigate('/')}>Airi Studio</span>,
                    }}
                    footer={
                        <div className="flex items-center gap-2">
                            <Button
                                theme="borderless"
                                icon={<IconBell />}
                                style={{ color: 'var(--semi-color-text-2)' }}
                            />
                            <Button
                                theme="borderless"
                                icon={<IconHelpCircle />}
                                style={{ color: 'var(--semi-color-text-2)' }}
                            />
                            <Dropdown
                                trigger="click"
                                position="bottomRight"
                                render={
                                    <Dropdown.Menu>
                                        <Dropdown.Item icon={<IconUser />}>
                                            个人信息
                                        </Dropdown.Item>
                                        <Dropdown.Divider />
                                        <Dropdown.Item icon={<IconExit />} onClick={handleLogout}>
                                            退出登录
                                        </Dropdown.Item>
                                    </Dropdown.Menu>
                                }
                            >
                                <Avatar
                                    size="small"
                                    src={user?.avatar_url}
                                    style={{ cursor: 'pointer' }}
                                >
                                    {user?.nick_name?.[0] || user?.name?.[0] || 'U'}
                                </Avatar>
                            </Dropdown>
                        </div>
                    }
                >
                    <Nav.Item itemKey="workspace" text="工作室" icon={<IconHome />} />
                    <Nav.Item itemKey="knowledge" text="知识库" icon={<IconFile />} />
                    <Nav.Item itemKey="models" text="模型" icon={<IconSetting />} />
                    <Nav.Item itemKey="tool" text="工具" icon={<IconWrench />} />
                </Nav>
            </Header>
            <Content style={{ backgroundColor: 'var(--semi-color-bg-0)', flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
                <Outlet />
            </Content>
        </Layout>
    );
};