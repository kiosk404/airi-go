// @ts-nocheck - Semi UI Form has known type recursion issues
import React, { useState, useEffect, useRef } from 'react';
import {
    Modal,
    Form,
    Button,
    Toast,
    Divider,
    Collapsible,
    Space,
    AutoComplete,
    Typography,
} from '@douyinfe/semi-ui';
import type { FormApi } from '@douyinfe/semi-ui/lib/es/form';
import {
    ModelClass,
    MODEL_PROVIDERS,
    EMBEDDING_PROVIDERS,
    RERANK_PROVIDERS,
    updateModel,
    createModel,
    type ModelListItem,
    type CreateModelParams,
} from '@/services/models';

const { Text } = Typography;

// 模型类型选项
const MODEL_TYPE_OPTIONS = [
    { value: 0, label: 'LLM', desc: 'LLM 模型' },
    { value: 1, label: 'Embedding', desc: 'Embedding 模型'},
    { value: 2, label: 'Rerank', desc: 'Rerank 模型'},
];

// 思考模式选项
const THINKING_TYPE_OPTIONS = [
    { value: 0, label: '默认' },
    { value: 1, label: '启用' },
    { value: 2, label: '禁用' },
    { value: 3, label: '自动' },
];

interface AddModelDialogProps {
    visible: boolean;
    editModel?: ModelListItem | null;
    onClose: () => void;
    onSuccess: () => void;
}

