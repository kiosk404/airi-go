import React, { useState, useEffect } from "react";
import { Spin, Toast } from "@douyinfe/semi-ui";
import CreateCharacterCard from "./CreateCharacterCard";
import CharacterCardSpace from "./CharacterCardList";
import { listCharacters, type CharacterItem } from "@/services/character";

interface Character {
    id: number;
    name: string;
    description: string;
    isStarred: boolean;
    updateAt: number;
}

const WorkSpaceContent: React.FC = () => {
    const [characters, setCharacters] = useState<Character[]>([]);
    const [loading, setLoading] = useState(true);

    // 从后端加载角色列表
    const loadCharacters = async () => {
        setLoading(true);
        try {
            const response = await listCharacters({ page: 1, page_size: 50 });
            if (response.code === 0 && response.data?.items !== undefined) {
                const list = response.data.items.map((item: CharacterItem) => ({
                    id: Number(item.id),
                    name: item.name,
                    description: item.description || '',
                    isStarred: false, // TODO: 后续从后端获取收藏状态
                    updateAt: item.updated_at,
                }));
                setCharacters(list);
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

    const handleStarToggle = (id: number) => {
        setCharacters(characters.map(c => c.id === id ? {...c, isStarred: !c.isStarred} : c));
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

                    {characters.map((character) => (
                        <CharacterCardSpace
                            key={character.id}
                            id={character.id}
                            name={character.name}
                            description={character.description}
                            isStarred={character.isStarred}
                            onStarToggle={handleStarToggle}
                            updatedAt={character.updateAt}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

export default WorkSpaceContent;
