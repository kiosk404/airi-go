import React from 'react';
import {Layout, Typography, Space, Button } from 'tdesign-react'; 
import ChatPanel from './ChatPanel';
import EditorPanel from './EditorPanel';
import SkillPanel from './SkillPanel';
import { AppIcon } from 'tdesign-icons-react';


const { Header } = Layout;
const { Title } = Typography;

const PlaygroundHeader: React.FC = () => {
    return (
        <Header style={{ background: '#fff', display: 'flex', justifyContent: 'space-bewtween', alignItems: 'center', padding: '0 20px', borderBottom: '1px solid #e0e0e0' }}>
            <Space>
                <AppIcon size='16px' />
                <Title level="h6" style={{margin: 0}}> Airi</Title>
            </Space>
            <Space style={{ marginLeft: 'auto' }}>
                <Button shape='round' theme='primary'>发布</Button>
            </Space>
        </Header>
    )
    
}


const PlaygroundPage: React.FC = () => {
    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden',
            height: '91vh',
            width: '100vw',
        }}>
            <PlaygroundHeader />
            <div style={{ display: 'flex', flex: 1, overflow: 'auto' }}>
                <div style={{ flex: 1, overflow: 'auto' }}><EditorPanel /></div>
                <div style={{ flex: 1, borderRight: '1px solid #e0e0e0',overflow: 'auto' }}><SkillPanel /></div>
                <div style={{ flex: 1, overflow: 'auto' }}><ChatPanel/></div>
            </div>
        </div>
    )
};

export default PlaygroundPage;
