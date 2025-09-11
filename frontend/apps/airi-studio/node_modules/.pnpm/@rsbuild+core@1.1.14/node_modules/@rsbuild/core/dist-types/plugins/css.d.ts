import type { CSSLoaderOptions, RsbuildPlugin } from '../types';
export declare const normalizeCssLoaderOptions: (options: CSSLoaderOptions, exportOnlyLocals: boolean) => CSSLoaderOptions;
export declare const pluginCss: () => RsbuildPlugin;
