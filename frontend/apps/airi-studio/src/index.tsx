import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './app';
import './index.css';
import 'antd/dist/reset.css';

const main = () => {
    const $root = document.getElementById('root');
    if (!$root) {
        throw new Error('root element not found');
    }
    const root = ReactDOM.createRoot($root);

    root.render(<App />);
};

main();