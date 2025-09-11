import type { JsModule } from "@rspack/binding";
import type { Compiler } from "../..";
import { type LazyCompilationDefaultBackendOptions } from "./backend";
export default class LazyCompilationPlugin {
    cacheable: boolean;
    entries: boolean;
    imports: boolean;
    test?: RegExp | ((m: JsModule) => boolean);
    backend?: LazyCompilationDefaultBackendOptions;
    constructor(cacheable: boolean, entries: boolean, imports: boolean, test?: RegExp | ((m: JsModule) => boolean), backend?: LazyCompilationDefaultBackendOptions);
    apply(compiler: Compiler): void;
}
export { LazyCompilationPlugin };
