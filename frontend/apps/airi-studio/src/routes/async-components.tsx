import { lazy } from 'react';

export const Redirect = lazy(() => import('../pages/redirect'));
export const Dashboard = lazy(() => import('../pages/DashboardPage'));
export const Develop = lazy(() => import('../pages/develop'));
export const Library = lazy(() => import('../pages/library'));

export const PluginLayout = lazy(() => import('../pages/plugin/layout'));
export const PluginPage = lazy(() => import('../pages/plugin/page'));
export const PluginToolPage = lazy(() => import('../pages/plugin/tool/page'));

export const ExplorePluginPage = lazy(() => import('../pages/explore-plugin'));
export const ExploreTemplatePage = lazy(() => import('../pages/explore-template'));
