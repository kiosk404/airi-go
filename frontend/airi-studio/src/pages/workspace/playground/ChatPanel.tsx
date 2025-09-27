import { Layout, Typography, Input, Button, Avatar, Space } from "tdesign-react";
import { BrowseIcon, Delete1Icon, RollbackIcon, SendIcon } from "tdesign-icons-react";
const { Title } = Typography
const { Header, Content, Footer } = Layout;

const kioskAvatarURL = "https://kiosk007.top/images/avatar.jpg"
const airiAvatarURL = "https://raw.githubusercontent.com/moeru-ai/airi/refs/heads/main/docs/content/public/maskable_icon_x192.avif"

const mockMessages = [
    {
        id: 1,
        author: "Airi",
        content: 'Hello, Airi here ~',
        avatar: airiAvatarURL,
    },
    {
        id: 2,
        author: "kiosk",
        content: '你好，我是kiosk ~',
        avatar: kioskAvatarURL,
    },    
    {
        id: 3,
        author: "Airi",
        content: 'Hello, Airi here ~',
        avatar: airiAvatarURL,
    },   
    {
        id: 4,
        author: "kiosk",
        content: '你好，我是kiosk ~',
        avatar: kioskAvatarURL,
    },   
    {
        id: 5,
        author: "Airi",
        content: 'Hello, Airi here ~',
        avatar: airiAvatarURL,
    },   
    {
        id: 6,
        author: "kiosk",
        content: '你好，我是kiosk ~',
        avatar: kioskAvatarURL,
    },    
]




const ChatPanel: React.FC = () => {
    return (
        <Layout style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Header style={{ backgroundColor: 'white', padding: '0 16px', borderBottom: '1px rgb(255, 255, 255)'}}>
                <Title level="h5" style={{ margin: 0, lineHeight: '64px'}}> 预览与调试 </Title>
            </Header>
            <Content style={{ flex: 1, padding: '0 16px', overflowY: 'scroll', backgroundColor: 'white' }}>
                {mockMessages.map((msg)=>(
                    <div
                        key={msg.id}
                        style={{
                            display: 'flex',
                            justifyContent: msg.author === 'Airi' ? 'flex-start' : 'flex-end',
                            marginBottom: '16px',
                        }}
                    >
                        <div style={{ display: 'flex', alignItems: msg.author === 'Airi' ? 'row' : 'row-reverse' }}>
                            <Avatar image={msg.avatar} style={{margin: '0 8px'}} />
                            <div>
                                <div style={{ color: '#888', fontSize: '12px', marginBottom: '4px' }}>{msg.author}</div>
                                <div style={{
                                    background: msg.author === 'Airi' ? '#f3f3f3' : '#e0e0fff',
                                    padding: '8px 12px',
                                    borderRadius: '8px',
                                    maxWidth: '300px',
                                }}>
                                    {msg.content}
                                </div>
                                {msg.author !== 'Airi' && (
                                    <Space style={{ marginTop: '4px'}}>
                                        <Button shape="circle" size="small" variant="text" icon={<RollbackIcon/>}/>
                                        <Button shape="circle" size="small" variant="text" icon={<Delete1Icon/>}/>
                                    </Space>
                                )}
                            </div>
                        </div>
                    </div>
                ))}
            </Content>
            <Footer style={{ backgroundColor: 'white', padding: '8px 16px', borderTop: '1px solid #e0e0e0'}}>
                <Input 
                    placeholder="请继续对话..."
                    size="medium"
                    label={<Button shape="square" variant="text" icon={<BrowseIcon />} />}
                    suffix={<Button variant="text" icon={<SendIcon />} />}
                />
            </Footer>
        </Layout>
    )
}

export default ChatPanel;