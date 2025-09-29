import { Dropdown, Card } from "@douyinfe/semi-ui";
import { useNavigate } from "react-router-dom";
import { StarIcon, TagIcon, MoreIcon } from "tdesign-icons-react"
import { Space, Button, Avatar } from "tdesign-react"

interface CharacterCardProps {
    id: number;                     // 角色id
    name: string;                   // 角色名
    description: string;              // 简介
    isStarred: boolean;         // 是否收藏
    onStarToggle: (id: number) => void;     // 收藏状态切换回调
}

function DropDownMenu({ name }: { name: string }) {
    // 处理编辑按钮的点击事件
    const handleEditClick = () => {
        console.log(`Edit ${name}`);
    }

    // 处理删除按钮的点击事件
    const handleDeleteClick = () => {
        console.log(`Delete ${name}`);
    }

    return (
        <>
            <Dropdown position="bottomLeft" render={
                <Dropdown.Menu>
                    <Dropdown.Title>操作</Dropdown.Title>
                    <Dropdown.Item icon={<StarIcon />}>编辑</Dropdown.Item>
                    <Dropdown.Item icon={<StarIcon />}>删除</Dropdown.Item>
                </Dropdown.Menu>
            }>
                <MoreIcon />
            </Dropdown>
        </>
    )
}

const CharacterCardSpace: React.FC<CharacterCardProps> = ({
    id, name, description, isStarred, onStarToggle }) => {

    const navigate = useNavigate();

    const handleStarClick = (e: React.MouseEvent) => {
        e.stopPropagation(); // 防止事件冒泡到卡片
        onStarToggle(id);
    }

    const handleAddTagClick = (e: React.MouseEvent) => {
        e.stopPropagation(); // 防止事件冒泡到卡片
        console.log(`Add tag to ${name}`);
    }

    return (
        <>
            <div onClick={() => navigate(`/workspace/playground/${id}`)} >
                <Card shadows="hover" bordered headerLine={true}>
                    {/* 卡片头部 */}
                    <div style={{ display: 'flex', alignItems: 'center', padding: '12px 14px', borderBottom: '1px solid #f0f0f0' }}>
                        <Avatar 
                            hideOnLoadFailed
                            size="40px"
                            shape="round"
                            image="https://raw.githubusercontent.com/moeru-ai/airi/refs/heads/main/docs/content/public/logo-dark.avif"
                            style={{ marginRight: '30px' }} 
                        />
                        <div>
                            <div style={{ fontSize: '14px', fontWeight: '600', color: '#333333' }}>{name}</div>
                            <div style={{ fontSize: '14px', color: '#999999' }}> kiosk404 编辑于 2025/09/24 16:13 </div>
                        </div>
                        {/* 右侧操作按钮 */}
                        <Space size="small">
                            <Button
                                variant="text" 
                                shape="circle"
                                size="small"
                                style={{ color: '#888888' }}
                                onClick={handleStarClick} >
                                {isStarred ? <StarIcon size="16" style={{ color: '#faad14', fill: '#faad14' }} /> : <StarIcon size="16" />}
                            </Button>
                            <DropDownMenu name={name} />
                        </Space>
                    </div>


                    {/* 卡片内容 */}
                    <div style={{ padding: '12px 16px', fontSize: '13px', color: '#666666', minHeight: '60px', display: 'flex', alignItems: 'center', justifyContent: 'flex-start'}}>
                        <div style={{ display: '-webkit-box', WebkitLineClamp: '2', WebkitBoxOrient: 'vertical', overflow: 'hidden', textOverflow: 'ellipsis', width: '100%' }}>
                            {description}
                        </div>
                    </div>
                    
                    {/* 卡片底部 */}
                    <div style={{ display: 'flex', alignItems: 'center', padding: '2px 5px'}}>
                        <Button 
                            variant="text"
                            size="small"
                            onClick={handleAddTagClick}
                            icon={<TagIcon size={12} />}  
                            style={{ color: '#888888', border: '1px dashed #d9d9d9', borderRadius: '4px', padding: '2px 6px', fontSize: '12px'}}>
                            添加标签
                        </Button>
                    </div>
                </Card>
            </div>
        </>
    );
}

export default CharacterCardSpace;