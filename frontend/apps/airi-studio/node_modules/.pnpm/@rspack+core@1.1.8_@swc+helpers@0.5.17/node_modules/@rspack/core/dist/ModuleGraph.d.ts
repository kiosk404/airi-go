import type { JsModuleGraph } from "@rspack/binding";
import { Dependency } from "./Dependency";
import { Module } from "./Module";
export default class ModuleGraph {
    #private;
    static __from_binding(binding: JsModuleGraph): ModuleGraph;
    private constructor();
    getModule(dependency: Dependency): Module | null;
    getIssuer(module: Module): Module | null;
}
