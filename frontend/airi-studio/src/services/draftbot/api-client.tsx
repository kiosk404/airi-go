/**
 * DraftBot API Client
 * 对应后端接口: /api/draftbot/* 和 /api/playground_api/draftbot/*
 */

import { httpClient } from '@/services/core';

const API_BASE = '/api/draftbot';
const PLAYGROUND_API_BASE = '/api/playground_api/draftbot';

// ============ 类型定义 ============

/** Bot 显示信息请求 */
export interface GetDisplayInfoRequest {
    bot_id: string;
}

/** Tab 显示信息 */
export interface TabDisplayItems {
    plugin_tab_status?: number;
    workflow_tab_status?: number;
    knowledge_tab_status?: number;
    database_tab_status?: number;
}

/** 显示信息数据 */
export interface DraftBotDisplayInfoData {
    tab_display_info?: TabDisplayItems;
}

/** Bot 显示信息响应 */
export interface GetDisplayInfoResponse {
    code: number;
    msg: string;
    data?: DraftBotDisplayInfoData;
}

/** 模型信息 */
export interface ModelInfo {
    model_id?: string;
    model_name?: string;
    provider?: string;
}

/** Prompt 信息 */
export interface PromptInfo {
    system_prompt?: string;
    user_prompt?: string;
}

/** 插件信息 */
export interface PluginInfo {
    api_id?: string;
    name?: string;
    description?: string;
    icon?: string;
}

/** 工作流信息 */
export interface WorkflowInfo {
    workflow_id?: string;
    name?: string;
}

/** 开场白信息 */
export interface OnboardingInfo {
    prologue?: string;
    suggested_questions?: string[];
}

/** 知识库信息 */
export interface Knowledge {
    knowledge_ids?: string[];
}

/** 变量信息 */
export interface Variable {
    name?: string;
    type?: string;
    default_value?: string;
}

/** 数据库信息 */
export interface Database {
    database_id?: string;
    name?: string;
}

/** 建议回复信息 */
export interface SuggestReplyInfo {
    enabled?: boolean;
}

/** Bot 完整信息 */
export interface BotInfo {
    BotId: string;
    Name: string;
    Description: string;
    IconUri: string;
    IconUrl?: string;
    CreatorId: string;
    CreateTime: number;
    UpdateTime: number;
    Version?: string;
    ModelInfo?: ModelInfo;
    PromptInfo?: PromptInfo;
    PluginInfoList?: PluginInfo[];
    WorkflowInfoList?: WorkflowInfo[];
    OnboardingInfo?: OnboardingInfo;
    Knowledge?: Knowledge;
    VariableList?: Variable[];
    DatabaseList?: Database[];
    SuggestReplyInfo?: SuggestReplyInfo;
    BotMode?: number;
    BackgroundImageInfoList?: unknown[];
    ShortcutSort?: string[];
}

/** 获取 Bot 信息请求 */
export interface GetDraftBotInfoRequest {
    bot_id: string;
}

/** 获取 Bot 信息数据 */
export interface GetDraftBotInfoData {
    bot_info?: BotInfo;
    editable?: boolean;
    has_unpublished_change?: boolean;
}

/** 获取 Bot 信息响应 */
export interface GetDraftBotInfoResponse {
    code: number;
    msg: string;
    data?: GetDraftBotInfoData;
}

// ============ API 方法 ============

/**
 * 获取 Bot 显示信息
 * POST /api/draftbot/get_display_info
 */
export async function getDisplayInfo(req: GetDisplayInfoRequest): Promise<GetDisplayInfoResponse> {
    return httpClient.post<GetDisplayInfoResponse>(`${API_BASE}/get_display_info`, req);
}

/**
 * 获取 Bot 完整信息
 * POST /api/playground_api/draftbot/get_draft_bot_info
 */
export async function getDraftBotInfo(req: GetDraftBotInfoRequest): Promise<GetDraftBotInfoResponse> {
    return httpClient.post<GetDraftBotInfoResponse>(`${PLAYGROUND_API_BASE}/get_draft_bot_info`, req);
}