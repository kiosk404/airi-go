import type { ConfigChainWithContext, RsbuildPlugin, Rspack } from '@rsbuild/core';
export declare const isPlainObject: (obj: unknown) => obj is Record<string, any>;
export declare const PLUGIN_LESS_NAME = "rsbuild:less";
export type LessLoaderOptions = {
    lessOptions?: import('../compiled/less/index.js').default.Options;
    additionalData?: string | ((content: string, loaderContext: Rspack.LoaderContext<LessLoaderOptions>) => string | Promise<string>);
    sourceMap?: boolean;
    webpackImporter?: boolean;
    implementation?: unknown;
};
export type PluginLessOptions = {
    /**
     * Options passed to less-loader.
     * @see https://github.com/webpack-contrib/less-loader
     */
    lessLoaderOptions?: ConfigChainWithContext<LessLoaderOptions, {
        /**
         * @deprecated
         * use `exclude` option instead.
         */
        addExcludes: (items: string | RegExp | Array<string | RegExp>) => void;
    }>;
    /**
     * Include some `.less` files, they will be transformed by less-loader.
     * @default /\.less$/
     */
    include?: Rspack.RuleSetCondition;
    /**
     * Exclude some `.less` files, they will not be transformed by less-loader.
     * @default undefined
     */
    exclude?: Rspack.RuleSetCondition;
};
export declare const pluginLess: (pluginOptions?: PluginLessOptions) => RsbuildPlugin;
