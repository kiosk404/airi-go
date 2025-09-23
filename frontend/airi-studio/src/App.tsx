import './App.css';
import { RouterProvider } from 'react-router-dom';
import { Suspense } from 'react';
import { Loading } from 'tdesign-react';
import { router } from './routes';

const App = () => {
  return (
    <Suspense
    fallback={
      <div className='w-full h-full flex justify-center items-center'><Loading size='large'/></div>
    }
    >
      <RouterProvider router={router} ></RouterProvider>
    </Suspense>
  );
};

export default App;