const AddModelDialog: React.FC<AddModelDialogProps> = ({
                                                           visible,
                                                           editModel,
                                                           onClose,
                                                           onSuccess,
                                                       }) => {
    const formRef = useRef<FormApi>();
    const [loading, setLoading] = useState(false);
    const [modelType, setModelType] = useState<number>(0)
    const [selectedProvider, setSelectedProvider] = useState<ModelClass>(ModelClass.GPT);
    const [advancedOpen, setAdvancedOpen] = useState(false);
    const [modelId, setModelId] = useState('');

    // 根据模型类型获取对应的提供商列表
    const getProvidersByType = (type: number) => {
        switch (type) {
            case 1: return EMBEDDING_PROVIDERS;
            case 2: return RERANK_PROVIDERS;
            default: return MODEL_PROVIDERS;
        }
    }

    const currentProviders = getProvidersByType(modelType)
    const currentProvider = currentProviders.find(p => p.class === selectedProvider)

    // 初始化/重置表单
    useEffect(() => {
        if (visible) {
            // 使用 setTimeout 确保 Form 已经渲染并且 formRef 已经设置
            setTimeout(() => {
                if (!formRef.current) {
                    console.warn('formRef.current is null');
                    return;
                }

                if (editModel) {
                    // 编辑模式
                    setModelType(editModel.type || 0);
                    setSelectedProvider(editModel.modelClass);
                    setModelId(editModel.modelId);
                    formRef.current.setValues({
                        model_name: editModel.name || '',
                        base_url: editModel.baseUrl || '',
                        api_key: '*****', // 安全考虑，不回显
                        thinking_type: editModel.thinkingType ?? 0,
                        enable_base64_url: editModel.enableBase64Url ?? false,
                        by_azure: false,
                        api_version: '',
                        gemini_backend: 1,
                        gemini_project: '',
                        gemini_location: '',
                    });
                } else {
                    // 新增模式
                    setModelType(0);
                    setSelectedProvider(ModelClass.GPT);
                    const defaultModel = MODEL_PROVIDERS[0].models[0] || '';
                    setModelId(defaultModel);
                    formRef.current.reset();
                    formRef.current.setValues({
                        base_url: MODEL_PROVIDERS[0].defaultBaseUrl,
                        thinking_type: 0,
                        enable_base64_url: false,
                        by_azure: false,
                        gemini_backend: 1,
                    });
                }
            }, 50); // 增加延迟确保 Form 完全渲染
        }
    }, [visible, editModel]);

    // 切换模型类型时更新提供商列表
    const handleModelTypeChange = (value: number) => {
        setModelType(value);
        const provider = getProvidersByType(value)[0] || MODEL_PROVIDERS[0];
        if (provider.length > 0) {
            const firstProvider = provider[0];
            setSelectedProvider(firstProvider.class);
            const defaultModel = firstProvider.models[0] || '';
            setModelId(defaultModel);
            if (formRef.current) {
                formRef.current.setValues({
                    base_url: firstProvider.defaultBaseUrl,
                });
            }
        }
    };

    // 切换提供商时更新默认 URL 和模型
    const handleProviderChange = (value: ModelClass) => {
        setSelectedProvider(value);
        const provider = currentProviders.find(p => p.class === value);
        if (provider && formRef.current) {
            const defaultModel = provider.models[0] || '';
            setModelId(defaultModel);
            formRef.current.setValues({
                base_url: provider.defaultBaseUrl,
            });
        }
    };

    // 模型 ID 变化处理
    const handleModelIdChange = (value: string) => {
        console.log('handleModelIdChange:', value);
        setModelId(value || '');
    };

    // 构建请求数据
    const buildRequest = (): CreateModelParams => {
        const values = formRef.current?.getValues() || {};

        const params: CreateModelParams = {
            modelClass: selectedProvider,
            modelName: values.model_name,
            baseUrl: values.base_url,
            apiKey: values.api_key,
            model: modelId,
            modelType: modelType,
            thinkingType: values.thinking_type ?? 0,
            enableBase64Url: values.enable_base64_url || false,
        };

        // 添加提供商特殊配置
        if (selectedProvider === ModelClass.GPT) {
            params.openai = {
                byAzure: values.by_azure || false,
                apiVersion: values.api_version || '',
            };
        } else if (selectedProvider === ModelClass.Gemini) {
            params.gemini = {
                backend: values.gemini_backend || 1,
                project: values.gemini_project || '',
                location: values.gemini_location || '',
            };
        }

        return params;
    };

    // 提交表单
    const handleSubmit = async () => {
        // 手动验证 modelId
        if (!modelId || !modelId.trim()) {
            Toast.error('请输入模型 ID');
            return;
        }

        // Semi UI Form validate() 返回 Promise<values>，验证失败会 reject
        try {
            await formRef.current?.validate(['model_name', 'base_url']);
        } catch (errors) {
            console.log('validate errors:', errors);
            if (errors && typeof errors === 'object') {
                const errorFields = Object.keys(errors);
                Toast.error(`请填写必填项: ${errorFields.join(', ')}`);
            } else {
                Toast.error('请填写必填项');
            }
            return;
        }

        setLoading(true);
        try {
            const params = buildRequest();
            console.log('createModel params:', params);
            if (editModel) {
                await updateModel({...params, id: editModel.id});
                Toast.success('模型更新成功');
            } else {
                await createModel(params);
                Toast.success('模型添加成功');
            }
            onSuccess();
            onClose();
        } catch (error) {
            console.error('createModel error:', error);
            Toast.error(`操作失败: ${error instanceof Error ? error.message : '未知错误'}`);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal
            visible={visible}
            title={editModel ? '编辑模型' : '添加模型'}
            width={640}
            onCancel={onClose}
            footer={
                <Space>
                    <Button onClick={onClose}>取消</Button>
                    <Button
                        theme="solid"
                        type="primary"
                        loading={loading}
                        onClick={handleSubmit}
                    >
                        {editModel ? '保存' : '添加'}
                    </Button>
                </Space>
            }
        >
            <Form
                key={visible ? 'open' : 'closed'}
                getFormApi={(api) => { formRef.current = api; }}
                labelPosition="top"
                style={{ padding: '0 12px' }}
                initValues={{
                    base_url: MODEL_PROVIDERS[0].defaultBaseUrl,
                    thinking_type: 0,
                    enable_base64_url: false,
                    by_azure: false,
                    gemini_backend: 1,
                }}
            >
                {/* 模型类型选择 */}
                <Form.Slot label="模型类型">
                    <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                        {MODEL_TYPE_OPTIONS.map(type => (
                            <Button
                                key={type.value}
                                theme={modelType === type.value ? 'solid' : 'light'}
                                type={modelType === type.value ? 'primary' : 'tertiary'}
                                onClick={() => handleModelTypeChange(type.value)}
                                style={{ minWidth: 120, height: 56, padding: '8px 16px'}}
                                disabled={!!editModel}
                            >
                                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                                    <span style={{ fontWeight: 600 }}>{type.label}</span>
                                    <span style={{ fontSize: 12, color: '#999', opacity: 0.7}}>{type.desc}</span>
                                </div>
                            </Button>
                        ))}
                    </div>
                </Form.Slot>
                {/* 模型提供商选择 */}
                <Form.Slot label="模型提供商">
                    <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                        {currentProviders.map(provider => (
                            <Button
                                key={provider.class}
                                theme={selectedProvider === provider.class ? 'solid' : 'light'}
                                type={selectedProvider === provider.class ? 'primary' : 'tertiary'}
                                onClick={() => handleProviderChange(provider.class)}
                                style={{ minWidth: 100 }}
                            >
                                <img
                                    src={provider.icon}
                                    alt={provider.name}
                                    style={{ width: 16, height: 16, marginRight: 6, objectFit: 'contain', verticalAlign: 'middle' }}
                                />
                                {provider.name}
                            </Button>
                        ))}
                    </div>
                </Form.Slot>

                {/* 模型名称 */}
                <Form.Input
                    field="model_name"
                    label="模型名称"
                    placeholder="输入显示名称，如 My GPT-4o"
                    rules={[{ required: true, message: '请输入模型名称' }]}
                />

                <Divider margin={16}>连接配置</Divider>

                {/* Base URL */}
                <Form.Input
                    field="base_url"
                    label="Base URL"
                    placeholder={currentProvider?.defaultBaseUrl || '输入 API Base URL'}
                    rules={[{ required: true, message: '请输入 Base URL' }]}
                />

                {/* API Key */}
                <Form.Input
                    field="api_key"
                    label="API Key"
                    mode="password"
                    placeholder="输入 API Key（可选）"
                />

                {/* 模型 ID - 使用 AutoComplete 支持输入和选择 */}
                <Form.Slot label={<span>模型 ID <Text type="danger">*</Text></span>}>
                    <AutoComplete
                        key={`model-autocomplete-${selectedProvider}`}
                        data={currentProvider?.models || []}
                        placeholder={currentProvider?.models?.length ? "输入或选择模型 ID" : "输入模型 ID"}
                        value={modelId}
                        onChange={handleModelIdChange}
                        onSelect={handleModelIdChange}
                        showClear
                        style={{ width: '100%' }}
                    />
                </Form.Slot>

                {/* 高级配置 */}
                {modelType === 0 && (
                    <>
                        <Collapsible
                            isOpen={advancedOpen}
                            collapseHeight={0}
                            keepDOM
                        >
                            <div style={{ marginTop: 16 }}>
                                {/* 思考模式 */}
                                <Form.Select
                                    field="thinking_type"
                                    label="思考模式"
                                    placeholder="选择思考模式"
                                    optionList={THINKING_TYPE_OPTIONS.map(opt => ({ label: opt.label, value: opt.value }))}
                                />

                                {/* Base64 URL */}
                                <Form.Switch
                                    field="enable_base64_url"
                                    label="启用 Base64 URL"
                                />

                                {/* OpenAI 特殊配置 */}
                                {selectedProvider === ModelClass.GPT && (
                                    <>
                                        <Divider margin={16}>OpenAI 配置</Divider>
                                        <Form.Switch
                                            field="by_azure"
                                            label="使用 Azure OpenAI"
                                        />
                                        <Form.Input
                                            field="api_version"
                                            label="API Version"
                                            placeholder="如 2024-02-15-preview"
                                        />
                                    </>
                                )}

                                {/* Gemini 特殊配置 */}
                                {selectedProvider === ModelClass.Gemini && (
                                    <>
                                        <Divider margin={16}>Gemini 配置</Divider>
                                        <Form.Select
                                            field="gemini_backend"
                                            label="Backend"
                                            optionList={[
                                                { label: 'Gemini API', value: 1 },
                                                { label: 'Vertex AI', value: 2 },
                                            ]}
                                        />
                                        <Form.Input
                                            field="gemini_project"
                                            label="Project"
                                            placeholder="GCP Project ID"
                                        />
                                        <Form.Input
                                            field="gemini_location"
                                            label="Location"
                                            placeholder="如 us-central1"
                                        />
                                    </>
                                )}
                            </div>
                        </Collapsible>

                        <Button
                            type="tertiary"
                            block
                            style={{ marginTop: 12 }}
                            onClick={() => setAdvancedOpen(!advancedOpen)}
                        >
                            {advancedOpen ? '收起高级配置' : '展开高级配置'}
                        </Button>
                    </>
                )}
            </Form>
        </Modal>
    );
};

export default AddModelDialog;
