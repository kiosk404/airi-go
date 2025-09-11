import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import { pluginLess } from '@rsbuild/plugin-less';
import { pluginSass } from '@rsbuild/plugin-sass';
import { pluginSvgr } from '@rsbuild/plugin-svgr';

export const defineRsbuildConfig = (config: any = {}) => {
  return defineConfig({
    plugins: [
      pluginReact(),
      pluginLess(),
      pluginSass(),
      pluginSvgr(),
    ],
    source: {
      define: {
        'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'development'),
      },
    },
    tools: {
      postcss: (opts, { addPlugins }) => {
        // Add Tailwind CSS
        addPlugins([require('tailwindcss')]);
      },
    },
    html: {
      title: 'Airi Studio',
      favicon: './assets/favicon.png',
    },
    server: {
      port: 3000,
      host: 'localhost',
    },
    ...config,
  });
};






