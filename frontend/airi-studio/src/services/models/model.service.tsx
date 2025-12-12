import * as apiClient from './api-client';
import { ApiError } from '@/services/core';
import type {
    IProviderModelListArgs,
    IConnectionArgs,
} from '@/api/generated/modelapi';
import { ModelClass } from '@/api/generated/app/developer_api/ModelClass';

// Re-export types for UI layer
export { ModelClass };
export type { IProviderModelListArgs, IConnectionArgs };

// ============ 业务类型定义 ============

/** 模型列表项 - UI 层使用的简化类型 */
export interface ModelListItem {
    id: number;
    name: string;
    modelId: string;
    providerName: string;
    providerIcon: string;
    modelClass: ModelClass;
    baseUrl: string;
    type: number;
    maxTokens?: number;
    minTokens?: number;
    capabilities: {
        functionCall?: boolean;
        imageUnderstanding?: boolean;
        videoUnderstanding?: boolean;
        audioUnderstanding?: boolean;
        cotDisplay?: boolean;
        multiModal?: boolean;
        prefillResp?: boolean;
    };
}

/** 创建模型参数 */
export interface CreateModelParams {
    modelClass: ModelClass;
    modelName: string;
    baseUrl: string;
    apiKey: string;
    model: string;
    thinkingType?: number;
    enableBase64Url?: boolean;
    // OpenAI 特殊配置
    openai?: {
        byAzure?: boolean;
        apiVersion?: string;
    };
    // Gemini 特殊配置
    gemini?: {
        backend?: number;
        project?: string;
        location?: string;
    };
    // DeepSeek 配置(预留)
    deepseek?: Record<string, unknown>;
    // Qwen 配置(预留)
    qwen?: Record<string, unknown>;
    // Ollama 配置(预留)
    ollama?: Record<string, unknown>;
    // Claude 配置(预留)
    claude?: Record<string, unknown>;
}

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
        defaultBaseUrl: 'http://localhost:11434/v1',
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

/** 获取提供商配置 */
export function getProviderConfig(modelClass: ModelClass) {
    return MODEL_PROVIDERS.find(p => p.class === modelClass) || MODEL_PROVIDERS[MODEL_PROVIDERS.length - 1];
}

// ============ Service 方法 ============

/**
 * 获取模型列表
 */
export async function fetchModelList(): Promise<ModelListItem[]> {
    try {
        const resp = await apiClient.getModelList();

        if (!resp.provider_model_list) {
            return [];
        }

        // 扁平化并转换为 UI 友好的格式
        const models: ModelListItem[] = [];

        for (const providerList of resp.provider_model_list) {
            const provider = providerList.provider;
            const modelList = providerList.model_list || [];

            for (const model of modelList) {
                const providerConfig = getProviderConfig(provider?.model_class ?? ModelClass.Other);

                models.push({
                    id: Number(model.model_id) || 0,
                    name: model.name || '未知模型',
                    modelId: model.protocol_config?.model || model.name ||'',
                    providerName: provider?.name?.zh_cn || providerConfig.name,
                    providerIcon: provider?.icon_url || model.icon_url || providerConfig.icon,
                    modelClass: provider?.model_class ?? ModelClass.Other,
                    baseUrl: model.protocol_config?.base_url || '',
                    type: 0,
                    maxTokens: Number(model.ability?.max_context_tokens) || 0,
                    capabilities: {
                        functionCall: model.ability?.function_call,
                        multiModal: model.ability?.multi_modal,
                    },
                });
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
export async function createModel(params: CreateModelParams): Promise<number> {
    const connection: IConnectionArgs = {
        base_conn_info: {
            base_url: params.baseUrl,
            api_key: params.apiKey,
            model: params.model,
            thinking_type: params.thinkingType ?? 0,
        },
    };

    // 添加提供商特殊配置
    switch (params.modelClass) {
        case ModelClass.GPT:
            connection.openai = {
                by_azure: params.openai.byAzure,
                api_version: params.openai.apiVersion,
            };
            break;
        case ModelClass.Gemini:
            connection.gemini = {
                backend: params.gemini.backend,
                project: params.gemini.project,
                location: params.gemini.location,
            };
            break;
        case ModelClass.Ollama:
            connection.ollama = {};
            break;
        case ModelClass.Claude:
            connection.claude = {};
            break;
        case ModelClass.DeepSeek:
            connection.deepseek = {};
            break;
        case ModelClass.QWen:
            connection.qwen = {};
            break;
    }

    const resp = await apiClient.createModel({
        model_class: params.modelClass,
        model_name: params.modelName,
        connection,
        enable_base64_url: params.enableBase64Url,
    });

    return Number(resp.id) || 0;
}

/**
 * 删除模型
 */
export async function deleteModel(id: number): Promise<void> {
    await apiClient.deleteModel({ id });
}

/**
 * 测试模型连接
 */
export async function testConnection(params: CreateModelParams): Promise<boolean> {
    try {
        await createModel(params);
        return true;
    } catch {
        return false;
    }
}
