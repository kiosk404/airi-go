import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Layout as AntLayout, Button } from 'antd';
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons';
import { useAppStore } from './stores/appStore';

const { Header, Sider, Content } = AntLayout;

const Layout = () => {
    const { sidebarCollapsed, toggleSidebar } = useAppStore();
    const navigate = useNavigate();
    const location = useLocation();

    const menuItems = [
        { key: 'dashboard', label: '仪表板', path: '/space/demo/dashboard' },
        { key: 'develop', label: '开发', path: '/space/demo/develop' },
        { key: 'library', label: '资源库', path: '/space/demo/library' },
        { key: 'explore', label: '探索', path: '/explore' },
    ];

    const handleMenuClick = (path: string) => {
        navigate(path);
    };

    const getCurrentKey = () => {
        const currentPath = location.pathname;
        const currentItem = menuItems.find(item => currentPath === item.path);
        return currentItem?.key || 'develop';
    };

    return (
        <AntLayout className="h-screen">
            <Sider
                width={240}
                collapsed={sidebarCollapsed}
                className="bg-gray-50 border-r border-gray-200"
            >
                <div className="p-4">
                    <h1 className="text-xl font-bold text-gray-800">
                        {sidebarCollapsed ? 'AS' : 'Airi Studio'}
                    </h1>
                </div>
                <div className="px-4 space-y-2">
                    <div className="text-sm text-gray-500">导航菜单</div>
                    <div className="space-y-1">
                        {menuItems.map((item) => (
                            <div
                                key={item.key}
                                onClick={() => handleMenuClick(item.path)}
                                className={`px-3 py-2 text-sm rounded cursor-pointer transition-colors ${
                                    getCurrentKey() === item.key
                                        ? 'bg-blue-100 text-blue-700 border-r-2 border-blue-500'
                                        : 'text-gray-700 hover:bg-gray-100'
                                }`}
                            >
                                {item.label}
                            </div>
                        ))}
                    </div>
                </div>
            </Sider>
            <AntLayout>
                <Header className="bg-white shadow-sm border-b border-gray-200 px-6">
                    <div className="flex items-center justify-between h-full">
                        <div className="flex items-center space-x-4">
                            <Button
                                type="text"
                                icon={sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                                onClick={toggleSidebar}
                                className="text-gray-600 hover:text-gray-800"
                            />
                            <h2 className="text-lg font-semibold text-gray-800">欢迎使用 Airi Studio</h2>
                        </div>
                        <div className="flex items-center space-x-4">
                            <span className="text-sm text-gray-600">用户菜单</span>
                        </div>
                    </div>
                </Header>
                <Content className="p-6 bg-gray-50">
                    <Outlet />
                </Content>
            </AntLayout>
        </AntLayout>
    );
};

export default Layout;
