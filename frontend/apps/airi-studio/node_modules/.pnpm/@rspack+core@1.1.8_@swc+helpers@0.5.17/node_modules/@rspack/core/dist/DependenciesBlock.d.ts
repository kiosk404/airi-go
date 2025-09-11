import type { JsDependenciesBlock } from "@rspack/binding";
import { Dependency } from "./Dependency";
export declare class DependenciesBlock {
    #private;
    readonly dependencies: Dependency[];
    readonly blocks: DependenciesBlock[];
    static __from_binding(binding: JsDependenciesBlock): DependenciesBlock;
    static __to_binding(block: DependenciesBlock): JsDependenciesBlock;
    private constructor();
}
