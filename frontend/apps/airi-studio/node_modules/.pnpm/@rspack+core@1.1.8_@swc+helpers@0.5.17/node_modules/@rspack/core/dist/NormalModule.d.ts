import * as liteTapable from "@rspack/lite-tapable";
import { Compilation } from "./Compilation";
import type { Module } from "./Module";
import type { LoaderContext } from "./config";
export declare class NormalModule {
    static getCompilationHooks(compilation: Compilation): {
        loader: liteTapable.SyncHook<[LoaderContext, Module]>;
        readResourceForScheme: any;
        readResource: liteTapable.HookMap<liteTapable.AsyncSeriesBailHook<[LoaderContext], string | Buffer>>;
    };
}
export default NormalModule;
