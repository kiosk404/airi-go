import * as apiClient from './api-client';
import { ApiError } from '@/services/core';
import type {
    ICreateModelReqArgs,
    IUpdateModelReqArgs,
    IConnectionArgs,
    IModelArgs,
    IGetModelListRespArgs,
} from '@/api/generated/modelapi';
import { ModelClass } from '@/api/generated/app/developer_api/ModelClass';

// Re-export types for UI layer
export { ModelClass };
export type { IConnectionArgs, ICreateModelReqArgs, IUpdateModelReqArgs };

// ============ 业务类型定义 ============

/** 模型列表项 - UI 层使用的类型 (只要求 id 必需) */
export type ModelListItem = Required<Pick<IModelArgs, 'id'>> &
    Omit<IModelArgs, 'id'>;

// ============ 提供商配置 ============

export const MODEL_PROVIDERS = [
    {
        class: ModelClass.GPT,
        name: 'OpenAI',
        icon: 'https://cdn.oaistatic.com/assets/favicon-miwirzcw.ico',
        defaultBaseUrl: 'https://api.openai.com/v1',
        models: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-4', 'gpt-5'],
    },
    {
        class: ModelClass.DeepSeek,
        name: 'DeepSeek',
        icon: 'https://chat.deepseek.com/favicon.svg',
        defaultBaseUrl: 'https://api.deepseek.com',
        models: ['deepseek-v3', 'deepseek-coder', 'deepseek-r1'],
    },
    {
        class: ModelClass.QWen,
        name: 'Qwen',
        icon: 'https://assets.alicdn.com/g/qwenweb/qwen-webui-fe/0.0.191/static/favicon.png',
        defaultBaseUrl: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
        models: ['qwen-max', 'qwen-plus', 'qwen-turbo', 'qwen-long'],
    },
    {
        class: ModelClass.Gemini,
        name: 'Gemini',
        icon: 'https://www.gstatic.com/lamda/images/gemini_favicon_f069958c85030456e93de685481c559f160ea06b.png',
        defaultBaseUrl: 'https://generativelanguage.googleapis.com/v1beta',
        models: ['gemini-2.0-flash', 'gemini-2.5-pro', 'gemini-3.0-pro'],
    },
    {
        class: ModelClass.Ollama,
        name: 'Ollama',
        icon: 'https://ollama.com/public/ollama.png',
        defaultBaseUrl: 'http://localhost:11434/',
        models: ['llama3', 'deepseek-v3', 'gpt-oss', 'qwen2'],
    },
    {
        class: ModelClass.Claude,
        name: 'Claude',
        icon: 'https://www.anthropic.com/favicon.ico',
        defaultBaseUrl: 'https://api.anthropic.com/v1',
        models: ['claude-3-5-sonnet', 'claude-3-opus', 'claude-3-haiku'],
    },
    {
        class: ModelClass.Other,
        name: 'Other',
        icon: 'https://cdn-icons-png.flaticon.com/512/4616/4616734.png',
        defaultBaseUrl: '',
        models: [],
    },
] as const;


// Embedding 提供商配置
export const EMBEDDING_PROVIDERS = [
    {
        class: ModelClass.GPT,
        name: 'OpenAI',
        icon: 'https://cdn.oaistatic.com/assets/favicon-miwirzcw.ico',
        defaultBaseUrl: 'https://api.openai.com/v1',
        models: ['text-embedding-3-large', 'text-embedding-3-small', 'text-embedding-ada-002'],
    },
    {
        class: ModelClass.QWen,
        name: 'Qwen',
        icon: 'https://assets.alicdn.com/g/qwenweb/qwen-webui-fe/0.0.191/static/favicon.png',
        defaultBaseUrl: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
        models: ['text-embedding-v3', 'text-embedding-v2'],
    },
    {
        class: ModelClass.Ollama,
        name: 'Ollama',
        icon: 'https://ollama.com/public/ollama.png',
        defaultBaseUrl: 'http://localhost:11434/',
        models: ['nomic-embed-text', 'mxbai-embed-large', 'all-minilm'],
    },
    {
        class: ModelClass.Other,
        name: 'Other',
        icon: 'https://cdn-icons-png.flaticon.com/512/4616/4616734.png',
        defaultBaseUrl: '',
        models: [],
    },
] as const;

// Rerank 提供商配置
export const RERANK_PROVIDERS = [
    {
        class: ModelClass.QWen,
        name: 'Qwen',
        icon: 'https://assets.alicdn.com/g/qwenweb/qwen-webui-fe/0.0.191/static/favicon.png',
        defaultBaseUrl: 'https://dashscope.aliyuncs.com/api/v1',
        models: ['gte-rerank', 'gte-rerank-hybrid'],
    },
    {
        class: ModelClass.Other,
        name: 'Other',
        icon: 'https://cdn-icons-png.flaticon.com/512/4616/4616734.png',
        defaultBaseUrl: '',
        models: [],
    },
] as const;


/** 获取提供商配置 */
export function getProviderConfig(modelClass: ModelClass) {
    return MODEL_PROVIDERS.find(p => p.class === modelClass) || MODEL_PROVIDERS[MODEL_PROVIDERS.length - 1];
}

/**
 * 获取模型列表
 */
export async function fetchModelList(): Promise<ModelListItem[]> {
    try {
        const resp = await apiClient.getModelList() as unknown as IGetModelListRespArgs;

        if (!resp.provider_model_list) {
            return [];
        }

        // 扁平化 ProviderModelList 为 Model 列表
        const models: ModelListItem[] = [];

        for (const providerList of resp.provider_model_list) {
            const modelList = providerList.model_list || [];
            for (const model of modelList) {
                // 只要有 id 即可
                if (model.id) {
                    models.push(model as ModelListItem);
                }
            }
        }

        return models;
    } catch (error) {
        // 接口未实现时返回空数组
        if (error instanceof ApiError && error.code === 404) {
            console.warn('模型列表接口未实现');
            return [];
        }
        throw error;
    }
}

/**
 * 创建模型
 */
export async function createModel(params: ICreateModelReqArgs): Promise<string> {
    const resp = await apiClient.createModel(params);
    return String(resp.id) || '0';
}

/**
 * 更新模型
 */
export async function updateModel(params: IUpdateModelReqArgs): Promise<void> {
    await apiClient.updateModel(params);
}

/**
 * 删除模型
 */
export async function deleteModel(id: string): Promise<void> {
    await apiClient.deleteModel({ id: id });
}

/**
 * 设置默认模型
 */
export async function setDefaultModel(id: string): Promise<void> {
    await apiClient.setDefaultModel({ id: id });
}