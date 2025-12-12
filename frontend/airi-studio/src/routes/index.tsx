import { createBrowserRouter, Navigate } from "react-router-dom";
import AppLayout from '@/layout';
import WorkspacePage from '@/pages/workspace';
import KnowledgePage from "@/pages/knowledge";
import AppsPage from "@/pages/apps";
import ModelsPage from "@/pages/models";
import PlaygroundPage from "@/pages/workspace/playground";
import LoginPage from "@/pages/login";
import RequireAuth from "@/components/auth/RequireAuth";

const LoadingFallback = () => <div>loading</div>

export const router = createBrowserRouter([
    {
        path: "/login",
        Component: LoginPage,
    },
    {
        path: "/",
        element: (
            <RequireAuth>
                <AppLayout />
            </RequireAuth>
        ),
        HydrateFallback: LoadingFallback,
        children: [
            { index: true, element: <Navigate to="/apps" replace /> },
            { path: "apps", Component: AppsPage },
            { path: "workspace", Component: WorkspacePage },
            { path: "knowledge", Component: KnowledgePage },
            { path: "models", Component: ModelsPage },
            { path: "tool", Component: KnowledgePage },
        ]
    },
    {
        path: "/workspace/playground",
        element: (
            <RequireAuth>
                <AppLayout />
            </RequireAuth>
        ),
        HydrateFallback: LoadingFallback,
        children: [
            {
                path: ":id",
                Component: PlaygroundPage
            }
        ]
    }
]);