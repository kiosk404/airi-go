import React, { useState, useEffect, useRef } from 'react';
import { Typography, Space, Button } from 'tdesign-react';
import { EditorView } from '@codemirror/view';
import CodeMirror from "@uiw/react-codemirror"
import { EditIcon, HelpCircleIcon } from 'tdesign-icons-react';
import { markdown } from '@codemirror/lang-markdown';
import type { BotInfo, } from '@/services/draftbot';
import * as draftBotService from '@/services/draftbot';

const { Title } = Typography

const customTheme = EditorView.theme({
    "&": {
        color: "#999",
        backgroundColor: "#fafafa",
        fontFamily: "Fira Code, monospace",
        fontSize: "14px",
    },
    ".cm-content": {    
        caretColor: "#ff0000",
    },
    "&.cm-editor.cm-focused": {
        outline: "none"   // 去掉默认的蓝色
    }
})

interface EditorPanelProps {
    botInfo?: BotInfo | null;
}

const ensureString = (val: unknown): string => {
    if (val === null || val === undefined) {
        return '';
    }
    return String(val);
}

const EditorPanel: React.FC<EditorPanelProps> = ({ botInfo }) => {
    // 从 botInfo 获取 prompt，如果没有则使用默认值
    const initialPrompt = ensureString(botInfo?.prompt_info?.prompt);
    const [value, setValue] = useState(initialPrompt);
    const saveTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    // 当 botInfo 变化时更新 prompt
    useEffect(() => {
        if (botInfo?.prompt_info?.prompt) {
            setValue(ensureString(botInfo?.prompt_info?.prompt));
        }
    }, [botInfo]);

    const updateDraftBotInfo = async (promptContent: string) => {
        console.log(botInfo);
        if (!botInfo?.bot_id) {
            console.error("Bot ID is missing, pass update draft bot info")
            return;
        }

        try {
            await draftBotService.updateDraftBotInfo({
                bot_info: {
                    bot_id: botInfo.bot_id,
                    prompt_info: {
                        prompt: promptContent,
                    }
                }
            });
        } catch (error) {
            console.error("Update draft bot info failed")
        }
    }
  
    const onChange = React.useCallback((value: string) => {
        setValue(value);

        console.log("kiosk---> 123: ", value);

        // 清除定时器
        if (saveTimeoutRef.current) {
            clearTimeout(saveTimeoutRef.current);
        }

        // 500ms 后保存
        saveTimeoutRef.current = setTimeout(() => {
            updateDraftBotInfo(value).then(() => {});
        }, 500);
    }, [botInfo?.bot_id]);

    // 组件卸载时取消定时器
    useEffect(() => {
        return () => {
            if (saveTimeoutRef.current) {
                clearTimeout(saveTimeoutRef.current);
            }
        }
    }, []);

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%', background:"#fff" }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0 16px', height: '56px', borderBottom: '1px solid #e0e0e0' }}>
                <Title level="h5" style={{ margin: 0, lineHeight: '50px' }}>人设与回复逻辑</Title>
                <Space>
                    <Button shape='square' variant='text' icon={<EditIcon />} />
                    <Button shape='square' variant='text' icon={<HelpCircleIcon />} />
                </Space>
            </div>
            <CodeMirror
                value={value}
                theme="light"
                extensions={[
                    markdown(),
                    EditorView.lineWrapping,
                    customTheme,
                ]}
                onChange={onChange}
                basicSetup={{ lineNumbers: false }}
                style={{ flex: 1, overflowY: 'auto'}}
            />
        </div>
    )
};

export default EditorPanel;