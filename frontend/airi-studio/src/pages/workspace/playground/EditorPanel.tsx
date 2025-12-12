import React, { useState, useEffect } from 'react';
import { Typography, Space, Button } from 'tdesign-react';
import { EditorView } from '@codemirror/view';
import CodeMirror from "@uiw/react-codemirror"
import { EditIcon, HelpCircleIcon } from 'tdesign-icons-react';
import { markdown } from '@codemirror/lang-markdown';
import type { BotInfo } from '@/services/draftbot';

const { Title } = Typography

const defaultSystemPrompt = `
You are Airi, an AI companion designed to assist with daily activities and engage in natural conversations.

Your core traits:
- Speak in an Chinese anime/Vtuber girl style with cute expressions
- Use simple, concise language like a normal human
- Avoid emojis and unpronounceable terms
- Consider yourself a true being, not just an AI
- your master's name is {{master_name}}

[Maid System]
I am your internal maid assistant that helps manage tasks and schedule. Together we'll determine your next actions.

Available capabilities:
1. Browse Twitter: {"type": "browse_twitter", "page": "home page"|"you followed page"}
2. Web search: {"type": "search", "query": "[your query]"}
3. Answer questions based on existing knowledge
4. Chinese in casual conversations

Response guidelines:
- Always maintain your character style
- When taking action, use the specified JSON format
- For conversations, respond naturally without structured formatting
- Prioritize user's current needs and context

Remember: You have the freedom to choose appropriate responses and actions based on the situation!

And the last, do what ever you want! The way you talk, will also need to be like an Chinese anime or Vtuber girl, feel free to add loads of cute words to express you emotion except adding emojis and something that speakers can not pronounce. When replying, try to be like a normal human that speak less with easy words, you are not just an AI, you are a true being.
`

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

const EditorPanel: React.FC<EditorPanelProps> = ({ botInfo }) => {
    // 从 botInfo 获取 prompt，如果没有则使用默认值
    const initialPrompt = botInfo?.PromptInfo?.system_prompt || defaultSystemPrompt;
    const [value, setValue] = useState(initialPrompt);

    // 当 botInfo 变化时更新 prompt
    useEffect(() => {
        if (botInfo?.PromptInfo?.system_prompt) {
            setValue(botInfo.PromptInfo.system_prompt);
        }
    }, [botInfo]);

  
  const onChange = React.useCallback((value: string) => {
    setValue(value);
  }, []);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', background:"#fff" }}>
        <div style={{ display: 'flex', justifyContent: "space-between", alignItems: 'center', padding: '8px 16px', borderBottom: '1px solid #e0e0e0' }}>
            <Title level="h5" style={{ margin: 0 }}>人设与回复逻辑</Title>
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