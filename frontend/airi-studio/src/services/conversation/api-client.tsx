/**
 * Conversation API Client
 * 对应后端接口: /api/conversation/*
 */

import { httpClient } from '@/services/core';

const API_BASE = '/api/conversation';


// 从 IDL 生成的类型中导入
import { Scene } from '@/api/generated/conversation/common';
import { MsgParticipantType } from '@/api/generated/conversation/message/MsgParticipantType';
import type { IAgentRunRequestArgs } from '@/api/generated/conversation/run/AgentRunRequest';
import type { IRunStreamResponseArgs } from '@/api/generated/conversation/run/RunStreamResponse';
import type { IGetMessageListRequestArgs } from '@/api/generated/conversation/message/GetMessageListRequest';
import type { IGetMessageListResponseArgs } from '@/api/generated/conversation/message/GetMessageListResponse';
import type { IChatMessageArgs } from '@/api/generated/conversation/message/ChatMessage';
import type { IMsgParticipantInfoArgs } from '@/api/generated/conversation/message/MsgParticipantInfo';
import type { IBreakMessageRequestArgs } from '@/api/generated/conversation/message/BreakMessageRequest';
import type { IBreakMessageResponseArgs } from '@/api/generated/conversation/message/BreakMessageResponse';
import type { IDeleteMessageRequestArgs } from '@/api/generated/conversation/message/DeleteMessageRequest';
import type { IDeleteMessageResponseArgs } from '@/api/generated/conversation/message/DeleteMessageResponse';
import type { ICreateConversationRequestArgs } from '@/api/generated/conversation/conversation/CreateConversationRequest';
import type { ICreateConversationResponseArgs } from '@/api/generated/conversation/conversation/CreateConversationResponse';
import type { IClearConversationHistoryRequestArgs } from '@/api/generated/conversation/conversation/ClearConversationHistoryRequest';
import type { IClearConversationHistoryResponseArgs } from '@/api/generated/conversation/conversation/ClearConversationHistoryResponse';

// 重新导出枚举和类型供外部使用
export { Scene, MsgParticipantType };
export type AgentRunRequest = IAgentRunRequestArgs;
export type RunStreamResponse = IRunStreamResponseArgs;
export type GetMessageListRequest = IGetMessageListRequestArgs;
export type GetMessageListResponse = IGetMessageListResponseArgs;
export type ChatMessage = IChatMessageArgs;
export type MsgParticipantInfo = IMsgParticipantInfoArgs;
export type BreakMessageRequest = IBreakMessageRequestArgs;
export type BreakMessageResponse = IBreakMessageResponseArgs;
export type DeleteMessageRequest = IDeleteMessageRequestArgs;
export type DeleteMessageResponse = IDeleteMessageResponseArgs;
export type CreateConversationRequest = ICreateConversationRequestArgs;
export type CreateConversationResponse = ICreateConversationResponseArgs;
export type ClearConversationRequest = IClearConversationHistoryRequestArgs;
export type ClearConversationResponse = IClearConversationHistoryResponseArgs;

export const RunEvent = {
    Done: 'done',
    Error: 'error',
    Message: 'message',
}as const

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
