import { Layout, Menu, Button, Avatar } from 'tdesign-react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { HelpCircleIcon, EarthIcon,
     BookIcon, RocketIcon, NotificationIcon, ToolsIcon } from 'tdesign-icons-react';

const {Header, Content} = Layout;
const {HeadMenu, MenuItem} = Menu;


export default () => {
    const navigate = useNavigate();
    const location = useLocation();
    return (
        <Layout>
            <Header>
                <HeadMenu
                    defaultValue={location.pathname}
                    logo={<RocketIcon size="36px" />}
                    operations={
                        <div className="t-menu__operations">
                            <Button variant='text' shape='square' icon={<NotificationIcon />}></Button>
                            <Button variant='text' shape='square' icon={<HelpCircleIcon />}></Button>
                            <Avatar size="medium" hideOnLoadFailed={false} image="https://kiosk007.top/images/avatar.jpg">Kiosk</Avatar>
                        </div>
                    }
                >
                    <MenuItem 
                        value="/workspace" icon={<EarthIcon/>} 
                        onClick={() => navigate('/workspace')} content={
                           <span className="font-mono">工作室</span>}/>
                    <MenuItem 
                        value="/knowledge" icon={<BookIcon/>} 
                        onClick={() => navigate('/knowledge')} content={
                            <span className="font-mono">知识库</span>}/> 
                    <MenuItem 
                        value="/tool" icon={<ToolsIcon/>} 
                        onClick={() => navigate('/tool')} content={
                            <span className="font-mono">工具</span>}/>           
                </HeadMenu>
            </Header>
            <Content style={{ backgroundColor: 'white'}}>
                <Outlet />
            </Content>
        </Layout>
    );
};