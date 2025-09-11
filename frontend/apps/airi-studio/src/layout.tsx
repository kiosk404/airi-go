import React from 'react';
import { Outlet } from 'react-router-dom';

export const Layout: React.FC<{ children?: React.ReactNode }> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="h-14 border-b bg-white flex items-center px-4">
        <div className="font-semibold text-primary-600">Airi Studio</div>
      </header>
      <div className="flex-1">
        {children}
        <Outlet />
      </div>
    </div>
  );
};

export default Layout;
