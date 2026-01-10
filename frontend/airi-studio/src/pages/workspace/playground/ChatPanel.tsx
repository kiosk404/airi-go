import { Layout, Typography, Spin, Avatar, Button, Tooltip } from '@douyinfe/semi-ui';
import { IconDelete } from '@douyinfe/semi-icons';
import React, { useCallback, useState, useRef, useEffect } from 'react';
import { AIChatInput } from '@douyinfe/semi-ui';
import { useParams } from 'react-router-dom';
import * as conversationApi from '@/services/conversation';
import { Scene } from '@/services/conversation';
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

interface ChatPanelProps {
    botInfo?: BotInfo;
    conversationId?: string;
}

const ChatPanel: React.FC<ChatPanelProps> = ({ botInfo, conversationId: propConversationId }) => {
    const params = useParams<{ id: string }>();

    // 优先使用 props，否则从 URL 获取
    const currentBotId = botInfo && botInfo.BotId ? botInfo.BotId : (params.id || '');
    const currentConversationId = useRef<string>(propConversationId ?? '');

    const [messages, setMessages] = useState<Message[]>([]);
    const [isStreaming, setIsStreaming] = useState(false);
    const [isLoadingHistory, setIsLoadingHistory] = useState(false);
    const abortControllerRef = useRef<AbortController | null>(null);

    // 将后端消息转换为前端消息格式
    const convertToMessage = useCallback((msg: conversationApi.ChatMessage): Message => {
        const isUser = msg.role === 'user';
        return {
            id: msg.message_id || `msg-${Date.now()}-${Math.random()}`,
            role: isUser ? 'user' : 'assistant',
            content: msg.content || '',
            createAt: typeof msg.content_time === 'number' ? msg.content_time : Date.now(),
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
                bot_id: String(currentBotId),
                conversation_id: currentConversationId.current || undefined,
                count: 50,
                scene: Scene.Playground,
                cursor: "0"
            });

            if (resp.code === 0) {
                // 如果返回了 conversation_id，保存起来
                if (resp.conversation_id) {
                    currentConversationId.current = resp.conversation_id;
                }

                if (resp.message_list && resp.message_list.length > 0) {
                    const historyMessages = resp.message_list.map(convertToMessage);
                    // 按时间正序排列
                    historyMessages.sort((a: Message, b: Message) => a.createAt - b.createAt);
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

    // 发送消息核心逻辑
    const sendMessage = useCallback((content: string, _attachment?: unknown[]) => {
        if (!currentBotId || isStreaming) return;

        console.log('[ChatPanel] sendMessage called with:', {
            botId: currentBotId,
            conversationId: currentConversationId.current,
            content,
        });

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
                draft_mode: true,
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

    // 处理消息发送 - 适配 AIChatInput
    // MessageContent 结构: { references, attachments, inputContents, setup }
    // inputContents 是数组，每个元素有 type 和对应内容（如 { type: 'text', text: '...' }）
    const handleMessageSend = useCallback((messageContent: any) => {
        // 从 inputContents 数组中提取文本内容
        const inputContents = messageContent?.inputContents || [];
        const textContent = inputContents
            .filter((item: { type: string }) => item.type === 'text')
            .map((item: { text: string }) => item.text || '')
            .join('');

        if (!currentBotId || isStreaming || !textContent.trim()) return;

        sendMessage(textContent);
    }, [currentBotId, isStreaming, sendMessage]);

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
                await conversationApi.clearConversation({ conversation_id: currentConversationId.current });
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

    // 消息列表滚动到底部
    const messagesEndRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    // 渲染单条消息
    const renderMessage = (msg: Message) => {
        const config = roleConfig[msg.role];
        const isUser = msg.role === 'user';
        const isLoading = msg.status === 'loading' && !msg.content;

        return (
            <div
                key={msg.id}
                style={{
                    display: 'flex',
                    flexDirection: isUser ? 'row-reverse' : 'row',
                    alignItems: 'flex-start',
                    marginBottom: 16,
                    gap: 12,
                }}
            >
                <Avatar
                    src={isUser ? config.avatar : (botInfo?.IconUrl || config.avatar)}
                    size="small"
                    style={{ flexShrink: 0 }}
                />
                <div
                    style={{
                        maxWidth: '70%',
                        padding: '10px 14px',
                        borderRadius: 12,
                        backgroundColor: isUser
                            ? 'var(--semi-color-primary)'
                            : 'var(--semi-color-fill-0)',
                        color: isUser ? '#fff' : 'var(--semi-color-text-0)',
                        wordBreak: 'break-word',
                    }}
                >
                    {isLoading ? (
                        <Spin size="small" />
                    ) : (
                        <div style={{ whiteSpace: 'pre-wrap' }}>
                            {msg.content}
                        </div>
                    )}
                    {msg.status === 'error' && (
                        <Text type="danger" size="small" style={{ display: 'block', marginTop: 4 }}>
                            发送失败
                        </Text>
                    )}
                </div>
            </div>
        );
    };

    return (
        <Layout style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Header style={{
                backgroundColor: 'var(--semi-color-bg-1)',
                padding: '0 16px',
                height: '56px',
                borderBottom: '1px solid var(--semi-color-border)',
                flexShrink: 0,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
            }}>
                <Title heading={5} style={{ margin: 0 }}>
                    预览与调试
                </Title>
                <div style={{ display: 'flex', gap: 8 }}>
                    {messages.length > 0 && (
                        <>
                            <Tooltip content="清空对话">
                                <Button
                                    icon={<IconDelete />}
                                    theme="borderless"
                                    onClick={handleClearContext}
                                    disabled={isStreaming}
                                />
                            </Tooltip>
                        </>
                    )}
                </div>
            </Header>
            <Content style={{
                flex: 1,
                overflow: 'hidden',
                backgroundColor: 'var(--semi-color-bg-0)',
                display: 'flex',
                flexDirection: 'column',
            }}>
                {isLoadingHistory ? (
                    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
                        <Spin size="large" tip="加载历史消息..." />
                    </div>
                ) : (
                    <>
                        {/* 消息列表区域 */}
                        <div style={{
                            flex: 1,
                            overflow: 'auto',
                            padding: '16px 24px',
                            position: 'relative',
                        }}>
                            {messages.length === 0 ? (
                                <div style={{
                                    height: '100%',
                                    display: 'flex',
                                    flexDirection: 'column',
                                    justifyContent: 'center',
                                    alignItems: 'center',
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
                            ) : (
                                <>
                                    {messages.map(renderMessage)}
                                    <div ref={messagesEndRef} />
                                </>
                            )}
                        </div>

                        {/* 输入框区域 */}
                        <div style={{
                            padding: '12px 24px 16px',
                            borderTop: '1px solid var(--semi-color-border)',
                            backgroundColor: 'var(--semi-color-bg-1)',
                        }}>
                            <div className="chat-input-wrapper">
                                <AIChatInput
                                    placeholder="输入消息开始对话..."
                                    onMessageSend={handleMessageSend}
                                    onStopGenerate={handleStopGenerate}
                                    generating={isStreaming}
                                    uploadProps={{
                                        action: '#',
                                    }}
                                    showUploadFile={false}
                                />
                            </div>
                        </div>
                    </>
                )}
            </Content>
        </Layout>
    );
};

export default ChatPanel;