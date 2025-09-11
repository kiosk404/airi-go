import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import { pluginLess } from '@rsbuild/plugin-less';
import { pluginSass } from '@rsbuild/plugin-sass';
import { pluginSvgr } from '@rsbuild/plugin-svgr';

export default defineConfig({
  plugins: [
    pluginReact(),
    pluginLess(),
    pluginSass(),
    pluginSvgr(),
  ],
  server: {
    port: 3000,
    proxy: [
      {
        context: ['/api'],
        target: 'http://localhost:8888/',
        secure: false,
        changeOrigin: true,
      },
    ],
  },
  html: {
    title: 'Airi Studio',
    template: './index.html',
  },
  source: {
    define: {
      'process.env.IS_REACT18': JSON.stringify(true),
    },
  },
});
