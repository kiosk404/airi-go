import { Dropdown, Space, Button, Avatar, Card } from "@douyinfe/semi-ui";
import { IconStar, IconMore, IconPriceTag, IconEdit, IconDelete } from "@douyinfe/semi-icons";
import { useNavigate } from "react-router-dom";

interface CharacterCardProps {
    id: number;                     // 角色id
    name: string;                   // 角色名
    description: string;            // 简介
    isStarred: boolean;             // 是否收藏
    onStarToggle: (id: number) => void;     // 收藏状态切换回调
    updatedAt: number;              // 更新时间戳（毫秒）
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


// 卡片头部样式
const cardHeaderStyle: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    padding: '12px 16px',
    borderBottom: '1px solid #f0f0f0',
};

// 头像和信息之间的间距容器
const headerInfoStyle: React.CSSProperties = {
    marginRight: '30px',
};

// 名称文本样式
const nameTextStyle: React.CSSProperties = {
    fontSize: '14px',
    fontWeight: '600',
    color: '#333333',
};

// 时间文本样式
const timeTextStyle: React.CSSProperties = {
    fontSize: '14px',
    color: '#999999',
};

// 卡片内容容器样式
const cardContentStyle: React.CSSProperties = {
    padding: '12px 16px',
    fontSize: '13px',
    color: '#666666',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-start',
};

// 描述文本样式 (【关键：固定高度】)
// 2.8em 保证了两行 13px 文本的空间
const descriptionTextStyle: React.CSSProperties = {
    display: '-webkit-box',
    WebkitLineClamp: 2, // 注意：在 TSX 中数字不加引号
    WebkitBoxOrient: 'vertical',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: '100%',
    height: '3.6em',
    minHeight: '3.6em',
    maxHeight: '3.6em',
};

// 卡片底部样式
const cardFooterStyle: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    padding: '2px 5px',
};

// 底部添加标签按钮样式
const addTagButtonStyle: React.CSSProperties = {
    color: '#888888',
    border: '1px dashed #d9d9d9',
    borderRadius: '4px',
    padding: '2px 6px',
    fontSize: '12px',
};

// 【手形光标样式】
const clickableWrapperStyle: React.CSSProperties = {
    cursor: 'pointer',
    height: '100%',
};

// 【圆形头像容器样式】 - 强制 40x40 裁剪成圆形
const avatarWrapperStyle: React.CSSProperties = {
    marginRight: '30px',
};

interface DropDownMenuProps {
    name: string;
    onEdit: () => void;
    onDelete: () => void;
}

function DropDownMenu({ onEdit, onDelete }: DropDownMenuProps) {
    const handleClick = (e: React.MouseEvent) => {
        e.stopPropagation(); // 阻止事件冒泡到卡片
    };

    const handleEdit = (e: React.MouseEvent) => {
        e.stopPropagation();
        onEdit();
    };

    const handleDelete = (e: React.MouseEvent) => {
        e.stopPropagation();
        onDelete();
    };

    return (
        <Dropdown
            position="bottomLeft"
            render={
                <Dropdown.Menu>
                    <Dropdown.Title>操作</Dropdown.Title>
                    <Dropdown.Item icon={<IconEdit />} onClick={handleEdit}>编辑</Dropdown.Item>
                    <Dropdown.Item icon={<IconDelete />} onClick={handleDelete}>删除</Dropdown.Item>
                </Dropdown.Menu>
            }
        >
            <Button
                type="tertiary"
                theme="borderless"
                icon={<IconMore />}
                size="small"
                onClick={handleClick}
            />
        </Dropdown>
    );
}

const CharacterCardSpace: React.FC<CharacterCardProps> = ({
                                                              id, name, description, isStarred, onStarToggle, updatedAt }) => {

    const navigate = useNavigate();

    const handleStarClick = (e: React.MouseEvent) => {
        e.stopPropagation(); // 防止事件冒泡到卡片
        onStarToggle(id);
    }

    const handleAddTagClick = (e: React.MouseEvent) => {
        e.stopPropagation(); // 防止事件冒泡到卡片
        console.log(`Add tag to ${name}`);
    }

    const handleEdit = () => {
        console.log(`编辑角色: ${name} (ID: ${id})`);
        // TODO: 实现编辑逻辑
    };

    const handleDelete = () => {
        console.log(`删除角色: ${name} (ID: ${id})`);
        // TODO: 实现删除逻辑
    };

    // 动态样式：星标颜色
    const starIconStyle = { color: '#888888' };
    const starredIconStyle = { color: '#faad14', fill: '#faad14' };

    // 右侧操作按钮，使用 margin-left: auto 推到最右边
    const actionSpaceStyle = { marginLeft: 'auto' };

    return (
        <div onClick={() => navigate(`/workspace/playground/${id}`)} style={clickableWrapperStyle}>
            <Card
                bordered
                shadows="hover"
                style={{ height: '100%', boxSizing: 'border-box' }}
            >
                {/* 卡片头部 */}
                <div style={cardHeaderStyle}>
                    {/* 【应用圆形头像容器】 */}
                    <div style={avatarWrapperStyle}>
                        <Avatar
                            size="default"
                            shape="square"
                            color="red"
                            src="https://raw.githubusercontent.com/moeru-ai/airi/refs/heads/main/docs/content/public/logo-dark.avif"
                        />
                    </div>

                    {/* 信息区 */}
                    <div style={headerInfoStyle}>
                        <div style={nameTextStyle}>{name}</div>
                        <div style={timeTextStyle}>编辑于 {formatTime(updatedAt)}</div>
                    </div>

                    {/* 右侧操作按钮 */}
                    <Space spacing={4} style={actionSpaceStyle}>
                        <Button
                            type="tertiary"
                            theme="borderless"
                            size="small"
                            style={starIconStyle}
                            onClick={handleStarClick}
                            icon={
                                isStarred
                                    ? <IconStar style={starredIconStyle} />
                                    : <IconStar />
                            }
                        />
                        <DropDownMenu name={name} onEdit={handleEdit} onDelete={handleDelete} />
                    </Space>
                </div>

                {/* 卡片内容 */}
                <div style={cardContentStyle}>
                    {/* 【应用固定高度描述文本】 */}
                    <div style={descriptionTextStyle}>
                        {description}
                    </div>
                </div>

                {/* 卡片底部 */}
                <div style={cardFooterStyle}>
                    <Button
                        type="tertiary"
                        theme="borderless"
                        size="small"
                        onClick={handleAddTagClick}
                        icon={<IconPriceTag size="small" />}
                        style={addTagButtonStyle}
                    >
                        添加标签
                    </Button>
                </div>
            </Card>
        </div>
    );
}

export default CharacterCardSpace;