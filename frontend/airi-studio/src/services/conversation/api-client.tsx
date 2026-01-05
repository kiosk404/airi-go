/**
 * Conversation API Client
 * 对应后端接口: /api/conversation/*
 */

import { httpClient } from '@/services/core';

const API_BASE = '/api/conversation';

// ============ 枚举定义 ============

/** 场景类型 */
export enum Scene {
    Default = 0,
    Explore = 1,
    BotStore = 2,
    AiriHome = 3,
    Playground = 4,
    Evaluation = 5,
    AgentAPP = 6,
    PromptOptimize = 7,
    GenerateAgentInfo = 8,
    SceneOpenApi = 9,
    SceneWorkflow = 50,
}

/** 消息参与者类型 */
export enum MsgParticipantType {
    Bot = 1,
    User = 2,
}

// ============ 类型定义 ============

/** 聊天请求 */
export interface AgentRunRequest {
    bot_id: string;
    conversation_id: string;
    query: string;
    scene?: Scene;
    content_type?: string;
    draft_mode?: boolean;
    extra?: Record<string, string>;
    custom_variables?: Record<string, string>;
}

/** SSE 消息事件 */
export interface RunStreamResponse {
    message: ChatMessage;
    is_finish?: boolean;
    index: number;
    conversation_id: string;
    seq_id: number;
}

/** 聊天消息 */
export interface ChatMessage {
    id?: string;
    content?: string;
    role?: string;
    type?: number;
    sender?: MsgParticipantInfo;
    create_time?: number;
    status?: number;
}

/** 消息参与者信息 */
export interface MsgParticipantInfo {
    type?: MsgParticipantType;
    id?: string;
    name?: string;
    avatar?: string;
}

/** 获取消息列表请求 */
export interface GetMessageListRequest {
    bot_id: string;
    conversation_id?: string;
    page_num?: number;
    page_size?: number;
    cursor?: string;
    scene?: number;
}

/** 获取消息列表响应 */
export interface GetMessageListResponse {
    code: number;
    msg: string;
    conversation_id?: string;
    messages?: ChatMessage[];
    has_more?: boolean;
    cursor?: string;
}

/** 创建会话请求 */
export interface CreateConversationRequest {
    BotId?: string;
    MetaData?: Record<string, string>;
}

/** 会话数据 */
export interface ConversationData {
    Id: string;
    CreatedAt: number;
    MetaData?: Record<string, string>;
    CreatorID?: string;
    LastSectionID?: string;
    Name?: string;
}

/** 创建会话响应 */
export interface CreateConversationResponse {
    code: number;
    msg: string;
    ConversationData?: ConversationData;
}

/** 清除会话历史请求 */
export interface ClearConversationRequest {
    ConversationID: string;
}

/** 清除会话历史响应 */
export interface ClearConversationResponse {
    code: number;
    msg: string;
    data?: {
        id: string;
        conversation_id: string;
    };
}

// ============ API 方法 ============

/**
 * 发送聊天消息 (SSE 流式)
 * POST /api/conversation/chat
 * 返回 EventSource 用于接收流式响应
 */
export function chat(
    req: AgentRunRequest,
    onMessage: (data: RunStreamResponse) => void,
    onError?: (error: Error) => void,
    onDone?: () => void
): AbortController {
    const controller = new AbortController();

    fetch(`${API_BASE}/chat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'text/event-stream',
        },
        body: JSON.stringify(req),
        signal: controller.signal,
        credentials: 'include',
    }).then(async (response) => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const reader = response.body?.getReader();
        if (!reader) {
            throw new Error('No response body');
        }

        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split('\n');
            buffer = lines.pop() || '';

            for (const line of lines) {
                if (line.startsWith('data:')) {
                    const dataStr = line.slice(5).trim();
                    if (dataStr) {
                        try {
                            const data = JSON.parse(dataStr);
                            onMessage(data);
                        } catch (e) {
                            console.warn('Failed to parse SSE data:', dataStr);
                        }
                    }
                } else if (line.startsWith('event:')) {
                    const event = line.slice(6).trim();
                    if (event === 'done') {
                        onDone?.();
                    } else if (event === 'error') {
                        // 错误会在下一个 data 行
                    }
                }
            }
        }

        onDone?.();
    }).catch((error) => {
        if (error.name !== 'AbortError') {
            onError?.(error);
        }
    });

    return controller;
}

/**
 * 获取消息列表
 * POST /api/conversation/get_message_list
 */
export async function getMessageList(req: GetMessageListRequest): Promise<GetMessageListResponse> {
    return httpClient.post<GetMessageListResponse>(`${API_BASE}/get_message_list`, req);
}

/**
 * 创建会话
 * POST /api/conversation/create
 */
export async function createConversation(req: CreateConversationRequest): Promise<CreateConversationResponse> {
    return httpClient.post<CreateConversationResponse>(`${API_BASE}/create`, req);
}

/**
 * 清除会话历史
 * POST /api/conversation/clear
 */
export async function clearConversation(req: ClearConversationRequest): Promise<ClearConversationResponse> {
    return httpClient.post<ClearConversationResponse>(`${API_BASE}/clear`, req);
}
