import {lazy} from 'react';

// 页面组件懒加载
export const Playground = lazy(() => import('../pages/playground/index'));

// 布局
export const Layout = lazy(() => import('../layout'));