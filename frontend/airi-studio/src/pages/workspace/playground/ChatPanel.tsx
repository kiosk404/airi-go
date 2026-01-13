import { Layout, Typography, Spin, Avatar, Button, Tooltip, Toast } from '@douyinfe/semi-ui';
import { IconDelete, IconCopy, IconQuote } from '@douyinfe/semi-icons';
import React, { useCallback, useState, useRef, useEffect, useMemo } from 'react';
import { AIChatInput, AIChatDialogue } from '@douyinfe/semi-ui';
import { useParams } from 'react-router-dom';
import * as conversationApi from '@/services/conversation';
import { Scene } from '@/services/conversation';
import { BotInfo } from '@/services/draftbot';

const { Title, Text } = Typography;
const { Header, Content } = Layout;

// 消息类型定义
interface ChatMessage {
    id: string | number;
    role: 'user' | 'assistant' | 'system';
    content: string;
    createdAt?: number;
    status?: 'in_progress' | 'failed' | 'completed' | 'cancelled' | 'queued' | 'incomplete';
}

// 角色配置类型
interface RoleMetadata {
    name?: string;
    avatar?: string;
}

interface ChatRoleConfig {
    user?: RoleMetadata;
    assistant?: RoleMetadata;
    system?: RoleMetadata;
}

interface ChatPanelProps {
    botInfo?: BotInfo;
    conversationId?: string;
}

// 引用类型
interface Reference {
    id: string;
    type ?: 'text' | 'file' | 'image' | 'code';
    content?: string;
    name?: string;
    url?: string;
}

