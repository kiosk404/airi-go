import React, {useState} from "react";
import {Button, Card, Typography} from "tdesign-react";
import {AddRectangleFilledIcon, EditIcon, StarIcon} from "tdesign-icons-react";

const { Title, Paragraph } = Typography;

const mockCharacter = [
    { id: 1, name: "Airis", description: "基于大语言模型的AI管家", isStarred: true },
    { id: 2, name: "Atlas", description: "基于大语言模型的AI小笨蛋", isStarred: true }
];

// ----- WorkSpaceContent -----
type WorkSpaceContentProps = {
    activeTab: string;
};

const WorkSpaceContent: React.FC<WorkSpaceContentProps> = ({ activeTab }) => {
    const [characters, setCharacters] = useState(mockCharacter);
    const handleStar = (id: number) => {
        setCharacters(characters.map(c => c.id === id ? { ...c, isStarred: !c.isStarred } : c));
    };

    const filteredCharacters = activeTab === "starred"
        ? characters.filter(c => c.isStarred)
        : characters;

    return (
        <div style={{ padding: "0 20px", marginTop: 20 }}>
            <div style={{ marginBottom: 16, fontSize: 16, fontWeight: 600 }}>应用列表</div>
            <div style={{ fontSize: 14, color: "#666", marginBottom: 16 }}>
                共 {filteredCharacters.length} 个应用
            </div>

            <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(300px, 1fr))", gap: 16 }}>
                {/* 创建空白应用 */}
                <Card style={{ height: 220, cursor: "pointer", backgroundColor: "#f0f7ff", border: "1px solid #bae0ff" }}>
                    <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", height: "100%" }}>
                        <div style={{ width: 64, height: 64, borderRadius: "50%", backgroundColor: "#e6f7ff", display: "flex", alignItems: "center", justifyContent: "center", marginBottom: 16 }}>
                            <AddRectangleFilledIcon size={28} style={{ color: "#1890ff" }} />
                        </div>
                        <Title level="h5" style={{ margin: 4 }}>创建空白应用</Title>
                        <Paragraph style={{ margin: 0, color: "#666", textAlign: "center" }}>从零创建全新对话应用</Paragraph>
                    </div>
                </Card>

                {/* 动态应用卡片 */}
                {characters.filter(
                    character => activeTab !== 'started' || character.isStarred)
                    .map(character => (
                        <Card key={character.id}
                              style={{ height: '220px', boxShadow: '0 1px 3px rgba(0,0,0,0,1)' }}> <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                            <div style={{ display: 'flex', alignItems: 'center' }}>
                                <Title level='h5' style={{ margin: 0 }}>{character.name}</Title>
                                {character.isStarred && (<StarIcon size={16} style={{ marginLeft: '8px', color: '#faad14', fill: '#faad14' }} />)}
                                <Button
                                    variant="text"
                                    shape="circle"
                                    size="small"
                                    onClick={() => handleStar(character.id)}>
                                    {character.isStarred ? (<StarIcon size={16} style={{ color: '#faad14', fill: '#faad14' }}></StarIcon>) : (<StarIcon size={16}></StarIcon>)}
                                </Button>
                            </div>
                            <div>
                                <Paragraph style={{ margin: '8px 0', color: '#666' }}>
                                    {character.description}
                                </Paragraph>
                                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', marginBottom: '16px' }}>
                                    <div style={{ fontSize: '12px', gap: '4px' }}>
                                        <Button variant="text" size="small" icon={<EditIcon size={16} />} />
                                    </div>
                                </div>
                            </div>
                        </div>
                        </Card>
                    ))}
            </div>
        </div>
    );
};

export default WorkSpaceContent;