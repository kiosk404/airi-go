import * as __WEBPACK_EXTERNAL_MODULE_node_module__ from "node:module";
import * as __WEBPACK_EXTERNAL_MODULE_node_path__ from "node:path";
import * as __WEBPACK_EXTERNAL_MODULE_node_url__ from "node:url";
import * as __WEBPACK_EXTERNAL_MODULE_deepmerge__ from "deepmerge";
import * as __WEBPACK_EXTERNAL_MODULE_reduce_configs__ from "reduce-configs";
import * as __WEBPACK_EXTERNAL_MODULE__compiled_resolve_url_loader_index_js__ from "../compiled/resolve-url-loader/index.js";
let GLOBAL_PATCHED_SYMBOL = Symbol('GLOBAL_PATCHED_SYMBOL');
function patchGlobalLocation() {
    if (!global.location) {
        let location = Object.freeze({
            [GLOBAL_PATCHED_SYMBOL]: !0,
            href: (0, __WEBPACK_EXTERNAL_MODULE_node_url__.pathToFileURL)(process.cwd()).href + __WEBPACK_EXTERNAL_MODULE_node_path__.default.sep
        });
        global.location = location;
    }
}
function unpatchGlobalLocation() {
    var _global_location;
    (null === (_global_location = global.location) || void 0 === _global_location ? void 0 : _global_location[GLOBAL_PATCHED_SYMBOL]) && delete global.location;
}
let getResolveUrlJoinFn = async ()=>{
    let { createJoinFunction, asGenerator, createJoinImplementation, defaultJoinGenerator } = __WEBPACK_EXTERNAL_MODULE__compiled_resolve_url_loader_index_js__.default;
    return createJoinFunction('rsbuild-resolve-join-fn', createJoinImplementation(asGenerator((item, ...rest)=>item.uri.startsWith('.') ? defaultJoinGenerator(item, ...rest) : [
            null
        ])));
}, src_dirname = __WEBPACK_EXTERNAL_MODULE_node_path__.default.dirname((0, __WEBPACK_EXTERNAL_MODULE_node_url__.fileURLToPath)(import.meta.url)), src_require = (0, __WEBPACK_EXTERNAL_MODULE_node_module__.createRequire)(import.meta.url), PLUGIN_SASS_NAME = 'rsbuild:sass', getSassLoaderOptions = (userOptions, isUseCssSourceMap)=>{
    let excludes = [], mergedOptions = (0, __WEBPACK_EXTERNAL_MODULE_reduce_configs__.reduceConfigsWithContext)({
        initial: {
            sourceMap: isUseCssSourceMap,
            api: 'modern-compiler',
            implementation: src_require.resolve('sass-embedded'),
            sassOptions: {
                quietDeps: !0
            }
        },
        config: userOptions,
        ctx: {
            addExcludes: (items)=>{
                excludes.push(...Array.isArray(items) ? items : [
                    items
                ]);
            }
        },
        mergeFn: (defaults, userOptions)=>({
                ...defaults,
                ...userOptions,
                sassOptions: defaults.sassOptions && userOptions.sassOptions ? (0, __WEBPACK_EXTERNAL_MODULE_deepmerge__.default)(defaults.sassOptions, userOptions.sassOptions) : userOptions.sassOptions || defaults.sassOptions
            })
    });
    return mergedOptions.sassOptions ||= {}, !mergedOptions.sassOptions.silenceDeprecations && (mergedOptions.sassOptions.silenceDeprecations = [
        'import'
    ], 'legacy' === mergedOptions.api && mergedOptions.sassOptions.silenceDeprecations.push('legacy-js-api')), {
        options: mergedOptions,
        excludes
    };
}, findRuleId = (chain, defaultId)=>{
    let id = defaultId, index = 0;
    for(; chain.module.rules.has(id);)id = `${defaultId}-${++index}`;
    return id;
}, pluginSass = (pluginOptions = {})=>({
        name: PLUGIN_SASS_NAME,
        setup (api) {
            api.onAfterCreateCompiler(({ compiler })=>{
                var compiler1;
                (compiler1 = compiler).hooks.run.tap('PatchGlobalLocation', patchGlobalLocation), compiler1.hooks.watchRun.tap('PatchGlobalLocation', patchGlobalLocation), compiler1.hooks.watchClose.tap('PatchGlobalLocation', unpatchGlobalLocation), compiler1.hooks.done.tap('PatchGlobalLocation', unpatchGlobalLocation);
            }), api.modifyBundlerChain(async (chain, { CHAIN_ID })=>{
                let { excludes, options } = getSassLoaderOptions(pluginOptions.sassLoaderOptions, !0), ruleId = findRuleId(chain, CHAIN_ID.RULE.SASS), rule = chain.module.rule(ruleId).test(pluginOptions.include ?? /\.s(?:a|c)ss$/).merge({
                    sideEffects: !0
                }).resolve.preferRelative(!0).end();
                for (let item of excludes)rule.exclude.add(item);
                pluginOptions.exclude && rule.exclude.add(pluginOptions.exclude);
                let cssRule = chain.module.rules.get(CHAIN_ID.RULE.CSS);
                for (let id of (rule.dependency(cssRule.get('dependency')), Object.keys(cssRule.uses.entries()))){
                    let loader = cssRule.uses.get(id), options = loader.get('options') ?? {}, clonedOptions = (0, __WEBPACK_EXTERNAL_MODULE_deepmerge__.default)({}, options);
                    id === CHAIN_ID.USE.CSS && (clonedOptions.importLoaders += 2), rule.use(id).loader(loader.get('loader')).options(clonedOptions);
                }
                rule.use(CHAIN_ID.USE.RESOLVE_URL).loader(__WEBPACK_EXTERNAL_MODULE_node_path__.default.join(src_dirname, '../compiled/resolve-url-loader/index.js')).options({
                    join: await getResolveUrlJoinFn(),
                    sourceMap: !1
                }).end().use(CHAIN_ID.USE.SASS).loader(__WEBPACK_EXTERNAL_MODULE_node_path__.default.join(src_dirname, '../compiled/sass-loader/index.js')).options(options);
            });
        }
    });
export { PLUGIN_SASS_NAME, pluginSass };
