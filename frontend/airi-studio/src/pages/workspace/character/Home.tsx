import React, {useState} from "react";
import CreateCharacterCard from "./CreateCharacterCard";
import CharacterCardSpace from "./CharacterCardList";


const mockCharacter = [
    { id: 1, name: "Airis", description: "基于大语言模型的AI管家", isStarred: true },
    { id: 2, name: "Atlas", description: "基于大语言模型的AI小笨蛋", isStarred: true }
];


const WorkSpaceContent: React.FC = () => {
    const [characters, setCharacters] = useState(mockCharacter);

    const handleStarToggle = (id: number) => {
        setCharacters(characters.map(c => c.id === id ? {...c, isStarred: !c.isStarred} : c))
    }


    return (
        <div style={{ padding: "0 20px", marginTop: 20 }}> 
            <div style={{ marginBottom: 16, fontSize: 16, fontWeight: 300 }}>角色列表</div>

            <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(300px, 1fr))", gap: 24 }}>
            <CreateCharacterCard />
            
      
            {characters.map((character) => (
                <CharacterCardSpace 
                key={character.id} 
                id={character.id} 
                name={character.name}
                description={character.description}
                isStarred={character.isStarred}
                onStarToggle={handleStarToggle}
                />
            ))}
            </div>
            
        </div>
    );
};

export default WorkSpaceContent;