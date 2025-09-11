import { BuiltinPluginName, type JsModule } from "@rspack/binding";
export declare const BuiltinLazyCompilationPlugin: {
    new (module: (args: {
        module: string;
        path: string;
    }) => {
        active: boolean;
        data: string;
        client: string;
    }, cacheable: boolean, entries: boolean, imports: boolean, test?: RegExp | ((m: JsModule) => boolean) | undefined): {
        name: BuiltinPluginName;
        _args: [module: (args: {
            module: string;
            path: string;
        }) => {
            active: boolean;
            data: string;
            client: string;
        }, cacheable: boolean, entries: boolean, imports: boolean, test?: RegExp | ((m: JsModule) => boolean) | undefined];
        affectedHooks: "done" | "make" | "compile" | "emit" | "afterEmit" | "invalid" | "thisCompilation" | "afterDone" | "compilation" | "normalModuleFactory" | "contextModuleFactory" | "initialize" | "shouldEmit" | "infrastructureLog" | "beforeRun" | "run" | "assetEmitted" | "failed" | "shutdown" | "watchRun" | "watchClose" | "environment" | "afterEnvironment" | "afterPlugins" | "afterResolvers" | "beforeCompile" | "afterCompile" | "finishMake" | "entryOption" | undefined;
        raw(compiler: import("../..").Compiler): import("@rspack/binding").BuiltinPlugin;
        apply(compiler: import("../..").Compiler): void;
    };
};
