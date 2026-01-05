import { Layout, Typography, Spin, Avatar } from '@douyinfe/semi-ui';
import React, { useCallback, useState, useRef, useEffect } from 'react';
import { Chat } from '@douyinfe/semi-ui';
import { useParams } from 'react-router-dom';
import * as conversationApi from '@/services/conversation';
import { Scene, MsgParticipantType } from '@/services/conversation';
import { BotInfo } from '@/services/draftbot';

const { Title, Text } = Typography;
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
    botInfo?: BotInfo;
    conversationId?: string;
}

const ChatPanel: React.FC<ChatPanelProps> = ({ botInfo, conversationId: propConversationId }) => {
    const params = useParams<{ id: string }>();

    console.log("kiosk: ", botInfo)

    // 优先使用 props，否则从 URL 获取
    const currentBotId = botInfo && botInfo.BotId ? botInfo.BotId : (params.id || '');
    const currentConversationId = useRef<string>(propConversationId ?? '');

    const [messages, setMessages] = useState<Message[]>([]);
    const [isStreaming, setIsStreaming] = useState(false);
    const [isLoadingHistory, setIsLoadingHistory] = useState(false);
    const abortControllerRef = useRef<AbortController | null>(null);

    // 将后端消息转换为前端消息格式
    const convertToMessage = useCallback((msg: conversationApi.ChatMessage): Message => {
        const isUser = msg.sender?.type === MsgParticipantType.User;
        return {
            id: msg.id || `msg-${Date.now()}-${Math.random()}`,
            role: isUser ? 'user' : 'assistant',
            content: msg.content || '',
            createAt: msg.create_time || Date.now(),
            status: 'complete',
        };
    }, []);

    // 加载历史消息
    const loadMessageHistory = useCallback(async () => {
        if (!currentBotId) return;

        setIsLoadingHistory(true);
        try {
            // 直接获取消息列表，后端会根据 bot_id 返回或创建会话
            const resp = await conversationApi.getMessageList({
                bot_id: currentBotId,
                conversation_id: currentConversationId.current || undefined,
                page_size: 50,
                scene: Scene.Playground,
                cursor: "0"
            });

            if (resp.code === 0) {
                // 如果返回了 conversation_id，保存起来
                if (resp.conversation_id) {
                    currentConversationId.current = resp.conversation_id;
                }

                if (resp.messages && resp.messages.length > 0) {
                    const historyMessages = resp.messages.map(convertToMessage);
                    // 按时间正序排列
                    historyMessages.sort((a, b) => a.createAt - b.createAt);
                    setMessages(historyMessages);
                }
            }
        } catch (error) {
            console.error('Failed to load message history:', error);
        } finally {
            setIsLoadingHistory(false);
        }
    }, [currentBotId, convertToMessage]);

    // 组件加载时获取历史消息
    useEffect(() => {
        loadMessageHistory();
    }, [loadMessageHistory]);

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
                height: '50px',
                borderBottom: '1px solid var(--semi-color-border)',
                flexShrink: 0,
                display: 'flex',
                alignItems: 'center',
            }}>
                <Title heading={5} style={{ margin: 0, lineHeight: '56px' }}>
                    预览与调试
                </Title>
            </Header>
            <Content style={{
                flex: 1,
                overflow: 'hidden',
                backgroundColor: 'var(--semi-color-bg-0)',
                position: 'relative',
            }}>
                {isLoadingHistory ? (
                    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
                        <Spin size="large" tip="加载历史消息..." />
                    </div>
                ) : (
                    <>
                        {messages.length === 0 && (
                            <div style={{
                                position: 'absolute',
                                top: 0,
                                left: 0,
                                right: 0,
                                bottom: 120,
                                display: 'flex',
                                flexDirection: 'column',
                                justifyContent: 'center',
                                alignItems: 'center',
                                pointerEvents: 'none',
                                zIndex: 1,
                            }}>
                                <Avatar
                                    src={botInfo?.IconUrl || roleConfig.assistant.avatar}
                                    size="extra-large"
                                    style={{ width: 80, height: 80, marginBottom: 16 }}
                                />
                                <Text strong style={{ fontSize: 18, marginBottom: 8 }}>
                                    {botInfo?.Name || 'Airi'}
                                </Text>
                                {botInfo?.Description && (
                                    <Text type="tertiary" style={{ textAlign: 'center', maxWidth: 300 }}>
                                        {botInfo.Description}
                                    </Text>
                                )}
                            </div>
                        )}
                        <Chat
                            chats={messages}
                            onChatsChange={handleChatsChange as any}
                            onMessageSend={handleMessageSend}
                            onMessageReset={handleMessageReset}
                            onClear={handleClearContext}
                            onStopGenerator={handleStopGenerate}
                            onHintClick={handleHintClick}
                            roleConfig={roleConfig}
                            hints={messages.length > 0 ? defaultHints : []}
                            showStopGenerate={isStreaming}
                            placeholder="输入消息开始对话..."
                            sendHotKey="enter"
                            align="leftRight"
                            mode="bubble"
                            style={{ height: '100%' }}
                            action="#"
                            uploadProps={{ action: '', disabled: true}}
                        />
                    </>
                )}
            </Content>
        </Layout>
    );
};

export default ChatPanel;