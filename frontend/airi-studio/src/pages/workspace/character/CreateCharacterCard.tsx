import React, { useState, useRef } from "react";
import { Card, Input, Button, Modal, TextArea, Toast } from '@douyinfe/semi-ui';
import { IconPlus, IconFolder } from "@douyinfe/semi-icons";
import { createCharacter } from "@/services/character";

type CreateCharacterModalProps = {
    dialogVisible: boolean;
    setDialogVisible: (visible: boolean) => void;
    onSuccess?: () => void;
}

function CreateCharacterModal({dialogVisible, setDialogVisible, onSuccess} : CreateCharacterModalProps) {
    const [formData, setFormData] = useState({
        avatar: '',
        name: '',
        description: ''
    })
    const [loading, setLoading] = useState(false);
    const [uploading, setUploading] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    /** 处理创建角色的表单提交 **/
    const handleSubmit = async () => {
        if (!formData.name.trim()) {
            Toast.error('请输入角色名称');
            return;
        }

        setLoading(true);
        try {
            const response = await createCharacter({
                name: formData.name,
                description: formData.description,
                icon_uri: formData.avatar || 'default_icon/default_agent_icon.png', // 使用默认图标
            });

            if (response.code === 0) {
                Toast.success('创建成功');
                setDialogVisible(false);
                setFormData({ avatar: '', name: '', description: '' });
                onSuccess?.();
            } else {
                Toast.error(response.msg || '创建失败');
            }
        } catch (error) {
            console.error('Create character failed:', error);
            Toast.error('创建失败');
        } finally {
            setLoading(false);
        }
    }


    const handleAvatarUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        setUploading(true);
        const reader = new FileReader();

        reader.onload = async () => {
            try {
                const base64Data = (reader.result as string).split(',')[1];
                const fileType = file.name.split('.').pop() || 'png';

                const response = await fetch('/api/bot/upload_file', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        file_head: {
                            file_type: fileType,
                            biz_type: 1,
                        },
                        data: base64Data,
                    }),
                });

                const result = await response.json();
                if (result.code === 0 && result.url) {
                    setFormData(prev => ({ ...prev, avatar: result.url }));
                    Toast.success('头像上传成功');
                } else {
                    Toast.error(result.msg || '上传失败');
                }
            } catch (error) {
                console.error('Upload failed:', error);
                Toast.error('上传失败');
            } finally {
                setUploading(false);
            }
        };

        reader.onerror = () => {
            Toast.error('文件读取失败');
            setUploading(false);
        };

        reader.readAsDataURL(file);
    }

    const handleAvatarClick = () => {
        fileInputRef.current?.click();
    }

    const handleInputChange = (name: string, value: string) => {
        setFormData(prev => ({ ...prev, [name]: value }));
    }

    const footer = (
        <div
            style={{
                display: 'flex',
                justifyContent: 'flex-end',
                padding: '0 16px 16px 0',
                borderTop: '2px solid #f0f0f0',
                gap: '24px'
            }}>
            <Button onClick={() => setDialogVisible(false)}>取消</Button>
            <Button theme="solid" type="primary" loading={loading} onClick={handleSubmit}>确定</Button>
        </div>
    );

    return (
        <Modal
            title="创建角色"
            visible={dialogVisible}
            onCancel={() => setDialogVisible(false)}
            footer={footer}
            confirmLoading={loading}
            maskClosable={false}
            width="500px"
        >
            <div style={{ padding: 0 }}>
                <div style={{ padding: '32px' }}>
                    <div>
                        <div style={{ marginBottom: '24px' }}>
                            <div style={{ fontSize: '16px', marginBottom: '16px'}}> 角色名称 & 头像</div>

                            <div style={{ display: 'flex', gap: '24px',  alignItems: 'center' }}>
                                {/** 头像上传 **/}
                                <input
                                    type="file"
                                    ref={fileInputRef}
                                    accept="image/*"
                                    style={{ display: 'none' }}
                                    onChange={handleAvatarUpload}
                                />
                                <div
                                    style={{
                                        width: '50px',
                                        height: '50px',
                                        border: '1px dashed #d9d9d9',
                                        cursor: uploading ? 'wait' : 'pointer',
                                        borderRadius: '50%',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        overflow: 'hidden',
                                        opacity: uploading ? 0.6 : 1
                                    }}
                                    onClick={handleAvatarClick}
                                >
                                    {formData.avatar ? (
                                        <img src={formData.avatar} alt="角色头像" style={{ width: '100%', height: '100%', objectFit: 'cover'}} />
                                    ) : (
                                        <span style={{ color: '#999'}}>{uploading ? '...' : '🤖'}</span>
                                    )}
                                </div>

                                {/** 角色名称输入 **/}
                                <div style={{ flex: 1 }}>
                                    <Input placeholder="给你的角色起个名字吧" value={formData.name} onChange={(val) => handleInputChange('name', val)} />
                                </div>
                            </div>
                        </div>

                        <div style={{ marginBottom: '24px' }}>
                            <div style={{ fontSize: '16px', marginBottom: '16px'}}> 角色描述</div>

                            {/** 描述输入 **/}
                            <TextArea
                                placeholder="请输入角色的描述"
                                value={formData.description}
                                onChange={(val) => handleInputChange('description', val)}
                                maxCount={100}
                            />
                        </div>
                    </div>
                </div>
            </div>

        </Modal>
    )
}

interface CreateCardProps {
    onSuccess?: () => void;
}

const CreateCharacterCard: React.FC<CreateCardProps> = ({ onSuccess }) => {
    const [dialogVisible, setDialogVisible] = useState(false);

    const handleCreateCharacter = () => {
        setDialogVisible(true);
    }

    return (
        <Card
            bordered
            style={{ padding: '16px' }}
        >
            {/* 标题*/}
            <div style={{ fontSize: '16px', fontWeight: '600', marginBottom: '16px', color: '#333333'}}> 创建角色</div>

            <div style={{ display: 'flex', gap: '8px', flexDirection: 'column' }}>
                {/* 创建角色 */}
                <div style={{ display: 'flex', alignItems: 'center', padding: '10px 12px', backgroundColor: '#f6f8fa', borderRadius: '6px', cursor: 'pointer', transition: 'all .3s ease' }} onClick={handleCreateCharacter}>
                    <div style={{ width: '20px', height: '20px', borderRadius: '50%', backgroundColor: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center', marginRight: '10px' }}>
                        <IconPlus size="small" style={{ color: '#1890ff'}} />
                    </div>
                    <span style={{ fontSize: '13px', color: '#333333'}}>创建角色</span>
                </div>

                {/* 从模板创建角色 */}
                <div style={{ display: 'flex', alignItems: 'center', padding: '10px 12px', backgroundColor: '#f6f8fa', borderRadius: '6px', cursor: 'pointer', transition: 'all .3s ease' }} onClick={handleCreateCharacter}>
                    <div style={{ width: '20px', height: '20px', borderRadius: '50%', backgroundColor: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center', marginRight: '10px' }}>
                        <IconFolder size="small" style={{ color: '#fa8c16'}} />
                    </div>
                    <span style={{ fontSize: '13px', color: '#333333'}}>从模板创建角色</span>
                </div>
            </div>

            <CreateCharacterModal dialogVisible={dialogVisible} setDialogVisible={setDialogVisible} onSuccess={onSuccess}/>
        </Card>
    )
}

export default CreateCharacterCard;
