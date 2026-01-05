import { Layout, Typography } from '@douyinfe/semi-ui';
import React, { useCallback, useState, useRef } from 'react';
import { Chat } from '@douyinfe/semi-ui';
import { useParams } from 'react-router-dom';
import * as conversationApi from '@/services/conversation';
import { Scene } from '@/services/conversation';

const { Title } = Typography;
const { Header, Content } = Layout;

// 消息类型定义
interface Message {
    id: string;
    role: 'user' | 'assistant' | 'system';
    content: string;
    createAt: number;
    status?: 'loading' | 'complete' | 'error';
}

const roleConfig = {
    user: {
        name: 'User',
        avatar: 'https://lf3-static.bytednsdoc.com/obj/eden-cn/ptlz_zlp/ljhwZthlaukjlkulzlp/docs-icon.png',
    },
    assistant: {
        name: 'Airi',
        avatar: 'https://raw.githubusercontent.com/moeru-ai/airi/refs/heads/main/docs/content/public/maskable_icon_x192.avif',
    },
    system: {
        name: 'System',
        avatar: 'https://lf3-static.bytednsdoc.com/obj/eden-cn/ptlz_zlp/ljhwZthlaukjlkulzlp/other/logo.png',
    },
};

const defaultHints = [
    '你能帮我分析一下这个问题吗？',
    '如何使用这个功能?',
    '请介绍一下你自己',
];

interface ChatPanelProps {
    botId?: string;
    conversationId?: string;
}

const ChatPanel: React.FC<ChatPanelProps> = ({ botId, conversationId: propConversationId }) => {
    const params = useParams<{ id: string }>();

    // 优先使用 props，否则从 URL 获取
    const currentBotId = botId ?? (params.id || '');
    const currentConversationId = useRef<string>(propConversationId ?? '');

    const [messages, setMessages] = useState<Message[]>([]);
    const [isStreaming, setIsStreaming] = useState(false);
    const abortControllerRef = useRef<AbortController | null>(null);

    // 消息变化
    const handleChatsChange = useCallback((chats?: Message[]) => {
        if (chats) {
            setMessages(chats);
        }
    }, []);

    // 处理消息发送
    const handleMessageSend = useCallback((content: string, _attachment?: unknown[]) => {
        if (!currentBotId || isStreaming) return;

        // 添加用户消息
        const userMessage: Message = {
            id: `user-${Date.now()}`,
            role: 'user',
            content: content,
            createAt: Date.now(),
            status: 'complete',
        };

        // 添加 AI 消息占位
        const assistantMessage: Message = {
            id: `assistant-${Date.now()}`,
            role: 'assistant',
            content: '',
            createAt: Date.now(),
            status: 'loading',
        };

        setMessages(prev => [...prev, userMessage, assistantMessage]);
        setIsStreaming(true);

        // 发送请求
        abortControllerRef.current = conversationApi.chat(
            {
                bot_id: currentBotId,
                conversation_id: currentConversationId.current,
                query: userMessage.content,
                scene: Scene.Playground,
                content_type: 'text',
            },
            // onMessage
            (data) => {
                setMessages(prev => {
                    const newMessages = [...prev];
                    const lastIdx = newMessages.length - 1;
                    if (lastIdx >= 0 && newMessages[lastIdx].role === 'assistant') {
                        newMessages[lastIdx] = {
                            ...newMessages[lastIdx],
                            content: newMessages[lastIdx].content + (data.message?.content || ''),
                            status: data.is_finish ? 'complete' : 'loading',
                        };
                    }
                    return newMessages;
                });
            },
            // onError
            (error) => {
                console.error('Chat error:', error);
                setMessages(prev => {
                    const newMessages = [...prev];
                    const lastIdx = newMessages.length - 1;
                    if (lastIdx >= 0 && newMessages[lastIdx].role === 'assistant') {
                        newMessages[lastIdx] = {
                            ...newMessages[lastIdx],
                            content: newMessages[lastIdx].content || '抱歉，发生了错误，请重试。',
                            status: 'error',
                        };
                    }
                    return newMessages;
                });
                setIsStreaming(false);
            },
            // onDone
            () => {
                setMessages(prev => {
                    const newMessages = [...prev];
                    const lastIdx = newMessages.length - 1;
                    if (lastIdx >= 0 && newMessages[lastIdx].role === 'assistant') {
                        newMessages[lastIdx] = {
                            ...newMessages[lastIdx],
                            status: 'complete',
                        };
                    }
                    return newMessages;
                });
                setIsStreaming(false);
            }
        );
    }, [currentBotId, isStreaming]);

    // 处理消息重发
    const handleMessageReset = useCallback(() => {
        if (messages.length < 2) return;

        // 获取最后一条用户消息
        let lastUserMsgIdx = -1;
        for (let i = messages.length - 1; i >= 0; i--) {
            if (messages[i].role === 'user') {
                lastUserMsgIdx = i;
                break;
            }
        }
        if (lastUserMsgIdx === -1) return;

        const userMessage = messages[lastUserMsgIdx];

        // 移除最后的 AI 回复，重新发送
        setMessages(prev => prev.slice(0, lastUserMsgIdx + 1));

        // 重新发送
        setTimeout(() => {
            handleMessageSend(userMessage.content);
        }, 100);
    }, [messages, handleMessageSend]);

    // 处理点击提示
    const handleHintClick = useCallback((hint: string) => {
        handleMessageSend(hint);
    }, [handleMessageSend]);

    // 处理清除上下文
    const handleClearContext = useCallback(async () => {
        // 中止当前流式请求
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
            abortControllerRef.current = null;
        }
        setIsStreaming(false);

        // 调用后端清除会话
        if (currentConversationId.current) {
            try {
                await conversationApi.clearConversation({ ConversationID: currentConversationId.current });
            } catch (error) {
                console.error('Failed to clear conversation:', error);
            }
        }

        setMessages([]);
    }, []);

    // 停止生成
    const handleStopGenerate = useCallback(() => {
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
            abortControllerRef.current = null;
        }
        setIsStreaming(false);

        // 标记最后一条消息为完成
        setMessages(prev => {
            const newMessages = [...prev];
            const lastIdx = newMessages.length - 1;
            if (lastIdx >= 0 && newMessages[lastIdx].status === 'loading') {
                newMessages[lastIdx] = {
                    ...newMessages[lastIdx],
                    status: 'complete',
                };
            }
            return newMessages;
        });
    }, []);

    return (
        <Layout style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Header style={{
                backgroundColor: 'var(--semi-color-bg-1)',
                padding: '0 16px',
                borderBottom: '1px solid var(--semi-color-border)',
                flexShrink: 0,
            }}>
                <Title heading={5} style={{ margin: 0, lineHeight: '56px' }}>
                    预览与调试
                </Title>
            </Header>
            <Content style={{
                flex: 1,
                overflow: 'hidden',
                backgroundColor: 'var(--semi-color-bg-0)',
            }}>
                <Chat
                    chats={messages}
                    onChatsChange={handleChatsChange as any}
                    onMessageSend={handleMessageSend}
                    onMessageReset={handleMessageReset}
                    onClear={handleClearContext}
                    onStopGenerator={handleStopGenerate}
                    onHintClick={handleHintClick}
                    roleConfig={roleConfig}
                    hints={defaultHints}
                    showStopGenerate={isStreaming}
                    placeholder="输入消息开始对话..."
                    sendHotKey="enter"
                    align="leftRight"
                    mode="bubble"
                    style={{ height: '100%' }}
                    action="#"
                    uploadProps={{ action: '', disabled: true}}
                />
            </Content>
        </Layout>
    );
};

export default ChatPanel;