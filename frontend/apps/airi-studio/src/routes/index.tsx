import { createBrowserRouter, Navigate } from 'react-router-dom';

import { Layout } from '../layout';
import {
  Redirect,
  Dashboard,
  Develop,
  Library,
  PluginLayout,
  PluginPage,
  PluginToolPage,
  ExplorePluginPage,
  ExploreTemplatePage,
} from './async-components';
import SpaceLayout from '../pages/space-layout';

export const router = createBrowserRouter([
  // docs redirects
  {
    path: '/open/docs/*',
    Component: Redirect,
    loader: () => ({ hasSider: false, requireAuth: false }),
  },
  {
    path: '/docs/*',
    Component: Redirect,
    loader: () => ({ hasSider: false, requireAuth: false }),
  },
  // main app
  {
    path: '/',
    Component: Layout,
    children: [
      { index: true, element: <Navigate to="/space/demo/develop" replace /> },

      // workspace
      {
        path: 'space',
        Component: SpaceLayout,
        loader: () => ({ hasSider: true, requireAuth: false }),
        children: [
          { index: true, element: <Navigate to="demo/develop" replace /> },
          {
            path: ':space_id',
            children: [
              { index: true, element: <Navigate to="develop" replace /> },
              { path: 'dashboard', Component: Dashboard },
              { path: 'develop', Component: Develop },
              { path: 'library', Component: Library },

              // plugin resources
              {
                path: 'plugin/:plugin_id',
                Component: PluginLayout,
                children: [
                  { index: true, Component: PluginPage },
                  {
                    path: 'tool/:tool_id',
                    children: [{ index: true, Component: PluginToolPage }],
                  },
                ],
              },
            ],
          },
        ],
      },

      // explore
      {
        path: 'explore',
        Component: () => null,
        children: [
          { index: true, element: <Navigate to="plugin" replace /> },
          { path: 'plugin', element: <ExplorePluginPage /> },
          { path: 'template', element: <ExploreTemplatePage /> },
        ],
      },
    ],
  },
]);
