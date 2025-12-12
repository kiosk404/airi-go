import ReactDOM from 'react-dom/client';
import App from './App';
import 'tdesign-react/dist/tdesign.css'; // import TDesign styles
import '@douyinfe/semi-ui/dist/css/semi.css'; // import Semi UI styles

// 忽略 Semi UI 的警告
const originalConsoleError = console.error;
console.error = (...args) => {
  if (typeof args[0] === 'string' && args[0].includes('Failed prop type: The prop `action` is marked as required in `Upload`')) {
    return; // 忽略这个警告
  }
  originalConsoleError.apply(console, args);
};

const rootEl = document.getElementById('root');
if (rootEl) {
  const root = ReactDOM.createRoot(rootEl);
  root.render(
    <App />,
  );
}
