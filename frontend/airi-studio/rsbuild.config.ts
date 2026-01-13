import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import { pluginNodePolyfill } from '@rsbuild/plugin-node-polyfill';

export default defineConfig({
  plugins: [pluginReact(), pluginNodePolyfill()],
  html: {
    title: 'Airi Studio', // You can change this to your desired website title
    favicon: './public/favicon.ico',
    tags: [
      {
        tag: 'link',
        attrs: {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap',
        }
      }
    ],
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: process.env.API_NASE_URL || 'http://localhost:9527',
        changeOrigin: true,
      },
    },
  },
  output: {
    assetPrefix: '/',
  },
});