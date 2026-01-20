import * as apiClient from './api-client';
import { ApiError } from '@/services/core';
import type {
    IProviderModelListArgs,
    IConnectionArgs,
    IModelProviderArgs,
    IDisplayInfoArgs,
    IBaseConnectionInfoArgs,
} from '@/api/generated/modelapi';
import type { IModelAbilityArgs } from '@/api/generated/app/developer_api';
import { ModelClass } from '@/api/generated/app/developer_api/ModelClass';

// Re-export types for UI layer
export { ModelClass };
export type { IProviderModelListArgs, IConnectionArgs };

// ============ 后端返回的 Model 类型定义 (来自 model_api.thrift) ============

/** 后端返回的 Model 结构 (modelapi.Model) */
interface IModelApiModel {
    id?: number | string;
    provider?: IModelProviderArgs;
    display_info?: IDisplayInfoArgs;
    capability?: IModelAbilityArgs;
    connection?: {
        base_conn_info?: IBaseConnectionInfoArgs;
        openai?: unknown;
        deepseek?: unknown;
        gemini?: unknown;
        qwen?: unknown;
        ollama?: unknown;
        claude?: unknown;
    };
    type?: number;
    parameters?: unknown[];
    status?: number;
    enable_base64_url?: boolean;
    delete_at_ms?: number | string;
}

/** 后端返回的 ProviderModelList 结构 */
interface IModelApiProviderModelList {
    provider?: IModelProviderArgs;
    model_list?: IModelApiModel[];
}

/** 后端返回的 GetModelListResp 结构 */
interface IModelApiGetModelListResp {
    provider_model_list?: IModelApiProviderModelList[];
    code?: number;
    msg?: string;
}

// ============ 业务类型定义 ============

/** 模型列表项 - UI 层使用的简化类型 */
export interface ModelListItem {
    id: string; 
    name: string;
    modelId: string;
    providerName: string;
    providerIcon: string;
    modelClass: ModelClass;
    baseUrl: string;
    type: number;
    maxTokens?: number;
    minTokens?: number;
    thinkingType?: number;
    enableBase64Url?: boolean;
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

// ============ Service 方法 ============

/**
 * 获取模型列表
 * 后端返回的 Model 结构来自 modelapi (model_api.thrift):
 * - id, provider, display_info, capability, connection, type, parameters, status, enable_base64_url
 */
export async function fetchModelList(): Promise<ModelListItem[]> {
    try {
        const resp = await apiClient.getModelList() as unknown as IModelApiGetModelListResp;

        if (!resp.provider_model_list) {
            return [];
        }

        // 扁平化并转换为 UI 友好的格式
        const models: ModelListItem[] = [];

        for (const providerList of resp.provider_model_list) {
            const provider = providerList.provider;
            const modelList = providerList.model_list || [];

            console.log('provider:', provider);

            for (const model of modelList) {
                // 优先使用 model 自带的 provider，其次使用 providerList 的 provider
                const modelProvider = model.provider || provider;
                const providerConfig = getProviderConfig(modelProvider?.model_class ?? ModelClass.Other);

                models.push({
                    id: String(model.id) || '0',
                    name: model.display_info?.name || model.connection?.base_conn_info?.model || '未知模型',
                    modelId: model.connection?.base_conn_info?.model || model.display_info?.name || '',
                    providerName: modelProvider?.name?.zh_cn || providerConfig.name,
                    providerIcon: modelProvider?.icon_url || providerConfig.icon,
                    modelClass: modelProvider?.model_class ?? ModelClass.Other,
                    baseUrl: model.connection?.base_conn_info?.base_url || '',
                    type: Number(model.type) || 0,
                    maxTokens: Number(model.display_info?.max_tokens) || 0,
                    thinkingType: model.connection?.base_conn_info?.thinking_type ?? 0,
                    enableBase64Url: model.enable_base64_url ?? false,
                    capabilities: {
                        functionCall: model.capability?.function_call,
                        imageUnderstanding: model.capability?.image_understanding,
                        videoUnderstanding: model.capability?.video_understanding,
                        audioUnderstanding: model.capability?.audio_understanding,
                        cotDisplay: model.capability?.cot_display,
                        multiModal: model.capability?.support_multi_modal,
                        prefillResp: model.capability?.prefill_resp,
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
export async function createModel(params: CreateModelParams): Promise<string> {
    const connection: IConnectionArgs = {
        base_conn_info: {
            base_url: params.baseUrl,
            api_key: params.apiKey,
            model: params.model,
            thinking_type: params.thinkingType ?? 0,
        },
    };

    // 添加提供商特殊配置
    if (params.modelClass === ModelClass.GPT && params.openai) {
        connection.openai = {
            by_azure: params.openai.byAzure,
            api_version: params.openai.apiVersion,
        };
    } else if (params.modelClass === ModelClass.Gemini && params.gemini) {
        connection.gemini = {
            backend: params.gemini.backend,
            project: params.gemini.project,
            location: params.gemini.location,
        };
    } else if (params.modelClass === ModelClass.Ollama && params.ollama) {
        connection.ollama = {};
    } else if (params.modelClass === ModelClass.Claude && params.claude) {
        connection.claude = {};
    } else if (params.modelClass === ModelClass.DeepSeek && params.deepseek) {
        connection.deepseek = {};
    } else if (params.modelClass === ModelClass.QWen && params.qwen) {
        connection.qwen = {};
    }

    const resp = await apiClient.createModel({
        model_class: params.modelClass,
        model_name: params.modelName,
        connection,
        enable_base64_url: params.enableBase64Url,
    });

    return String(resp.id) || '0';
}

export interface UpdateModelParams extends CreateModelParams {
    id: string; // 使用字符串避免大整数精度丢失
}

/**
 * 更新模型
 */
export async function updateModel(params: UpdateModelParams): Promise<void> {
    const connection: IConnectionArgs = {
        base_conn_info: {
            base_url: params.baseUrl,
            api_key: params.apiKey,
            model: params.model,
            thinking_type: params.thinkingType ?? 0,
        },
    };

    // 添加提供商特殊配置
    if (params.modelClass === ModelClass.GPT && params.openai) {
        connection.openai = {
            by_azure: params.openai.byAzure,
            api_version: params.openai.apiVersion,
        };
    } else if (params.modelClass === ModelClass.Gemini && params.gemini) {
        connection.gemini = {
            backend: params.gemini.backend,
            project: params.gemini.project,
            location: params.gemini.location,
        };
    }

    await apiClient.updateModel({
        model: {
            id: params.id,
            model_name: params.modelName,
            connection,
            enable_base64_url: params.enableBase64Url
        } as any,
    });
}

/**
 * 删除模型
 */
export async function deleteModel(id: string): Promise<void> {
    await apiClient.deleteModel({ id: id });
}