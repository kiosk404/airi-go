import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import 'tdesign-react/dist/tdesign.css'; // import TDesign styles
import '@douyinfe/semi-ui/dist/css/semi.css'; // import Semi UI styles

const rootEl = document.getElementById('root');
if (rootEl) {
  const root = ReactDOM.createRoot(rootEl);
  root.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  );
}
