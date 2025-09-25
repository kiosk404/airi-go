import { createBrowserRouter, Navigate } from "react-router-dom";
import AppLayout from '../layout';
import PlaygroundPage from '../pages/playground';
import WorkspacePage from '../pages/workspace';

const LoadingFallback = () => <div>loading</div>

export const router = createBrowserRouter([
    {
        path: "/",
        Component: AppLayout,
        HydrateFallback: LoadingFallback,
        children: [
            { index: true, element: <Navigate to="/apps" replace /> },
            { path: "apps", Component: WorkspacePage },
            { path: "workspace", Component: WorkspacePage },
            { path: "playground", Component: PlaygroundPage },
            { path: "knowledge", Component: WorkspacePage },
            { path: "tool", Component: WorkspacePage },
        ]
    },
    
]);