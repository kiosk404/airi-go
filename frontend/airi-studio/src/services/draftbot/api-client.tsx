/**
 * DraftBot API Client
 * 对应后端接口: /api/draftbot/* 和 /api/playground_api/draftbot/*
 *
 * 类型定义复用自 @/api/generated
 */

import { httpClient } from '@/services/core';

// ============ 从 IDL 生成的类型中导入 ============

// Playground - 获取 Bot 信息
import type { IGetDraftBotInfoAgwRequestArgs } from '@/api/generated/playground/GetDraftBotInfoAgwRequest';
import type { IGetDraftBotInfoAgwResponseArgs } from '@/api/generated/playground/GetDraftBotInfoAgwResponse';
import type { IGetDraftBotInfoAgwDataArgs } from '@/api/generated/playground/GetDraftBotInfoAgwData';
import type { IBotOptionDataArgs } from '@/api/generated/playground/BotOptionData';
import type { IModelDetailArgs } from '@/api/generated/playground/ModelDetail';

// Bot 通用类型
import type { IBotInfoArgs } from '@/api/generated/app/bot_common/BotInfo';
import type { IModelInfoArgs } from '@/api/generated/app/bot_common/ModelInfo';
import type { IPluginInfoArgs } from '@/api/generated/app/bot_common/PluginInfo';
import type { IWorkflowInfoArgs } from '@/api/generated/app/bot_common/WorkflowInfo';
import type { IDatabaseArgs } from '@/api/generated/app/bot_common/Database';
import type { IVariableArgs } from '@/api/generated/app/intelligence/common/Variable';

// Developer API - 显示信息
import type { IGetDraftBotDisplayInfoRequestArgs } from '@/api/generated/app/developer_api/GetDraftBotDisplayInfoRequest';
import type { IGetDraftBotDisplayInfoResponseArgs } from '@/api/generated/app/developer_api/GetDraftBotDisplayInfoResponse';
import type { IDraftBotDisplayInfoDataArgs } from '@/api/generated/app/developer_api/DraftBotDisplayInfoData';
import type { ITabDisplayItemsArgs } from '@/api/generated/app/developer_api/TabDisplayItems';

// ============ 重新导出类型（简化外部使用） ============

// 请求/响应类型
export type GetDraftBotInfoRequest = IGetDraftBotInfoAgwRequestArgs;
export type GetDraftBotInfoResponse = IGetDraftBotInfoAgwResponseArgs;
export type GetDraftBotInfoData = IGetDraftBotInfoAgwDataArgs;

export type GetDisplayInfoRequest = IGetDraftBotDisplayInfoRequestArgs;
export type GetDisplayInfoResponse = IGetDraftBotDisplayInfoResponseArgs;
export type DraftBotDisplayInfoData = IDraftBotDisplayInfoDataArgs;
export type TabDisplayItems = ITabDisplayItemsArgs;

// Bot 相关类型
export type BotInfo = IBotInfoArgs;
export type BotOptionData = IBotOptionDataArgs;
export type ModelDetail = IModelDetailArgs;
export type ModelInfo = IModelInfoArgs;
export type PluginInfo = IPluginInfoArgs;
export type WorkflowInfo = IWorkflowInfoArgs;
export type Database = IDatabaseArgs;
export type Variable = IVariableArgs;

// ============ API 端点 ============

const API_BASE = '/api/draftbot';
const PLAYGROUND_API_BASE = '/api/playground_api/draftbot';

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