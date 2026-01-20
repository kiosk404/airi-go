import { httpClient } from '@/services/core';

// 从 IDL 生成的类型中导入
import type {
    IGetModelListRespArgs,
    ICreateModelReqArgs,
    ICreateModelRespArgs,
    IDeleteModelReqArgs,
    IDeleteModelRespArgs,
    IUpdateModelReqArgs,
    IUpdateModelRespArgs,
    ISetDefaultModelReqArgs,
    ISetDefaultModelRespArgs
} from '@/api/generated/modelapi';

const API_BASE = '/api/admin/model';

/**
 * 获取模型列表
 * GET /api/admin/model/list
 */
export async function getModelList(): Promise<IGetModelListRespArgs> {
    return httpClient.get<IGetModelListRespArgs>(`${API_BASE}/list`, {});
}

/**
 * 创建模型
 * POST /api/admin/model/create
 */
export async function createModel(
    req: ICreateModelReqArgs
): Promise<ICreateModelRespArgs> {
    return httpClient.post<ICreateModelRespArgs>(`${API_BASE}/create`, req);
}

/**
 * 更新模型
 * POST /api/admin/model/update
 */
export async function updateModel(
    req: IUpdateModelReqArgs
): Promise<IUpdateModelRespArgs> {
    return httpClient.post<IUpdateModelRespArgs>(`${API_BASE}/update`, req);
}

/**
 * 删除模型
 * POST /api/admin/model/delete
 */
export async function deleteModel(
    req: IDeleteModelReqArgs
): Promise<IDeleteModelRespArgs> {
    return httpClient.post<IDeleteModelRespArgs>(`${API_BASE}/delete`, req);
}

/**
 * 设置默认模型
 * GET /api/admin/model/set_default
 */
export async function setDefaultModel(
    req: ISetDefaultModelReqArgs
): Promise<ISetDefaultModelRespArgs> {
    return httpClient.post<ISetDefaultModelRespArgs>(`${API_BASE}/set_default`, req);
}