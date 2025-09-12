import { RouterProvider } from 'react-router-dom';
import { Suspense } from 'react';
import { Spin } from 'antd';
import { router } from './routes';

function App() {
  return (
    <Suspense
        fallback={
            <div className="w-full h-full flex items-center justify-center">
                <Spin spinning style={{ height: '100%', width: '100%' }} />
            </div>
        }
    >
        <RouterProvider router={router} fallbackElement={<div>loading...</div>} />
    </Suspense>
  );
}

export default App;

