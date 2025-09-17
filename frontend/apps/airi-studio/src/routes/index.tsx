import { createBrowserRouter, Navigate } from 'react-router-dom';
import { Layout, Dashboard, Develop, Library, Explore } from './async-components';

export const router = createBrowserRouter([
  {
    path: '/',
    Component: Layout,
    children: [
      { index: true, element: <Navigate to="/space/demo/develop" replace /> },
      {
        path: 'space',
        children: [
          { index: true, element: <Navigate to="demo/develop" replace /> },
          {
            path: ':space_id',
            children: [
              { index: true, element: <Navigate to="develop" replace /> },
              { path: 'dashboard', Component: Dashboard },
              { path: 'develop', Component: Develop },
              { path: 'library', Component: Library },
            ],
          },
        ],
      },
      {
        path: 'explore',
        Component: Explore,
      },
    ],
  },
]);