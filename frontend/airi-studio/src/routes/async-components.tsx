import {lazy} from 'react';

// 页面组件懒加载
export const Playground = lazy(() => import('@/pages/workspace/playground/index'));
export const Workspace = lazy(() => import('../pages/workspace/index'));
export const Knowledge = lazy(() => import("../pages/knowledge/index"));


// 布局
export const Layout = lazy(() => import('../layout'));