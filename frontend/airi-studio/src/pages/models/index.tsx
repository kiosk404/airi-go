import React, { useState, useEffect, useCallback, useRef } from 'react';
import {
    Button,
    Input,
    RadioGroup,
    Radio,
    Toast,
    Spin,
    Empty,
    Row,
    Col,
    Typography,
} from '@douyinfe/semi-ui';
import ModelCard from './ModelCard';
import AddModelDialog from './AddModelDialog';
import {
    fetchModelList,
    deleteModel,
    setDefaultModel,
    type ModelListItem,
} from '@/services/models';

const { Title, Text } = Typography;

// 模型类型选项
const MODEL_TYPE_OPTIONS = [
    { value: 0, label: 'LLM' },
    { value: 1, label: 'Embedding' },
    { value: 2, label: 'Rerank' },
];

const ModelsPage: React.FC = () => {
    const [models, setModels] = useState<ModelListItem[]>([]);
    const [loading, setLoading] = useState(false);
    const [searchKeyword, setSearchKeyword] = useState('');
    const [filterType, setFilterType] = useState<number | 'all'>('all');
    const [dialogVisible, setDialogVisible] = useState(false);
    const [editingModel, setEditingModel] = useState<ModelListItem | null>(null);
    const mountedRef = useRef(false);

    // 加载模型列表
    const loadModels = useCallback(async () => {
        setLoading(true);
        try {
            const data = await fetchModelList();
            setModels(data);
        } catch (error) {
            Toast.error(`加载失败: ${error instanceof Error ? error.message : '未知错误'}`);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        if (!mountedRef.current) {
            mountedRef.current = true;
            loadModels().then(r => console.log('loadModels ->', r))
        }
    }, [loadModels]);

    // 过滤模型
    const filteredModels = models.filter(model => {
        console.log('model', model);
        // 类型过滤
        if (filterType !== 'all' && model.type !== filterType) {
            return false;
        }
        // 关键词搜索
        if (searchKeyword) {
            const keyword = searchKeyword.toLowerCase();
            const name = (model.display_info?.name || '').toLowerCase();
            const modelId = (model.connection?.base_conn_info?.model || '').toLowerCase();
            const providerName = (model.provider?.name?.zh_cn || '').toLowerCase();
            return name.includes(keyword) || modelId.includes(keyword) || providerName.includes(keyword);
        }
        return true;
    });

    // 打开添加弹窗
    const handleAdd = () => {
        setEditingModel(null);
        setDialogVisible(true);
    };

    // 打开编辑弹窗
    const handleEdit = (model: ModelListItem) => {
        setEditingModel(model);
        setDialogVisible(true);
    };

    // 删除模型
    const handleDelete = async (id: string) => {
        try {
            await deleteModel(id);
            Toast.success('删除成功');
            await loadModels();
        } catch (error) {
            Toast.error(`删除失败: ${error instanceof Error ? error.message : '未知错误'}`);
        }
    };

    // 模型选择
    const handleModelSelect = async (id: string) => {
        console.log('handleModelSelect', id);
        try{
            await setDefaultModel(id);
            Toast.success('设为默认成功');
            await loadModels();
        } catch (error) {
            Toast.error(`设为默认失败: ${error instanceof Error ? error.message : '未知错误'}`);
        }
    };

    // 弹窗关闭
    const handleDialogClose = () => {
        setDialogVisible(false);
        setEditingModel(null);
    };

    // 操作成功后刷新
    const handleSuccess = () => {
        loadModels().then(r => console.log('loadModels ->', r))
    };

    return (
        <div style={{ padding: '24px', minHeight: '100vh', background: '#f5f7fa' }}>
            {/* 页面标题和操作栏 */}
            <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: 24
            }}>
                <div>
                    <Title heading={3} style={{ margin: 0 }}>大模型管理</Title>
                    <Text type="tertiary" style={{ marginTop: 8 }}>
                        管理和配置 AI 大语言模型
                    </Text>
                </div>
                <Button theme="solid" type="primary" onClick={handleAdd}>
                    + 添加模型
                </Button>
            </div>

            {/* 筛选栏 */}
            <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: 20,
                padding: '16px 20px',
                background: '#fff',
                borderRadius: 8,
                boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
            }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                    {/* 类型筛选 */}
                    <RadioGroup
                        type="button"
                        buttonSize="middle"
                        value={filterType}
                        onChange={(e) => setFilterType(e.target.value as number | 'all')}
                    >
                        <Radio value="all">全部</Radio>
                        {MODEL_TYPE_OPTIONS.map(opt => (
                            <Radio key={opt.value} value={opt.value}>
                                {opt.label}
                            </Radio>
                        ))}
                    </RadioGroup>
                </div>

                <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                    {/* 搜索框 */}
                    <Input
                        prefix="🔍"
                        placeholder="搜索模型名称、ID..."
                        value={searchKeyword}
                        onChange={(val) => setSearchKeyword(val)}
                        style={{ width: 240 }}
                        showClear
                    />
                    {/* 刷新按钮 */}
                    <Button
                        type="tertiary"
                        onClick={loadModels}
                        loading={loading}
                    >
                        刷新
                    </Button>
                </div>
            </div>

            {/* 模型列表 */}
            <Spin spinning={loading} tip="加载中...">
                {filteredModels.length > 0 ? (
                    <Row gutter={[16, 16]}>
                        {filteredModels.map(model => (
                            <Col key={model.id} xs={24} sm={24} md={12} lg={8} xl={6}>
                                <ModelCard
                                    model={model}
                                    onEdit={handleEdit}
                                    onDelete={handleDelete}
                                    onSelect={handleModelSelect}
                                />
                            </Col>
                        ))}
                    </Row>
                ) : (
                    <div style={{
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'center',
                        alignItems: 'center',
                        minHeight: 400,
                        background: '#fff',
                        borderRadius: 8,
                        padding: '60px 20px',
                    }}>
                        <Empty
                            description={
                                searchKeyword || filterType !== 'all'
                                    ? '没有找到匹配的模型'
                                    : '暂无模型，点击上方按钮添加'
                            }
                        />
                        {!searchKeyword && filterType === 'all' && (
                            <Button
                                theme="solid"
                                type="primary"
                                onClick={handleAdd}
                                style={{ marginTop: 16 }}
                            >
                                添加第一个模型
                            </Button>
                        )}
                    </div>
                )}
            </Spin>

            {/* 添加/编辑弹窗 */}
            <AddModelDialog
                visible={dialogVisible}
                editModel={editingModel}
                onClose={handleDialogClose}
                onSuccess={handleSuccess}
            />
        </div>
    );
};

export default ModelsPage;
