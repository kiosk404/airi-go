import React, { useState } from "react";
import {Card, Form, Input, Upload, Button } from 'tdesign-react';
import { PlusIcon, FolderIcon, FileIcon } from "tdesign-icons-react";
import { Modal, TextArea } from '@douyinfe/semi-ui';

type CreateCharacterModalProps = {
    dialogVisible: boolean;
    setDialogVisible: (visible: boolean) => void;
}

function CreateCharacterModal({dialogVisible, setDialogVisible} : CreateCharacterModalProps) {
    const [formData, setFormData] = useState({
        avatar: '',
        name: '',
        description: ''
    })
    const [loading, setLoading] = useState(false);

    /** 处理创建角色的表单提交 **/
    const handleSubmit = () => {
        console.log('fromData 提交的数据', formData)
        setLoading(false);
        setDialogVisible(false);
    }

   
    const handleAvatarUpload = (files: any) => {
        console.log('上传头像', files)
    }

    const handleInputChange = (name: string, value: string) => {
        console.log(name, value)
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
            <Button theme="primary" loading={loading} onClick={handleSubmit}>确定</Button>
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
                    <Form layout="vertical">
                        <div style={{ marginBottom: '24px' }}>
                            <div style={{ fontSize: '16px', marginBottom: '16px'}}> 角色名称 & 头像</div>

                            <div style={{ display: 'flex', gap: '24px',  alignItems: 'center' }}>
                                {/** 头像上传 **/}
                                <Upload action="/api/upload" accept="image/*" multiple={false} onSuccess={(files)=> handleAvatarUpload(files)}>
                                    <div style={{ width: '50px', height: '50px', border: '1px dashed #d9d9d9', cursor: 'pointer', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', overflow: 'hidden'}}>
                                        {formData.avatar ? (
                                            <img src={formData.avatar} alt="角色头像" style={{ width: '100%', height: '100%', objectFit: 'cover'}} /> ) 
                                            : (
                                                <span style={{ color: '#999'}}>🤖</span>
                                        )}
                                    </div>
                                </Upload>

                                {/** 角色名称输入 **/}
                                <div style={{ flex: 1 }}>
                                    <Input placeholder="给你的角色起个名字吧" value={formData.name} onChange={(val) => handleInputChange('name', val)}></Input>
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
                    </Form>
                </div>
           </div>

        </Modal>
    )
}

interface CreateCardProps {};

const CreateCharacterCard: React.FC<CreateCardProps> = () => {
    const [dialogVisible, setDialogVisible] = useState(false);

    const handleCreateCharacter = () => {
        setDialogVisible(true);
    }

    return (
        <Card bordered headerBordered hoverShadow shadow size="medium" theme="normal">
            {/* 标题*/}
            <div style={{ fontSize: '16px', fontWeight: '600', marginBottom: '16px', color: '#333333'}}> 创建角色</div>

            <div style={{ display: 'flex', gap: '8px', flexDirection: 'column' }}>
                {/* 创建角色 */}
                <div style={{ display: 'flex', alignItems: 'center', padding: '10px 12px', backgroundColor: '#f6f8fa', borderRadius: '6px', cursor: 'pointer', transition: 'all .3s ease' }} onClick={handleCreateCharacter}>
                    <div style={{ width: '20px', height: '20px', borderRadius: '50%', backgroundColor: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center', marginRight: '10px' }}>
                        <PlusIcon size={12} style={{ color: '#1890ff'}} />
                    </div>
                    <span style={{ fontSize: '13px', color: '#333333'}}>创建角色</span>
                </div>

                {/* 创建角色 */}
                <div style={{ display: 'flex', alignItems: 'center', padding: '10px 12px', backgroundColor: '#f6f8fa', borderRadius: '6px', cursor: 'pointer', transition: 'all .3s ease' }} onClick={handleCreateCharacter}>
                    <div style={{ width: '20px', height: '20px', borderRadius: '50%', backgroundColor: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center', marginRight: '10px' }}>
                        <FolderIcon size={12} style={{ color: '#fa8c16'}} />
                    </div>
                    <span style={{ fontSize: '13px', color: '#333333'}}>从模板创建角色</span>
                </div>

                {/* 创建角色 */}
                <div style={{ display: 'flex', alignItems: 'center', padding: '10px 12px', backgroundColor: '#f6f8fa', borderRadius: '6px', cursor: 'pointer', transition: 'all .3s ease' }} onClick={handleCreateCharacter}>
                    <div style={{ width: '20px', height: '20px', borderRadius: '50%', backgroundColor: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center', marginRight: '10px' }}>
                        <FileIcon size={12} style={{ color: '#52c41a'}} />
                    </div>
                    <span style={{ fontSize: '13px', color: '#333333'}}>导入DSL文件</span>
                </div>
            </div>

            <CreateCharacterModal dialogVisible={dialogVisible} setDialogVisible={setDialogVisible}/>
        </Card>
    )
}

export default CreateCharacterCard;

