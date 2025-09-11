import type { IncomingMessage } from 'node:http';
import type { Socket } from 'node:net';
import type { EnvironmentContext, DevConfig as OriginDevConfig, Rspack, ServerConfig } from '../types';
import { type DevMiddlewareAPI } from './devMiddleware';
type Options = {
    publicPaths: string[];
    environments: Record<string, EnvironmentContext>;
    dev: OriginDevConfig;
    server: ServerConfig;
    compiler: Rspack.Compiler | Rspack.MultiCompiler;
};
/**
 * Setup compiler-related logic:
 * 1. setup rsbuild-dev-middleware
 * 2. establish webSocket connect
 */
export declare class CompilerDevMiddleware {
    middleware: DevMiddlewareAPI;
    private devConfig;
    private serverConfig;
    private compiler;
    private publicPaths;
    private socketServer;
    constructor({ dev, server, compiler, publicPaths, environments }: Options);
    init(): Promise<void>;
    upgrade(req: IncomingMessage, sock: Socket, head: any): void;
    close(): Promise<void>;
    sockWrite(type: string, data?: Record<string, any> | string | boolean): void;
    private setupDevMiddleware;
}
export {};
