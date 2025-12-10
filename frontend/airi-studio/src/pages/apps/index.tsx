import React, { useState } from "react";
import { LightRays } from "@/components/base/bit";
import {Card, Button, Radio, Typography,} from "tdesign-react";
import { Gamepad1Icon } from "tdesign-icons-react";

const { Title, Paragraph } = Typography;


// ----- NavSpaceHeader -----
type NavSpaceHeaderProps = {
    activeTab: string;
    setActiveTab: React.Dispatch<React.SetStateAction<string>>;
};

const NavSpaceHeader: React.FC<NavSpaceHeaderProps> = ({ activeTab, setActiveTab }) => {
    return (
        <div style={{padding: "10px 30px",
            display: "flex",
                justifyContent: "center",
        alignItems: "center" }}>
            <Radio.Group
                variant="default-filled"
                value={activeTab}
                onChange={(val) => setActiveTab(val as string)}
            >
                <Radio.Button value="airis"><Gamepad1Icon style={{ marginRight: 6 }} />
                    <span className="font-mono">Airis</span>
                </Radio.Button>
                <Radio.Button value="playground"><Gamepad1Icon style={{ marginRight: 6 }} />
                    <span className="font-mono">Playground</span>
                </Radio.Button>
                <Radio.Button value="chatflow"><Gamepad1Icon style={{ marginRight: 6 }} />
                    <span className="font-mono">ChatFlow</span>
                </Radio.Button>
            </Radio.Group>
        </div>
    );
};


// ----- ChatFlowPage (简单的临时页面) -----
const ChatFlow: React.FC = () => {
    return (
        <div style={{ padding: "20px 30px", textAlign: "center" }}>
            <Title style={{ marginBottom: 16 }}>
                💬 ChatFlow
            </Title>
            <Paragraph style={{ fontSize: 16, color: "#666", marginBottom: 20 }}>
                这里是 ChatFlow 页面，用于管理对话流程。
            </Paragraph>
            <Card style={{ maxWidth: 600, margin: "0 auto", padding: 20 }}>
                <Title style={{ marginBottom: 16 }}>流程管理</Title>
                <div style={{ display: "flex", gap: 12, justifyContent: "center", flexWrap: "wrap" }}>
                    <Button theme="primary">创建流程</Button>
                    <Button theme="default">编辑流程</Button>
                    <Button theme="default">流程模板</Button>
                    <Button theme="default">导入导出</Button>
                </div>
                <Paragraph style={{ marginTop: 20, color: "#888" }}>
                    设计和管理复杂的对话流程
                </Paragraph>
            </Card>
        </div>
    );
};

// ----- PlaygroundPage (简单的临时页面) -----
const Playground: React.FC = () => {
    return (
        <div style={{ padding: "20px 30px", textAlign: "center" }}>
            <Title style={{ marginBottom: 16 }}>
                🤖 Palygourd
            </Title>
            <Paragraph style={{ fontSize: 16, color: "#666", marginBottom: 20 }}>
                这里是 Palygourd 页面，用于管理对话流程。
            </Paragraph>
            <Card style={{ maxWidth: 600, margin: "0 auto", padding: 20 }}>
                <Title style={{ marginBottom: 16 }}>流程管理</Title>
                <div style={{ display: "flex", gap: 12, justifyContent: "center", flexWrap: "wrap" }}>
                    <Button theme="primary">创建流程</Button>
                    <Button theme="default">编辑流程</Button>
                    <Button theme="default">流程模板</Button>
                    <Button theme="default">导入导出</Button>
                </div>
                <Paragraph style={{ marginTop: 20, color: "#888" }}>
                    设计和管理复杂的对话流程
                </Paragraph>
            </Card>
        </div>
    );
};


// ----- WorkSpaceContent (简单的临时页面) -----
const WorkSpace: React.FC = () => {
    return (
        <div style={{ padding: "20px 30px", textAlign: "center" }}>
            <Title style={{ marginBottom: 16 }}>
                🛸 WorkSpace
            </Title>
            <Paragraph style={{ fontSize: 16, color: "#666", marginBottom: 20 }}>
                这里是 WorkSpace 页面，用于管理对话流程。
            </Paragraph>
            <Card style={{ maxWidth: 600, margin: "0 auto", padding: 20 }}>
                <Title style={{ marginBottom: 16 }}>流程管理</Title>
                <div style={{ display: "flex", gap: 12, justifyContent: "center", flexWrap: "wrap" }}>
                    <Button theme="primary">创建流程</Button>
                    <Button theme="default">编辑流程</Button>
                    <Button theme="default">流程模板</Button>
                    <Button theme="default">导入导出</Button>
                </div>
                <Paragraph style={{ marginTop: 20, color: "#888" }}>
                    设计和管理复杂的对话流程
                </Paragraph>
            </Card>
        </div>
    );
};


// ----- 主页面内容渲染函数 -----
const renderPageContent = (activeTab: string) => {
    switch (activeTab) {
        case "airis":
            return <WorkSpace />;
        case "playground":
            return <Playground />;
        case "chatflow":
            return <ChatFlow />;
        default:
            return <WorkSpace />;
    }
};


// ----- AppPage -----
const AppPage: React.FC = () => {
    const [activeTab, setActiveTab] = useState("airis");

    return (
        <div style={{ width: "100%", height: "100%", position: "relative", overflow: "hidden", flex: 1 }}>
            {/* 背景 */}
            <div style={{ position: "absolute", top: 0, left: 0, width: "100%", height: "100%", zIndex: 0 }}>
                <LightRays
                    raysOrigin="bottom-center"
                    raysColor="#f5f5f5"
                    raysSpeed={1.5}
                    lightSpread={0.8}
                    rayLength={1.2}
                    followMouse={true}
                    mouseInfluence={0.1}
                    noiseAmount={0.1}
                    distortion={0.05}
                    className="custom-rays"
                />
            </div>

            {/* 内容 */}
            <div style={{ position: "relative", zIndex: 1, display: "flex", flexDirection: "column", height: "100%" }}>
                <NavSpaceHeader activeTab={activeTab} setActiveTab={setActiveTab} />
                <div style={{ flex: 1, overflowY: "auto" }}>
                    {renderPageContent(activeTab)}
                </div>
            </div>
        </div>
    );
};

export default AppPage;
