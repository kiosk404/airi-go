import React from 'react';
import { Layout, Typography, Space, Button } from 'tdesign-react';
import ChatPanel from './ChatPanel.tsx';
import EditorPanel from './EditorPanel.tsx';
import SkillPanel from './SkillPanel.tsx';
import { Shuffle } from "@/components/base/bit"


const { Header } = Layout;
const { Title } = Typography;

const PlaygroundHeader: React.FC = () => {
    return (
        <Header style={{ background: '#fff', display: 'flex', justifyContent: 'space-bewtween', alignItems: 'center', padding: '0 20px', borderBottom: '1px solid #e0e0e0' }}>
            <Space>
                <Title style={{ marginTop: -20 , marginLeft: 10, padding: 0, margin: 0, fontSize: '16px', fontWeight: 'bold'}}>
                    <Shuffle
                        text=" Airi Here ~"
                        shuffleDirection="right"
                        duration={0.35}
                        animationMode="evenodd"
                        shuffleTimes={1}
                        ease="power3.out"
                        stagger={0.03}
                        threshold={0.1}
                        triggerOnce={true}
                        triggerOnHover={true}
                        respectReducedMotion={true}
                        style={{ fontSize: '16px', textShadow: '2px 2px 4px rgba(0,0,0,0.5)' }}
                    />
                </Title>
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
                <div style={{ flex: 1, borderRight: '1px solid #e0e0e0', overflow: 'auto' }}><SkillPanel /></div>
                <div style={{ flex: 1, overflow: 'auto' }}><ChatPanel /></div>
            </div>
        </div>
    )
};

export default PlaygroundPage;
