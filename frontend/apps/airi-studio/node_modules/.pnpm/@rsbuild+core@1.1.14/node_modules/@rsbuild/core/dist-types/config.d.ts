import type { ChokidarOptions } from '../compiled/chokidar/index.js';
import type { InspectConfigOptions, InspectConfigResult, NormalizedConfig, PluginManager, PublicDir, PublicDirOptions, RsbuildConfig, RsbuildEntry } from './types';
export declare function getDefaultEntry(root: string): RsbuildEntry;
export declare const withDefaultConfig: (rootPath: string, config: RsbuildConfig) => Promise<RsbuildConfig>;
/** #__PURE__
 * 1. May used by multiple plugins.
 * 2. Object value that should not be empty.
 * 3. Meaningful and can be filled by constant value.
 */
export declare const normalizeConfig: (config: RsbuildConfig) => NormalizedConfig;
export type ConfigParams = {
    env: string;
    command: string;
    envMode?: string;
    meta?: Record<string, unknown>;
};
export type RsbuildConfigAsyncFn = (env: ConfigParams) => Promise<RsbuildConfig>;
export type RsbuildConfigSyncFn = (env: ConfigParams) => RsbuildConfig;
export type RsbuildConfigExport = RsbuildConfig | RsbuildConfigSyncFn | RsbuildConfigAsyncFn;
/**
 * This function helps you to autocomplete configuration types.
 * It accepts a Rsbuild config object, or a function that returns a config.
 */
export declare function defineConfig(config: RsbuildConfig): RsbuildConfig;
export declare function defineConfig(config: RsbuildConfigSyncFn): RsbuildConfigSyncFn;
export declare function defineConfig(config: RsbuildConfigAsyncFn): RsbuildConfigAsyncFn;
export declare function defineConfig(config: RsbuildConfigExport): RsbuildConfigExport;
export declare function watchFilesForRestart(files: string[], root: string, isBuildWatch: boolean, watchOptions?: ChokidarOptions): Promise<void>;
export type LoadConfigOptions = {
    /**
     * The root path to resolve the config file.
     * @default process.cwd()
     */
    cwd?: string;
    /**
     * The path to the config file, can be a relative or absolute path.
     * If `path` is not provided, the function will search for the config file in the `cwd`.
     */
    path?: string;
    /**
     * A custom meta object to be passed into the config function of `defineConfig`.
     */
    meta?: Record<string, unknown>;
    /**
     * The `envMode` passed into the config function of `defineConfig`.
     * @default process.env.NODE_ENV
     */
    envMode?: string;
};
export type LoadConfigResult = {
    /**
     * The loaded configuration object.
     */
    content: RsbuildConfig;
    /**
     * The path to the loaded configuration file.
     * Return `null` if the configuration file is not found.
     */
    filePath: string | null;
};
export declare function loadConfig({ cwd, path, envMode, meta, }?: LoadConfigOptions): Promise<LoadConfigResult>;
export declare const getRsbuildInspectConfig: ({ normalizedConfig, inspectOptions, pluginManager, }: {
    normalizedConfig: NormalizedConfig;
    inspectOptions: InspectConfigOptions;
    pluginManager: PluginManager;
}) => {
    rawRsbuildConfig: string;
    rsbuildConfig: InspectConfigResult["origin"]["rsbuildConfig"];
    rawEnvironmentConfigs: Array<{
        name: string;
        content: string;
    }>;
    environmentConfigs: InspectConfigResult["origin"]["environmentConfigs"];
};
export declare function outputInspectConfigFiles({ rawBundlerConfigs, rawEnvironmentConfigs, inspectOptions, configType, }: {
    configType: string;
    rawEnvironmentConfigs: Array<{
        name: string;
        content: string;
    }>;
    rawBundlerConfigs: Array<{
        name: string;
        content: string;
    }>;
    inspectOptions: InspectConfigOptions & {
        outputPath: string;
    };
}): Promise<void>;
export declare function stringifyConfig(config: unknown, verbose?: boolean): string;
export declare const normalizePublicDirs: (publicDir?: PublicDir) => Required<PublicDirOptions>[];
