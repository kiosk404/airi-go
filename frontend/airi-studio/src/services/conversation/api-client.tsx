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

/** 消息内容类型 */
export const ContentType = {
    Text: 'text',
    Image: 'image',
    Audio: 'audio',
    Video: 'video',
    Link: 'link',
    Music: 'music',
    Card: 'card',
    App: 'app',
    File: 'file',
    Mix: 'mix',
    MixApi: 'object_string',
} as const;

/** SSE 事件类型 */
export const RunEvent = {
    Message: 'message',
    Done: 'done',
    Error: 'error',
} as const;

// ============ 类型定义 ============

/** 额外信息 */
export interface ExtraInfo {
    local_message_id?: string;
    input_tokens?: string;
    output_tokens?: string;
    token?: string;
    plugin_status?: string;
    time_cost?: string;
    workflow_tokens?: string;
    bot_state?: string;
    plugin_request?: string;
    tool_name?: string;
    plugin?: string;
    log_id?: string;
    stream_id?: string;
    message_title?: string;
    stream_plugin_running?: string;
    new_section_id?: string;
    remove_query_id?: string;
    execute_display_name?: string;
    task_type?: string;
    refer_format?: string;
    call_id?: string;
}

/** 消息参与者信息 */
export interface MsgParticipantInfo {
    id?: string;
    type?: MsgParticipantType;
    name?: string;
    desc?: string;
    avatar_url?: string;
    user_id?: string;
    user_name?: string;
    allow_mention?: boolean;
    access_path?: string;
    is_fav?: boolean;
    allow_share?: boolean;
}

/** 聊天消息 */
export interface ChatMessage {
    id?: string;
    role?: string;
    type?: string;
    content?: string;
    content_type?: string;
    message_id?: string;
    reply_id?: string;
    section_id?: string;
    extra_info?: ExtraInfo;
    status?: string;
    broken_pos?: number;
    sender_id?: string;
    sender?: MsgParticipantInfo;
    mention_list?: MsgParticipantInfo[];
    content_time?: number;
    create_time?: number;
    message_index?: number;
    source?: number;
    reply_message?: ChatMessage;
    reasoning_content?: string;
    card_status?: Record<string, string>;
}

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
    regen_message_id?: string;
    local_message_id?: string;
    preset_bot?: string;
    device_id?: string;
    mention_list?: MsgParticipantInfo[];
}

/** SSE 消息事件 */
export interface RunStreamResponse {
    message: ChatMessage;
    is_finish?: boolean;
    index: number;
    conversation_id: string;
    seq_id: number;
}

/** 获取消息列表请求 */
export interface GetMessageListRequest {
    bot_id: string;
    conversation_id?: string;
    cursor?: string;
    count?: number;
    page_size?: number;
    scene?: Scene;
    draft_mode?: boolean;
    preset_bot?: string;
    biz_kind?: string;
    load_direction?: number;
    must_append?: boolean;
    share_id?: string;
}

/** 获取消息列表响应 */
export interface GetMessageListResponse {
    code: number;
    msg: string;
    conversation_id?: string;
    messages?: ChatMessage[];
    message_list?: ChatMessage[];
    has_more?: boolean;
    hasmore?: boolean;
    cursor?: string;
    next_cursor?: string;
    next_has_more?: boolean;
    last_section_id?: string;
    participant_info_map?: Record<string, MsgParticipantInfo>;
    read_message_index?: number;
    connector_conversation_id?: string;
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

/** 删除消息请求 */
export interface DeleteMessageRequest {
    conversation_id: string;
    message_id: string;
    scene?: Scene;
    bot_id?: string;
}

/** 删除消息响应 */
export interface DeleteMessageResponse {
    code: number;
    msg: string;
}

/** 中断消息请求 */
export interface BreakMessageRequest {
    conversation_id: string;
    query_message_id: string;
    answer_message_id?: string;
    broken_pos?: number;
    scene?: Scene;
}

/** 中断消息响应 */
export interface BreakMessageResponse {
    code: number;
    msg: string;
}

// ============ API 方法 ============

/**
 * 发送聊天消息 (SSE 流式)
 * POST /api/conversation/chat
 * 返回 AbortController 用于取消请求
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
        let currentEvent = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split('\n');
            buffer = lines.pop() || '';

            for (const line of lines) {
                const trimmedLine = line.trim();

                if (trimmedLine.startsWith('event:')) {
                    currentEvent = trimmedLine.slice(6).trim();
                } else if (trimmedLine.startsWith('data:')) {
                    const dataStr = trimmedLine.slice(5).trim();

                    if (currentEvent === RunEvent.Done || dataStr === '[DONE]') {
                        onDone?.();
                        currentEvent = '';
                        continue;
                    }

                    if (currentEvent === RunEvent.Error) {
                        try {
                            const errorData = JSON.parse(dataStr);
                            onError?.(new Error(errorData.msg || errorData.message || 'Unknown error'));
                        } catch {
                            onError?.(new Error(dataStr || 'Unknown error'));
                        }
                        currentEvent = '';
                        continue;
                    }

                    if (dataStr) {
                        try {
                            const data = JSON.parse(dataStr) as RunStreamResponse;
                            onMessage(data);
                        } catch (e) {
                            console.warn('Failed to parse SSE data:', dataStr);
                        }
                    }
                    currentEvent = '';
                } else if (trimmedLine === '') {
                    // 空行，重置事件类型
                    currentEvent = '';
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

/**
 * 删除消息
 * POST /api/conversation/delete_message
 */
export async function deleteMessage(req: DeleteMessageRequest): Promise<DeleteMessageResponse> {
    return httpClient.post<DeleteMessageResponse>(`${API_BASE}/delete_message`, req);
}

/**
 * 中断消息生成
 * POST /api/conversation/break_message
 */
export async function breakMessage(req: BreakMessageRequest): Promise<BreakMessageResponse> {
    return httpClient.post<BreakMessageResponse>(`${API_BASE}/break_message`, req);
}
