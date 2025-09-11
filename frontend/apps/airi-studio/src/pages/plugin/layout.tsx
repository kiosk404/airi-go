import React from 'react';
import { Outlet } from 'react-router-dom';

const PluginLayout: React.FC = () => {
  return (
    <div className="p-6">
      <h2 className="text-xl font-semibold mb-4">插件</h2>
      <Outlet />
    </div>
  );
};

export default PluginLayout;
