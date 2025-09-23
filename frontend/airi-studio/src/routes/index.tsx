import { createBrowserRouter, Navigate } from "react-router-dom";
import AppLayout from '../layout';
import PlaygroundPage from '../pages/playground';

const LoadingFacllback = () => <div>loading</div>

export const router = createBrowserRouter([
    {
        path: "/",
        Component: AppLayout,
        HydrateFallback: LoadingFacllback,
        children: [
            { index: true, element: <Navigate to="/playground" replace /> },
            { path: "playground", Component: PlaygroundPage },
        ]
    },
    
]);