const ChatPanel: React.FC<ChatPanelProps> = ({ botInfo, conversationId: propConversationId }) => {
    const params = useParams<{ id: string }>();

    const currentBotId = botInfo && botInfo.bot_id ? botInfo.bot_id : (params.id || '');
    const currentConversationId = useRef<string>(propConversationId ?? '');

    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [hints, setHints] = useState<string[]>([]);
    const [isStreaming, setIsStreaming] = useState(false);
    const [isLoadingHistory, setIsLoadingHistory] = useState(false);
    const [references, setReferences] = useState<Reference[]>([]);

    const abortControllerRef = useRef<AbortController | null>(null);
    const dialogueRef = useRef<any>(null);
    const chatInputRef = useRef<any>(null);

    // 角色配置
    const roleConfig: ChatRoleConfig = {
        user: {
            name: 'User',
            avatar: 'https://lf3-static.bytednsdoc.com/obj/eden-cn/ptlz_zlp/ljhwZthlaukjlkulzlp/docs-icon.png',
        },
        assistant: {
            name: botInfo?.name || 'Airi',
            avatar: botInfo?.icon_url || '🤖',
        },
        system: {
            name: 'System',
            avatar: '⚙️',
        },
    };

    // 将后端消息转换为 AIChatDialogue 消息格式
    const convertToMessage = useCallback((msg: conversationApi.ChatMessage): ChatMessage | null => {
        if (msg.type === 'generate_answer_finish' || msg.type === 'follow_up' || msg.type === 'verbose') {
            return null;
        }
        const isUser = msg.role === 'user';
        return {
            id: msg.message_id || `msg-${Date.now()}-${Math.random()}`,
            role: isUser ? 'user' : 'assistant',
            content: msg.content || '',
            createdAt: typeof msg.content_time === 'number' ? msg.content_time : Date.now(),
            status: 'completed',
        };
    }, []);

    // 从消息列表中提取 hints（建议问题）
    const extractHints = useCallback((messageList: conversationApi.ChatMessage[]): string[] => {
        return messageList
            .filter(msg => msg.type === 'follow_up' && msg.content)
            .map(msg => msg.content || '');
    }, []);

    // 加载历史消息
    const loadMessageHistory = useCallback(async () => {
        if (!currentBotId) return;

        setIsLoadingHistory(true);
        try {
            const resp = await conversationApi.getMessageList({
                bot_id: String(currentBotId),
                conversation_id: currentConversationId.current || undefined,
                count: 50,
                scene: Scene.Playground,
                cursor: "0"
            });

            if (resp.code === 0) {
                if (resp.conversation_id) {
                    currentConversationId.current = resp.conversation_id;
                }

                if (resp.message_list && resp.message_list.length > 0) {
                    // 转换普通消息
                    const historyMessages = resp.message_list
                        .map(convertToMessage)
                        .filter((msg): msg is ChatMessage => msg !== null);
                    historyMessages.sort((a, b) => (a.createdAt || 0) - (b.createdAt || 0));
                    setMessages(historyMessages);
                }
            }
        } catch (error) {
            console.error('Failed to load message history:', error);
        } finally {
            setIsLoadingHistory(false);
        }
    }, [currentBotId, convertToMessage, extractHints]);

    // 组件加载时获取历史消息
    useEffect(() => {
        void loadMessageHistory();
    }, [loadMessageHistory]);

    // 发送消息核心逻辑
    const sendMessage = useCallback((content: string) => {
        if (!currentBotId || isStreaming) return;

        // 添加用户消息
        const userMessage: ChatMessage = {
            id: `user-${Date.now()}`,
            role: 'user',
            content: content,
            createdAt: Date.now(),
            status: 'completed',
        };

        // 添加 AI 消息占位
        const assistantMessage: ChatMessage = {
            id: `assistant-${Date.now()}`,
            role: 'assistant',
            content: '',
            createdAt: Date.now(),
            status: 'in_progress',
        };

        setMessages(prev => [...prev, userMessage, assistantMessage]);
        setIsStreaming(true);

        // 滚动到底部
        setTimeout(() => dialogueRef.current?.scrollToBottom(true), 100);

        // 发送请求
        abortControllerRef.current = conversationApi.chat(
            {
                bot_id: currentBotId,
                conversation_id: currentConversationId.current,
                query: content,
                scene: Scene.Playground,
                content_type: 'text',
                draft_mode: true,
            },
            // onMessage
            (data) => {
                // 忽略 ack 类型的消息
                if (data.message?.type === 'ack') {
                    return;
                }

                // 忽略 verbose 类型的消息
                if (data.message?.type === 'verbose') {
                    return;
                }

                setMessages(prev => {
                    const newMessages = [...prev];
                    const lastIdx = newMessages.length - 1;
                    if (lastIdx >= 0 && newMessages[lastIdx].role === 'assistant') {
                        newMessages[lastIdx] = {
                            ...newMessages[lastIdx],
                            content: (newMessages[lastIdx].content || '') + (data.message?.content || ''),
                            status: data.is_finish ? 'completed' : 'in_progress',
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
                            content: (newMessages[lastIdx].content as string) || '抱歉，发生了错误，请重试。',
                            status: 'failed',
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
                            status: 'completed',
                        };
                    }
                    return newMessages;
                });
                setIsStreaming(false);
            }
        );
    }, [currentBotId, isStreaming]);

    // 处理消息发送 - 适配 AIChatInput
    const handleMessageSend = useCallback((messageContent: any) => {
        const inputContents = messageContent?.inputContents || [];
        const textContent = inputContents
            .filter((item: { type: string }) => item.type === 'text')
            .map((item: { text: string }) => item.text || '')
            .join('');

        if (!currentBotId || isStreaming || !textContent.trim()) return;

        // 发送消息后清空 hints
        setHints([]);
        sendMessage(textContent);
    }, [currentBotId, isStreaming, sendMessage]);

    // 处理 hint 点击
    const handleHintClick = useCallback((hint: string) => {
        if (!currentBotId || isStreaming || !hint.trim()) return;

        // 点击 hint 后清空 hints 并发送消息
        setHints([]);
        sendMessage(hint);
    }, [currentBotId, isStreaming, sendMessage]);

    // 处理清除上下文
    const handleClearContext = useCallback(async () => {
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
            abortControllerRef.current = null;
        }
        setIsStreaming(false);

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

        setMessages(prev => {
            const newMessages = [...prev];
            const lastIdx = newMessages.length - 1;
            if (lastIdx >= 0 && newMessages[lastIdx].status === 'in_progress') {
                newMessages[lastIdx] = {
                    ...newMessages[lastIdx],
                    status: 'completed',
                };
            }
            return newMessages;
        });
    }, []);

    // 消息变化时滚动到底部
    useEffect(() => {
        dialogueRef.current?.scrollToBottom(true);
    }, [messages]);

    // 处理引用删除
    const handleReferenceDelete = useCallback((item: Reference) => {
        setReferences(prev => prev.filter(ref => ref.id !== item.id));
    }, []);

    // 处理引用点击
    const handleReferenceClick = useCallback((item: Reference) => {
        console.log('Reference clicked:', item);
        // 可以在这里实现点击引用后的逻辑，比如预览文件等
    }, []);

    // 引用消息到输入框
    const handleQuoteMessage = useCallback((message: ChatMessage) => {
        const contentPreview = message.content.length > 100
            ? message.content.substring(0, 100) + '...'
            : message.content;
        const ref: Reference = {
            id: `quote-${message.id}-${Date.now()}`,
            type: 'text',
            content: contentPreview,
        };
        setReferences(prev => {
            // 避免重复添加同一消息的引用
            if (prev.some(r => r.content === contentPreview)) {
                return prev;
            }
            return [...prev, ref];
        });
    }, []);

    // 复制消息内容
    const handleCopyMessage = useCallback((message: ChatMessage) => {
        navigator.clipboard.writeText(message.content).then(() => {
            Toast.success('已复制到剪贴板');
        }).catch(() => {
            Toast.error('复制失败');
        });
    }, []);

    // 对话框渲染配置 - 自定义操作按钮
    const dialogueRenderConfig = useMemo(() => ({
        renderDialogueAction: (props: any) => {
            const { message, className } = props;
            // 只为已完成的消息显示操作按钮
            if (message?.status === 'loading') {
                return null;
            }
            return (
                <div className={className} style={{ display: 'flex', gap: 4 }}>
                    <Tooltip content="复制">
                        <Button
                            icon={<IconCopy size="small" />}
                            size="small"
                            theme="borderless"
                            type="tertiary"
                            onClick={() => handleCopyMessage(message)}
                        />
                    </Tooltip>
                    <Tooltip content="引用">
                        <Button
                            icon={<IconQuote size="small" />}
                            size="small"
                            theme="borderless"
                            type="tertiary"
                            onClick={() => handleQuoteMessage(message)}
                        />
                    </Tooltip>
                </div>
            );
        },
    }), [handleCopyMessage, handleQuoteMessage]);

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
                        <Tooltip content="清空对话">
                            <Button
                                icon={<IconDelete />}
                                theme="borderless"
                                onClick={handleClearContext}
                                disabled={isStreaming}
                            />
                        </Tooltip>
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
                        {/* 对话区域 */}
                        <div style={{ flex: 1, overflow: 'hidden', position: 'relative' }}>
                            {messages.length === 0 ? (
                                <div style={{
                                    height: '100%',
                                    display: 'flex',
                                    flexDirection: 'column',
                                    justifyContent: 'center',
                                    alignItems: 'center',
                                }}>
                                    <Avatar
                                        src={botInfo?.icon_url || roleConfig.assistant?.avatar}
                                        size="extra-large"
                                        style={{ width: 80, height: 80, marginBottom: 16 }}
                                    />
                                    <Text strong style={{ fontSize: 18, marginBottom: 8 }}>
                                        {botInfo?.name || 'Airi'}
                                    </Text>
                                    {botInfo?.description && (
                                        <Text type="tertiary" style={{ textAlign: 'center', maxWidth: 300 }}>
                                            {botInfo.description}
                                        </Text>
                                    )}
                                </div>
                            ) : (
                                    <AIChatDialogue
                                        ref={dialogueRef}
                                        chats={messages as any}
                                        roleConfig={roleConfig as any}
                                        showReference
                                        onChatsChange={(chats) => chats && setMessages(chats as ChatMessage[])}
                                        mode="userBubble"
                                        align="leftRight"
                                        hints={hints}
                                        dialogueRenderConfig={dialogueRenderConfig}
                                        onHintClick={handleHintClick}
                                        style={{ flex: 1, overflow: 'auto' }}
                                    />

                            )}
                        </div>

                        {/* 输入框区域 */}
                        <div style={{
                            padding: '12px 24px 16px',
                            borderTop: '1px solid var(--semi-color-border)',
                            backgroundColor: 'var(--semi-color-bg-1)',
                        }}>
                            <AIChatInput
                                placeholder="输入消息开始对话..."
                                onMessageSend={handleMessageSend}
                                onStopGenerate={handleStopGenerate}
                                generating={isStreaming}
                                uploadProps={{ action: '#' }}
                                showUploadFile={false}
                                ref={chatInputRef}
                                onReferenceDelete={handleReferenceDelete as any}
                                onReferenceClick={handleReferenceClick as any}
                                references={references as any}
                            />
                        </div>
                    </>
                )}
            </Content>
        </Layout>
    );
};

export default ChatPanel;