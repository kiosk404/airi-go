import React from 'react';
import { Link, Outlet, useLocation, useParams } from 'react-router-dom';
import classNames from 'classnames';

const SpaceLayout: React.FC = () => {
  const { space_id = 'airi' } = useParams();
  const location = useLocation();
  const base = `/space/${space_id}`;

  const tabs = [
    { name: 'Dashboard', path: `${base}/dashboard` },
    { name: 'Develop', path: `${base}/develop` },
    { name: 'Library', path: `${base}/library` },
  ];

  return (
    <div className="flex">
      <aside className="w-48 border-r bg-white">
        <div className="p-4 font-medium text-gray-700">空间：{space_id}</div>
        <nav className="space-y-1 p-2">
          {tabs.map(item => (
            <Link
              key={item.path}
              to={item.path}
              className={classNames(
                'block px-3 py-2 rounded-md text-sm',
                location.pathname === item.path
                  ? 'bg-primary-50 text-primary-700'
                  : 'text-gray-600 hover:bg-gray-50'
              )}
            >
              {item.name}
            </Link>
          ))}
        </nav>
      </aside>
      <section className="flex-1">
        <Outlet />
      </section>
    </div>
  );
};

export default SpaceLayout;





