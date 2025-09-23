import { Layout, Menu, Button, Avatar } from 'tdesign-react';
import { Outlet, useNavigate} from 'react-router-dom';
import { HelpCircleIcon, HomeIcon, LogoGithubIcon, NotificationIcon } from 'tdesign-icons-react';

const {Header, Content} = Layout;
const {HeadMenu, MenuItem} = Menu;


export default () => {
    const navigate = useNavigate();
    return (
        <Layout>
            <Header>
                <HeadMenu 
                    defaultValue={"main"}
                    logo={<LogoGithubIcon size="36px" />}
                    operations={
                        <div className="t-menu__operations">
                            <Button variant='text' shape='square' icon={<NotificationIcon />}></Button>
                            <Button variant='text' shape='square' icon={<HelpCircleIcon />}></Button>
                            <Avatar size="medium" hideOnLoadFailed={false} image="https://kiosk007.top/images/avatar.jpg">Kiosk</Avatar>
                        </div>
                    }
                >

                    <MenuItem value="main" icon={<HomeIcon/>} onClick={() => navigate('/')} content="工作空间"/>
                    <MenuItem value="playground" icon={<HomeIcon/>} onClick={() => navigate('/playground')} content="Playground"/>
                
                </HeadMenu>
            </Header>
            <Content>
                <Outlet />
            </Content>
        </Layout>
    );
};