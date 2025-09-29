import { Dropdown,  } from "@douyinfe/semi-ui";
import { useNavigate } from "react-router-dom";
import { StarIcon, TagIcon, MoreIcon } from "tdesign-icons-react"
import { Space, Button, Avatar, Card } from "tdesign-react"

interface CharacterCardProps {
    id: number;                     // 角色id
    name: string;                   // 角色名
    description: string;              // 简介
    isStarred: boolean;         // 是否收藏
    onStarToggle: (id: number) => void;     // 收藏状态切换回调
}


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
};

// 【圆形头像容器样式】 - 强制 40x40 裁剪成圆形
const avatarWrapperStyle: React.CSSProperties = {
    marginRight: '30px',
};

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

    // 动态样式：星标颜色
    const starIconStyle = { color: '#888888' };
    const starredIconStyle = { color: '#faad14', fill: '#faad14' };
    
    // 右侧操作按钮，使用 margin-left: auto 推到最右边
    const actionSpaceStyle = { marginLeft: 'auto' };

    return (
        <>
            <div onClick={() => navigate(`/workspace/playground/${id}`)} style={clickableWrapperStyle}>
            
            <Card bordered headerBordered hoverShadow shadow size="medium" theme="normal" >
                
                {/* 卡片头部 */}
                <div style={cardHeaderStyle}>
                    
                    {/* 【应用圆形头像容器】 */}
                    <div style={avatarWrapperStyle}>
                        <Avatar 
                            hideOnLoadFailed
                            size="40px"
                            shape="round" // 保持 round，但由外部样式控制最终形状
                            image="https://raw.githubusercontent.com/moeru-ai/airi/refs/heads/main/docs/content/public/logo-dark.avif"
                            // 这里的 Avatar 必须放在一个 40x40 且 overflow: hidden 的容器内才能实现裁剪
                        />
                    </div>
                    
                    {/* 信息区 */}
                    <div style={headerInfoStyle}> 
                        <div style={nameTextStyle}>{name}</div>
                        <div style={timeTextStyle}> kiosk404 编辑于 2025/09/24 16:13 </div>
                    </div>
                    
                    {/* 右侧操作按钮 */}
                    <Space size="small" style={actionSpaceStyle}>
                        <Button
                            variant="text" 
                            shape="circle"
                            size="small"
                            style={starIconStyle}
                            onClick={handleStarClick} >
                            {isStarred 
                                ? <StarIcon size="16" style={starredIconStyle} /> 
                                : <StarIcon size="16" />}
                        </Button>
                        <DropDownMenu name={name} />
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
                        variant="text"
                        size="small"
                        onClick={handleAddTagClick}
                        icon={<TagIcon size={12} />}  
                        style={addTagButtonStyle}>
                        添加标签
                    </Button>
                </div>
            </Card>
        </div>
        </>
    );
}

export default CharacterCardSpace;