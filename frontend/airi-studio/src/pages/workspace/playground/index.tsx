import React, { useEffect, useState } from 'react';
import { Layout, Typography, Space, Button, MessagePlugin } from 'tdesign-react';
import { ChevronLeftIcon } from 'tdesign-icons-react';
import ChatPanel from './ChatPanel.tsx';
import EditorPanel from './EditorPanel.tsx';
import SkillPanel from './SkillPanel.tsx';
import { Shuffle } from "@/components/base/bit"
import { useNavigate, useParams } from 'react-router-dom';
import * as draftbotApi from '@/services/draftbot';
import type { BotInfo, DraftBotDisplayInfoData, BotOptionData } from '@/services/draftbot';


const { Header } = Layout;
const { Title } = Typography;

interface PlaygroundHeaderProps {
    botName?: string;
}

const PlaygroundHeader: React.FC<PlaygroundHeaderProps> = ({ botName }) => {

    const navigate = useNavigate();

    // 处理返回按钮点击事件
    const handleBackButtonClick = () => {
        navigate(-1); // 跳转到工作台页面
    }

    return (
        <Header style={{ background: '#fff', display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0 20px', height: '50px', borderBottom: '1px solid #e0e0e0' }}>
            <Space>
                <Title style={{ padding: 0, margin: 0, display: 'flex', alignItems: 'center'}}>
                    <Button shape='square' variant='text' size='large' onClick={handleBackButtonClick} icon={<ChevronLeftIcon />}/>
                    <Shuffle
                        text={botName ? ` ${botName} - Playground` : " Playgroud ~"}
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
                        style={{ fontSize: '16px', alignItems: 'center', marginTop: '4px', textShadow: '2px 2px 4px rgba(0,0,0,0.5)' }}
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
    const params = useParams<{ id: string }>();
    const botId = params.id || '';

    const [botInfo, setBotInfo] = useState<BotInfo | null>(null);
    const [displayInfo, setDisplayInfo] = useState<DraftBotDisplayInfoData | null>(null);
    const [botOptionData, setBotOptionData] = useState<BotOptionData | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!botId) {
            setLoading(false);
            return;
        }

        const fetchBotData = async () => {
            try {
                setLoading(true);
                // 并行请求两个 API
                const [displayInfoResp, botInfoResp] = await Promise.all([
                    draftbotApi.getDisplayInfo({ bot_id: botId }),
                    draftbotApi.getDraftBotInfo({ bot_id: botId }),
                ]);

                if (displayInfoResp.data) {
                    setDisplayInfo(displayInfoResp.data);
                }

                if (botInfoResp.data?.bot_info) {
                    setBotInfo(botInfoResp.data.bot_info);
                }

                if (botInfoResp.data?.bot_option_data) {
                    setBotOptionData(botInfoResp.data.bot_option_data);
                }
            } catch (error) {
                console.error('Failed to fetch bot data:', error);
                void MessagePlugin.error('获取 Bot 信息失败');
            } finally {
                setLoading(false);
            }
        };

        void fetchBotData();
    }, [botId]);

    if (loading) {
        return (
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                加载中...
            </div>
        );
    }

    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden',
            height: '100vh',
            width: '100vw',
        }}>
            <PlaygroundHeader botName={botInfo?.name} />
            <div style={{ display: 'flex', flex: 1, overflow: 'auto' }}>
                <div style={{ flex: 1, overflow: 'auto' }}>
                    <EditorPanel botInfo={botInfo} />
                </div>
                <div style={{ flex: 1, borderRight: '1px solid #e0e0e0', overflow: 'auto' }}>
                    <SkillPanel botInfo={botInfo} displayInfo={displayInfo} botOptionData={botOptionData} />
                </div>
                <div style={{ flex: 1, overflow: 'auto' }}>
                    <ChatPanel botInfo={botInfo ?? undefined} />
                </div>
            </div>
        </div>
    )
};

export default PlaygroundPage;
