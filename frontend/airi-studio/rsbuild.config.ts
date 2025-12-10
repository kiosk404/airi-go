import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import { pluginNodePolyfill } from '@rsbuild/plugin-node-polyfill';

export default defineConfig({
  plugins: [pluginReact(), pluginNodePolyfill()],
  html: {
    title: 'Airi Studio', // You can change this to your desired website title
    favicon: './public/favicon.ico',
  },
  server: {
    port: 3000,
  },
  output: {
    assetPrefix: '/',
  },
});