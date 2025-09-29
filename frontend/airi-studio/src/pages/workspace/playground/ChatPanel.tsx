import { Layout, Typography } from "tdesign-react";
import React, {useCallback, useState} from 'react';
import { Chat, Avatar } from '@douyinfe/semi-ui';
const { Title } = Typography
const { Header, Content } = Layout;


const roleConfig = {
    user: {
        name: 'kiosk',
        avatar: "https://kiosk007.top/images/avatar.jpg",
    },
    assistant: {
        name: 'Airi',
        avatar: "https://raw.githubusercontent.com/moeru-ai/airi/refs/heads/main/docs/content/public/maskable_icon_x192.avif",
    }
}

const mockMessages = [
    {
        id: '1',
        role: 'assistant',
        content: 'Hello, Airi here ~',
        createAt: Date.now() - 3600000,        
    },
    {
        id: '2',
        role: 'user',
        content: '你好，我是kiosk ~',
        createAt: Date.now() - 3600000,    
    },    
    {
        id: '3',
        role: 'assistant',
        content: 'Hello, Airi here ~',
        createAt: Date.now() - 3600000,   
    },   
    {
        id: '4',
        role: 'user',
        content: '你好，我是kiosk ~',
        createAt: Date.now() - 3600000,   
    },   
    {
        id: '5',
        role: 'assistant',
        content: 'Hello, Airi here ~',
        createAt: Date.now() - 3600000,   
    },   
    {
        id: '6',
        role: 'user',
        content: '你好，我是kiosk ~',
        createAt: Date.now() - 3600000,   
    },    
];

const hints = [
    '你能帮我分析一下这个问题吗？',
    '如何使用这个功能?',
]



const ChatPanel: React.FC = () => {
    // 
    const [message, setMessage] = useState(mockMessages);

    // 消息变化
    const handleChatsChange = useCallback((chats: any) =>{
        setMessage(chats)
    }, [])


    // 处理消息发送
    const handleMessageSend = useCallback((props: any) => {
        const newAssistantMessage = {
            role: 'assistant',
            id: Date.now().toString(),
            createAt: Date.now(),
            content: "this is a mock assistant message",
        }
        setTimeout(() => {
            setMessage((message) => ([...message, newAssistantMessage]));
        }, 200);
    }, []);

    // 处理消息重发
    const handleMessageReset = useCallback(() => {
        setMessage((message) => {
            const lastMessage = message[message.length - 1];
            const newLastMessage = {
                ...lastMessage,
                status: 'complete',
                content: 'This is a mock reset message.',
            }
            return [...message.slice(0, -1), newLastMessage]
        });
    }, [])

    // 处理点击提示
    const handleHintClick = (hint: string) => {
        handleMessageSend({ content: hint, attachment: [] });
    }

    // 处理清除上下文
    const handleClearContext = () => {
        setMessage([]);
    };

    return (
        <Layout style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Header style={{ backgroundColor: 'white', padding: '0 16px', borderBottom: '1px rgb(255, 255, 255)'}}>
                <Title level="h5" style={{ margin: 0, lineHeight: '64px'}}> 预览与调试 </Title>
            </Header>
            <Content style={{ flex: 1, padding: '0 16px', overflowY: 'scroll', backgroundColor: 'white' }}>
                <Chat
                    chats={message}
                    onChatsChange={handleChatsChange}
                    onMessageSend={handleMessageSend}
                    onMessageReset={handleMessageReset}
                    onClear={handleClearContext}
                    onHintClick={handleHintClick}
                    roleConfig={roleConfig}
                    hints={hints}
                    timeFormat="YYYY-MM-DD HH:mm:ss"
                    placeholder="请继续对话..."
                    sendHotKey="enter"
                    align="leftRight"
                    mode="bubble"
                    style={{ height: '100%' }}
                />
            </Content>
        
        </Layout>
    )
}

export default ChatPanel;