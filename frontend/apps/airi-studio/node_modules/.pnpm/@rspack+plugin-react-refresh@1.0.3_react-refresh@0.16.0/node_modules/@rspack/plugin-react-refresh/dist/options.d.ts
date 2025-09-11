import type { RuleSetCondition } from '@rspack/core';
import type { IntegrationType } from './utils/getSocketIntegration';
interface OverlayOptions {
    entry: string;
    module: string;
    sockIntegration: IntegrationType | false;
    sockHost?: string;
    sockPath?: string;
    sockPort?: string;
    sockProtocol?: string;
}
export type PluginOptions = {
    /**
     * Include files to be processed by the plugin.
     * The value is the same as the `rule.test` option in Rspack.
     * @default /\.([cm]js|[jt]sx?|flow)$/i
     */
    include?: RuleSetCondition | null;
    /**
     * Exclude files from being processed by the plugin.
     * The value is the same as the `rule.exclude` option in Rspack.
     * @default /node_modules/
     */
    exclude?: RuleSetCondition | null;
    /**
     * Sets a namespace for the React Refresh runtime.
     * It is most useful when multiple instances of React Refresh is running
     * together simultaneously.
     * @default `output.uniqueName || output.library`
     */
    library?: string;
    /**
     * Whether to force enable the plugin.
     * By default, the plugin will not be enabled in non-development environments.
     * If you want to force enable the plugin, you can set this option to `true`.
     * @default false
     */
    forceEnable?: boolean;
    /**
     * Modify the behavior of the error overlay.
     * @default false
     */
    overlay?: boolean | Partial<OverlayOptions>;
};
export interface NormalizedPluginOptions extends Required<PluginOptions> {
    overlay: false | OverlayOptions;
}
export declare function normalizeOptions(options: PluginOptions): NormalizedPluginOptions;
export {};
