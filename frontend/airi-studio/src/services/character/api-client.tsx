/**
 * Character (Agent) API Client
 * 对应后端接口: /api/draftbot/*
 */

import { httpClient } from '@/services/core';

const API_BASE = '/api/draftbot';

// ============ 类型定义 ============

/** 角色/Agent 项 */
export interface CharacterItem {
    id: string;
    name: string;
    description: string;
    icon_uri: string;
    created_at: number;
    updated_at: number;
}

/** 获取角色列表请求参数 */
export interface ListCharacterParams {
    page?: number;
    page_size?: number;
}

/** 获取角色列表响应 */
export interface ListCharacterResponse {
    code: number;
    msg: string;
    data: {
        items: CharacterItem[];
        total: number;
        page: number;
        page_size: number;
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

// ============ API 方法 ============

/**
 * 获取角色列表
 * GET /api/draftbot/list
 */
export async function listCharacters(params?: ListCharacterParams): Promise<ListCharacterResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.page_size) query.set('page_size', String(params.page_size));

    const queryStr = query.toString();
    const url = queryStr ? `${API_BASE}/list?${queryStr}` : `${API_BASE}/list`;

    return httpClient.get<ListCharacterResponse>(url);
}

/**
 * 创建角色
 * POST /api/draftbot/create
 */
export async function createCharacter(params: CreateCharacterParams): Promise<CreateCharacterResponse> {
    return httpClient.post<CreateCharacterResponse>(`${API_BASE}/create`, params);
}
