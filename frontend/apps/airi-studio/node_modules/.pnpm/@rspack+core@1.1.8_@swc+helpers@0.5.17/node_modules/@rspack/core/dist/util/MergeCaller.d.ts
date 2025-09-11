type CallFn<D> = (args: D[]) => void;
export default class MergeCaller<D> {
    private callArgs;
    private callFn;
    constructor(fn: CallFn<D>);
    private finalCall;
    push(...data: D[]): void;
}
export {};
