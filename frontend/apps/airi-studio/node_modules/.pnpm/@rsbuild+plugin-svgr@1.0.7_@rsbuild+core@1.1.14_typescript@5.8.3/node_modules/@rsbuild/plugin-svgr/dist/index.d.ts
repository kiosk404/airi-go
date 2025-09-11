import type { RsbuildPlugin, Rspack } from '@rsbuild/core';
import type { Config as SvgrOptions } from '@svgr/core';
export type SvgDefaultExport = 'component' | 'url';
export type PluginSvgrOptions = {
    /**
     * Configure SVGR options.
     * @see https://react-svgr.com/docs/options/
     */
    svgrOptions?: SvgrOptions;
    /**
     * Whether to allow the use of default import and named import at the same time.
     * @default false
     */
    mixedImport?: boolean;
    /**
     * Custom query suffix to match SVGR transformation.
     * @default /react/
     */
    query?: RegExp;
    /**
     * Exclude some SVG files, they will not be transformed by SVGR.
     */
    exclude?: Rspack.RuleSetCondition;
    /**
     * Exclude some modules, the SVGs imported by these modules will not be transformed by SVGR.
     */
    excludeImporter?: Rspack.RuleSetCondition;
};
export declare const PLUGIN_SVGR_NAME = "rsbuild:svgr";
export declare const pluginSvgr: (options?: PluginSvgrOptions) => RsbuildPlugin;
