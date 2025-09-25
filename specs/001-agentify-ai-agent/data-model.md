# Data Model: Agentify - AI Agent 开发调试平台

## AI Agent

代表用户创建的私人 AI 助手。

**属性**:
- id: string (唯一标识符)
- name: string (Agent 名称)
- description: string (Agent 描述)
- systemPrompt: string (系统提示词)
- plugins: Plugin[] (插件列表)
- workflows: Workflow[] (工作流列表)
- createdAt: Date (创建时间)
- updatedAt: Date (更新时间)

**验证规则**:
- name: 必填，长度 1-100 字符
- description: 可选，长度 0-500 字符
- systemPrompt: 必填，长度 1-10000 字符

## Workspace

代表三种不同的工作环境。

**枚举值**:
- Develop: 开发工作空间，用于数据分析
- Playground: 演练场，用于提示词编排和对话
- Explore: 探索空间，用于创建工作流和插件

## Plugin/Tool

代表可以添加到 AI Agent 的功能扩展。

**属性**:
- id: string (唯一标识符)
- name: string (插件名称)
- description: string (插件描述)
- icon: string (插件图标)
- config: object (插件配置参数)
- enabled: boolean (是否启用)
- createdAt: Date (创建时间)

**验证规则**:
- name: 必填，长度 1-50 字符
- description: 可选，长度 0-200 字符

## Workflow

代表用户创建的自动化工作流程。

**属性**:
- id: string (唯一标识符)
- name: string (工作流名称)
- description: string (工作流描述)
- steps: object[] (步骤序列)
- createdAt: Date (创建时间)
- updatedAt: Date (更新时间)

**验证规则**:
- name: 必填，长度 1-50 字符
- description: 可选，长度 0-200 字符

## Prompt

代表 AI Agent 的系统提示词。

**属性**:
- id: string (唯一标识符)
- content: string (提示词内容)
- version: number (版本号)
- createdAt: Date (创建时间)
- updatedAt: Date (更新时间)

**验证规则**:
- content: 必填，长度 1-10000 字符

## AnalyticsData

代表 Develop 工作空间中的分析数据。

**属性**:
- agentId: string (关联的 Agent ID)
- tokenConsumption: number (Token 消耗量)
- pluginCount: number (插件数量)
- workflowCount: number (工作流数量)
- lastActive: Date (最后活跃时间)

## ChatMessage

代表 Playground 中的对话消息。

**属性**:
- id: string (唯一标识符)
- agentId: string (关联的 Agent ID)
- role: "user" | "assistant" (消息角色)
- content: string (消息内容)
- timestamp: Date (时间戳)

**验证规则**:
- content: 必填，长度 1-10000 字符
- role: 必须是 "user" 或 "assistant"