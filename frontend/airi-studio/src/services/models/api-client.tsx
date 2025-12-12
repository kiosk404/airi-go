
import { httpClient } from '@/services/core';

// 从 IDL 生成的类型中导入
import type {
    IGetModelListRespArgs,
    ICreateModelReqArgs,
    ICreateModelRespArgs,
    IDeleteModelReqArgs,
    IDeleteModelRespArgs,
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
 * 删除模型
 * POST /api/admin/model/delete
 */
export async function deleteModel(
    req: IDeleteModelReqArgs
): Promise<IDeleteModelRespArgs> {
    return httpClient.post<IDeleteModelRespArgs>(`${API_BASE}/delete`, req);
}
