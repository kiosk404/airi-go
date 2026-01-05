import {Avatar, Card, Dropdown, Button, Toast} from "@douyinfe/semi-ui";
import { IconStar, IconMore, IconEdit, IconDelete } from "@douyinfe/semi-icons";
import { useNavigate } from "react-router-dom";
import { deleteCharacter } from "@/services/character/api-client";

interface CharacterCardProps {
    id: string;                     // 角色id
    name: string;                   // 角色名
    avatar: string;                 // 头像url
    description: string;            // 简介
    updatedAt: number;              // 更新时间戳（毫秒）
    isStarred?: boolean;            // 是否收藏
    onStarToggle?: (id: string) => void;  // 收藏状态切换回调
    onDelete?: (id: string) => void;  // 删除成功回调
}

// 格式化时间戳为可读字符串
const formatTime = (timestamp: number): string => {
    const date = new Date(timestamp);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}/${month}/${day} ${hours}:${minutes}`;
};

const CharacterCardSpace: React.FC<CharacterCardProps> = ({
                                                              id, name, avatar, description, updatedAt, isStarred = false, onStarToggle, onDelete }) => {

    const navigate = useNavigate();

    const handleStarClick = (e: React.MouseEvent) => {
        e.stopPropagation();
        onStarToggle?.(id);
    };

    const handleMoreClick = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const handleEdit = (e: React.MouseEvent) => {
        e.stopPropagation();
        console.log(`编辑角色: ${name} (ID: ${id})`);
    };

    const handleDelete = async (e: React.MouseEvent) => {
        e.stopPropagation();
        try {
            await deleteCharacter({
                bot_id: id
            });
            Toast.success('删除成功');
            onDelete?.(id)  // 删除成功后调用回调函数
        } catch (error) {
            Toast.error('删除失败');
            console.error(`删除角色失败: ${name} (ID: ${id})`, error);
        }
    };

    return (
        <div
            onClick={() => navigate(`/workspace/playground/${id}`)}
            style={{ cursor: 'pointer', height: '100%' }}
        >
            <Card
                bordered
                shadows="hover"
                style={{ height: '100%', boxSizing: 'border-box' }}
                bodyStyle={{ padding: '16px 20px', height: '100%', boxSizing: 'border-box' }}
            >
                {/* 整体布局容器 */}
                <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
                    {/* 上半部分：名称、描述、头像 */}
                    <div style={{ display: 'flex', gap: '14px', marginBottom: '10px' }}>
                        {/* 左侧：名称和描述 */}
                        <div style={{ flex: 1, minWidth: 0, paddingTop: '12px' }}>
                            <div style={{ fontSize: '17px', fontWeight: '600', color: '#333', marginBottom: '24px' }}>
                                {name}
                            </div>
                            <div style={{
                                fontSize: '13px',
                                color: '#666',
                                lineHeight: '1.5',
                                overflow: 'hidden',
                                display: '-webkit-box',
                                WebkitLineClamp: 2,
                                WebkitBoxOrient: 'vertical',
                            }}>
                                {description}
                            </div>
                        </div>

                        {/* 右侧：头像 */}
                        <div style={{ flexShrink: 0, marginTop: '20px'}}>
                            <Avatar
                                shape="square"
                                style={{ width: '55px', height: '55px' }}
                                src={avatar}
                            />
                        </div>
                    </div>

                    {/* 标签 - 使用 marginTop: auto 推到底部时间上方 */}
                    <div style={{ marginTop: 'auto', marginBottom: '8px' }}>
                        <span style={{
                            color: '#666',
                            backgroundColor: '#f5f5f5',
                            borderRadius: '4px',
                            padding: '3px 10px',
                            fontSize: '12px'
                        }}>
                            智能体
                        </span>
                    </div>

                    {/* 底部：时间和按钮 */}
                    <div style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                    }}>
                        <span style={{ fontSize: '12px', color: '#999' }}>
                            最近编辑 {formatTime(updatedAt)}
                        </span>

                        <div style={{ display: 'flex', gap: '8px' }}>
                            <Button
                                type="tertiary"
                                theme="light"
                                size="small"
                                style={{ color: '#666', backgroundColor: '#f5f5f5', borderRadius: '6px' }}
                                onClick={handleStarClick}
                                icon={<IconStar style={isStarred ? { color: '#faad14' } : undefined} />}
                            />
                            <Dropdown
                                position="bottomRight"
                                render={
                                    <Dropdown.Menu>
                                        <Dropdown.Item icon={<IconEdit />} onClick={handleEdit}>编辑</Dropdown.Item>
                                        <Dropdown.Item icon={<IconDelete />} onClick={handleDelete}>删除</Dropdown.Item>
                                    </Dropdown.Menu>
                                }
                            >
                                <Button
                                    type="tertiary"
                                    theme="light"
                                    size="small"
                                    style={{ color: '#666', backgroundColor: '#f5f5f5', borderRadius: '6px' }}
                                    onClick={handleMoreClick}
                                    icon={<IconMore />}
                                />
                            </Dropdown>
                        </div>
                    </div>
                </div>
            </Card>
        </div>
    );
}

export default CharacterCardSpace;