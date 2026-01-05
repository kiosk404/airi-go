/**
 * Character (Agent) API Client
 * 对应后端接口: /api/intelligence_api/* 和 /api/draftbot/*
 */

import { httpClient } from '@/services/core';

const API_BASE = '/api/draftbot';
const INTELLIGENCE_API_BASE = '/api/intelligence_api/search';

// ============ 类型定义 ============

/** Intelligence 基础信息 */
export interface IntelligenceBasicInfo {
    id: string;
    name: string;
    description: string;
    icon_uri: string;
    icon_url: string;
    owner_id: string;
    create_time: string;
    update_time: string;
    status: number;
    publish_time?: string;
}

/** Intelligence 数据项 */
export interface IntelligenceData {
    basic_info: IntelligenceBasicInfo;
    type: number;
    publish_info?: {
        publish_time?: string;
    };
    permission_info?: unknown;
    owner_info?: unknown;
    favorite_info?: {
        is_fav?: boolean;
    };
}

/** 获取角色列表请求参数 */
export interface ListCharacterParams {
    name?: string;
    size?: number;
    cursor_id?: string;
}

/** 获取角色列表响应 */
export interface ListCharacterResponse {
    code: number;
    msg: string;
    data: {
        intelligences: IntelligenceData[];
        total: number;
        has_more: boolean;
        next_cursor_id?: string;
    };
}

/** 创建角色请求参数 */
export interface CreateCharacterParams {
    name: string;
    description?: string;
    icon_uri: string;
}

/** 创建角色响应 */
export interface CreateCharacterResponse {
    code: number;
    msg: string;
    data: {
        bot_id: string;
    };
}

export interface DeleteCharacterParams {
    bot_id: string;
}

export interface DeleteCharacterResponse {
    code: number;
    msg: string;
}

// ============ API 方法 ============

/**
 * 获取角色列表
 * POST /api/intelligence_api/search/get_draft_intelligence_list
 */
export async function listCharacters(params: ListCharacterParams): Promise<ListCharacterResponse> {
    return httpClient.post<ListCharacterResponse>(`${INTELLIGENCE_API_BASE}/get_draft_intelligence_list`, {
        name: params.name,
        size: params.size || 50,
        cursor_id: params.cursor_id,
        types: [1], // IntelligenceType.Bot = 1
    });
}

/**
 * 创建角色
 * POST /api/draftbot/create
 */
export async function createCharacter(params: CreateCharacterParams): Promise<CreateCharacterResponse> {
    return httpClient.post<CreateCharacterResponse>(`${API_BASE}/create`, params);
}

/**
 * 删除角色
 * POST /api/draftbot/delete
 */
export async function deleteCharacter(params: DeleteCharacterParams): Promise<DeleteCharacterResponse> {
    return httpClient.post<DeleteCharacterResponse>(`${API_BASE}/delete`, params);
}
