import { lazy } from 'react';

// 页面组件懒加载
export const Dashboard = lazy(() => import('../pages/dashboard/DashboardPage'));
export const Develop = lazy(() => import('../pages/develop/DevelopPage'));
export const Library = lazy(() => import('../pages/library/LibraryPage'));
export const Explore = lazy(() => import('../pages/explore/ExplorePage'));

// 布局组件
export const Layout = lazy(() => import('../layout'));