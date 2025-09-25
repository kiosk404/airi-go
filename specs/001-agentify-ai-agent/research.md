# Research Findings: Agentify - AI Agent 开发调试平台

## Ant Design Integration with rsbuild + React + TypeScript

**Decision**: 使用 @ant-design/react 中的组件库，并配置 rsbuild 支持 Ant Design 的样式和主题定制。

**Rationale**: 
- Ant Design 是成熟的企业级 UI 设计语言和 React 组件库
- 与 React 和 TypeScript 有良好的类型支持
- rsbuild 官方文档提供了 Ant Design 的集成示例
- 社区支持广泛，文档完善

**Alternatives considered**:
- Material-UI: 功能强大但与项目要求的 Ant Design 主题不一致
- 自定义组件库: 开发成本高，不符合快速原型的要求

## Tailwind CSS + Ant Design Integration

**Decision**: 采用 Tailwind CSS 作为基础样式系统，Ant Design 作为高级组件库，通过配置解决样式冲突。

**Rationale**:
- Tailwind CSS 提供原子化 CSS 类，便于快速开发和维护
- Ant Design 提供复杂的业务组件
- 通过配置可以实现两者共存，Tailwind 处理布局和基础样式，Ant Design 处理复杂交互组件

**Alternatives considered**:
- 仅使用 Ant Design: 缺乏灵活性，难以处理定制化样式需求
- 仅使用 Tailwind CSS: 需要重新实现复杂的交互组件

## Local JSON Mock Data Patterns for Frontend Development

**Decision**: 使用 Mock Service Worker (MSW) 作为本地 JSON Mock 解决方案。

**Rationale**:
- MSW 可以拦截网络请求并返回模拟数据
- 支持 REST 和 GraphQL API
- 与真实 API 兼容，便于后期切换
- 不需要修改应用代码即可切换 mock 和真实数据

**Alternatives considered**:
- 手动创建 mock 数据文件: 难以模拟真实 API 行为
- JSON Server: 需要额外的服务进程

## UI/UX Patterns for AI Agent Development Platforms

**Decision**: 参考 Coze Studio 的设计模式，采用三栏布局：左侧导航，中间主内容区，右侧配置面板。

**Rationale**:
- Coze Studio 是同类产品的成熟解决方案
- 三栏布局适合复杂配置和实时预览场景
- 工作空间分离符合用户心智模型

**Alternatives considered**:
- 单页应用: 不适合复杂配置场景
- 标签页切换: 无法同时查看多个工作区内容

## Implementation Technologies

**Decision**: 
- 使用 rsbuild 作为构建工具，支持快速开发和构建
- 使用 React Hooks 和 TypeScript 确保代码质量
- 使用 React Router 处理页面路由
- 使用 Zustand 或 Context API 管理全局状态

**Rationale**:
- rsbuild 是现代化的构建工具，性能优于 Webpack
- TypeScript 提供类型安全，减少运行时错误
- React Router 是 React 生态中的标准路由解决方案
- Zustand 轻量且易于使用，适合中小型项目

## Styling Approach

**Decision**: Tailwind CSS + Ant Design + CSS Modules

**Rationale**:
- Tailwind CSS 用于快速布局和基础样式
- Ant Design 用于复杂组件
- CSS Modules 用于局部样式覆盖和定制

## Data Management

**Decision**: 使用本地存储 (localStorage) 模拟后端数据持久化

**Rationale**:
- 符合项目约束（无后端依赖）
- localStorage API 简单可靠
- 便于演示和测试

## Testing Strategy

**Decision**: 使用 Jest + React Testing Library 进行单元测试和集成测试

**Rationale**:
- Jest 是 JavaScript 生态中的标准测试框架
- React Testing Library 鼓励以用户视角测试组件
- 与 Create React App 和类似工具集成良好