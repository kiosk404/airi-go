import type { JsDependency } from "@rspack/binding";
export declare class Dependency {
    #private;
    readonly type: string;
    readonly category: string;
    readonly request: string | undefined;
    critical: boolean;
    static __from_binding(binding: JsDependency): Dependency;
    static __to_binding(data: Dependency): JsDependency;
    private constructor();
}
