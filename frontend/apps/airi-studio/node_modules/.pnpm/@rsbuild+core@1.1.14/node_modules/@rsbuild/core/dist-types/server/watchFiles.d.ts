import type { FSWatcher } from '../../compiled/chokidar/index.js';
import type { ChokidarOptions, DevConfig, ServerConfig } from '../types';
import type { CompileMiddlewareAPI } from './getDevMiddlewares';
type WatchFilesOptions = {
    dev: DevConfig;
    server: ServerConfig;
    compileMiddlewareAPI?: CompileMiddlewareAPI;
    root: string;
};
export declare function setupWatchFiles(options: WatchFilesOptions): Promise<{
    close(): Promise<void>;
} | undefined>;
export declare function createChokidar(pathOrGlobs: string[], root: string, options: ChokidarOptions): Promise<FSWatcher>;
export {};
