import React, { useState, useEffect } from "react";
import { Spin, Toast } from "@douyinfe/semi-ui";
import CreateCharacterCard from "./CreateCharacterCard";
import CharacterCardSpace from "./CharacterCardList";
import { listCharacters, type IntelligenceData } from "@/services/character";

const WorkSpaceContent: React.FC = () => {
    const [characters, setCharacters] = useState<IntelligenceData[]>([]);
    const [loading, setLoading] = useState(true);

    // 从后端加载角色列表
    const loadCharacters = async () => {
        setLoading(true);
        try {
            const response = await listCharacters({ size: 50 });
            if (response.code === 0 && response.data?.intelligences) {
                setCharacters(response.data.intelligences);
            } else {
                Toast.error(response.msg || '获取角色列表失败');
            }
        } catch (error) {
            console.error('Failed to load characters:', error);
            Toast.error('获取角色列表失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadCharacters().catch(console.error);
    }, []);

    const handleStarToggle = (id: string) => {
        setCharacters(characters.map(c =>
            c.basic_info?.id === id
                ? { ...c, favorite_info: { is_fav: !c.favorite_info?.is_fav } }
                : c
        ));
    };

    return (
        <div style={{ padding: "0 20px", marginTop: 20 }}>
            <div style={{ marginBottom: 16, fontSize: 16, fontWeight: 300 }}>角色列表</div>

            {loading ? (
                <div style={{ display: 'flex', justifyContent: 'center', padding: 40 }}>
                    <Spin size="large" />
                </div>
            ) : (
                <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))", gap: 24 }}>
                    <CreateCharacterCard onSuccess={loadCharacters} />
                    {characters.map((item) => (
                        <CharacterCardSpace
                            key={String(item.basic_info?.id || '')}
                            id={String(item.basic_info?.id || '')}
                            name={item.basic_info?.name || ''}
                            avatar={item.basic_info?.icon_url || item.basic_info?.icon_uri || ''}
                            description={item.basic_info?.description || ''}
                            isStarred={item.favorite_info?.is_fav || false}
                            onStarToggle={handleStarToggle}
                            onDelete={loadCharacters}
                            updatedAt={parseInt(String(item.basic_info?.update_time)) || 0}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

export default WorkSpaceContent;
