/**
 * Character (Agent) API Client
 * 对应后端接口: /api/intelligence_api/* 和 /api/draftbot/*
 */

import { httpClient } from '@/services/core';

// ============ API 端点 ============

const API_BASE = '/api/draftbot';
const INTELLIGENCE_API_BASE = '/api/intelligence_api/search';


// ============ 从 IDL 生成的类型中导入 ============

// Intelligence 列表相关
import type { IGetDraftIntelligenceListRequestArgs } from '@/api/generated/app/intelligence/GetDraftIntelligenceListRequest';
import type { IGetDraftIntelligenceListResponseArgs } from '@/api/generated/app/intelligence/GetDraftIntelligenceListResponse';
import type { IDraftIntelligenceListDataArgs } from '@/api/generated/app/intelligence/DraftIntelligenceListData';
import type { IIntelligenceDataArgs } from '@/api/generated/app/intelligence/IntelligenceData';
import type { IIntelligenceBasicInfoArgs } from '@/api/generated/app/intelligence/common/IntelligenceBasicInfo';

// Bot 创建/删除相关
import type { IDraftBotCreateRequestArgs } from '@/api/generated/app/developer_api/DraftBotCreateRequest';
import type { IDraftBotCreateResponseArgs } from '@/api/generated/app/developer_api/DraftBotCreateResponse';
import type { IDeleteDraftBotRequestArgs } from '@/api/generated/app/developer_api/DeleteDraftBotRequest';
import type { IDeleteDraftBotResponseArgs } from '@/api/generated/app/developer_api/DeleteDraftBotResponse';

// ============ 重新导出类型（简化外部使用） ============

export type GetDraftIntelligenceListRequest = IGetDraftIntelligenceListRequestArgs;
export type GetDraftIntelligenceListResponse = IGetDraftIntelligenceListResponseArgs;
export type DraftIntelligenceListData = IDraftIntelligenceListDataArgs;
export type IntelligenceData = IIntelligenceDataArgs;
export type IntelligenceBasicInfo = IIntelligenceBasicInfoArgs;

export type DraftBotCreateRequest = IDraftBotCreateRequestArgs;
export type DraftBotCreateResponse = IDraftBotCreateResponseArgs;
export type DeleteDraftBotRequest = IDeleteDraftBotRequestArgs;
export type DeleteDraftBotResponse = IDeleteDraftBotResponseArgs;

// ============ API 方法 ============

/**
 * 获取角色列表
 * POST /api/intelligence_api/search/get_draft_intelligence_list
 */
export async function listCharacters(params: GetDraftIntelligenceListRequest): Promise<GetDraftIntelligenceListResponse> {
    return httpClient.post<GetDraftIntelligenceListResponse>(`${INTELLIGENCE_API_BASE}/get_draft_intelligence_list`, {
        name: params.name,
        size: params.size || 50,
        cursor_id: params.cursor_id,
        types: params.types || [1], // IntelligenceType.Bot = 1
    });
}

/**
 * 创建角色
 * POST /api/draftbot/create
 */
export async function createCharacter(params: DraftBotCreateRequest): Promise<DraftBotCreateResponse> {
    return httpClient.post<DraftBotCreateResponse>(`${API_BASE}/create`, params);
}

/**
 * 删除角色
 * POST /api/draftbot/delete
 */
export async function deleteCharacter(params: DeleteDraftBotRequest): Promise<DeleteDraftBotResponse> {
    return httpClient.post<DeleteDraftBotResponse>(`${API_BASE}/delete`, params);
}
