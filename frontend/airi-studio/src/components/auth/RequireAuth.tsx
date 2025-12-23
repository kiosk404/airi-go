import { Navigate, useLocation } from 'react-router-dom';
import { Spin } from '@douyinfe/semi-ui';
import { useAuth } from '@/contexts/AuthContext';

interface RequireAuthProps {
  children: React.ReactNode;
}

export default function RequireAuth({ children }: RequireAuthProps) {
    const { isAuthenticated, isLoading } = useAuth();
    const location = useLocation();

   if (isLoading) {
       return (
           <div className="w-full h-screen flex justify-center items-center">
               <Spin size="large"/>
           </div>
       )
   }

   if (!isAuthenticated) {
       return <Navigate to="/login" state={{ from: location }} />;
   }

    return children;
}
